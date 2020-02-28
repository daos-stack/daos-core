/**
 * (C) Copyright 2019 Intel Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
 * The Government's rights to use, modify, reproduce, release, perform, display,
 * or disclose this software are subject to the terms of the Apache License as
 * provided in Contract No. B609815.
 * Any reproduction of computer software, computer software documentation, or
 * portions thereof marked with this legend must also reproduce the markings.
 */
/**
  * Implementation for aggregation and discard
 */
#define D_LOGFAC	DD_FAC(vos)

#include <daos_srv/vos.h>
#include <daos/object.h>	/* for daos_unit_oid_compare() */
#include <daos/checksum.h>
#include "vos_internal.h"
#include "evt_priv.h"

/*
 * EV tree sorted iterator returns logical entry in extent start order, and
 * the information like: physical entry it belongs to, visibility, is it the
 * last entry, etc. will be returned as well. One thing worth noting is that
 * sorted iterator is working on an out-of-band auxiliary index, so it won't
 * be affected by any in-tree modifications.
 *
 * EV tree aggregation is driven by sorted iterator: As the iterator moves on,
 * the visible logical entries will be queued to form a merge window, and all
 * the physical entries (no matter if it's fully covered or not) within the
 * window will be queued too. When the window size reaches certain threshold,
 * a procedure called merge window flush will be triggered to replace the old
 * physical entries by the new coalesced physical entries.
 *
 * On merge window flush, the queued visible logical entries will be used to
 * facilitate the data transfer from old physical entries to new coalesced
 * physical entries. If any old physical entry (not fully covered in current
 * window) straddles window end, it has to be head-truncated on window flush,
 * and the remaining part will be processed in next merge window.
 */

/* Merge window for evtree aggregation */
struct agg_merge_window {
	/* Record size */
	daos_size_t			 mw_rsize;
	/* Threshold for merge window flush */
	daos_size_t			 mw_flush_thresh;
	/* Merge window extent */
	struct evt_extent		 mw_ext;
	/* Physical entries in merge window */
	d_list_t			 mw_phy_ents;
	unsigned int			 mw_phy_cnt;
	/* Visible logical entries in merge window */
	struct vos_agg_lgc_ent		*mw_lgc_ents;
	unsigned int			 mw_lgc_max;
	unsigned int			 mw_lgc_cnt;
	/* I/O context for transfering data on flush */
	struct vos_agg_io_context	 mw_io_ctxt;
};

struct vos_agg_param {
	uint32_t		ap_credits_max; /* # of tight loops to yield */
	uint32_t		ap_credits;	/* # of tight loops */
	daos_handle_t		ap_coh;		/* container handle */
	daos_unit_oid_t		ap_oid;		/* current object ID */
	daos_key_t		ap_dkey;	/* current dkey */
	daos_key_t		ap_akey;	/* current akey */
	unsigned int		ap_discard:1;
	struct umem_instance	*ap_umm;
	/* SV tree: Max epoch in specified iterate epoch range */
	daos_epoch_t		 ap_max_epoch;
	/* EV tree: Merge window for evtree aggregation */
	struct agg_merge_window	 ap_window;
};

static bool cont_has_csums;
static bool unit_test;

static inline void
mark_yield(bio_addr_t *addr, unsigned int *acts)
{
	/*
	 * When read/write or reserve/delete a NVMe record, the BIO or VEA
	 * call might yield (BIO read/write yield and wait for NVMe DMA done,
	 * VEA reserve/free may trigger free extents reclaiming then yield
	 * and wait on blob unmap done).
	 *
	 * But we can't tell if it really yield or not (BIO read/write could
	 * skip DMA transfer on certain cases, free extents reclaiming isn't
	 * necessarily being triggered on every VEA call), to ensure the
	 * correctness, we always inform vos_iterate() yield, which may result
	 * in some unnecessary re-probe.
	 */
	if (addr->ba_type == DAOS_MEDIA_NVME)
		*acts |= VOS_ITER_CB_YIELD;
}

static int
agg_del_entry(daos_handle_t ih, struct umem_instance *umm,
	      vos_iter_entry_t *entry, unsigned int *acts)
{
	int	rc;

	D_ASSERT(umm != NULL);
	D_ASSERT(acts != NULL);

	rc = umem_tx_begin(umm, NULL);
	if (rc)
		return rc;

	mark_yield(&entry->ie_biov.bi_addr, acts);

	rc = vos_iter_delete(ih, NULL);
	if (rc != 0)
		rc = umem_tx_abort(umm, rc);
	else
		rc = umem_tx_commit(umm);

	if (rc) {
		D_ERROR("Failed to delete entry: "DF_RC"\n", DP_RC(rc));
		return rc;
	}

	*acts |= VOS_ITER_CB_DELETE;
	return rc;
}

static inline void
reset_agg_pos(vos_iter_type_t type, struct vos_agg_param *agg_param)
{
	switch (type) {
	case VOS_ITER_OBJ:
		memset(&agg_param->ap_oid, 0, sizeof(agg_param->ap_oid));
		break;
	case VOS_ITER_DKEY:
		memset(&agg_param->ap_dkey, 0, sizeof(agg_param->ap_dkey));
		break;
	case VOS_ITER_AKEY:
		memset(&agg_param->ap_akey, 0, sizeof(agg_param->ap_akey));
		break;
	default:
		break;
	}
}

static int
vos_agg_obj(daos_handle_t ih, vos_iter_entry_t *entry,
	    struct vos_agg_param *agg_param, unsigned int *acts)
{
	D_ASSERT(agg_param != NULL);
	if (daos_unit_oid_compare(agg_param->ap_oid, entry->ie_oid)) {
		agg_param->ap_oid = entry->ie_oid;
		reset_agg_pos(VOS_ITER_DKEY, agg_param);
		reset_agg_pos(VOS_ITER_AKEY, agg_param);
	} else {
		/*
		 * When recursive vos_iterate() yield in sub tree, re-probe
		 * is required when it returns back to upper level tree, if
		 * the just processed object is found on re-probe, we need
		 * to notify vos_iterate() to not iterate into to sub tree
		 * again.
		 */
		D_DEBUG(DB_EPC, "Skip oid:"DF_UOID" aggregation on re-probe\n",
			DP_UOID(agg_param->ap_oid));
		*acts |= VOS_ITER_CB_SKIP;
	}

	return 0;
}

static inline int
vos_agg_key_compare(daos_key_t key1, daos_key_t key2)
{
	if (key1.iov_len != key2.iov_len)
		return 1;

	return memcmp(key1.iov_buf, key2.iov_buf, key1.iov_len);
}

static int
vos_agg_dkey(daos_handle_t ih, vos_iter_entry_t *entry,
	     struct vos_agg_param *agg_param, unsigned int *acts)
{
	D_ASSERT(agg_param != NULL);
	if (vos_agg_key_compare(agg_param->ap_dkey, entry->ie_key)) {
		agg_param->ap_dkey = entry->ie_key;
		reset_agg_pos(VOS_ITER_AKEY, agg_param);
	} else {
		D_DEBUG(DB_EPC, "Skip dkey: "DF_KEY" aggregation on re-probe\n",
			DP_KEY(&entry->ie_key));
		*acts |= VOS_ITER_CB_SKIP;
	}

	return 0;
}

static inline bool
ext1_covers_ext2(struct evt_extent *ext1, struct evt_extent *ext2)
{
	D_ASSERT(ext1->ex_lo <= ext1->ex_hi);
	D_ASSERT(ext2->ex_lo <= ext2->ex_hi);

	return (ext1->ex_lo <= ext2->ex_lo && ext1->ex_hi >= ext2->ex_hi);
}

enum {
	MW_CLOSED	= 0,	/* Window closed, resource released */
	MW_FLUSHED,		/* Window flushed, no logical entries */
	MW_OPENED,		/* Window opened, has logical entries */
};

