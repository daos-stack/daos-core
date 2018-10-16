/*
 * (C) Copyright 2018 Intel Corporation.
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
 * \file
 *
 * server: VOS-Related Utilities
 */

#define D_LOGFAC DD_FAC(server)

#include <daos_srv/daos_server.h>
#include <daos_srv/vos.h>

/**
 * Iterate VOS entries (i.e., containers, objects, dkeys, etc.) and call \a
 * cb(\a arg) for each entry.
 *
 * If \a cb returns a nonzero (either > 0 or < 0) value that is not
 * -DER_NONEXIST, this function stops the iteration and returns that nonzero
 * value from \a cb. If \a cb returns -DER_NONEXIST, this function completes
 * the iteration and returns 0. If \a cb returns 0, the iteration continues.
 *
 * \param[in]		type	entry type
 * \param[in]		param	parameters for \a type
 * \param[in,out]	anchor	[in]: where to begin; [out]: where stopped
 * \param[in]		cb	callback called for each entry
 * \param[in]		arg	callback argument
 *
 * \retval		0	iteration complete
 * \retval		> 0	callback return value
 * \retval		-DER_*	error (but never -DER_NONEXIST)
 */
int
dss_vos_iterate(vos_iter_type_t type, vos_iter_param_t *param,
		daos_anchor_t *anchor, dss_vos_iterate_cb_t cb, void *arg)
{
	daos_anchor_t		*probe_anchor = NULL;
	vos_iter_entry_t	key_ent;
	daos_handle_t		ih;
	int			rc;

	rc = vos_iter_prepare(type, param, &ih);
	if (rc != 0) {
		if (rc == -DER_NONEXIST) {
			daos_anchor_set_eof(anchor);
			rc = 0;
		} else {
			D_ERROR("failed to prepare iterator (type=%d): %d\n",
				type, rc);
		}
		D_GOTO(out, rc);
	}

	if (!daos_anchor_is_zero(anchor))
		probe_anchor = anchor;
	rc = vos_iter_probe(ih, probe_anchor);
	if (rc != 0) {
		if (rc == -DER_NONEXIST || rc == -DER_AGAIN) {
			daos_anchor_set_eof(anchor);
			rc = 0;
		} else {
			D_ERROR("failed to probe iterator (type=%d anchor=%p): "
				"%d\n", type, probe_anchor, rc);
		}
		D_GOTO(out_iter_fini, rc);
	}

	while (1) {
		rc = vos_iter_fetch(ih, &key_ent, anchor);
		if (rc != 0) {
			D_ERROR("failed to fetch iterator (type=%d): %d\n",
				type, rc);
			break;
		}

		rc = cb(ih, &key_ent, type, param, arg);
		if (rc != 0)
			break;

		rc = vos_iter_next(ih);
		if (rc) {
			if (rc != -DER_NONEXIST)
				D_ERROR("failed to iterate next (type=%d): "
					"%d\n", type, rc);
			break;
		}
	}

	if (rc == -DER_NONEXIST) {
		daos_anchor_set_eof(anchor);
		rc = 0;
	}

out_iter_fini:
	vos_iter_finish(ih);
out:
	return rc;
}

/* obj_enum_rec.rec_flags */
#define RECX_INLINE	(1U << 0)

struct obj_enum_rec {
	daos_recx_t		rec_recx;
	daos_epoch_range_t	rec_epr;
	uuid_t			rec_cookie;
	uint64_t		rec_size;
	uint32_t		rec_version;
	uint32_t		rec_flags;
};

static int
fill_recxs(daos_handle_t ih, vos_iter_entry_t *key_ent,
	   struct dss_enum_arg *arg, vos_iter_type_t type)
{
	/* check if recxs is full */
	if (arg->recxs_len >= arg->recxs_cap) {
		D_DEBUG(DB_IO, "recx_len %d recx_cap %d\n",
			arg->recxs_len, arg->recxs_cap);
		return 1;
	}

	arg->eprs[arg->eprs_len].epr_lo = key_ent->ie_epoch;
	arg->eprs[arg->eprs_len].epr_hi = DAOS_EPOCH_MAX;
	arg->eprs_len++;

	arg->recxs[arg->recxs_len] = key_ent->ie_recx;
	arg->recxs_len++;
	if (arg->rsize == 0) {
		arg->rsize = key_ent->ie_rsize;
	} else if (arg->rsize != key_ent->ie_rsize) {
		D_ERROR("different size "DF_U64" != "DF_U64"\n", arg->rsize,
			key_ent->ie_rsize);
		return -DER_INVAL;
	}

	D_DEBUG(DB_IO, "Pack recxs_eprs "DF_U64"/"DF_U64" recxs_len %d size "
		DF_U64"\n", key_ent->ie_recx.rx_idx, key_ent->ie_recx.rx_nr,
		arg->recxs_len, arg->rsize);

	arg->rnum++;
	return 0;
}

static int
fill_recxs_cb(daos_handle_t ih, vos_iter_entry_t *key_ent,
	      vos_iter_type_t type, vos_iter_param_t *param, void *arg)
{
	return fill_recxs(ih, key_ent, arg, type);
}