static int
merge_window_status(struct agg_merge_window *mw)
{
	struct vos_agg_io_context	*io = &mw->mw_io_ctxt;

	D_ASSERT(io->ic_seg_cnt == 0);
	D_ASSERT(io->ic_scm_cnt == 0);
	D_ASSERT(d_list_empty(&io->ic_nvme_exts));

	D_ASSERT(mw->mw_ext.ex_lo <= mw->mw_ext.ex_hi);

	if (mw->mw_lgc_cnt != 0) {
		D_ASSERT(mw->mw_rsize != 0);
		D_ASSERT(mw->mw_phy_cnt != 0);
		D_ASSERT(!d_list_empty(&mw->mw_phy_ents));

		return MW_OPENED;
	}

	D_ASSERT(mw->mw_ext.ex_lo == 0 && mw->mw_ext.ex_hi == 0);

	if (mw->mw_lgc_ents != NULL) {
		/*
		 * Even there if there isn't any logical entries,
		 * there could be some truncated physical entries.
		 */
		D_ASSERT(mw->mw_rsize != 0);
		return MW_FLUSHED;
	}

	/* Window closed, all resource should have been released */
	D_ASSERT(mw->mw_phy_cnt == 0);
	D_ASSERT(d_list_empty(&mw->mw_phy_ents));
	D_ASSERT(mw->mw_lgc_max == 0);

	D_ASSERT(io->ic_buf_len == 0);
	D_ASSERT(io->ic_buf == NULL);
	D_ASSERT(io->ic_seg_max == 0);
	D_ASSERT(io->ic_segs == NULL);
	D_ASSERT(io->ic_scm_max == 0);
	D_ASSERT(io->ic_scm_exts == NULL);

	return MW_CLOSED;
}

static int
vos_agg_akey(daos_handle_t ih, vos_iter_entry_t *entry,
	     struct vos_agg_param *agg_param, unsigned int *acts)
{
	D_ASSERT(agg_param != NULL);
	if (vos_agg_key_compare(agg_param->ap_akey, entry->ie_key)) {
		agg_param->ap_akey = entry->ie_key;
	} else {
		D_DEBUG(DB_EPC, "Skip akey: "DF_KEY" aggregation on re-probe\n",
			DP_KEY(&entry->ie_key));
		*acts |= VOS_ITER_CB_SKIP;
	}

	if (agg_param->ap_discard) {
		/* No merge window for discard path so bypass checks below. */
		return 0;
	}

	/* Reset the max epoch for low-level SV tree iteration */
	agg_param->ap_max_epoch = 0;
	/* The merge window for EV tree aggregation should have been closed */
	if (merge_window_status(&agg_param->ap_window) != MW_CLOSED)
		D_ASSERTF(false, "Merge window isn't closed.\n");

	return 0;
}

static int
vos_agg_sv(daos_handle_t ih, vos_iter_entry_t *entry,
	   struct vos_agg_param *agg_param, unsigned int *acts)
{
	int	rc;

	D_ASSERT(agg_param != NULL);
	D_ASSERT(entry->ie_epoch != 0);

	/* Discard */
	if (agg_param->ap_discard)
		goto delete;

	/* If entry is covered, the key or object is punched */
	if (entry->ie_vis_flags & VOS_VIS_FLAG_COVERED)
		goto delete;
	/*
	 * Aggregate: preserve the first recx which has highest epoch, because
	 * of re-probe, the highest epoch could be iterated multiple times.
	 */
	if (agg_param->ap_max_epoch == 0 ||
	    agg_param->ap_max_epoch == entry->ie_epoch) {
		agg_param->ap_max_epoch = entry->ie_epoch;
		return 0;
	}

	D_ASSERTF(entry->ie_epoch < agg_param->ap_max_epoch,
		  "max:"DF_U64", cur:"DF_U64"\n",
		  agg_param->ap_max_epoch, entry->ie_epoch);

delete:
	rc = agg_del_entry(ih, agg_param->ap_umm, entry, acts);
	if (rc) {
		D_ERROR("Failed to delete SV entry: "DF_RC"\n", DP_RC(rc));
	} else if (vos_iter_empty(ih) == 1 && agg_param->ap_discard) {
		/* Trigger re-probe in akey iteration */
		*acts |= VOS_ITER_CB_YIELD;
	}

	return rc;
}

static int
prepare_segments(struct agg_merge_window *mw)
{
	struct vos_agg_io_context	*io = &mw->mw_io_ctxt;
	struct vos_agg_phy_ent		*phy_ent;
	struct vos_agg_lgc_ent		*lgc_ent;
	struct vos_agg_lgc_seg		*lgc_seg;
	struct evt_entry_in		*ent_in;
	struct evt_extent		 ext;
	unsigned int			 i, seg_max, cs_total = 0;
	bool				 hole = false, coalesce;
	int				 rc = 0;

	/*
	 * Allocate large enough segments array to hold all the coalesced
	 * segments (at most mw_lgc_cnt) and truncated segments (at most
	 * mw_phy_cnt).
	 */
	D_ASSERT(mw->mw_lgc_cnt > 0);
	D_ASSERT(mw->mw_phy_cnt > 0);
	seg_max = MAX((mw->mw_lgc_cnt + mw->mw_phy_cnt), 200);
	if (io->ic_seg_max < seg_max) {
		D_REALLOC(lgc_seg, io->ic_segs, seg_max * sizeof(*lgc_seg));
		if (lgc_seg == NULL)
			return -DER_NOMEM;

		io->ic_segs = lgc_seg;
		io->ic_seg_max = seg_max;
	}
	memset(io->ic_segs, 0, io->ic_seg_max * sizeof(*lgc_seg));
	io->ic_seg_cnt = 0;

	/* Generate coalesced segments according to visible logical entries */
	for (i = 0; i < mw->mw_lgc_cnt; i++) {
		lgc_ent = &mw->mw_lgc_ents[i];
		phy_ent = lgc_ent->le_phy_ent;

		ext = lgc_ent->le_ext;
		D_ASSERT(ext1_covers_ext2(&mw->mw_ext, &ext));

		if (i == 0) {
			coalesce = false;
		} else if (hole != bio_addr_is_hole(&phy_ent->pe_addr)) {
			coalesce = false;
			io->ic_seg_cnt++;
			D_ASSERT(io->ic_seg_cnt < io->ic_seg_max);
		} else {
			coalesce = true;
		}

		hole = bio_addr_is_hole(&phy_ent->pe_addr);
		lgc_seg = &io->ic_segs[io->ic_seg_cnt];
		ent_in = &lgc_seg->ls_ent_in;

		if (!coalesce) {
			lgc_seg->ls_phy_ent = NULL;
			lgc_seg->ls_idx_start = i;
			ent_in->ei_inob = mw->mw_rsize;
			ent_in->ei_rect.rc_ex.ex_lo = ext.ex_lo;
			bio_addr_set_hole(&ent_in->ei_addr, hole);
			if (hole) {
				bio_addr_set(&ent_in->ei_addr, DAOS_MEDIA_SCM,
					     0);
				ent_in->ei_inob = 0;
			}
		} else {
			D_ASSERT(ext.ex_lo == ent_in->ei_rect.rc_ex.ex_hi + 1);
		}

		lgc_seg->ls_idx_end = i;
		ent_in->ei_rect.rc_ex.ex_hi = ext.ex_hi;
		/* Merge to highest epoch */
		if (ent_in->ei_rect.rc_epc < phy_ent->pe_rect.rc_epc)
			ent_in->ei_rect.rc_epc = phy_ent->pe_rect.rc_epc;
		/* Merge to highest pool map version */
		if (ent_in->ei_ver < phy_ent->pe_ver)
			ent_in->ei_ver = phy_ent->pe_ver;
		if (cont_has_csums && ent_in->ei_inob != 0)
			cs_total += vos_csum_prepare_ent(ent_in, phy_ent);
	}

	io->ic_seg_cnt++;
	D_ASSERT(io->ic_seg_cnt < io->ic_seg_max);

	/* Generate truncated segments according to physical entries */
	d_list_for_each_entry(phy_ent, &mw->mw_phy_ents, pe_link) {

		lgc_seg = &io->ic_segs[io->ic_seg_cnt];
		ent_in = &lgc_seg->ls_ent_in;

		ext = phy_ent->pe_rect.rc_ex;

		/* The physical entry was truncated on prev window flush */
		if (phy_ent->pe_off != 0)
			ext.ex_lo += phy_ent->pe_off;

		D_ASSERT(ext.ex_lo <= ext.ex_hi);
		D_ASSERT(ext.ex_lo <= mw->mw_ext.ex_hi);
		D_ASSERT(ext.ex_hi >= mw->mw_ext.ex_lo);

		/*
		 * Physical entry is in window, or it's fully covered (not
		 * visible) in current window.
		 */
		if (ext.ex_hi <= mw->mw_ext.ex_hi || phy_ent->pe_ref == 0)
			continue;

		lgc_seg->ls_phy_ent = phy_ent;
		lgc_seg->ls_idx_start = 0;
		lgc_seg->ls_idx_end = 0;

		ent_in->ei_inob = mw->mw_rsize;
		ent_in->ei_rect.rc_ex.ex_lo = mw->mw_ext.ex_hi + 1;
		ent_in->ei_rect.rc_ex.ex_hi = ext.ex_hi;
		ent_in->ei_rect.rc_epc = phy_ent->pe_rect.rc_epc;
		ent_in->ei_ver = phy_ent->pe_ver;

		hole = bio_addr_is_hole(&phy_ent->pe_addr);
		bio_addr_set_hole(&ent_in->ei_addr, hole);
		if (hole) {
			bio_addr_set(&ent_in->ei_addr, DAOS_MEDIA_SCM, 0);
			ent_in->ei_inob = 0;
		}

		if (cont_has_csums && ent_in->ei_inob != 0)
			cs_total += vos_csum_prepare_ent(ent_in, phy_ent);

		io->ic_seg_cnt++;
		D_ASSERT(io->ic_seg_cnt <= io->ic_seg_max);
	}
	if (cont_has_csums && cs_total) {
		rc = vos_csum_prepare_buf(io->ic_segs, io->ic_seg_cnt,
				      &io->ic_csum_buf, io->ic_csum_buf_len,
				      cs_total);
		io->ic_csum_buf_len += cs_total;
	}

	return rc;
}

static int
reserve_segment(struct vos_object *obj, struct vos_agg_io_context *io,
		daos_size_t size, bio_addr_t *addr)
{
	struct vea_space_info	*vsi;
	struct vea_hint_context	*hint_ctxt;
	struct vea_resrvd_ext	*nvme_ext;
	uint32_t		 blk_cnt;
	uint16_t		 media;
	int			 rc;

	memset(addr, 0, sizeof(*addr));
	media = vos_media_select(obj->obj_cont, DAOS_IOD_ARRAY, size);

	if (media == DAOS_MEDIA_SCM) {
		struct pobj_action	*scm_ext;
		umem_off_t		 umoff;

		D_ASSERT(io->ic_scm_max > io->ic_scm_cnt);
		D_ASSERT(io->ic_scm_exts != NULL);
		scm_ext = &io->ic_scm_exts[io->ic_scm_cnt];

		if (vos_obj2umm(obj)->umm_ops->mo_reserve != NULL)
			umoff = umem_reserve(vos_obj2umm(obj), scm_ext, size);
		else
			umoff = umem_alloc(vos_obj2umm(obj), size);

		if (UMOFF_IS_NULL(umoff)) {
			D_ERROR("Reserve "DF_U64" bytes on SCM failed.\n",
				size);
			return -DER_NOSPACE;
		}

		io->ic_scm_cnt++;
		bio_addr_set(addr, media, umoff);
		return 0;
	}

	D_ASSERT(media == DAOS_MEDIA_NVME);

	vsi = obj->obj_cont->vc_pool->vp_vea_info;
	D_ASSERT(vsi);
	hint_ctxt = obj->obj_cont->vc_hint_ctxt[VOS_IOS_AGGREGATION];
	D_ASSERT(hint_ctxt);
	blk_cnt = vos_byte2blkcnt(size);

	rc = vea_reserve(vsi, blk_cnt, hint_ctxt, &io->ic_nvme_exts);
	if (rc) {
		D_ERROR("Reserve %u blocks on NVMe failed: "DF_RC"\n", blk_cnt,
		DP_RC(rc));
		return rc;
	}

	nvme_ext = d_list_entry(io->ic_nvme_exts.prev, struct vea_resrvd_ext,
				vre_link);
	D_ASSERTF(nvme_ext->vre_blk_cnt == blk_cnt, "%u != %u\n",
		  nvme_ext->vre_blk_cnt, blk_cnt);
	D_ASSERT(nvme_ext->vre_blk_off != 0);

	bio_addr_set(addr, media, nvme_ext->vre_blk_off << VOS_BLK_SHIFT);
	return 0;
}

static inline daos_size_t
merge_window_size(struct agg_merge_window *mw)
{
	D_ASSERT(mw->mw_ext.ex_hi >= mw->mw_ext.ex_lo);
	D_ASSERT(mw->mw_rsize != 0);
	return evt_extent_width(&mw->mw_ext) * mw->mw_rsize;
}

static int
fill_one_segment(daos_handle_t ih, struct agg_merge_window *mw,
		 struct vos_agg_lgc_seg *lgc_seg, unsigned int *acts)
{
	struct vos_obj_iter		*oiter = vos_hdl2oiter(ih);
	struct vos_object		*obj = oiter->it_obj;
	struct vos_agg_io_context	*io = &mw->mw_io_ctxt;
	struct evt_entry_in		*ent_in = &lgc_seg->ls_ent_in;
	struct vos_agg_phy_ent		*phy_ent;
	struct bio_io_context		*bio_ctxt;
	struct bio_sglist		 bsgl;
	struct csum_recalc		*csum_recalcs = NULL;
	d_sg_list_t			 sgl;
	d_iov_t				 iov;
	bio_addr_t			 addr_dst, addr_src;
	daos_size_t			 seg_size, copy_size, buf_max;
	struct evt_extent		 ext = { 0 };
	daos_off_t			 phy_lo = 0;
	unsigned int			 i, seg_count, biov_idx = 0;
	size_t				 buf_add = 0;
	unsigned int			 added_csum_segs = 0;
	int				rc;

	D_ASSERT(obj != NULL);
	D_ASSERT(mw->mw_rsize > 0);

	if (bio_addr_is_hole(&ent_in->ei_addr))
		return 0;

	seg_size = evt_rect_width(&ent_in->ei_rect) * mw->mw_rsize;
	D_ASSERTF(seg_size > 0, "seg_size:"DF_U64"\n", seg_size);

	buf_max = MAX(seg_size, merge_window_size(mw));
	buf_max = MAX(buf_max, VOS_MW_FLUSH_THRESH);

	/* Copy data from old logical entries into new segment */
	D_ASSERT(lgc_seg->ls_idx_start <= lgc_seg->ls_idx_end);
	D_ASSERT(lgc_seg->ls_idx_end < mw->mw_lgc_cnt);

	bio_ctxt = obj->obj_cont->vc_pool->vp_io_ctxt;
	D_ASSERT(bio_ctxt != NULL);

	seg_count = lgc_seg->ls_idx_end - lgc_seg->ls_idx_start + 1;
	rc = bio_sgl_init(&bsgl, seg_count);
	if (rc) {
		D_ERROR("Init bsgl error: "DF_RC"\n", DP_RC(rc));
		return rc;
	}

	if (cont_has_csums) {
		D_ALLOC_ARRAY(csum_recalcs, seg_count);
		if (csum_recalcs == NULL) {
			rc = -DER_NOMEM;
			goto out;
		}
	}