static int
is_sgl_kds_full(struct dss_enum_arg *arg, daos_size_t size)
{
	d_sg_list_t *sgl = arg->sgl;

	/* Find avaible iovs in sgl
	 * XXX this is buggy because key descriptors require keys are stored
	 * in sgl in the same order as descriptors, but it's OK for now because
	 * we only use one IOV.
	 */
	while (arg->sgl_idx < sgl->sg_nr) {
		daos_iov_t *iovs = sgl->sg_iovs;

		if (iovs[arg->sgl_idx].iov_len + size >=
		    iovs[arg->sgl_idx].iov_buf_len) {
			D_DEBUG(DB_IO, "current %dth iov buf is full"
				" iov_len %zd size "DF_U64" buf_len %zd\n",
				arg->sgl_idx, iovs[arg->sgl_idx].iov_len, size,
				iovs[arg->sgl_idx].iov_buf_len);
			arg->sgl_idx++;
			continue;
		}
		break;
	}

	/* Update sg_nr_out */
	if (arg->sgl_idx < sgl->sg_nr && sgl->sg_nr_out < arg->sgl_idx + 1)
		sgl->sg_nr_out = arg->sgl_idx + 1;

	/* Check if the sgl is full */
	if (arg->sgl_idx >= sgl->sg_nr || arg->kds_len >= arg->kds_cap) {
		D_DEBUG(DB_IO, "sgl or kds full sgl %d/%d kds %d/%d size "
			DF_U64"\n", arg->sgl_idx, sgl->sg_nr,
			arg->kds_len, arg->kds_cap, size);
		return 1;
	}

	return 0;
}

static int
fill_obj(daos_handle_t ih, vos_iter_entry_t *entry, struct dss_enum_arg *arg,
	 vos_iter_type_t type)
{
	daos_iov_t *iovs = arg->sgl->sg_iovs;

	D_ASSERTF(type == VOS_ITER_OBJ, "%d\n", type);

	if (is_sgl_kds_full(arg, sizeof(entry->ie_oid)))
		return 1;

	/* Append a new descriptor to kds. */
	D_ASSERT(arg->kds_len < arg->kds_cap);
	memset(&arg->kds[arg->kds_len], 0, sizeof(arg->kds[arg->kds_len]));
	arg->kds[arg->kds_len].kd_key_len = sizeof(entry->ie_oid);
	arg->kds[arg->kds_len].kd_val_types = type;
	arg->kds_len++;

	/* Append the object ID to iovs. */
	D_ASSERT(iovs[arg->sgl_idx].iov_len + sizeof(entry->ie_oid) <
		 iovs[arg->sgl_idx].iov_buf_len);
	memcpy(iovs[arg->sgl_idx].iov_buf + iovs[arg->sgl_idx].iov_len,
	       &entry->ie_oid, sizeof(entry->ie_oid));
	iovs[arg->sgl_idx].iov_len += sizeof(entry->ie_oid);

	D_DEBUG(DB_IO, "Pack obj "DF_UOID" iov_len %zu kds_len %d\n",
		DP_UOID(entry->ie_oid), iovs[arg->sgl_idx].iov_len,
		arg->kds_len);
	return 0;
}

static int
fill_obj_cb(daos_handle_t ih, vos_iter_entry_t *key_ent, vos_iter_type_t type,
	    vos_iter_param_t *param, void *arg)
{
	return fill_obj(ih, key_ent, arg, type);
}

static int
fill_key(daos_handle_t ih, vos_iter_entry_t *key_ent, struct dss_enum_arg *arg,
	 vos_iter_type_t type)
{
	daos_iov_t	*iovs = arg->sgl->sg_iovs;
	daos_size_t	 size;

	D_ASSERT(type == VOS_ITER_DKEY || type == VOS_ITER_AKEY);
	size = key_ent->ie_key.iov_len;

	if (is_sgl_kds_full(arg, size)) {
		if (arg->kds_len == 0) {
			arg->kds[0].kd_key_len = size;
			return -DER_KEY2BIG;
		} else {
			return 1;
		}
	}

	D_ASSERT(arg->kds_len < arg->kds_cap);
	arg->kds[arg->kds_len].kd_key_len = size;
	arg->kds[arg->kds_len].kd_csum_len = 0;
	arg->kds[arg->kds_len].kd_val_types = type;
	arg->kds_len++;

	if (arg->eprs != NULL) {
		arg->eprs[arg->eprs_len].epr_lo = key_ent->ie_epoch;
		arg->eprs[arg->eprs_len].epr_hi = DAOS_EPOCH_MAX;
		arg->eprs_len++;
	}

	D_ASSERT(iovs[arg->sgl_idx].iov_len + key_ent->ie_key.iov_len <
		 iovs[arg->sgl_idx].iov_buf_len);
	memcpy(iovs[arg->sgl_idx].iov_buf + iovs[arg->sgl_idx].iov_len,
	       key_ent->ie_key.iov_buf, key_ent->ie_key.iov_len);

	iovs[arg->sgl_idx].iov_len += key_ent->ie_key.iov_len;
	D_DEBUG(DB_IO, "Pack key %d %s iov total %zd"
		" kds len %d\n", (int)key_ent->ie_key.iov_len,
		(char *)key_ent->ie_key.iov_buf,
		iovs[arg->sgl_idx].iov_len, arg->kds_len - 1);

	return 0;
}

static int
fill_key_cb(daos_handle_t ih, vos_iter_entry_t *key_ent, vos_iter_type_t type,
	    vos_iter_param_t *param, void *arg)
{
	return fill_key(ih, key_ent, arg, type);
}