	iov.iov_buf_len = buf_max; /* for sanity check */
	i = lgc_seg->ls_idx_start;
	while (i <= lgc_seg->ls_idx_end) {
		if (lgc_seg->ls_phy_ent != NULL) {
			phy_ent = lgc_seg->ls_phy_ent;
			ext = ent_in->ei_rect.rc_ex;
		} else {
			struct vos_agg_lgc_ent *lgc_ent = &mw->mw_lgc_ents[i];

			phy_ent = lgc_ent->le_phy_ent;
			ext = lgc_ent->le_ext;
		}
		i++;

		D_ASSERT(ext1_covers_ext2(&ent_in->ei_rect.rc_ex, &ext));
		D_ASSERT(ext1_covers_ext2(&phy_ent->pe_rect.rc_ex, &ext));

		phy_lo = phy_ent->pe_rect.rc_ex.ex_lo;
		if (phy_ent->pe_off != 0)
			phy_lo += phy_ent->pe_off;

		D_ASSERT(phy_lo <= phy_ent->pe_rect.rc_ex.ex_hi);
		D_ASSERT(ext.ex_lo >= phy_lo);

		copy_size = evt_extent_width(&ext) * ent_in->ei_inob;

		addr_src = phy_ent->pe_addr;
		addr_src.ba_off += (ext.ex_lo - phy_lo) * ent_in->ei_inob;

		D_ASSERT(!bio_addr_is_hole(&addr_src));
		D_ASSERT(iov.iov_buf_len >= copy_size);

		mark_yield(&addr_src, acts);
		D_ASSERT(biov_idx < bsgl.bs_nr);
		bio_iov_set(&bsgl.bs_iovs[biov_idx], addr_src, copy_size);

		if (cont_has_csums) {
			unsigned int wider = 0; /* length of per-ext csum add */
			unsigned int add_cnt =
				vos_csum_widen_biov(&bsgl.bs_iovs[biov_idx],
						    phy_ent, &ext, ent_in->ei_inob,
						    phy_lo, &wider);

			/* add_cnt is the number of data segments to add. */
			if (add_cnt) {
				buf_add += wider;
				iov.iov_buf_len += wider;
				copy_size += wider;
				added_csum_segs += add_cnt;
			}
			vos_csum_add_recalcs(&csum_recalcs, phy_ent, &ext,
					     &bsgl, biov_idx);
		}
		biov_idx++;
		D_ASSERT(iov.iov_buf_len >= copy_size);
		iov.iov_buf_len -= copy_size;
	}

	D_ASSERT(seg_size == buf_max - iov.iov_buf_len);

	/* Moved read buf allocation to after loop, to allow inclusion
	 * of additional data needed for verification of prior checksums.
	 */
	if (io->ic_buf_len < buf_max + buf_add) {
		void *buffer;

		D_REALLOC(buffer, io->ic_buf, buf_max + buf_add);
		if (buffer == NULL) {
			rc = -DER_NOMEM;
			goto out;
		}
		io->ic_buf = buffer;
		io->ic_buf_len = buf_max + buf_add;
	}

	if (added_csum_segs) {
		rc = vos_csum_append_added_segs(&bsgl,added_csum_segs);
		if (rc) {
			D_ERROR("Extend bsgl error: "DF_RC"\n", DP_RC(rc));
			goto out;
		}
	}

	iov.iov_buf = io->ic_buf;
	iov.iov_len = 0;
	iov.iov_buf_len = io->ic_buf_len;
	sgl.sg_nr = 1;
	sgl.sg_iovs = &iov;
	rc = bio_readv(bio_ctxt, &bsgl, &sgl);
	if (rc) {
		D_ERROR("Readv for "DF_RECT" error: "DF_RC"\n",
			DP_RECT(&ent_in->ei_rect), DP_RC(rc));
		goto out;
	}
	D_ASSERT(iov.iov_len == seg_size + buf_add);

	if (cont_has_csums) {
		rc = vos_csum_recalc(io, &bsgl, &sgl, ent_in, csum_recalcs,
				     seg_count, seg_size, unit_test);
		if (rc) {
			D_ERROR("CSUM verify error: "DF_RC"\n", DP_RC(rc));
			goto out;
		}
	}

	/* For csum support, this has moved reserve to after read, in case
	 * there's a csum mismatch on the verification of the read data
	 */
	rc = reserve_segment(obj, io, seg_size, &ent_in->ei_addr);
	if (rc) {
		D_ERROR("Reserve "DF_U64" segment error: "DF_RC"\n", seg_size,
			DP_RC(rc));
		goto out;
	}

	addr_dst = ent_in->ei_addr;
	D_ASSERT(!bio_addr_is_hole(&addr_dst));
	mark_yield(&addr_dst, acts);

	iov.iov_buf = io->ic_buf;
	iov.iov_buf_len = io->ic_buf_len;
	iov.iov_len = seg_size;
	rc = bio_write(bio_ctxt, addr_dst, &iov);
	if (rc)
		D_ERROR("Write "DF_RECT" error: "DF_RC"\n",
			DP_RECT(&ent_in->ei_rect), DP_RC(rc));
out:
	bio_sgl_fini(&bsgl);
	if (cont_has_csums)
		D_FREE(csum_recalcs);
	return rc;
}

static int
fill_segments(daos_handle_t ih, struct agg_merge_window *mw,
	      unsigned int *acts)
{
	struct vos_agg_io_context	*io = &mw->mw_io_ctxt;
	struct vos_agg_lgc_seg		*lgc_seg;
	unsigned int			 i, scm_max;
	int				 rc = 0;

	scm_max = MAX(io->ic_seg_cnt, 200);
	if (io->ic_scm_max < scm_max) {
		struct pobj_action *scm_exts;

		D_REALLOC(scm_exts, io->ic_scm_exts,
			  scm_max * sizeof(*scm_exts));
		if (scm_exts == NULL)
			return -DER_NOMEM;

		io->ic_scm_exts = scm_exts;
		io->ic_scm_max = scm_max;
	}
	memset(io->ic_scm_exts, 0, io->ic_scm_max * sizeof(*io->ic_scm_exts));
	D_ASSERT(io->ic_scm_cnt == 0);

	for (i = 0; i < io->ic_seg_cnt; i++) {
		lgc_seg = &io->ic_segs[i];

		D_DEBUG(DB_EPC, "Fill segment: %u-%u "DF_RECT"\n",
			lgc_seg->ls_idx_start, lgc_seg->ls_idx_end,
			DP_RECT(&lgc_seg->ls_ent_in.ei_rect));

		rc = fill_one_segment(ih, mw, lgc_seg, acts);
		if (rc) {
			D_ERROR("Fill seg %u-%u %p "DF_RECT" error: "DF_RC"\n",
				lgc_seg->ls_idx_start, lgc_seg->ls_idx_end,
				lgc_seg->ls_phy_ent,
				DP_RECT(&lgc_seg->ls_ent_in.ei_rect),
					DP_RC(rc));
			/* continue if -DER_CSUM */
			if (rc == -DER_CSUM)
				rc = 0;
			else
				break;
		}
	}

	return rc;
}