/* Copy the data of the recx into buf. */
/* TODO: Use entry->ie_biov when implemented. */
static void
copy_data(vos_iter_type_t type, vos_iter_param_t *param,
	  vos_iter_entry_t *entry, void *buf, size_t len)
{
	daos_iod_t	iod;
	daos_iov_t	iov;
	daos_sg_list_t	sgl;
	int		rc;

	D_ASSERT(type == VOS_ITER_SINGLE || type == VOS_ITER_RECX);

	memset(&iod, 0, sizeof(iod));
	iod.iod_name = param->ip_akey;
	if (type == VOS_ITER_SINGLE)
		iod.iod_type = DAOS_IOD_SINGLE;
	else
		iod.iod_type = DAOS_IOD_ARRAY;
	iod.iod_nr = 1;
	iod.iod_recxs = &entry->ie_recx;
	iod.iod_eprs = NULL;

	iov.iov_buf = buf;
	iov.iov_buf_len = len;
	iov.iov_len = 0;
	sgl.sg_nr = 1;
	sgl.sg_nr_out = 0;
	sgl.sg_iovs = &iov;

	rc = vos_obj_fetch(param->ip_hdl, param->ip_oid, entry->ie_epoch,
			   &param->ip_dkey, 1 /* iod_nr */, &iod, &sgl);
	/* This vos_obj_fetch call is a workaround anyway. */
	D_ASSERTF(rc == 0, "%d\n", rc);
	D_ASSERT(iod.iod_size == entry->ie_rsize);
}

/* Callers are responsible for incrementing arg->kds_len. See iter_akey_cb. */
static int
fill_rec(daos_handle_t ih, vos_iter_entry_t *key_ent, struct dss_enum_arg *arg,
	 vos_iter_type_t type, vos_iter_param_t *param)
{
	daos_iov_t		*iovs = arg->sgl->sg_iovs;
	struct obj_enum_rec	*rec;
	daos_size_t		 data_size;
	daos_size_t		 size = sizeof(*rec);
	bool			 inline_data = false;

	D_ASSERT(type == VOS_ITER_SINGLE || type == VOS_ITER_RECX);

	/* Inline the data? A 0 threshold disables this completely. */
	data_size = key_ent->ie_rsize * key_ent->ie_recx.rx_nr;
	if (arg->inline_thres > 0 && data_size <= arg->inline_thres) {
		inline_data = true;
		size += data_size;
	}

	if (is_sgl_kds_full(arg, size))
		return 1;

	/* Grow the next new descriptor (instead of creating yet a new one). */
	arg->kds[arg->kds_len].kd_val_types = type;
	arg->kds[arg->kds_len].kd_key_len += sizeof(*rec);

	/* Append the recx record to iovs. */
	D_ASSERT(iovs[arg->sgl_idx].iov_len + sizeof(*rec) <
		 iovs[arg->sgl_idx].iov_buf_len);
	rec = iovs[arg->sgl_idx].iov_buf + iovs[arg->sgl_idx].iov_len;
	rec->rec_recx = key_ent->ie_recx;
	rec->rec_size = key_ent->ie_rsize;
	rec->rec_epr.epr_lo = key_ent->ie_epoch;
	rec->rec_epr.epr_hi = DAOS_EPOCH_MAX;
	uuid_copy(rec->rec_cookie, key_ent->ie_cookie);
	rec->rec_version = key_ent->ie_ver;
	rec->rec_flags = 0;
	iovs[arg->sgl_idx].iov_len += sizeof(*rec);

	/* If we've decided to inline the data, append the data to iovs. */
	if (inline_data) {
		arg->kds[arg->kds_len].kd_key_len += data_size;
		rec->rec_flags |= RECX_INLINE;
		/* Punched recxs do not have any data to copy. */
		if (data_size > 0)
			copy_data(type, param, key_ent,
				  iovs[arg->sgl_idx].iov_buf +
				  iovs[arg->sgl_idx].iov_len, data_size);
		iovs[arg->sgl_idx].iov_len += data_size;
	}

	D_DEBUG(DB_IO, "Pack rec "DF_U64"/"DF_U64
		" rsize "DF_U64" cookie "DF_UUID" ver %u"
		" kd_len "DF_U64" type %d sgl_idx %d kds_len %d inline "DF_U64
		" epr "DF_U64"/"DF_U64"\n", key_ent->ie_recx.rx_idx,
		key_ent->ie_recx.rx_nr, key_ent->ie_rsize,
		DP_UUID(rec->rec_cookie), rec->rec_version,
		arg->kds[arg->kds_len].kd_key_len, type, arg->sgl_idx,
		arg->kds_len, rec->rec_flags & RECX_INLINE ? data_size : 0,
		rec->rec_epr.epr_lo, rec->rec_epr.epr_hi);
	return 0;
}

static int
fill_rec_cb(daos_handle_t ih, vos_iter_entry_t *key_ent, vos_iter_type_t type,
	    vos_iter_param_t *param, void *arg)
{
	return fill_rec(ih, key_ent, arg, type, param);
}