static int
insert_segments(daos_handle_t ih, struct agg_merge_window *mw,
		unsigned int *acts)
{
	struct vos_obj_iter		*oiter = vos_hdl2oiter(ih);
	struct vos_object		*obj = oiter->it_obj;
	struct vos_agg_io_context	*io = &mw->mw_io_ctxt;
	struct vos_agg_phy_ent		*phy_ent, *tmp;
	struct vos_agg_lgc_ent		*lgc_ent;
	struct vos_agg_lgc_seg		*lgc_seg;
	struct evt_entry_in		*ent_in;
	struct evt_rect			 rect;
	unsigned int			 i, leftovers = 0;
	int				 rc;

	D_ASSERT(obj != NULL);
	rc = umem_tx_begin(vos_obj2umm(obj), NULL);
	if (rc)
		return rc;

	/* Publish SCM reservations */
	if (io->ic_scm_cnt) {
		rc = umem_tx_publish(vos_obj2umm(obj), io->ic_scm_exts,
				     io->ic_scm_cnt);
		io->ic_scm_cnt = 0;
		if (rc) {
			D_ERROR("Publish %u SCM extents error: "DF_RC"\n",
				io->ic_scm_cnt, DP_RC(rc));
			goto abort;
		}
	}

	/* Adjust logical entry queue */
	for (i = 0; i < mw->mw_lgc_cnt; i++) {
		lgc_ent = &mw->mw_lgc_ents[i];
		phy_ent = lgc_ent->le_phy_ent;

		D_ASSERT(ext1_covers_ext2(&mw->mw_ext, &lgc_ent->le_ext));
		D_ASSERT(phy_ent->pe_ref > 0);
		phy_ent->pe_ref--;
		phy_ent->pe_trunc_head = true;
	}
	mw->mw_lgc_cnt = 0;

	/* Adjust payload address of truncated physical entries */
	for (i = 0; i < io->ic_seg_cnt; i++) {
		lgc_seg = &io->ic_segs[i];
		ent_in = &io->ic_segs[i].ls_ent_in;
		phy_ent = lgc_seg->ls_phy_ent;

		if (phy_ent != NULL && !bio_addr_is_hole(&ent_in->ei_addr)) {
			phy_ent->pe_addr = ent_in->ei_addr;
			phy_ent->pe_csum_info = ent_in->ei_csum;
		}
	}

	/* Remove old physical entries from EV tree */
	d_list_for_each_entry_safe(phy_ent, tmp, &mw->mw_phy_ents, pe_link) {
		rect = phy_ent->pe_rect;

		D_ASSERT(phy_ent->pe_ref == 0);
		/* The physical entry was truncated on prev window
		 * flush
		 */
		if (phy_ent->pe_off != 0) {
			rect.rc_ex.ex_lo += phy_ent->pe_off;
		}

		D_ASSERT(rect.rc_ex.ex_lo <= rect.rc_ex.ex_hi);
		D_ASSERT(rect.rc_ex.ex_lo <= mw->mw_ext.ex_hi);
		D_ASSERT(rect.rc_ex.ex_hi >= mw->mw_ext.ex_lo);

		/*
		 * The physical entry spans window end, but is fully
		 * covered in current window, keep it intact.
		 */
		if (rect.rc_ex.ex_hi > mw->mw_ext.ex_hi &&
		!phy_ent->pe_trunc_head) {
			leftovers++;
			continue;
		}
		if (phy_ent->pe_csum_info.cs_not_valid)
			continue;

		mark_yield(&phy_ent->pe_addr, acts);

		rc = evt_delete(oiter->it_hdl, &rect, NULL);
		if (rc) {
			D_ERROR("Delete "DF_RECT" pe_off:"DF_U64" error: "
				""DF_RC"\n", DP_RECT(&rect),
				phy_ent->pe_off,
				DP_RC(rc));
			goto abort;
		}

		/* Physical entry is in window */
		if (rect.rc_ex.ex_hi <= mw->mw_ext.ex_hi) {
			d_list_del(&phy_ent->pe_link);
			D_FREE_PTR(phy_ent);
			D_ASSERT(mw->mw_phy_cnt > 0);
			mw->mw_phy_cnt--;
			continue;
		}

		/* Update extent start of truncated physical entry.
		 */

		rect.rc_ex.ex_lo = mw->mw_ext.ex_hi + 1;
		phy_ent->pe_off = rect.rc_ex.ex_lo -
				phy_ent->pe_rect.rc_ex.ex_lo;
		phy_ent->pe_trunc_head = false;
		leftovers++;
	}
	D_ASSERT(leftovers == mw->mw_phy_cnt);

	/* Clear window size */
	mw->mw_ext.ex_lo = mw->mw_ext.ex_hi = 0;

	/* Insert new segments into EV tree */
	for (i = 0; i < io->ic_seg_cnt; i++) {
		if (io->ic_segs[i].ls_has_csum_err)
			continue;
		ent_in = &io->ic_segs[i].ls_ent_in;

		rc = evt_insert(oiter->it_hdl, ent_in);
		if (rc) {
			D_ERROR("Insert segment "DF_RECT" error: "DF_RC"\n",
				DP_RECT(&ent_in->ei_rect), DP_RC(rc));
			goto abort;
		}
	}

	/* Publish NVMe reservations */
	rc = vos_publish_blocks(obj->obj_cont, &io->ic_nvme_exts, true,
				VOS_IOS_AGGREGATION);
	if (rc) {
		D_ERROR("Publish NVMe extents error: "DF_RC"\n", DP_RC(rc));
		goto abort;
	}
abort:
	if (rc)
		rc = umem_tx_abort(vos_obj2umm(obj), rc);
	else
		rc = umem_tx_commit(vos_obj2umm(obj));

	return rc;
}

static void
cleanup_segments(daos_handle_t ih, struct agg_merge_window *mw, int rc)
{
	struct vos_obj_iter		*oiter = vos_hdl2oiter(ih);
	struct vos_object		*obj = oiter->it_obj;
	struct vos_agg_io_context	*io = &mw->mw_io_ctxt;

	D_ASSERT(obj != NULL);
	if (rc) {
		if (io->ic_scm_cnt) {
			umem_cancel(vos_obj2umm(obj), io->ic_scm_exts,
				    io->ic_scm_cnt);
			io->ic_scm_cnt = 0;
		}
		if (!d_list_empty(&io->ic_nvme_exts))
			vos_publish_blocks(obj->obj_cont, &io->ic_nvme_exts,
					   false, VOS_IOS_AGGREGATION);
	}

	/* Reset io context */
	D_ASSERT(d_list_empty(&io->ic_nvme_exts));
	D_ASSERT(io->ic_scm_cnt == 0);
	io->ic_seg_cnt = 0;
}

static void
clear_merge_window(struct agg_merge_window *mw)
{
	struct vos_agg_phy_ent *phy_ent, *tmp;

	mw->mw_ext.ex_lo = mw->mw_ext.ex_hi = 0;
	mw->mw_lgc_cnt = 0;
	d_list_for_each_entry_safe(phy_ent, tmp, &mw->mw_phy_ents,
				   pe_link) {
		d_list_del(&phy_ent->pe_link);
		D_FREE_PTR(phy_ent);
	}
	mw->mw_phy_cnt = 0;
}

static bool
need_flush(struct agg_merge_window *mw)
{
	struct vos_agg_phy_ent	*phy_ent;
	struct vos_agg_lgc_ent	*lgc_ent;
	struct evt_extent	 lgc_ext, phy_ext;
	int			 i;
	bool			 hole = false;

	for (i = 0; i < mw->mw_lgc_cnt; i++) {
		lgc_ent = &mw->mw_lgc_ents[i];
		phy_ent = lgc_ent->le_phy_ent;
		phy_ext = phy_ent->pe_rect.rc_ex;
		lgc_ext = lgc_ent->le_ext;

		/*
		 * Any physical entry is partially covered, or appeared
		 * in other window
		 */
		if (lgc_ext.ex_lo != phy_ext.ex_lo ||
		    lgc_ext.ex_hi != phy_ext.ex_hi)
			return true;

		/* If any consecutive visible entries can be merged */
		if (i != 0 && hole == bio_addr_is_hole(&phy_ent->pe_addr))
			return true;

		hole = bio_addr_is_hole(&phy_ent->pe_addr);
	}

	/* Any invisible physical entries ? */
	if (mw->mw_lgc_cnt != mw->mw_phy_cnt)
		return true;

	clear_merge_window(mw);
	D_DEBUG(DB_EPC, "Skip window flush "DF_EXT"\n", DP_EXT(&mw->mw_ext));

	return false;
}

static int
flush_merge_window(daos_handle_t ih, struct agg_merge_window *mw,
		   unsigned int *acts)
{
	int	rc;

	/*
	 * If no new updates in an already aggregated window, window flush will
	 * be skipped, otherwise, all the data within the window will be
	 * migrated to a new location, such batch data migration is good for
	 * anti-fragmentaion.
	 */
	if (!need_flush(mw))
		return 0;

	/* Prepare the new segments to be inserted */
	rc = prepare_segments(mw);
	if (rc) {
		D_ERROR("Prepare segments "DF_EXT" error: "DF_RC"\n",
			DP_EXT(&mw->mw_ext), DP_RC(rc));
		goto out;
	}

	/* Transfer data from old logical records to reserved new segments */
	rc = fill_segments(ih, mw, acts);
	if (rc) {
		D_ERROR("Fill segments "DF_EXT" error: "DF_RC"\n",
			DP_EXT(&mw->mw_ext), DP_RC(rc));
		goto out;
	}

	/* Replace the old logical records with new segments in EV tree */
	rc = insert_segments(ih, mw, acts);
	if (rc) {
		D_ERROR("Insert segments "DF_EXT" error: "DF_RC"\n",
			DP_EXT(&mw->mw_ext), DP_RC(rc));
		goto out;
	}
out:
	cleanup_segments(ih, mw, rc);
	return rc;
}