static int
iter_akey_cb(daos_handle_t ih, vos_iter_entry_t *key_ent, vos_iter_type_t type,
	     vos_iter_param_t *param, void *varg)
{
	struct dss_enum_arg	*arg = varg;
	vos_iter_param_t	 iter_recx_param;
	daos_anchor_t		 single_anchor = { 0 };
	int			 rc;

	D_DEBUG(DB_IO, "enum key %d %s type %d\n",
		(int)key_ent->ie_key.iov_len,
		(char *)key_ent->ie_key.iov_buf, type);

	/* Fill the current key */
	rc = fill_key(ih, key_ent, arg, VOS_ITER_AKEY);
	if (rc)
		goto out;

	iter_recx_param = *param;
	iter_recx_param.ip_akey = key_ent->ie_key;

	/* iterate array record */
	rc = dss_vos_iterate(VOS_ITER_RECX, &iter_recx_param, &arg->recx_anchor,
			     fill_rec_cb, arg);

	if (arg->kds[arg->kds_len].kd_key_len > 0) {
		arg->kds_len++;
		/** This eprs will not be used during rebuild,
		 * because the epoch for each record will be returned
		 * through obj_enum_rec anyway, see fill_rec().
		 * This "empty" eprs is just for eprs and kds to
		 * be matched, so it would be easier for unpacking
		 * see dss_enum_unpack().
		 */
		if (arg->eprs != NULL) {
			arg->eprs[arg->eprs_len].epr_lo = DAOS_EPOCH_MAX;
			arg->eprs[arg->eprs_len].epr_hi = DAOS_EPOCH_MAX;
			arg->eprs_len++;
		}
	}

	/* Exit either failure or buffer is full */
	if (rc) {
		if (rc < 0)
			D_ERROR("failed to enumerate array recxs: %d\n", rc);
		goto out;
	}

	D_ASSERT(daos_anchor_is_eof(&arg->recx_anchor));
	daos_anchor_set_zero(&arg->recx_anchor);

	/* iterate single record */
	rc = dss_vos_iterate(VOS_ITER_SINGLE, &iter_recx_param, &single_anchor,
			     fill_rec_cb, arg);

	if (rc) {
		if (rc < 0)
			D_ERROR("failed to enumerate single recxs: %d\n", rc);
		goto out;
	}

	if (arg->kds[arg->kds_len].kd_key_len > 0) {
		arg->kds_len++;
		/** empty eprs, see comments above */
		if (arg->eprs != NULL) {
			arg->eprs[arg->eprs_len].epr_lo = DAOS_EPOCH_MAX;
			arg->eprs[arg->eprs_len].epr_hi = DAOS_EPOCH_MAX;
			arg->eprs_len++;
		}
	}
out:
	return rc;
}

static int
iter_dkey_cb(daos_handle_t ih, vos_iter_entry_t *key_ent, vos_iter_type_t type,
	     vos_iter_param_t *param, void *varg)
{
	struct dss_enum_arg	*arg = varg;
	vos_iter_param_t	 iter_akey_param;
	int			 rc;

	D_DEBUG(DB_IO, "enum key %d %s type %d\n",
		(int)key_ent->ie_key.iov_len,
		(char *)key_ent->ie_key.iov_buf, type);

	/* Fill the current dkey */
	rc = fill_key(ih, key_ent, arg, VOS_ITER_DKEY);
	if (rc != 0)
		return rc;

	/* iterate akey */
	iter_akey_param = *param;
	iter_akey_param.ip_dkey = key_ent->ie_key;
	rc = dss_vos_iterate(VOS_ITER_AKEY, &iter_akey_param, &arg->akey_anchor,
			     iter_akey_cb, arg);
	if (rc) {
		if (rc < 0)
			D_ERROR("failed to enumerate akeys: %d\n", rc);
		return rc;
	}

	D_ASSERT(daos_anchor_is_eof(&arg->akey_anchor));
	daos_anchor_set_zero(&arg->akey_anchor);
	daos_anchor_set_zero(&arg->recx_anchor);

	return rc;
}

static int
iter_obj_cb(daos_handle_t ih, vos_iter_entry_t *entry, vos_iter_type_t type,
	    vos_iter_param_t *param, void *varg)
{
	struct dss_enum_arg	*arg = varg;
	vos_iter_param_t	 iter_dkey_param;
	int			 rc;

	D_ASSERTF(type == VOS_ITER_OBJ, "%d\n", type);
	D_DEBUG(DB_IO, "enum obj "DF_UOID"\n", DP_UOID(entry->ie_oid));

	rc = fill_obj(ih, entry, arg, type);
	if (rc != 0)
		return rc;

	iter_dkey_param = *param;
	iter_dkey_param.ip_oid = entry->ie_oid;
	rc = dss_vos_iterate(VOS_ITER_DKEY, &iter_dkey_param, &arg->dkey_anchor,
			     iter_dkey_cb, arg);
	if (rc != 0) {
		if (rc < 0)
			D_ERROR("failed to enumerate dkeys: %d\n", rc);
		return rc;
	}

	D_ASSERT(daos_anchor_is_eof(&arg->dkey_anchor));
	daos_anchor_set_zero(&arg->dkey_anchor);
	daos_anchor_set_zero(&arg->akey_anchor);
	daos_anchor_set_zero(&arg->recx_anchor);

	return 0;
}

/**
 * Enumerate VOS objects, dkeys, akeys, and/or recxs and pack them into a set
 * of buffers.
 *
 * The buffers must be provided by the caller. They may contain existing data,
 * in which case this function appends to them.
 *
 * \param[in]		type	iteration type
 * \param[in,out]	arg	enumeration argument
 *
 * \retval		0	enumeration complete
 * \retval		1	buffer(s) full
 * \retval		-DER_*	error
 */
int
dss_enum_pack(vos_iter_type_t type, struct dss_enum_arg *arg)
{
	daos_anchor_t	       *anchor;
	dss_vos_iterate_cb_t	cb;
	int			rc;

	D_ASSERT(!arg->fill_recxs ||
		 type == VOS_ITER_SINGLE || type == VOS_ITER_RECX);
	switch (type) {
	case VOS_ITER_OBJ:
		anchor = &arg->obj_anchor;
		cb = arg->recursive ? iter_obj_cb : fill_obj_cb;
		break;
	case VOS_ITER_DKEY:
		anchor = &arg->dkey_anchor;
		cb = arg->recursive ? iter_dkey_cb : fill_key_cb;
		break;
	case VOS_ITER_AKEY:
		anchor = &arg->akey_anchor;
		cb = arg->recursive ? iter_akey_cb : fill_key_cb;
		break;
	case VOS_ITER_SINGLE:
	case VOS_ITER_RECX:
		anchor = &arg->recx_anchor;
		cb = arg->fill_recxs ? fill_recxs_cb : fill_rec_cb;
		break;
	default:
		D_ASSERTF(false, "unknown/unsupported type %d\n", type);
	}

	rc = dss_vos_iterate(type, &arg->param, anchor, cb, arg);

	D_DEBUG(DB_IO, "enum type %d tag %d rc %d\n", type,
		dss_get_module_info()->dmi_tid, rc);
	return rc;
}

static int
grow_array(void **arrayp, size_t elem_size, int old_len, int new_len)
{
	void *p;

	D_ASSERTF(old_len < new_len, "%d < %d\n", old_len, new_len);
	D_REALLOC(p, *arrayp, elem_size * new_len);
	if (p == NULL)
		return -DER_NOMEM;
	/* Until D_REALLOC does this, zero the new segment. */
	memset(p + elem_size * old_len, 0, elem_size * (new_len - old_len));
	*arrayp = p;
	return 0;
}

/* Parse recxs in <*data, len> and append them to iod and sgl. */
static int
unpack_recxs(daos_iod_t *iod, int *recxs_cap, daos_sg_list_t *sgl,
	     daos_key_t *akey, daos_key_desc_t *kds, void **data,
	     daos_size_t len, uuid_t cookie, uint32_t *version)
{
	int rc = 0;

	if (iod->iod_name.iov_len == 0)
		daos_iov_copy(&iod->iod_name, akey);
	else
		D_ASSERT(daos_key_match(&iod->iod_name, akey));

	if (kds == NULL)
		return 0;

	while (len > 0) {
		struct obj_enum_rec *rec = *data;

		D_DEBUG(DB_REBUILD, "data %p len "DF_U64"\n", *data, len);

		/* Every recx begins with an obj_enum_rec. */
		if (len < sizeof(*rec)) {
			D_ERROR("invalid recxs: <%p, %zu>\n", *data, len);
			rc = -DER_INVAL;
			break;
		}

		/* Check if the cookie or the version is changing. */
		if (uuid_is_null(cookie)) {
			uuid_copy(cookie, rec->rec_cookie);
			*version = rec->rec_version;
		} else if (uuid_compare(cookie, rec->rec_cookie) != 0 ||
			   *version != rec->rec_version) {
			D_DEBUG(DB_REBUILD, "different cookie or version"
				DF_UUIDF" "DF_UUIDF" %u != %u\n",
				DP_UUID(cookie), DP_UUID(rec->rec_cookie),
				*version, rec->rec_version);
			rc = 1;
			break;
		}

		/* Iteration might return multiple single record with same
		 * dkey/akeks but different epoch. But fetch & update only
		 * allow 1 SINGLE type record per IOD. Let's put these
		 * single records in different IODs.
		 */
		if (kds->kd_val_types == VOS_ITER_SINGLE && iod->iod_nr > 0) {
			rc = 1;
			break;
		}

		if (iod->iod_size != 0 && iod->iod_size != rec->rec_size)
			D_WARN("rsize "DF_U64" != "DF_U64" are different"
			       " under one akey\n", iod->iod_size,
			       rec->rec_size);

		/* If the arrays are full, grow them as if all the remaining
		 * recxs have no inline data.
		 */
		if (iod->iod_nr + 1 > *recxs_cap) {
			int cap = *recxs_cap + len / sizeof(*rec);

			rc = grow_array((void **)&iod->iod_recxs,
					sizeof(*iod->iod_recxs), *recxs_cap,
					cap);
			if (rc != 0)
				break;
			rc = grow_array((void **)&iod->iod_eprs,
					sizeof(*iod->iod_eprs), *recxs_cap,
					cap);
			if (rc != 0)
				break;
			if (sgl != NULL) {
				rc = grow_array((void **)&sgl->sg_iovs,
						sizeof(*sgl->sg_iovs),
						*recxs_cap, cap);
				if (rc != 0)
					break;
			}
			/* If we execute any of the three breaks above,
			 * *recxs_cap will be < the real capacities of some of
			 * the arrays. This is harmless, as it only causes the
			 * diff segment to be copied and zeroed unnecessarily
			 * next time we grow them.
			 */
			*recxs_cap = cap;
		}

		/* Append one more recx. */
		iod->iod_eprs[iod->iod_nr] = rec->rec_epr;
		/* Iteration does not fill the high epoch, so let's reset
		 * the high epoch with EPOCH_MAX to make vos fetch/update happy.
		 */
		iod->iod_eprs[iod->iod_nr].epr_hi = DAOS_EPOCH_MAX;
		iod->iod_recxs[iod->iod_nr] = rec->rec_recx;
		iod->iod_nr++;
		if (iod->iod_size == 0)
			iod->iod_size = rec->rec_size;
		*data += sizeof(*rec);
		len -= sizeof(*rec);

		/* Append the data, if inline. */
		if (sgl != NULL) {
			daos_iov_t *iov = &sgl->sg_iovs[sgl->sg_nr];

			if (rec->rec_flags & RECX_INLINE) {
				daos_iov_set(iov, *data,
					     rec->rec_size *
					     rec->rec_recx.rx_nr);
				D_DEBUG(DB_TRACE, "set data iov %p len %d\n",
					iov, (int)iov->iov_len);
			} else {
				daos_iov_set(iov, NULL, 0);
			}
			sgl->sg_nr++;
			D_ASSERTF(sgl->sg_nr == iod->iod_nr, "%u == %u\n",
				  sgl->sg_nr, iod->iod_nr);
			*data += iov->iov_len;
			len -= iov->iov_len;
		}

		D_DEBUG(DB_REBUILD, "pack %p idx/nr "DF_U64"/"DF_U64
			" epr lo/hi "DF_U64"/"DF_U64" size %zd inline %zu\n",
			*data, iod->iod_recxs[iod->iod_nr - 1].rx_idx,
			iod->iod_recxs[iod->iod_nr - 1].rx_nr,
			iod->iod_eprs[iod->iod_nr - 1].epr_lo,
			iod->iod_eprs[iod->iod_nr - 1].epr_hi, iod->iod_size,
			sgl != NULL ? sgl->sg_iovs[sgl->sg_nr - 1].iov_len : 0);
	}

	if (kds->kd_val_types == VOS_ITER_RECX)
		iod->iod_type = DAOS_IOD_ARRAY;
	else
		iod->iod_type = DAOS_IOD_SINGLE;

	D_DEBUG(DB_REBUILD, "pack nr %d cookie/version "DF_UUID"/%u rc %d\n",
		iod->iod_nr, DP_UUID(cookie), *version, rc);
	return rc;
}