static bool
trigger_flush(struct agg_merge_window *mw, struct evt_extent *lgc_ext)
{
	struct evt_extent *w_ext = &mw->mw_ext;

	D_ASSERT(w_ext->ex_lo <= lgc_ext->ex_lo);
	/* Empty or closed merge window */
	if (merge_window_status(mw) == MW_CLOSED ||
	    merge_window_status(mw) == MW_FLUSHED)
		return false;

	/*
	 * Window is formed by visible logical entries, must have no
	 * overlapping.
	 */
	D_ASSERTF(w_ext->ex_hi < lgc_ext->ex_lo,
		  "win:"DF_EXT", lgc_ent:"DF_EXT"\n",
		  DP_EXT(w_ext), DP_EXT(lgc_ext));

	/* Window is large enough */
	if (merge_window_size(mw) >= mw->mw_flush_thresh)
		return true;

	/* Trigger flush when entry is disjoint with window */
	return !((w_ext->ex_hi + 1) == lgc_ext->ex_lo);
}

static struct vos_agg_phy_ent *
enqueue_phy_ent(struct agg_merge_window *mw, struct evt_extent *phy_ext,
		daos_epoch_t epoch, bio_addr_t *addr,
		struct dcs_csum_info *csum_info, uint32_t ver)
{
	struct vos_agg_phy_ent *phy_ent;

	D_ALLOC_PTR(phy_ent);
	if (phy_ent == NULL)
		return NULL;

	phy_ent->pe_rect.rc_ex = *phy_ext;
	phy_ent->pe_rect.rc_epc = epoch;
	phy_ent->pe_addr = *addr;
	phy_ent->pe_csum_info = *csum_info;
	phy_ent->pe_off = 0;
	phy_ent->pe_ver = ver;
	phy_ent->pe_ref = 0;

	/* Sanity check */
	if (!d_list_empty(&mw->mw_phy_ents)) {
		struct vos_agg_phy_ent *prev;

		D_ASSERT(mw->mw_phy_cnt != 0);
		prev = d_list_entry(mw->mw_phy_ents.prev,
				    struct vos_agg_phy_ent, pe_link);
		D_ASSERTF(prev->pe_rect.rc_ex.ex_lo <= phy_ext->ex_lo,
			  "prev phy_ext: "DF_EXT", phy_ext: "DF_EXT"\n",
			  DP_EXT(&prev->pe_rect.rc_ex), DP_EXT(phy_ext));
	} else {
		D_ASSERT(mw->mw_phy_cnt == 0);
	}

	d_list_add_tail(&phy_ent->pe_link, &mw->mw_phy_ents);
	mw->mw_phy_cnt++;

	return phy_ent;
}

static int
enqueue_lgc_ent(struct agg_merge_window *mw, struct evt_extent *lgc_ext,
		struct vos_agg_phy_ent *phy_ent)
{
	struct vos_agg_lgc_ent	*lgc_ent;
	unsigned int		 max, cnt;

	max = mw->mw_lgc_max;
	cnt = mw->mw_lgc_cnt;
	/* Sanity check */
	if (cnt > 0) {
		lgc_ent = &mw->mw_lgc_ents[cnt - 1];
		D_ASSERTF(lgc_ext->ex_lo == lgc_ent->le_ext.ex_hi + 1 &&
			  lgc_ent->le_ext.ex_hi == mw->mw_ext.ex_hi,
			  "prev lgc_ext: "DF_EXT", lgc_ext: "DF_EXT"\n",
			  DP_EXT(&lgc_ent->le_ext), DP_EXT(lgc_ext));
	}

	if (cnt == max) {
		unsigned int new_max = max ? max * 2 : 10;

		D_REALLOC(lgc_ent, mw->mw_lgc_ents,
			  new_max * sizeof(*lgc_ent));
		if (lgc_ent == NULL)
			return -DER_NOMEM;

		mw->mw_lgc_max = new_max;
		mw->mw_lgc_ents = lgc_ent;
	}

	D_ASSERT(mw->mw_lgc_max > mw->mw_lgc_cnt);
	lgc_ent = &mw->mw_lgc_ents[cnt];
	lgc_ent->le_ext = *lgc_ext;
	phy_ent->pe_ref++;
	lgc_ent->le_phy_ent = phy_ent;
	mw->mw_lgc_cnt++;

	/*
	 * Extend window size. If the visible entry is a punched record, the
	 * window size could be very huge, but this is ok, because there won't
	 * be huge contiguous allocation on window flush, only lots of covered
	 * physical entries being deleted.
	 */
	if (mw->mw_lgc_cnt == 1)
		mw->mw_ext.ex_lo = lgc_ext->ex_lo;
	mw->mw_ext.ex_hi = lgc_ext->ex_hi;

	D_DEBUG(DB_EPC, "lgc_ext:"DF_EXT", phy_ext:"DF_RECT", mw:"DF_EXT", "
		"index:%u\n", DP_EXT(lgc_ext), DP_RECT(&phy_ent->pe_rect),
		DP_EXT(&mw->mw_ext), cnt);

	return 0;
}

static void
close_merge_window(struct agg_merge_window *mw, int rc)
{
	struct vos_agg_io_context *io = &mw->mw_io_ctxt;

	if (rc)
		clear_merge_window(mw);

	D_ASSERT(merge_window_status(mw) != MW_OPENED);

	mw->mw_rsize = 0;
	if (mw->mw_lgc_ents != NULL) {
		D_FREE(mw->mw_lgc_ents);
		mw->mw_lgc_ents = NULL;
		mw->mw_lgc_max = 0;
	}

	if (io->ic_buf != NULL) {
		D_FREE(io->ic_buf);
		io->ic_buf = NULL;
		io->ic_buf_len = 0;
	}

	if (io->ic_segs != NULL) {
		D_FREE(io->ic_segs);
		io->ic_segs = NULL;
		io->ic_seg_max = 0;
	}

	if (io->ic_scm_exts != NULL) {
		D_FREE(io->ic_scm_exts);
		io->ic_scm_exts = NULL;
		io->ic_scm_max = 0;
	}
}

static inline void
recx2ext(daos_recx_t *recx, struct evt_extent *ext)
{
	D_ASSERT(recx->rx_nr > 0);
	ext->ex_lo = recx->rx_idx;
	ext->ex_hi = recx->rx_idx + recx->rx_nr - 1;
}

static struct vos_agg_phy_ent *
lookup_phy_ent(struct agg_merge_window *mw, const struct evt_extent *phy_ext,
	       daos_epoch_t epoch)
{
	struct vos_agg_phy_ent *phy_ent;

	d_list_for_each_entry_reverse(phy_ent, &mw->mw_phy_ents, pe_link) {
		/* Physical entry list is sorted by extent start */
		if (phy_ent->pe_rect.rc_ex.ex_lo < phy_ext->ex_lo)
			break;

		if (phy_ent->pe_rect.rc_epc == epoch &&
		    phy_ent->pe_rect.rc_ex.ex_hi == phy_ext->ex_hi)
			return phy_ent;
	}

	return NULL;
}