/**
 * Initialize \a io with \a iods[\a iods_cap], \a recxs_caps[\a iods_cap], and
 * \a sgls[\a iods_cap].
 *
 * \param[in,out]	io		I/O descriptor
 * \param[in]		iods		daos_iod_t array
 * \param[in]		recxs_caps	recxs capacity array
 * \param[in]		sgls		optional sgl array for inline recxs
 * \param[in]		ephs		epoch array
 * \param[in]		iods_cap	maximal number of elements in \a iods,
 *					\a recxs_caps, \a sgls, and \a ephs
 */
static void
dss_enum_unpack_io_init(struct dss_enum_unpack_io *io, daos_iod_t *iods,
			int *recxs_caps, daos_sg_list_t *sgls,
			daos_epoch_t *ephs, int iods_cap)
{
	int i;

	memset(io, 0, sizeof(*io));

	io->ui_dkey_eph = DAOS_EPOCH_MAX;

	D_ASSERTF(iods_cap > 0, "%d\n", iods_cap);
	io->ui_iods_cap = iods_cap;

	D_ASSERT(iods != NULL);
	memset(iods, 0, sizeof(*iods) * iods_cap);
	io->ui_iods = iods;

	D_ASSERT(recxs_caps != NULL);
	memset(recxs_caps, 0, sizeof(*recxs_caps) * iods_cap);
	io->ui_recxs_caps = recxs_caps;

	if (sgls != NULL) {
		memset(sgls, 0, sizeof(*sgls) * iods_cap);
		io->ui_sgls = sgls;
	}

	for (i = 0; i < iods_cap; i++)
		ephs[i] = DAOS_EPOCH_MAX;

	io->ui_akey_ephs = ephs;
	uuid_clear(io->ui_cookie);
}

static void
clear_iod(daos_iod_t *iod, daos_sg_list_t *sgl, int *recxs_cap)
{
	daos_iov_free(&iod->iod_name);
	if (iod->iod_recxs != NULL)
		D_FREE(iod->iod_recxs);
	if (iod->iod_eprs != NULL)
		D_FREE(iod->iod_eprs);
	memset(iod, 0, sizeof(*iod));

	if (sgl != NULL) {
		if (sgl->sg_iovs != NULL)
			D_FREE(sgl->sg_iovs);
		memset(sgl, 0, sizeof(*sgl));
	}

	*recxs_cap = 0;
}

/**
 * Clear the iods/sgls in \a io.
 *
 * \param[in]	io	I/O descriptor
 */
static void
dss_enum_unpack_io_clear(struct dss_enum_unpack_io *io)
{
	int i;

	for (i = 0; i < io->ui_iods_len; i++) {
		daos_sg_list_t *sgl = NULL;

		if (io->ui_sgls != NULL)
			sgl = &io->ui_sgls[i];
		clear_iod(&io->ui_iods[i], sgl, &io->ui_recxs_caps[i]);
		if (io->ui_akey_ephs)
			io->ui_akey_ephs[i] = DAOS_EPOCH_MAX;
	}

	io->ui_dkey_eph = DAOS_EPOCH_MAX;
	io->ui_iods_len = 0;
	uuid_clear(io->ui_cookie);
	io->ui_version = 0;
}

/**
 * Finalize \a io. All iods/sgls must have already been cleard.
 *
 * \param[in]	io	I/O descriptor
 */
static void
dss_enum_unpack_io_fini(struct dss_enum_unpack_io *io)
{
	D_ASSERTF(io->ui_iods_len == 0, "%d\n", io->ui_iods_len);
	daos_iov_free(&io->ui_dkey);
}

/*
 * Close the current iod (i.e., io->ui_iods[io->ui_iods_len]). If it contains
 * recxs, append it to io by incrementing io->ui_iods_len. If it doesn't
 * contain any recx, clear it.
 */
static void
close_iod(struct dss_enum_unpack_io *io)
{
	D_ASSERTF(io->ui_iods_cap > 0, "%d > 0\n", io->ui_iods_cap);
	D_ASSERTF(io->ui_iods_len < io->ui_iods_cap, "%d < %d\n",
		  io->ui_iods_len, io->ui_iods_cap);
	if (io->ui_iods[io->ui_iods_len].iod_nr > 0) {
		io->ui_iods_len++;
	} else {
		daos_sg_list_t *sgl = NULL;

		D_DEBUG(DB_TRACE, "iod without recxs: %d\n", io->ui_iods_len);
		if (io->ui_sgls != NULL)
			sgl = &io->ui_sgls[io->ui_iods_len];
		clear_iod(&io->ui_iods[io->ui_iods_len], sgl,
			  &io->ui_recxs_caps[io->ui_iods_len]);
	}
}

/* Close io, pass it to cb, and clear it. */
static int
complete_io(struct dss_enum_unpack_io *io, dss_enum_unpack_cb_t cb, void *arg)
{
	int rc = 0;

	if (io->ui_iods_len == 0) {
		D_DEBUG(DB_TRACE, "io empty\n");
		goto out;
	}
	rc = cb(io, arg);
out:
	dss_enum_unpack_io_clear(io);
	return rc;
}

/**
 * Unpack the result of a dss_enum_pack enumeration into \a io, which can then
 * be used to issue a VOS update. \a arg->*_anchor are ignored currently. \a cb
 * will be called, for the caller to consume the recxs accumulated in \a io.
 *
 * \param[in]		type	enumeration type
 * \param[in]		arg	enumeration argument
 * \param[in]		cb	callback
 * \param[in]		cb_arg	callback argument
 */