static int
join_merge_window(daos_handle_t ih, struct agg_merge_window *mw,
		  vos_iter_entry_t *entry, unsigned int *acts)
{
	struct vos_obj_iter	*oiter = vos_hdl2oiter(ih);
	struct evt_extent	 phy_ext, lgc_ext;
	struct vos_agg_phy_ent	*phy_ent;
	bool			 visible, partial, last;
	int			 rc = 0;

	recx2ext(&entry->ie_recx, &lgc_ext);
	recx2ext(&entry->ie_orig_recx, &phy_ext);
	D_ASSERT(ext1_covers_ext2(&phy_ext, &lgc_ext));

	visible = (entry->ie_vis_flags & VOS_VIS_FLAG_VISIBLE);
	partial = (entry->ie_vis_flags & VOS_VIS_FLAG_PARTIAL);
	last = (entry->ie_vis_flags & VOS_VIS_FLAG_LAST);

	/* Just delete the fully covered intact physical entry */
	if (!visible && !partial) {
		struct evt_rect rect;

		D_ASSERTF(lgc_ext.ex_lo == phy_ext.ex_lo &&
			  lgc_ext.ex_hi == phy_ext.ex_hi,
			  ""DF_EXT" != "DF_EXT"\n",
			  DP_EXT(&lgc_ext), DP_EXT(&phy_ext));
		D_ASSERT(entry->ie_vis_flags & VOS_VIS_FLAG_COVERED);

		rect.rc_ex = phy_ext;
		rect.rc_epc = entry->ie_epoch;
		mark_yield(&entry->ie_biov.bi_addr, acts);

		rc = evt_delete(oiter->it_hdl, &rect, NULL);
		if (rc) {
			D_ERROR("Delete EV entry "DF_RECT" error: "DF_RC"\n",
				DP_RECT(&rect), DP_RC(rc));
			return rc;
		}

		goto out;
	}

	/* Trigger current window flush when reaching threshold */
	if (visible && trigger_flush(mw, &lgc_ext)) {
		rc = flush_merge_window(ih, mw, acts);
		if (rc) {
			D_ERROR("Flush window "DF_EXT" error: "DF_RC"\n",
				DP_EXT(&mw->mw_ext), DP_RC(rc));
			return rc;
		}
		D_ASSERT(merge_window_status(mw) == MW_FLUSHED);
	}

	/* Lookup physical entry, enqueue if it doesn't exist */
	phy_ent = lookup_phy_ent(mw, &phy_ext, entry->ie_epoch);
	if (phy_ent == NULL) {
		D_ASSERT(phy_ext.ex_lo == lgc_ext.ex_lo);

		phy_ent = enqueue_phy_ent(mw, &phy_ext, entry->ie_epoch,
					  &entry->ie_biov.bi_addr,
					  &entry->ie_csum, entry->ie_ver);
		if (phy_ent == NULL) {
			rc = -DER_NOMEM;
			D_ERROR("Enqueue phy_ent win:"DF_EXT", ent:"DF_EXT" "
				"error: "DF_RC"\n", DP_EXT(&mw->mw_ext),
				DP_EXT(&phy_ext), DP_RC(rc));
			return rc;
		}
	} else {
		/* Can't be the first logcial entry */
		D_ASSERT(phy_ext.ex_lo != lgc_ext.ex_lo);
	}

	/* Enqueue the visible logical entry */
	if (visible) {
		rc = enqueue_lgc_ent(mw, &lgc_ext, phy_ent);
		if (rc) {
			D_ERROR("Enqueue lgc_ent win: "DF_EXT", ent:"DF_EXT" "
				"error: "DF_RC"\n", DP_EXT(&mw->mw_ext),
				DP_EXT(&lgc_ext), DP_RC(rc));
			return rc;
		}
	} else {
		/* Fully covered physical entry must have been deleted */
		D_ASSERT(partial);
	}
out:
	/* Flush & close window on last entry */
	if (last) {
		rc = flush_merge_window(ih, mw, acts);
		if (rc)
			D_ERROR("Flush window "DF_EXT" error: "DF_RC"\n",
				DP_EXT(&mw->mw_ext), DP_RC(rc));
		close_merge_window(mw, rc);
	}

	return rc;
}

static int
set_window_size(struct agg_merge_window *mw, daos_size_t rsize)
{
	int	rc = 0;

	if (rsize == 0) {
		D_CRIT("EV tree 0 iod_size could be caused by inserting "
		       "punch records in an empty tree, this will be "
		       "disallowed in the future.\n");
		rc = -DER_INVAL;
	} else if (mw->mw_rsize == 0) {
		mw->mw_rsize = rsize;

		if (DAOS_FAIL_CHECK(DAOS_VOS_AGG_MW_THRESH)) {
			mw->mw_flush_thresh = daos_fail_value_get();
			D_INFO("Set flush threshold to: "DF_U64"\n",
			       mw->mw_flush_thresh);
		} else if (rsize < (VOS_MW_FLUSH_THRESH / 2)) {
			mw->mw_flush_thresh = VOS_MW_FLUSH_THRESH;
		} else {
			mw->mw_flush_thresh = (rsize < VOS_MW_FLUSH_THRESH) ?
						rsize * 2 : rsize;
			D_INFO("Bump flush threshold to: "DF_U64", rsize: "
			       ""DF_U64"\n", mw->mw_flush_thresh, rsize);
		}
	} else if (mw->mw_rsize != rsize) {
		D_CRIT("Mismatched iod_size "DF_U64" != "DF_U64"\n",
		       mw->mw_rsize, rsize);
		rc = -DER_INVAL;
	}

	return rc;
}

static int
vos_agg_ev(daos_handle_t ih, vos_iter_entry_t *entry,
	   struct vos_agg_param *agg_param, unsigned int *acts)
{
	struct agg_merge_window	*mw = &agg_param->ap_window;
	struct evt_extent	 phy_ext, lgc_ext;
	int			 rc = 0;

	D_ASSERT(agg_param != NULL);
	D_ASSERT(acts != NULL);
	recx2ext(&entry->ie_recx, &lgc_ext);
	recx2ext(&entry->ie_orig_recx, &phy_ext);

	/* Discard */
	if (agg_param->ap_discard) {
		struct vos_obj_iter	*oiter = vos_hdl2oiter(ih);
		struct evt_rect		 rect;

		/*
		 * Delete the physical entry when iterating to the first
		 * logical entry
		 */
		if (phy_ext.ex_lo == lgc_ext.ex_lo) {
			rect.rc_ex = phy_ext;
			rect.rc_epc = entry->ie_epoch;
			mark_yield(&entry->ie_biov.bi_addr, acts);

			rc = evt_delete(oiter->it_hdl, &rect, NULL);
			if (rc)
				D_ERROR("Delete EV entry "DF_RECT" error: "
					""DF_RC"\n", DP_RECT(&rect), DP_RC(rc));
		}

		/*
		 * Sorted iteration doesn't support tree empty check, so we
		 * always inform vos_iterate() to check if subtree is empty.
		 */
		if (entry->ie_vis_flags & VOS_VIS_FLAG_LAST) {
			/* Trigger re-probe in akey iteration */
			*acts |= VOS_ITER_CB_YIELD;
		}
		return rc;
	}

	/* Aggregation */
	D_DEBUG(DB_EPC, "oid:"DF_UOID", lgc_ext:"DF_EXT", "
		"phy_ext:"DF_EXT", epoch:"DF_U64", flags: %x\n",
		DP_UOID(agg_param->ap_oid), DP_EXT(&lgc_ext),
		DP_EXT(&phy_ext), entry->ie_epoch, entry->ie_vis_flags);

	rc = set_window_size(mw, entry->ie_rsize);
	if (rc)
		goto out;

	rc = join_merge_window(ih, mw, entry, acts);
	if (rc)
		D_ERROR("Join window "DF_EXT"/"DF_EXT" error: "DF_RC"\n",
			DP_EXT(&mw->mw_ext), DP_EXT(&phy_ext), DP_RC(rc));
out:
	if (rc)
		close_merge_window(mw, rc);
	return rc;
}

static int
vos_aggregate_pre_cb(daos_handle_t ih, vos_iter_entry_t *entry,
		     vos_iter_type_t type, vos_iter_param_t *param,
		     void *cb_arg, unsigned int *acts)
{
	struct vos_agg_param	*agg_param = cb_arg;
	struct vos_container	*cont;
	int			 rc;

	cont = vos_hdl2cont(param->ip_hdl);
	D_DEBUG(DB_EPC, DF_CONT": Aggregate pre, type:%d, is_discard:%d\n",
		DP_CONT(cont->vc_pool->vp_id, cont->vc_id), type,
		agg_param->ap_discard);

	switch (type) {
	case VOS_ITER_OBJ:
		rc = vos_agg_obj(ih, entry, agg_param, acts);
		break;
	case VOS_ITER_DKEY:
		rc = vos_agg_dkey(ih, entry, agg_param, acts);
		break;
	case VOS_ITER_AKEY:
		rc = vos_agg_akey(ih, entry, agg_param, acts);
		break;
	case VOS_ITER_SINGLE:
		rc = vos_agg_sv(ih, entry, agg_param, acts);
		break;
	case VOS_ITER_RECX:
		rc = vos_agg_ev(ih, entry, agg_param, acts);
		break;
	default:
		D_ASSERTF(false, "Invalid iter type\n");
		rc = -DER_INVAL;
		break;
	}

	if (rc < 0) {
		D_ERROR("VOS aggregation failed: "DF_RC"\n", DP_RC(rc));
		return rc;
	}

	if (cont->vc_abort_aggregation) {
		D_DEBUG(DB_EPC, "VOS aggregation aborted\n");
		return 1;
	}

	agg_param->ap_credits++;

	if (agg_param->ap_credits > agg_param->ap_credits_max ||
	    (DAOS_FAIL_CHECK(DAOS_VOS_AGG_RANDOM_YIELD) && (rand() % 2))) {
		agg_param->ap_credits = 0;
		*acts |= VOS_ITER_CB_YIELD;

		/*
		 * Reset position if we yield while iterating in object, dkey
		 * or akey level, so that subtree won't be skipped mistakenly,
		 * see the comment in vos_agg_obj().
		 */
		reset_agg_pos(type, agg_param);
		bio_yield();
	}

	return 0;
}

static int
vos_aggregate_post_cb(daos_handle_t ih, vos_iter_entry_t *entry,
		      vos_iter_type_t type, vos_iter_param_t *param,
		      void *cb_arg, unsigned int *acts)
{
	struct vos_agg_param	*agg_param = cb_arg;
	struct vos_container	*cont;
	int			 rc = 0;

	cont = vos_hdl2cont(param->ip_hdl);
	D_DEBUG(DB_EPC, DF_CONT": Aggregate post, type:%d, is_discard:%d\n",
		DP_CONT(cont->vc_pool->vp_id, cont->vc_id), type,
		agg_param->ap_discard);

	switch (type) {
	case VOS_ITER_OBJ:
		rc = oi_iter_aggregate(ih, agg_param->ap_discard);
		break;
	case VOS_ITER_DKEY:
	case VOS_ITER_AKEY:
		rc = vos_obj_iter_aggregate(ih, agg_param->ap_discard);
		break;
	case VOS_ITER_SINGLE:
		return 0;
	case VOS_ITER_RECX:
		return 0;
	default:
		D_ASSERTF(false, "Invalid iter type\n");
		return -DER_INVAL;
	}

	if (rc == 1) {
		/* Reprobe flag is set */
		*acts |= VOS_ITER_CB_YIELD;
		rc = 0;
	}

	if (rc != 0)
		D_ERROR("VOS aggregation failed: %d\n", rc);

	return rc;
}

static int
aggregate_enter(struct vos_container *cont, bool discard)
{
	if (cont->vc_in_aggregation) {
		D_ERROR(DF_CONT": Already in aggregation. discard:%d\n",
			DP_CONT(cont->vc_pool->vp_id, cont->vc_id), discard);
		return -DER_BUSY;
	}

	cont->vc_in_aggregation = 1;
	cont->vc_abort_aggregation = 0;
	return 0;
}

static void
aggregate_exit(struct vos_container *cont, bool discard)
{
	D_ASSERT(cont->vc_in_aggregation);
	cont->vc_in_aggregation = 0;
}

static void
merge_window_init(struct agg_merge_window *mw)
{
	struct vos_agg_io_context *io = &mw->mw_io_ctxt;

	memset(mw, 0, sizeof(*mw));
	D_INIT_LIST_HEAD(&mw->mw_phy_ents);
	D_INIT_LIST_HEAD(&io->ic_nvme_exts);
	io->ic_csum_buf = NULL;
	io->ic_csum_buf_len = 0;
}

int
vos_aggregate(daos_handle_t coh, daos_epoch_range_t *epr,
	      unsigned int agg_flags)
{
	struct vos_container	*cont = vos_hdl2cont(coh);
	vos_iter_param_t	 iter_param = { 0 };
	struct vos_agg_param	 agg_param = { 0 };
	struct vos_iter_anchors	 anchors = { 0 };
	int			 rc;

	cont_has_csums = agg_flags & VAF_CSUM;
	unit_test = agg_flags & VAF_UNIT_TEST;
	D_ASSERT(epr != NULL);
	D_ASSERTF(epr->epr_lo < epr->epr_hi && epr->epr_hi != DAOS_EPOCH_MAX,
		  "epr_lo:"DF_U64", epr_hi:"DF_U64"\n",
		  epr->epr_lo, epr->epr_hi);

	rc = aggregate_enter(cont, false);
	if (rc)
		return rc;

	/* Set iteration parameters */
	iter_param.ip_hdl = coh;
	iter_param.ip_epr = *epr;
	/*
	 * Iterate in epoch reserve order for SV tree, so that we can know for
	 * sure the first returned recx in SV tree has highest epoch and can't
	 * be aggregated.
	 */
	iter_param.ip_epc_expr = VOS_IT_EPC_RR;
	/* EV tree iterator returns all sorted logical rectangles */
	iter_param.ip_flags = VOS_IT_PUNCHED | VOS_IT_RECX_VISIBLE |
		VOS_IT_RECX_COVERED;

	/* Set aggregation parameters */
	agg_param.ap_umm = &cont->vc_pool->vp_umm;
	agg_param.ap_coh = coh;
	agg_param.ap_credits_max = VOS_AGG_CREDITS_MAX;
	agg_param.ap_credits = 0;
	agg_param.ap_discard = false;
	merge_window_init(&agg_param.ap_window);

	iter_param.ip_flags |= VOS_IT_FOR_PURGE;
	rc = vos_iterate(&iter_param, VOS_ITER_OBJ, true, &anchors,
			 vos_aggregate_pre_cb, vos_aggregate_post_cb,
			 &agg_param);
	if (rc != 0) {
		close_merge_window(&agg_param.ap_window, rc);
		goto exit;
	}

	/*
	 * Update HAE, when aggregating for snapshot deletion, the
	 * @epr->epr_hi could be smaller than the HAE
	 */
	if (cont->vc_cont_df->cd_hae < epr->epr_hi)
		cont->vc_cont_df->cd_hae = epr->epr_hi;
exit:
	aggregate_exit(cont, false);
	if (cont_has_csums)
		D_FREE(agg_param.ap_window.mw_io_ctxt.ic_csum_buf);

	if (merge_window_status(&agg_param.ap_window) != MW_CLOSED)
		D_ASSERTF(false, "Merge window resource leaked.\n");

	return rc;
}

int
vos_discard(daos_handle_t coh, daos_epoch_range_t *epr)
{
	struct vos_container	*cont = vos_hdl2cont(coh);
	vos_iter_param_t	 iter_param = { 0 };
	struct vos_agg_param	 agg_param = { 0 };
	struct vos_iter_anchors	 anchors = { 0 };
	int			 rc;

	D_ASSERT(epr != NULL);
	D_ASSERTF(epr->epr_lo <= epr->epr_hi,
		  "epr_lo:"DF_U64", epr_hi:"DF_U64"\n",
		  epr->epr_lo, epr->epr_hi);

	rc = aggregate_enter(cont, true);
	if (rc != 0)
		return rc;

	D_DEBUG(DB_EPC, "Discard epr "DF_U64"-"DF_U64"\n",
		epr->epr_lo, epr->epr_hi);

	/* Set iteration parameters */
	iter_param.ip_hdl = coh;
	iter_param.ip_epr = *epr;
	if (epr->epr_lo == epr->epr_hi)
		iter_param.ip_epc_expr = VOS_IT_EPC_EQ;
	else if (epr->epr_hi != DAOS_EPOCH_MAX)
		iter_param.ip_epc_expr = VOS_IT_EPC_RR;
	else
		iter_param.ip_epc_expr = VOS_IT_EPC_GE;
	/* EV tree iterator returns all sorted logical rectangles */
	iter_param.ip_flags = VOS_IT_PUNCHED | VOS_IT_RECX_VISIBLE |
		VOS_IT_RECX_COVERED;

	/* Set aggregation parameters */
	agg_param.ap_umm = &cont->vc_pool->vp_umm;
	agg_param.ap_coh = coh;
	agg_param.ap_credits_max = VOS_AGG_CREDITS_MAX;
	agg_param.ap_credits = 0;
	agg_param.ap_discard = true;

	iter_param.ip_flags |= VOS_IT_FOR_PURGE;
	rc = vos_iterate(&iter_param, VOS_ITER_OBJ, true, &anchors,
			 vos_aggregate_pre_cb, vos_aggregate_post_cb,
			 &agg_param);

	aggregate_exit(cont, true);
	return rc;
}