int
dss_enum_unpack(vos_iter_type_t type, struct dss_enum_arg *arg,
		dss_enum_unpack_cb_t cb, void *cb_arg)
{
	struct dss_enum_unpack_io	io;
	daos_iod_t			iods[DSS_ENUM_UNPACK_MAX_IODS];
	int				recxs_caps[DSS_ENUM_UNPACK_MAX_IODS];
	daos_epoch_t			ephs[DSS_ENUM_UNPACK_MAX_IODS];
	daos_sg_list_t			sgls[DSS_ENUM_UNPACK_MAX_IODS];
	daos_key_t			akey = {0};
	daos_epoch_range_t		*eprs = arg->eprs;
	void				*ptr;
	unsigned int			i;
	int				rc = 0;

	/* Currently, this function is only for unpacking recursive
	 * enumerations from arg->kds and arg->sgl.
	 */
	D_ASSERT(arg->recursive && !arg->fill_recxs);

	D_ASSERT(arg->kds_len > 0);
	D_ASSERT(arg->kds != NULL);
	if (arg->kds[0].kd_val_types != type) {
		D_ERROR("the first kds type %d != %d\n",
			arg->kds[0].kd_val_types, type);
		return -DER_INVAL;
	}

	dss_enum_unpack_io_init(&io, iods, recxs_caps, sgls, ephs,
				DSS_ENUM_UNPACK_MAX_IODS);
	if (type > VOS_ITER_OBJ)
		io.ui_oid = arg->param.ip_oid;

	D_ASSERTF(arg->sgl->sg_nr > 0, "%u\n", arg->sgl->sg_nr);
	D_ASSERT(arg->sgl->sg_iovs != NULL);
	ptr = arg->sgl->sg_iovs[0].iov_buf;

	for (i = 0; i < arg->kds_len; i++) {
		D_DEBUG(DB_REBUILD, "process %d type %d ptr %p len "DF_U64
			" total %zd\n", i, arg->kds[i].kd_val_types, ptr,
			arg->kds[i].kd_key_len, arg->sgl->sg_iovs[0].iov_len);

		D_ASSERT(arg->kds[i].kd_key_len > 0);
		if (arg->kds[i].kd_val_types == VOS_ITER_OBJ) {
			daos_unit_oid_t *oid = ptr;

			if (arg->kds[i].kd_key_len != sizeof(*oid)) {
				D_ERROR("Invalid object ID size: "DF_U64
					" != %zu\n", arg->kds[i].kd_key_len,
					sizeof(*oid));
				rc = -DER_INVAL;
				break;
			}
			if (daos_unit_oid_is_null(io.ui_oid)) {
				io.ui_oid = *oid;
			} else if (daos_unit_oid_compare(io.ui_oid, *oid) !=
				   0) {
				close_iod(&io);
				rc = complete_io(&io, cb, cb_arg);
				if (rc != 0)
					break;
				daos_iov_free(&io.ui_dkey);
				io.ui_oid = *oid;
			}
			D_DEBUG(DB_REBUILD, "process obj "DF_UOID"\n",
				DP_UOID(io.ui_oid));
		} else if (arg->kds[i].kd_val_types == VOS_ITER_DKEY) {
			daos_key_t tmp_key;

			tmp_key.iov_buf = ptr;
			tmp_key.iov_buf_len = arg->kds[i].kd_key_len;
			tmp_key.iov_len = arg->kds[i].kd_key_len;
			if (eprs != NULL)
				io.ui_dkey_eph = eprs[i].epr_lo;

			if (io.ui_dkey.iov_len == 0) {
				daos_iov_copy(&io.ui_dkey, &tmp_key);
			} else if (!daos_key_match(&io.ui_dkey, &tmp_key) ||
				   (eprs != NULL &&
				    io.ui_dkey_eph != eprs[i].epr_lo)) {
				close_iod(&io);
				rc = complete_io(&io, cb, cb_arg);
				if (rc != 0)
					break;

				if (!daos_key_match(&io.ui_dkey, &tmp_key)) {
					daos_iov_free(&io.ui_dkey);
					daos_iov_copy(&io.ui_dkey, &tmp_key);
				}
			}

			D_DEBUG(DB_REBUILD, "process dkey %d %s eph "DF_U64"\n",
				(int)io.ui_dkey.iov_len,
				(char *)io.ui_dkey.iov_buf,
				eprs ? io.ui_dkey_eph : 0);
		} else if (arg->kds[i].kd_val_types == VOS_ITER_AKEY) {
			daos_key_t *iod_akey;

			akey.iov_buf = ptr;
			akey.iov_buf_len = arg->kds[i].kd_key_len;
			akey.iov_len = arg->kds[i].kd_key_len;
			if (io.ui_dkey.iov_buf == NULL) {
				D_ERROR("No dkey for akey %*.s invalid buf.\n",
				      (int)akey.iov_len, (char *)akey.iov_buf);
				rc = -DER_INVAL;
				break;
			}
			D_DEBUG(DB_REBUILD, "process akey %d %s\n",
				(int)akey.iov_len, (char *)akey.iov_buf);

			if (io.ui_iods_len >= io.ui_iods_cap) {
				close_iod(&io);
				rc = complete_io(&io, cb, cb_arg);
				if (rc < 0)
					goto out;
			}

			/* If there are no records for akey(punched akey rec),
			 * then ui_iods_len might still point to the last dkey,
			 * i.e. close_iod are not being called.
			 */
			iod_akey = &io.ui_iods[io.ui_iods_len].iod_name;
			if (iod_akey->iov_len != 0 &&
			    !daos_key_match(iod_akey, &akey))
				io.ui_iods_len++;

			rc = unpack_recxs(&io.ui_iods[io.ui_iods_len],
					  NULL, NULL, &akey, NULL, NULL,
					  0, NULL, NULL);
			if (rc < 0)
				goto out;

			if (eprs)
				io.ui_akey_ephs[io.ui_iods_len] =
							eprs[i].epr_lo;
		} else if (arg->kds[i].kd_val_types == VOS_ITER_SINGLE ||
			   arg->kds[i].kd_val_types == VOS_ITER_RECX) {
			void *data = ptr;

			if (io.ui_dkey.iov_len == 0 || akey.iov_len == 0) {
				D_ERROR("invalid list buf for kds %d\n", i);
				rc = -DER_INVAL;
				break;
			}

			while (1) {
				daos_size_t	len;
				int		j = io.ui_iods_len;

				/* Because vos_obj_update only accept single
				 * cookie/version, let's go through the records
				 * to check different cookie and version, and
				 * queue rebuild.
				 */
				len = ptr + arg->kds[i].kd_key_len - data;
				rc = unpack_recxs(&io.ui_iods[j],
						  &io.ui_recxs_caps[j],
						  io.ui_sgls == NULL ?
						  NULL : &io.ui_sgls[j], &akey,
						  &arg->kds[i], &data, len,
						  io.ui_cookie,
						  &io.ui_version);
				if (rc < 0)
					goto out;

				/* All records referred by this kds has been
				 * packed, then it does not need to send
				 * right away, might pack more next round.
				 */
				if (rc == 0)
					break;

				/* Otherwise let's complete current io, and go
				 * next round.
				 */
				close_iod(&io);
				rc = complete_io(&io, cb, cb_arg);
				if (rc < 0)
					goto out;
			}
		} else {
			D_ERROR("unknow kds type %d\n",
				arg->kds[i].kd_val_types);
			rc = -DER_INVAL;
			break;
		}
		ptr += arg->kds[i].kd_key_len;
	}

	if (io.ui_iods[0].iod_nr > 0) {
		close_iod(&io);
		rc = complete_io(&io, cb, cb_arg);
	}

	D_DEBUG(DB_REBUILD, "process list buf "DF_UOID" rc %d\n",
		DP_UOID(io.ui_oid), rc);

out:
	dss_enum_unpack_io_fini(&io);
	return rc;
}
