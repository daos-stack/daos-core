/*
 *  (C) Copyright 2016-2020 Intel Corporation.
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
 * object shard operations.
 */
#define D_LOGFAC	DD_FAC(object)

#include <daos/container.h>
#include <daos/mgmt.h>
#include <daos/pool.h>
#include <daos/pool_map.h>
#include <daos/rpc.h>
#include <daos/checksum.h>
#include "obj_rpc.h"
#include "obj_internal.h"

static inline struct dc_obj_layout *
obj_shard2layout(struct dc_obj_shard *shard)
{
	return container_of(shard, struct dc_obj_layout,
			    do_shards[shard->do_shard]);
}

void
obj_shard_decref(struct dc_obj_shard *shard)
{
	struct dc_obj_layout	*layout;
	struct dc_object	*obj;
	bool			 release = false;

	D_ASSERT(shard != NULL);
	D_ASSERT(shard->do_ref > 0);
	D_ASSERT(shard->do_obj != NULL);

	obj = shard->do_obj;
	layout = obj_shard2layout(shard);

	D_SPIN_LOCK(&obj->cob_spin);
	if (--(shard->do_ref) == 0) {
		layout->do_open_count--;
		if (layout->do_open_count == 0 && layout != obj->cob_shards)
			release = true;
		shard->do_obj = NULL;
	}
	D_SPIN_UNLOCK(&obj->cob_spin);

	if (release)
		D_FREE(layout);
}

void
obj_shard_addref(struct dc_obj_shard *shard)
{
	D_ASSERT(shard->do_obj != NULL);
	D_SPIN_LOCK(&shard->do_obj->cob_spin);
	shard->do_ref++;
	D_SPIN_UNLOCK(&shard->do_obj->cob_spin);
}

int
dc_obj_shard_open(struct dc_object *obj, daos_unit_oid_t oid,
		  unsigned int mode, struct dc_obj_shard *shard)
{
	struct pool_target	*map_tgt;
	int			rc;

	D_ASSERT(obj != NULL && shard != NULL);
	D_ASSERT(shard->do_obj == NULL);

	rc = dc_cont_tgt_idx2ptr(obj->cob_coh, shard->do_target_id,
				 &map_tgt);
	if (rc)
		return rc;

	shard->do_id = oid;
	shard->do_target_rank = map_tgt->ta_comp.co_rank;
	shard->do_target_idx = map_tgt->ta_comp.co_index;
	shard->do_obj = obj;
	shard->do_co_hdl = obj->cob_coh;
	obj_shard_addref(shard);

	D_SPIN_LOCK(&obj->cob_spin);
	obj->cob_shards->do_open_count++;
	D_SPIN_UNLOCK(&obj->cob_spin);

	return 0;
}

void
dc_obj_shard_close(struct dc_obj_shard *shard)
{
	obj_shard_decref(shard);
}

struct rw_cb_args {
	crt_rpc_t		*rpc;
	daos_handle_t		*hdlp;
	d_sg_list_t		*rwaa_sgls;
	struct dc_obj_shard	*dobj;
	unsigned int		*map_ver;
	daos_iom_t		*maps;
	struct shard_rw_args	*shard_args;
};

static struct dcs_layout *
dc_rw_cb_singv_lo_get(daos_iod_t *iods, d_sg_list_t *sgls, uint32_t iod_nr,
		      struct obj_reasb_req *reasb_req)
{
	struct dcs_layout	*singv_lo, *singv_los;
	daos_iod_t		*iod;
	d_sg_list_t		*sgl;
	uint32_t		 i;

	if (reasb_req == NULL)
		return NULL;

	singv_los = reasb_req->orr_singv_los;
	for (i = 0; i < iod_nr; i++) {
		singv_lo = &singv_los[i];
		if (singv_lo->cs_even_dist == 0 || singv_lo->cs_bytes != 0)
			continue;
		/* the case of fetch singv with unknown rec size, now after the
		 * fetch need to re-calculate the singv_lo again
		 */
		iod = &iods[i];
		sgl = &sgls[i];
		D_ASSERT(iod->iod_size != DAOS_REC_ANY);
		if (obj_ec_singv_one_tgt(iod, sgl, reasb_req->orr_oca)) {
			singv_lo->cs_even_dist = 0;
			continue;
		}
		singv_lo->cs_bytes = obj_ec_singv_cell_bytes(iod->iod_size,
						reasb_req->orr_oca);
	}
	return singv_los;
}

int dc_rw_cb_csum_verify(const struct rw_cb_args *rw_args)
{
	struct daos_csummer	*csummer;
	d_sg_list_t		*sgls;
	struct obj_rw_in	*orw;
	struct obj_rw_out	*orwo;
	daos_iod_t		*iods;
	struct dcs_iod_csums	*iods_csums;
	daos_iom_t		*maps;
	struct dcs_layout	*singv_lo, *singv_los;
	uint32_t		 shard_idx;
	int			 i;
	int			 rc = 0;

	csummer = dc_cont_hdl2csummer(rw_args->dobj->do_co_hdl);
	if (!daos_csummer_initialized(csummer))
		return 0;

	orw = crt_req_get(rw_args->rpc);
	orwo = crt_reply_get(rw_args->rpc);
	sgls = rw_args->rwaa_sgls;
	iods = orw->orw_iod_array.oia_iods;
	iods_csums = orwo->orw_iod_csums.ca_arrays;
	maps = orwo->orw_maps.ca_arrays;

	if (daos_obj_is_echo(orw->orw_oid.id_pub))
		return 0; /** currently don't verify echo classes */

	if (sgls == NULL)
		return 0; /** no data to verify ... size only */
	D_ASSERTF(orwo->orw_maps.ca_count == orw->orw_iod_array.oia_iod_nr,
		  "orwo->orw_maps.ca_count(%lu) == "
		  "orw->orw_iod_array.oia_iod_nr(%d)",
		  orwo->orw_maps.ca_count, orw->orw_iod_array.oia_iod_nr);

	/** fault injection - corrupt data after getting from server and before
	 * verifying on client - simulates corruption over network
	 */
	if (DAOS_FAIL_CHECK(DAOS_CSUM_CORRUPT_FETCH))
		/** Got csum successfully from server. Now poison it!! */
		orwo->orw_iod_csums.ca_arrays->ic_data->cs_csum[0]++;

	shard_idx = rw_args->shard_args->auxi.shard -
		    rw_args->shard_args->auxi.start_shard;
	singv_los = dc_rw_cb_singv_lo_get(iods, sgls, orw->orw_nr,
					  rw_args->shard_args->reasb_req);
	for (i = 0; i < orw->orw_nr; i++) {
		daos_iod_t		*iod = &iods[i];
		struct dcs_iod_csums	*iod_csum = &iods_csums[i];
		daos_iom_t		*map = &maps[i];

		if (!csum_iod_is_supported(iod))
			continue;

		singv_lo = (singv_los == NULL) ? NULL : &singv_los[i];
		rc = daos_csummer_verify_iod(csummer, iod, &sgls[i],
					     iod_csum, singv_lo, shard_idx,
					     map);
		if (rc != 0) {
			if (iod->iod_type == DAOS_IOD_SINGLE) {
				D_ERROR("Data Verification failed (object: "
					DF_OID"): "DF_RC"\n",
					DP_OID(orw->orw_oid.id_pub),
					DP_RC(rc));
			} else  if (iod->iod_type == DAOS_IOD_ARRAY) {
				D_ERROR("Data Verification failed (object: "
						DF_OID" , extent: "DF_RECX"):"
						" "DF_RC"\n",
					DP_OID(orw->orw_oid.id_pub),
					DP_RECX(iod->iod_recxs[i]),
					DP_RC(rc));
			}

			break;
		}
	}

	return rc;
}

static void
daos_iom_copy(const daos_iom_t *src, daos_iom_t *dst)
{
	uint32_t	i;
	uint32_t	to_copy;

	dst->iom_type = src->iom_type;
	dst->iom_size = src->iom_size;
	dst->iom_recx_hi = src->iom_recx_hi;
	dst->iom_recx_lo = src->iom_recx_lo;
	dst->iom_nr_out = src->iom_nr_out;

	if (dst->iom_nr < dst->iom_nr_out)
		D_WARN("mapped recxs list will be truncated");

	to_copy = min(dst->iom_nr, dst->iom_nr_out);

	for (i = 0; i < to_copy ; i++)
		dst->iom_recxs[i] = src->iom_recxs[i];
}

static int
dc_rw_cb(tse_task_t *task, void *arg)
{
	struct rw_cb_args	*rw_args = arg;
	struct obj_rw_in	*orw;
	struct obj_rw_out	*orwo;
	daos_iod_t		*iods;
	uint64_t		*sizes;
	int			opc;
	int                     ret = task->dt_result;
	int			i, j;
	int			rc = 0;

	opc = opc_get(rw_args->rpc->cr_opc);
	D_DEBUG(DB_TRACE, "rpc %p opc:%d completed, dt_result %d.\n",
		rw_args->rpc, opc, ret);
	if (opc == DAOS_OBJ_RPC_FETCH &&
	    DAOS_FAIL_CHECK(DAOS_SHARD_OBJ_FETCH_TIMEOUT)) {
		D_ERROR("Inducing -DER_TIMEDOUT error on shard I/O fetch\n");
		D_GOTO(out, rc = -DER_TIMEDOUT);
	}
	if (opc == DAOS_OBJ_RPC_UPDATE &&
	    DAOS_FAIL_CHECK(DAOS_SHARD_OBJ_UPDATE_TIMEOUT)) {
		D_ERROR("Inducing -DER_TIMEDOUT error on shard I/O update\n");
		D_GOTO(out, rc = -DER_TIMEDOUT);
	}
	if (opc == DAOS_OBJ_RPC_UPDATE &&
	    DAOS_FAIL_CHECK(DAOS_OBJ_UPDATE_NOSPACE)) {
		D_ERROR("Inducing -DER_NOSPACE error on shard I/O update\n");
		D_GOTO(out, rc = -DER_NOSPACE);
	}

	orw = crt_req_get(rw_args->rpc);
	orwo = crt_reply_get(rw_args->rpc);
	D_ASSERT(orw != NULL && orwo != NULL);
	if (ret != 0) {
		/*
		 * If any failure happens inside Cart, let's reset failure to
		 * TIMEDOUT, so the upper layer can retry.
		 */
		D_ERROR("RPC %d failed: %d\n", opc, ret);
		D_GOTO(out, ret);
	}

	rc = obj_reply_get_status(rw_args->rpc);
	if (rc != 0) {
		if (rc == -DER_INPROGRESS) {
			D_DEBUG(DB_TRACE, "rpc %p opc %d to rank %d tag %d may "
				"need retry: "DF_RC"\n", rw_args->rpc, opc,
				rw_args->rpc->cr_ep.ep_rank,
				rw_args->rpc->cr_ep.ep_tag, DP_RC(rc));
		} else {
			D_ERROR("rpc %p opc %d to rank %d tag %d failed: "
				DF_RC"\n", rw_args->rpc, opc,
				rw_args->rpc->cr_ep.ep_rank,
				rw_args->rpc->cr_ep.ep_tag, DP_RC(rc));
			if (rc == -DER_REC2BIG && opc == DAOS_OBJ_RPC_FETCH) {
				/* update the sizes in iods */
				iods = orw->orw_iod_array.oia_iods;
				sizes = orwo->orw_iod_sizes.ca_arrays;
				for (i = 0; i < orw->orw_nr; i++)
					iods[i].iod_size = sizes[i];
			}
		}
		D_GOTO(out, rc);
	}
	*rw_args->map_ver = obj_reply_map_version_get(rw_args->rpc);

	if (opc == DAOS_OBJ_RPC_FETCH) {
		if (rw_args->maps != NULL && orwo->orw_maps.ca_count > 0) {
			/** Should have 1 map per iod */
			D_ASSERT(orwo->orw_maps.ca_count == orw->orw_nr);
			for (i = 0; i < orw->orw_nr; i++) {
				daos_iom_copy(&orwo->orw_maps.ca_arrays[i],
					      &rw_args->maps[i]);
			}
		}

		bool	is_ec_obj = false;

		rc = obj_ec_recov_add(rw_args->shard_args->reasb_req,
				      orwo->orw_rels.ca_arrays,
				      orwo->orw_rels.ca_count);
		if (rc) {
			D_ERROR("fail to add recov list for "DF_UOID",rc %d.\n",
				DP_UOID(orw->orw_oid), rc);
			goto out;
		}

		iods = orw->orw_iod_array.oia_iods;
		sizes = orwo->orw_iod_sizes.ca_arrays;

		if (rw_args->shard_args->reasb_req != NULL &&
		    DAOS_OC_IS_EC(rw_args->shard_args->reasb_req->orr_oca))
			is_ec_obj = true;

		if (orwo->orw_iod_sizes.ca_count != orw->orw_nr) {
			D_ERROR("out:%u != in:%u for "DF_UOID" with eph "
				DF_U64".\n",
				(unsigned)orwo->orw_iod_sizes.ca_count,
				orw->orw_nr, DP_UOID(orw->orw_oid),
				orw->orw_epoch);
			D_GOTO(out, rc = -DER_PROTO);
		}

		/* update the sizes in iods */
		for (i = 0; i < orw->orw_nr; i++)
			iods[i].iod_size = sizes[i];

		if (orwo->orw_sgls.ca_count > 0) {
			/* inline transfer */
			rc = daos_sgls_copy_data_out(rw_args->rwaa_sgls,
						     orw->orw_nr,
						     orwo->orw_sgls.ca_arrays,
						     orwo->orw_sgls.ca_count);
		} else if (rw_args->rwaa_sgls != NULL && !is_ec_obj) {
			/* XXX need some extra handling for EC obj */
			/* for bulk transfer it needs to update sg_nr_out */
			d_sg_list_t	*sgls = rw_args->rwaa_sgls;
			d_iov_t		*iov;
			uint32_t	*nrs;
			uint32_t	 nrs_count;
			daos_size_t	*replied_sizes;
			daos_size_t	 data_size, size_in_iod;
			daos_size_t	 buf_size;

			nrs = orwo->orw_nrs.ca_arrays;
			nrs_count = orwo->orw_nrs.ca_count;
			replied_sizes = orwo->orw_data_sizes.ca_arrays;
			if (nrs_count != orw->orw_nr) {
				D_ERROR("Invalid nrs %u != %u\n", nrs_count,
					orw->orw_nr);
				D_GOTO(out, rc = -DER_PROTO);
			}

			for (i = 0; i < orw->orw_nr; i++) {
				/* server returned bs_nr_out is only to check
				 * if it is empty record in that case just set
				 * sg_nr_out as zero, or will set sg_nr_out and
				 * iov_len by checking with iods as server
				 * filled the buffer from beginning.
				 */
				if (nrs[i] == 0) {
					sgls[i].sg_nr_out = 0;
					continue;
				}
				size_in_iod = daos_iods_len(&iods[i], 1);
				if (size_in_iod == -1) {
					/* only for echo mode */
					sgls[i].sg_nr_out = sgls[i].sg_nr;
					continue;
				}
				data_size = replied_sizes[i];
				D_ASSERT(data_size <= size_in_iod);
				buf_size = 0;
				for (j = 0; j < sgls[i].sg_nr; j++) {
					iov = &sgls[i].sg_iovs[j];
					buf_size += iov->iov_buf_len;
					if (buf_size < data_size) {
						iov->iov_len = iov->iov_buf_len;
						continue;
					}

					iov->iov_len = iov->iov_buf_len -
						       (buf_size - data_size);
					sgls[i].sg_nr_out = j + 1;
					break;
				}
			}
		}
		if (rc != 0)
			goto out;

		rc = dc_rw_cb_csum_verify(rw_args);
		if (rc != 0)
			goto out;

	}
out:
	crt_req_decref(rw_args->rpc);
	obj_shard_decref(rw_args->dobj);
	dc_pool_put((struct dc_pool *)rw_args->hdlp);

	if (ret == 0 || obj_retry_error(rc))
		ret = rc;
	return ret;
}

static struct dc_pool *
obj_shard_ptr2pool(struct dc_obj_shard *shard)
{
	daos_handle_t poh;

	poh = dc_cont_hdl2pool_hdl(shard->do_co_hdl);
	if (daos_handle_is_inval(poh))
		return NULL;

	return dc_hdl2pool(poh);
}

int
dc_obj_shard_rw(struct dc_obj_shard *shard, enum obj_rpc_opc opc,
		void *shard_args, struct daos_shard_tgt *fw_shard_tgts,
		uint32_t fw_cnt, tse_task_t *task)
{
	struct shard_rw_args	*args = shard_args;
	struct shard_auxi_args	*auxi = &args->auxi;
	daos_obj_rw_t		*api_args = args->api_args;
	struct dc_pool		*pool;
	daos_key_t		*dkey = api_args->dkey;
	unsigned int		 nr = api_args->nr;
	d_sg_list_t		*sgls = api_args->sgls;
	crt_rpc_t		*req = NULL;
	struct obj_rw_in	*orw;
	struct rw_cb_args	 rw_args;
	crt_endpoint_t		 tgt_ep;
	uuid_t			 cont_hdl_uuid;
	uuid_t			 cont_uuid;
	bool			 cb_registered = false;
	uint32_t		 flags = 0;
	int			 rc;

	obj_shard_addref(shard);

	if (DAOS_FAIL_CHECK(DAOS_SHARD_OBJ_UPDATE_TIMEOUT_SINGLE)) {
		if (auxi->shard == daos_fail_value_get()) {
			D_INFO("Set Shard %d update to return -DER_TIMEDOUT\n",
			       auxi->shard);
			daos_fail_loc_set(DAOS_SHARD_OBJ_UPDATE_TIMEOUT |
					  DAOS_FAIL_ONCE);
		}
	}
	if (DAOS_FAIL_CHECK(DAOS_OBJ_TGT_IDX_CHANGE)) {
		if (srv_io_mode == DIM_CLIENT_DISPATCH) {
			/* to trigger retry on all other shards */
			if (auxi->shard != daos_fail_value_get()) {
				D_INFO("complete shard %d update as "
				       "-DER_TIMEDOUT.\n", auxi->shard);
				D_GOTO(out_obj, rc = -DER_TIMEDOUT);
			}
		} else {
			flags = ORF_DTX_SYNC;
		}
	}

	if (auxi->epoch.oe_uncertain)
		flags |= ORF_EPOCH_UNCERTAIN;

	rc = dc_cont_hdl2uuid(shard->do_co_hdl, &cont_hdl_uuid, &cont_uuid);
	if (rc != 0)
		D_GOTO(out_obj, rc);

	pool = obj_shard_ptr2pool(shard);
	if (pool == NULL)
		D_GOTO(out_obj, rc = -DER_NO_HDL);

	tgt_ep.ep_grp = pool->dp_sys->sy_group;
	tgt_ep.ep_tag = shard->do_target_idx;
	tgt_ep.ep_rank = shard->do_target_rank;
	if ((int)tgt_ep.ep_rank < 0)
		D_GOTO(out_pool, rc = (int)tgt_ep.ep_rank);

	rc = obj_req_create(daos_task2ctx(task), &tgt_ep, opc, &req);
	D_DEBUG(DB_TRACE, "rpc %p opc:%d "DF_UOID" %d %s rank:%d tag:%d eph "
		DF_U64"\n", req, opc, DP_UOID(shard->do_id), (int)dkey->iov_len,
		(char *)dkey->iov_buf, tgt_ep.ep_rank, tgt_ep.ep_tag,
		auxi->epoch.oe_value);
	if (rc != 0)
		D_GOTO(out_pool, rc);

	if (DAOS_FAIL_CHECK(DAOS_SHARD_OBJ_FAIL))
		D_GOTO(out_req, rc = -DER_INVAL);

	orw = crt_req_get(req);
	D_ASSERT(orw != NULL);

	if (fw_shard_tgts != NULL) {
		D_ASSERT(fw_cnt >= 1);
		orw->orw_shard_tgts.ca_count = fw_cnt;
		orw->orw_shard_tgts.ca_arrays = fw_shard_tgts;
	} else {
		orw->orw_shard_tgts.ca_count = 0;
		orw->orw_shard_tgts.ca_arrays = NULL;
	}
	orw->orw_map_ver = auxi->map_ver;
	orw->orw_start_shard = auxi->start_shard;
	orw->orw_oid = shard->do_id;
	uuid_copy(orw->orw_pool_uuid, pool->dp_pool);
	uuid_copy(orw->orw_co_hdl, cont_hdl_uuid);
	uuid_copy(orw->orw_co_uuid, cont_uuid);
	daos_dti_copy(&orw->orw_dti, &args->dti);
	orw->orw_flags = auxi->flags | flags;
	if (obj_op_is_ec_fetch(auxi->obj_auxi) &&
	    (auxi->shard != (auxi->start_shard + auxi->ec_tgt_idx))) {
		orw->orw_flags |= ORF_EC_DEGRADED;
		orw->orw_tgt_idx = auxi->ec_tgt_idx;
	}
	orw->orw_dti_cos.ca_count = 0;
	orw->orw_dti_cos.ca_arrays = NULL;

	orw->orw_api_flags = api_args->flags;
	orw->orw_epoch = auxi->epoch.oe_value;
	orw->orw_dkey_hash = args->dkey_hash;
	orw->orw_nr = nr;
	orw->orw_dkey = *dkey;
	orw->orw_dkey_csum = args->dkey_csum;
	orw->orw_iod_array.oia_iod_nr = nr;
	orw->orw_iod_array.oia_iods = api_args->iods;
	orw->orw_iod_array.oia_iod_csums = args->iod_csums;
	orw->orw_iod_array.oia_oiods = args->oiods;
	orw->orw_iod_array.oia_oiod_nr = (args->oiods == NULL) ?
					 0 : nr;
	orw->orw_iod_array.oia_offs = args->offs;

	D_DEBUG(DB_TRACE, "opc %d "DF_UOID" %d %s rank %d tag %d eph "
		DF_U64", DTI = "DF_DTI"\n", opc, DP_UOID(shard->do_id),
		(int)dkey->iov_len, (char *)dkey->iov_buf, tgt_ep.ep_rank,
		tgt_ep.ep_tag, auxi->epoch.oe_value, DP_DTI(&orw->orw_dti));

	if (args->bulks != NULL) {
		orw->orw_sgls.ca_count = 0;
		orw->orw_sgls.ca_arrays = NULL;
		orw->orw_bulks.ca_count = nr;
		orw->orw_bulks.ca_arrays = args->bulks;
		if (fw_shard_tgts != NULL)
			orw->orw_flags |= ORF_BULK_BIND;
	} else {
		/* Transfer data inline */
		if (sgls != NULL)
			orw->orw_sgls.ca_count = nr;
		else
			orw->orw_sgls.ca_count = 0;
		orw->orw_sgls.ca_arrays = sgls;
		orw->orw_bulks.ca_count = 0;
		orw->orw_bulks.ca_arrays = NULL;
	}

	crt_req_addref(req);
	rw_args.rpc = req;
	rw_args.hdlp = (daos_handle_t *)pool;
	rw_args.map_ver = &auxi->map_ver;
	rw_args.dobj = shard;
	rw_args.shard_args = args;
	/* remember the sgl to copyout the data inline for fetch */
	rw_args.rwaa_sgls = (opc == DAOS_OBJ_RPC_FETCH) ? sgls : NULL;
	rw_args.maps = args->api_args->maps;
	if (opc == DAOS_OBJ_RPC_FETCH &&
	    (rw_args.maps != NULL || args->iod_csums != NULL))
		orw->orw_flags |= ORF_CREATE_MAP;

	if (DAOS_FAIL_CHECK(DAOS_SHARD_OBJ_RW_CRT_ERROR))
		D_GOTO(out_args, rc = -DER_HG);

	rc = tse_task_register_comp_cb(task, dc_rw_cb, &rw_args,
				       sizeof(rw_args));
	if (rc != 0)
		D_GOTO(out_args, rc);
	cb_registered = true;

	if (daos_io_bypass & IOBP_CLI_RPC) {
		rc = daos_rpc_complete(req, task);
	} else {
		rc = daos_rpc_send(req, task);
		if (rc != 0) {
			D_ERROR("update/fetch rpc failed rc "DF_RC"\n",
				DP_RC(rc));
			D_GOTO(out_args, rc);
		}
	}
	return rc;

out_args:
	crt_req_decref(req);
out_req:
	crt_req_decref(req);
out_pool:
	dc_pool_put(pool);
out_obj:
	if (!cb_registered)
		obj_shard_decref(shard);
	tse_task_complete(task, rc);
	return rc;
}

struct obj_punch_cb_args {
	crt_rpc_t	*rpc;
	unsigned int	*map_ver;
};

static int
obj_shard_punch_cb(tse_task_t *task, void *data)
{
	struct obj_punch_cb_args	*cb_args;
	crt_rpc_t			*rpc;

	cb_args = (struct obj_punch_cb_args *)data;
	rpc = cb_args->rpc;
	if (task->dt_result == 0) {
		task->dt_result = obj_reply_get_status(rpc);
		*cb_args->map_ver = obj_reply_map_version_get(rpc);
	}

	crt_req_decref(rpc);
	return task->dt_result;
}

int
dc_obj_shard_punch(struct dc_obj_shard *shard, enum obj_rpc_opc opc,
		   void *shard_args, struct daos_shard_tgt *fw_shard_tgts,
		   uint32_t fw_cnt, tse_task_t *task)
{
	struct shard_punch_args		*args = shard_args;
	daos_obj_punch_t		*obj_args = args->pa_api_args;
	daos_key_t			*dkey = obj_args->dkey;
	struct dc_pool			*pool;
	struct obj_punch_in		*opi;
	crt_rpc_t			*req;
	struct obj_punch_cb_args	 cb_args;
	daos_unit_oid_t			 oid;
	crt_endpoint_t			 tgt_ep;
	int				 rc;

	pool = obj_shard_ptr2pool(shard);
	if (pool == NULL)
		D_GOTO(out, rc = -DER_NO_HDL);

	oid = shard->do_id;
	tgt_ep.ep_grp	= pool->dp_sys->sy_group;
	tgt_ep.ep_tag	= shard->do_target_idx;
	tgt_ep.ep_rank = shard->do_target_rank;
	if ((int)tgt_ep.ep_rank < 0)
		D_GOTO(out, rc = (int)tgt_ep.ep_rank);

	D_DEBUG(DB_IO, "opc=%d, rank=%d tag=%d epoch "DF_U64".\n",
		 opc, tgt_ep.ep_rank, tgt_ep.ep_tag,
		 args->pa_auxi.epoch.oe_value);

	rc = obj_req_create(daos_task2ctx(task), &tgt_ep, opc, &req);
	if (rc != 0)
		D_GOTO(out, rc);

	crt_req_addref(req);
	cb_args.rpc = req;
	cb_args.map_ver = &args->pa_auxi.map_ver;
	rc = tse_task_register_comp_cb(task, obj_shard_punch_cb, &cb_args,
				       sizeof(cb_args));
	if (rc != 0)
		D_GOTO(out_req, rc);

	opi = crt_req_get(req);
	D_ASSERT(opi != NULL);

	opi->opi_map_ver	 = args->pa_auxi.map_ver;
	opi->opi_api_flags	 = obj_args->flags;
	opi->opi_epoch		 = args->pa_auxi.epoch.oe_value;
	opi->opi_dkey_hash	 = args->pa_dkey_hash;
	opi->opi_oid		 = oid;
	opi->opi_dkeys.ca_count  = (dkey == NULL) ? 0 : 1;
	opi->opi_dkeys.ca_arrays = dkey;
	opi->opi_akeys.ca_count	 = obj_args->akey_nr;
	opi->opi_akeys.ca_arrays = obj_args->akeys;
	if (fw_shard_tgts != NULL) {
		D_ASSERT(fw_cnt >= 1);
		opi->opi_shard_tgts.ca_count = fw_cnt;
		opi->opi_shard_tgts.ca_arrays = fw_shard_tgts;
	} else {
		opi->opi_shard_tgts.ca_count = 0;
		opi->opi_shard_tgts.ca_arrays = NULL;
	}
	uuid_copy(opi->opi_pool_uuid, pool->dp_pool);
	uuid_copy(opi->opi_co_hdl, args->pa_coh_uuid);
	uuid_copy(opi->opi_co_uuid, args->pa_cont_uuid);
	daos_dti_copy(&opi->opi_dti, &args->pa_dti);
	opi->opi_flags = args->pa_auxi.flags;
	if (args->pa_auxi.epoch.oe_uncertain)
		opi->opi_flags |= ORF_EPOCH_UNCERTAIN;
	opi->opi_dti_cos.ca_count = 0;
	opi->opi_dti_cos.ca_arrays = NULL;

	rc = daos_rpc_send(req, task);
	if (rc != 0) {
		D_ERROR("punch rpc failed rc "DF_RC"\n", DP_RC(rc));
		D_GOTO(out_req, rc);
	}

	dc_pool_put(pool);
	return 0;

out_req:
	crt_req_decref(req);
out:
	if (pool != NULL)
		dc_pool_put(pool);
	tse_task_complete(task, rc);
	return rc;
}

struct obj_enum_args {
	crt_rpc_t		*rpc;
	daos_handle_t		*hdlp;
	uint32_t		*eaa_nr;
	daos_key_desc_t		*eaa_kds;
	daos_anchor_t		*eaa_anchor;
	daos_anchor_t		*eaa_dkey_anchor;
	daos_anchor_t		*eaa_akey_anchor;
	struct dc_obj_shard	*eaa_obj;
	d_sg_list_t		*eaa_sgl;
	daos_recx_t		*eaa_recxs;
	daos_size_t		*eaa_size;
	unsigned int		*eaa_map_ver;
	d_iov_t			*csum;
};

/**
 * use iod/iod_csum as vehicle to verify data
 */
static int
csum_enum_verify_recx(struct daos_csummer *csummer,
		      struct obj_enum_rec *rec,
		      struct dcs_csum_info *csum_info,
		      d_iov_t *enum_type_val)
{
	daos_iod_t		 tmp_iod = {0};
	d_sg_list_t		 tmp_sgl = {0};
	struct dcs_iod_csums	 tmp_iod_csum = {0};

	tmp_iod.iod_size = rec->rec_size;
	tmp_iod.iod_type = DAOS_IOD_ARRAY;
	tmp_iod.iod_recxs = &rec->rec_recx;
	tmp_iod.iod_nr = 1;

	tmp_sgl.sg_nr = tmp_sgl.sg_nr_out = 1;
	tmp_sgl.sg_iovs = enum_type_val;

	tmp_iod_csum.ic_nr = 1;
	tmp_iod_csum.ic_data = csum_info;

	return daos_csummer_verify_iod(csummer,
				       &tmp_iod, &tmp_sgl, &tmp_iod_csum,
				       NULL, 0, NULL);
}

/**
 * use iod/iod_csum as vehicle to verify data
 */
static int
csum_enum_verify_sv(struct daos_csummer *csummer,
		 d_iov_t *enum_type_val, d_iov_t *csum_iov)
{
	daos_iod_t		 tmp_iod = {0};
	d_sg_list_t		 tmp_sgl = {0};
	struct dcs_iod_csums	 tmp_iod_csum = {0};
	struct dcs_csum_info	*tmp_csum_info = {0};

	tmp_iod.iod_size = enum_type_val->iov_len;
	tmp_iod.iod_type = DAOS_IOD_SINGLE;
	tmp_iod.iod_nr = 1;

	tmp_sgl.sg_nr = tmp_sgl.sg_nr_out = 1;
	tmp_sgl.sg_iovs = enum_type_val;

	ci_cast(&tmp_csum_info, csum_iov);
	ci_move_next_iov(tmp_csum_info, csum_iov);

	tmp_iod_csum.ic_nr = 1;
	tmp_iod_csum.ic_data = tmp_csum_info;

	return daos_csummer_verify_iod(csummer,
				       &tmp_iod, &tmp_sgl, &tmp_iod_csum,
				       NULL, 0, NULL);
}

static int
csum_enum_verify(const struct obj_enum_args *enum_args,
		 const struct obj_key_enum_out *oeo)
{
	struct daos_csummer	*csummer;
	uint64_t		 i;
	int			 rc = 0;
	struct daos_sgl_idx	 sgl_idx = {0};
	d_sg_list_t		 sgl = oeo->oeo_sgl;
	d_iov_t			 csum_iov = oeo->oeo_csum_iov;
	struct dcs_csum_info	 *tmp = NULL;

	if (enum_args->eaa_nr == NULL ||
	    *enum_args->eaa_nr == 0 ||
	    sgl.sg_nr_out == 0)
		return 0; /** no keys to verify */

	csummer = dc_cont_hdl2csummer(enum_args->eaa_obj->do_co_hdl);
	if (!daos_csummer_initialized(csummer))
		return 0; /** csums not enabled */

	if (csum_iov.iov_len == 0) {
		D_ERROR("CSUM is enabled but no  checksum provided.");
		return -DER_CSUM;
	}

	for (i = 0; i < *enum_args->eaa_nr; i++) {
		daos_key_desc_t		*kd = &enum_args->eaa_kds[i];
		void			*buf;
		d_iov_t			 enum_type_val;
		d_iov_t			 iov = sgl.sg_iovs[sgl_idx.iov_idx];

		buf = iov.iov_buf + sgl_idx.iov_offset;

		switch (kd->kd_val_type) {
		case OBJ_ITER_RECX: {
			struct obj_enum_rec *rec = buf;

			/**
			 * Even if don't use csum info at this point because
			 * the data isn't inline, still need to move to next
			 */
			ci_cast(&tmp, &csum_iov);
			ci_move_next_iov(tmp, &csum_iov);

			if (rec->rec_flags & RECX_INLINE) {
				buf += sizeof(*rec);
				d_iov_set(&enum_type_val, buf,
					  rec->rec_size * rec->rec_recx.rx_nr);
				rc = csum_enum_verify_recx(csummer, rec, tmp,
							   &enum_type_val);
				if (rc != 0)
					return rc;
			}
			break;
		}
		case OBJ_ITER_SINGLE: {
			d_iov_set(&enum_type_val, buf, kd->kd_key_len);

			ci_cast(&tmp, &csum_iov);
			ci_move_next_iov(tmp, &csum_iov);

			rc = csum_enum_verify_sv(csummer,
						 &enum_type_val, &csum_iov);

			if (rc != 0)
				return rc;
			break;
		}
		case OBJ_ITER_AKEY:
		case OBJ_ITER_DKEY:
			d_iov_set(&enum_type_val, buf, kd->kd_key_len);
			/**
			  * fault injection - corrupt keys before verifying -
			  * simulates corruption over network
			  */
			if (DAOS_FAIL_CHECK(DAOS_CSUM_CORRUPT_FETCH_AKEY) ||
			    DAOS_FAIL_CHECK(DAOS_CSUM_CORRUPT_FETCH_DKEY))
				((uint8_t *)buf)[0] += 2;

			ci_cast(&tmp, &csum_iov);
			ci_move_next_iov(tmp, &csum_iov);

			rc = daos_csummer_verify_key(csummer,
						     &enum_type_val, tmp);

			if (rc != 0)
				return rc;
			break;
		}

		sgl_idx.iov_offset += kd->kd_key_len;

		/** move to next iov if necessary */
		if (sgl_idx.iov_offset >= iov.iov_len) {
			sgl_idx.iov_idx++;
			sgl_idx.iov_offset = 0;
		}
	}
	return rc;
}

/**
 * If requested (dst iov is set) and there is csum info to copy, copy the
 * serialized csum. If not all of it will fit into the provided buffer, copy
 * what can and set the destination iov len to needed len and let caller
 * decide what to do.
 */
static int
dc_enumerate_copy_csum(d_iov_t *dst, const d_iov_t *src)
{
	if (dst != NULL && src->iov_len > 0) {
		memcpy(dst->iov_buf, src->iov_buf,
		       min(dst->iov_buf_len,
			   src->iov_len));
		dst->iov_len = src->iov_len;
		if (dst->iov_len > dst->iov_buf_len)
			return -DER_TRUNC;
	}
	return 0;
}

static int
dc_enumerate_cb(tse_task_t *task, void *arg)
{
	struct obj_enum_args	*enum_args = (struct obj_enum_args *)arg;
	struct obj_key_enum_in	*oei;
	struct obj_key_enum_out	*oeo;
	int			 opc = opc_get(enum_args->rpc->cr_opc);
	int			 ret = task->dt_result;
	int			 rc = 0;

	oei = crt_req_get(enum_args->rpc);
	D_ASSERT(oei != NULL);

	if (ret != 0) {
		/* If any failure happens inside Cart, let's reset
		 * failure to TIMEDOUT, so the upper layer can retry
		 **/
		D_ERROR("RPC %d failed: %d\n", opc, ret);
		D_GOTO(out, ret);
	}

	oeo = crt_reply_get(enum_args->rpc);
	rc = obj_reply_get_status(enum_args->rpc);

	if (rc != 0) {
		if (rc == -DER_KEY2BIG) {
			D_DEBUG(DB_IO, "key size "DF_U64" too big.\n",
				oeo->oeo_size);
			enum_args->eaa_kds[0].kd_key_len = oeo->oeo_size;
		} else if (rc == -DER_INPROGRESS) {
			D_DEBUG(DB_TRACE, "rpc %p RPC %d may need retry: "
				""DF_RC"\n", enum_args->rpc, opc, DP_RC(rc));
		} else {
			D_ERROR("rpc %p RPC %d failed: "DF_RC"\n",
				enum_args->rpc, opc, DP_RC(rc));
		}
		D_GOTO(out, rc);
	}

	rc = dc_enumerate_copy_csum(enum_args->csum, &oeo->oeo_csum_iov);
	if (rc != 0)
		D_GOTO(out, rc);

	*enum_args->eaa_map_ver = obj_reply_map_version_get(enum_args->rpc);

	if (enum_args->eaa_size)
		*enum_args->eaa_size = oeo->oeo_size;

	if (*enum_args->eaa_nr < oeo->oeo_num) {
		D_ERROR("key enumerate get %d > %d more kds, %d\n",
			oeo->oeo_num, *enum_args->eaa_nr, -DER_PROTO);
		D_GOTO(out, rc = -DER_PROTO);
	}

	*enum_args->eaa_nr = oeo->oeo_num;

	if (enum_args->eaa_kds && oeo->oeo_kds.ca_count > 0)
		memcpy(enum_args->eaa_kds, oeo->oeo_kds.ca_arrays,
		       sizeof(*enum_args->eaa_kds) *
		       oeo->oeo_kds.ca_count);

	if (enum_args->eaa_recxs && oeo->oeo_recxs.ca_count > 0) {
		D_ASSERT(*enum_args->eaa_nr >= oeo->oeo_recxs.ca_count);
		memcpy(enum_args->eaa_recxs, oeo->oeo_recxs.ca_arrays,
		       sizeof(*enum_args->eaa_recxs) *
		       oeo->oeo_recxs.ca_count);
	}

	if (enum_args->eaa_sgl && oeo->oeo_sgl.sg_nr > 0) {
		rc = daos_sgl_copy_data_out(enum_args->eaa_sgl, &oeo->oeo_sgl);
		if (rc)
			D_GOTO(out, rc);
	}

	/* Update dkey hash and tag */
	if (enum_args->eaa_dkey_anchor)
		enum_anchor_copy(enum_args->eaa_dkey_anchor,
				 &oeo->oeo_dkey_anchor);

	if (enum_args->eaa_akey_anchor)
		enum_anchor_copy(enum_args->eaa_akey_anchor,
				 &oeo->oeo_akey_anchor);

	if (enum_args->eaa_anchor)
		enum_anchor_copy(enum_args->eaa_anchor,
				 &oeo->oeo_anchor);
	rc = csum_enum_verify(enum_args, oeo);
	if (rc != 0)
		D_GOTO(out, rc);

out:
	if (enum_args->eaa_obj != NULL)
		obj_shard_decref(enum_args->eaa_obj);

	if (oei->oei_bulk != NULL)
		crt_bulk_free(oei->oei_bulk);
	if (oei->oei_kds_bulk != NULL)
		crt_bulk_free(oei->oei_kds_bulk);
	crt_req_decref(enum_args->rpc);
	dc_pool_put((struct dc_pool *)enum_args->hdlp);

	if (ret == 0 || obj_retry_error(rc))
		ret = rc;
	return ret;
}

#define KDS_BULK_LIMIT	128

int
dc_obj_shard_list(struct dc_obj_shard *obj_shard, enum obj_rpc_opc opc,
		  void *shard_args, struct daos_shard_tgt *fw_shard_tgts,
		  uint32_t fw_cnt, tse_task_t *task)
{
	struct shard_list_args *args = shard_args;
	daos_obj_list_t	       *obj_args = args->la_api_args;
	daos_key_desc_t	       *kds = obj_args->kds;
	d_sg_list_t	       *sgl = obj_args->sgl;
	crt_endpoint_t		tgt_ep;
	struct dc_pool	       *pool = NULL;
	crt_rpc_t	       *req;
	uuid_t			cont_hdl_uuid;
	uuid_t			cont_uuid;
	struct obj_key_enum_in	*oei;
	struct obj_enum_args	enum_args;
	daos_size_t		sgl_size = 0;
	bool			cb_registered = false;
	int			rc;

	D_ASSERT(obj_shard != NULL);
	obj_shard_addref(obj_shard);

	rc = dc_cont_hdl2uuid(obj_shard->do_co_hdl, &cont_hdl_uuid, &cont_uuid);
	if (rc != 0)
		D_GOTO(out_put, rc);

	pool = obj_shard_ptr2pool(obj_shard);
	if (pool == NULL)
		D_GOTO(out_put, rc = -DER_NO_HDL);

	tgt_ep.ep_grp = pool->dp_sys->sy_group;
	tgt_ep.ep_tag = obj_shard->do_target_idx;
	tgt_ep.ep_rank = obj_shard->do_target_rank;
	if ((int)tgt_ep.ep_rank < 0)
		D_GOTO(out_pool, rc = (int)tgt_ep.ep_rank);

	D_DEBUG(DB_IO, "opc %d "DF_UOID" rank %d tag %d\n",
		opc, DP_UOID(obj_shard->do_id), tgt_ep.ep_rank, tgt_ep.ep_tag);

	rc = obj_req_create(daos_task2ctx(task), &tgt_ep, opc, &req);
	if (rc != 0)
		D_GOTO(out_pool, rc);

	oei = crt_req_get(req);
	D_ASSERT(oei != NULL);

	if (obj_args->dkey != NULL)
		oei->oei_dkey = *obj_args->dkey;
	if (obj_args->akey != NULL)
		oei->oei_akey = *obj_args->akey;
	oei->oei_oid		= obj_shard->do_id;
	oei->oei_map_ver	= args->la_auxi.map_ver;
	if (args->la_auxi.epoch.oe_uncertain)
		oei->oei_flags |= ORF_EPOCH_UNCERTAIN;
	if (obj_args->eprs != NULL && opc == DAOS_OBJ_RPC_ENUMERATE) {
		oei->oei_epr = *obj_args->eprs;
		/*
		 * If an epoch range is specified, we shall not assume any
		 * epoch uncertainty.
		 */
		oei->oei_flags &= ~ORF_EPOCH_UNCERTAIN;
	} else {
		oei->oei_epr.epr_lo = 0;
		oei->oei_epr.epr_hi = args->la_auxi.epoch.oe_value;
	}

	oei->oei_nr		= *obj_args->nr;
	oei->oei_rec_type	= obj_args->type;
	uuid_copy(oei->oei_pool_uuid, pool->dp_pool);
	uuid_copy(oei->oei_co_hdl, cont_hdl_uuid);
	uuid_copy(oei->oei_co_uuid, cont_uuid);
	daos_dti_copy(&oei->oei_dti, &args->la_dti);

	if (obj_args->anchor != NULL)
		enum_anchor_copy(&oei->oei_anchor, obj_args->anchor);
	if (obj_args->dkey_anchor != NULL)
		enum_anchor_copy(&oei->oei_dkey_anchor, obj_args->dkey_anchor);
	if (obj_args->akey_anchor != NULL)
		enum_anchor_copy(&oei->oei_akey_anchor, obj_args->akey_anchor);

	if (sgl != NULL) {
		oei->oei_sgl = *sgl;
		sgl_size = daos_sgls_packed_size(sgl, 1, NULL);
		if (sgl_size >= OBJ_BULK_LIMIT) {
			/* Create bulk */
			rc = crt_bulk_create(daos_task2ctx(task),
					     sgl, CRT_BULK_RW,
					     &oei->oei_bulk);
			if (rc < 0)
				D_GOTO(out_req, rc);
		}
	}

	if (*obj_args->nr > KDS_BULK_LIMIT) {
		d_sg_list_t	tmp_sgl = { 0 };
		d_iov_t		tmp_iov = { 0 };

		tmp_iov.iov_buf_len = sizeof(*kds) * (*obj_args->nr);
		tmp_iov.iov_buf = kds;
		tmp_sgl.sg_nr_out = 1;
		tmp_sgl.sg_nr = 1;
		tmp_sgl.sg_iovs = &tmp_iov;

		rc = crt_bulk_create(daos_task2ctx(task),
				     &tmp_sgl, CRT_BULK_RW,
				     &oei->oei_kds_bulk);
		if (rc < 0)
			D_GOTO(out_req, rc);
	}

	crt_req_addref(req);
	enum_args.rpc = req;
	enum_args.hdlp = (daos_handle_t *)pool;
	enum_args.eaa_nr = obj_args->nr;
	enum_args.eaa_kds = kds;
	enum_args.eaa_anchor = obj_args->anchor;
	enum_args.eaa_dkey_anchor = obj_args->dkey_anchor;
	enum_args.eaa_akey_anchor = obj_args->akey_anchor;
	enum_args.eaa_obj = obj_shard;
	enum_args.eaa_size = obj_args->size;
	enum_args.eaa_sgl = sgl;
	enum_args.csum = obj_args->csum;
	enum_args.eaa_map_ver = &args->la_auxi.map_ver;
	enum_args.eaa_recxs = obj_args->recxs;
	rc = tse_task_register_comp_cb(task, dc_enumerate_cb, &enum_args,
				       sizeof(enum_args));
	if (rc != 0)
		D_GOTO(out_eaa, rc);
	cb_registered = true;

	rc = daos_rpc_send(req, task);
	if (rc != 0) {
		D_ERROR("enumerate rpc failed rc "DF_RC"\n", DP_RC(rc));
		D_GOTO(out_eaa, rc);
	}

	return rc;

out_eaa:
	crt_req_decref(req);
	if (sgl != NULL && sgl_size >= OBJ_BULK_LIMIT)
		crt_bulk_free(oei->oei_bulk);
out_req:
	crt_req_decref(req);
out_pool:
	dc_pool_put(pool);
out_put:
	if (!cb_registered)
		obj_shard_decref(obj_shard);
	tse_task_complete(task, rc);
	return rc;
}

struct obj_query_key_cb_args {
	crt_rpc_t		*rpc;
	unsigned int		*map_ver;
	daos_unit_oid_t		oid;
	daos_key_t		*dkey;
	daos_key_t		*akey;
	daos_recx_t		*recx;
	struct dc_object	*obj;
};

static int
obj_shard_query_key_cb(tse_task_t *task, void *data)
{
	struct obj_query_key_cb_args	*cb_args;
	struct obj_query_key_in		*okqi;
	struct obj_query_key_out	*okqo;
	uint32_t			flags;
	int				opc;
	int				ret = task->dt_result;
	int				rc = 0;
	crt_rpc_t			*rpc;

	cb_args = (struct obj_query_key_cb_args *)data;
	rpc = cb_args->rpc;

	okqi = crt_req_get(cb_args->rpc);
	D_ASSERT(okqi != NULL);

	flags = okqi->okqi_api_flags;
	opc = opc_get(cb_args->rpc->cr_opc);

	if (ret != 0) {
		D_ERROR("RPC %d failed: %d\n", opc, ret);
		D_GOTO(out, ret);
	}

	okqo = crt_reply_get(cb_args->rpc);
	rc = obj_reply_get_status(rpc);
	if (rc != 0) {
		if (rc == -DER_NONEXIST)
			D_GOTO(out, rc = 0);

		if (rc == -DER_INPROGRESS)
			D_DEBUG(DB_TRACE, "rpc %p RPC %d may need retry: %d\n",
				cb_args->rpc, opc, rc);
		else
			D_ERROR("rpc %p RPC %d failed: %d\n",
				cb_args->rpc, opc, rc);
		D_GOTO(out, rc);
	}
	*cb_args->map_ver = obj_reply_map_version_get(rpc);

	D_RWLOCK_WRLOCK(&cb_args->obj->cob_lock);

	bool check = true;
	bool changed = false;
	bool first = (cb_args->dkey->iov_len == 0);

	if (flags & DAOS_GET_DKEY) {
		uint64_t *val = (uint64_t *)okqo->okqo_dkey.iov_buf;
		uint64_t *cur = (uint64_t *)cb_args->dkey->iov_buf;

		if (okqo->okqo_dkey.iov_len != sizeof(uint64_t)) {
			D_ERROR("Invalid Dkey obtained\n");
			D_RWLOCK_UNLOCK(&cb_args->obj->cob_lock);
			D_GOTO(out, rc = -DER_IO);
		}

		/** for first cb, just set the dkey */
		if (first) {
			*cur = *val;
			cb_args->dkey->iov_len = okqo->okqo_dkey.iov_len;
		} else if (flags & DAOS_GET_MAX) {
			if (*val > *cur) {
				*cur = *val;
				/** set to change akey and recx */
				changed = true;
			} else {
				/** no change, don't check akey and recx */
				check = false;
			}
		} else if (flags & DAOS_GET_MIN) {
			if (*val < *cur) {
				*cur = *val;
				/** set to change akey and recx */
				changed = true;
			} else {
				/** no change, don't check akey and recx */
				check = false;
			}
		} else {
			D_ASSERT(0);
		}
	}

	if (check && flags & DAOS_GET_AKEY) {
		uint64_t *val = (uint64_t *)okqo->okqo_akey.iov_buf;
		uint64_t *cur = (uint64_t *)cb_args->akey->iov_buf;

		/** if first cb, or dkey changed, set akey */
		if (first || changed)
			*cur = *val;
		else
			D_ASSERT(0);
	}

	if (check && flags & DAOS_GET_RECX) {
		/** if first cb, set recx */
		if (first || changed) {
			cb_args->recx->rx_nr = okqo->okqo_recx.rx_nr;
			cb_args->recx->rx_idx = okqo->okqo_recx.rx_idx;
		} else {
			D_ASSERT(0);
		}
	}
	D_RWLOCK_UNLOCK(&cb_args->obj->cob_lock);

out:
	crt_req_decref(rpc);
	if (ret == 0 || obj_retry_error(rc))
		ret = rc;
	return ret;
}

int
dc_obj_shard_query_key(struct dc_obj_shard *shard, struct dc_obj_epoch *epoch,
		       uint32_t flags, struct dc_object *obj, daos_key_t *dkey,
		       daos_key_t *akey, daos_recx_t *recx,
		       const uuid_t coh_uuid, const uuid_t cont_uuid,
		       struct dtx_id *dti, unsigned int *map_ver,
		       tse_task_t *task)
{
	struct dc_pool			*pool = NULL;
	struct obj_query_key_in		*okqi;
	crt_rpc_t			*req;
	struct obj_query_key_cb_args	 cb_args;
	daos_unit_oid_t			 oid;
	crt_endpoint_t			 tgt_ep;
	uint64_t			 dkey_hash;
	int				 rc;

	tse_task_stack_pop_data(task, &dkey_hash, sizeof(dkey_hash));

	pool = obj_shard_ptr2pool(shard);
	if (pool == NULL)
		D_GOTO(out, rc = -DER_NO_HDL);

	oid = shard->do_id;
	tgt_ep.ep_grp	= pool->dp_sys->sy_group;
	tgt_ep.ep_tag	= shard->do_target_idx;
	tgt_ep.ep_rank = shard->do_target_rank;
	if ((int)tgt_ep.ep_rank < 0)
		D_GOTO(out, rc = (int)tgt_ep.ep_rank);

	D_DEBUG(DB_IO, "OBJ_QUERY_KEY_RPC, rank=%d tag=%d.\n",
		tgt_ep.ep_rank, tgt_ep.ep_tag);

	rc = obj_req_create(daos_task2ctx(task), &tgt_ep,
			    DAOS_OBJ_RPC_QUERY_KEY, &req);
	if (rc != 0)
		D_GOTO(out, rc);

	crt_req_addref(req);
	cb_args.rpc	= req;
	cb_args.map_ver = map_ver;
	cb_args.oid	= shard->do_id;
	cb_args.dkey	= dkey;
	cb_args.akey	= akey;
	cb_args.recx	= recx;
	cb_args.obj	= obj;

	rc = tse_task_register_comp_cb(task, obj_shard_query_key_cb, &cb_args,
				       sizeof(cb_args));
	if (rc != 0)
		D_GOTO(out_req, rc);

	okqi = crt_req_get(req);
	D_ASSERT(okqi != NULL);

	okqi->okqi_map_ver		= *map_ver;
	okqi->okqi_epoch		= epoch->oe_value;
	okqi->okqi_api_flags		= flags;
	okqi->okqi_oid			= oid;
	if (dkey != NULL)
		okqi->okqi_dkey		= *dkey;
	if (akey != NULL)
		okqi->okqi_akey		= *akey;
	if (epoch->oe_uncertain)
		okqi->okqi_flags	= ORF_EPOCH_UNCERTAIN;
	uuid_copy(okqi->okqi_pool_uuid, pool->dp_pool);
	uuid_copy(okqi->okqi_co_hdl, coh_uuid);
	uuid_copy(okqi->okqi_co_uuid, cont_uuid);
	daos_dti_copy(&okqi->okqi_dti, dti);

	rc = daos_rpc_send(req, task);
	if (rc != 0) {
		D_ERROR("query_key rpc failed rc "DF_RC"\n", DP_RC(rc));
		D_GOTO(out_req, rc);
	}

	dc_pool_put(pool);
	return 0;

out_req:
	crt_req_decref(req);
out:
	if (pool)
		dc_pool_put(pool);
	tse_task_complete(task, rc);
	return rc;
}

struct obj_shard_sync_cb_args {
	crt_rpc_t	*rpc;
	daos_epoch_t	*epoch;
	uint32_t	*map_ver;
};

static int
obj_shard_sync_cb(tse_task_t *task, void *data)
{
	struct obj_shard_sync_cb_args	*cb_args;
	struct obj_sync_out		*oso;
	int				 ret = task->dt_result;
	int				 rc = 0;
	crt_rpc_t			*rpc;

	cb_args = (struct obj_shard_sync_cb_args *)data;
	rpc = cb_args->rpc;

	if (ret != 0) {
		D_ERROR("OBJ_SYNC RPC failed: rc = %d\n", ret);
		D_GOTO(out, rc = ret);
	}

	oso = crt_reply_get(rpc);
	rc = oso->oso_ret;
	if (rc == -DER_NONEXIST)
		D_GOTO(out, rc = 0);

	if (rc == -DER_INPROGRESS) {
		D_DEBUG(DB_TRACE,
			"rpc %p OBJ_SYNC_RPC may need retry: rc = "DF_RC"\n",
			rpc, DP_RC(rc));
		D_GOTO(out, rc);
	}

	if (rc != 0) {
		D_ERROR("rpc %p OBJ_SYNC_RPC failed: rc = "DF_RC"\n", rpc,
			DP_RC(rc));
		D_GOTO(out, rc);
	}

	*cb_args->epoch = oso->oso_epoch;
	*cb_args->map_ver = oso->oso_map_version;

	D_DEBUG(DB_IO, "OBJ_SYNC_RPC reply: eph "DF_U64", version %u.\n",
		oso->oso_epoch, oso->oso_map_version);

out:
	crt_req_decref(rpc);
	return rc;
}

int
dc_obj_shard_sync(struct dc_obj_shard *shard, enum obj_rpc_opc opc,
		  void *shard_args, struct daos_shard_tgt *fw_shard_tgts,
		  uint32_t fw_cnt, tse_task_t *task)
{
	struct shard_sync_args		*args = shard_args;
	struct dc_pool			*pool = NULL;
	uuid_t				 cont_hdl_uuid;
	uuid_t				 cont_uuid;
	struct obj_sync_in		*osi;
	crt_rpc_t			*req;
	struct obj_shard_sync_cb_args	 cb_args;
	crt_endpoint_t			 tgt_ep;
	int				 rc;

	pool = obj_shard_ptr2pool(shard);
	if (pool == NULL)
		D_GOTO(out, rc = -DER_NO_HDL);

	rc = dc_cont_hdl2uuid(shard->do_co_hdl, &cont_hdl_uuid, &cont_uuid);
	if (rc != 0)
		D_GOTO(out, rc);

	tgt_ep.ep_grp	= pool->dp_sys->sy_group;
	tgt_ep.ep_tag	= shard->do_target_idx;
	tgt_ep.ep_rank	= shard->do_target_rank;
	if ((int)tgt_ep.ep_rank < 0)
		D_GOTO(out, rc = (int)tgt_ep.ep_rank);

	D_DEBUG(DB_IO, "OBJ_SYNC_RPC, rank=%d tag=%d.\n",
		tgt_ep.ep_rank, tgt_ep.ep_tag);

	rc = obj_req_create(daos_task2ctx(task), &tgt_ep, opc, &req);
	if (rc != 0)
		D_GOTO(out, rc);

	crt_req_addref(req);
	cb_args.rpc	= req;
	cb_args.epoch	= args->sa_epoch;
	cb_args.map_ver = &args->sa_auxi.map_ver;

	rc = tse_task_register_comp_cb(task, obj_shard_sync_cb, &cb_args,
				       sizeof(cb_args));
	if (rc != 0)
		D_GOTO(out_req, rc);

	osi = crt_req_get(req);
	D_ASSERT(osi != NULL);

	uuid_copy(osi->osi_co_hdl, cont_hdl_uuid);
	uuid_copy(osi->osi_pool_uuid, pool->dp_pool);
	uuid_copy(osi->osi_co_uuid, cont_uuid);
	osi->osi_oid		= shard->do_id;
	osi->osi_epoch		= args->sa_auxi.epoch.oe_value;
	osi->osi_map_ver	= args->sa_auxi.map_ver;

	rc = daos_rpc_send(req, task);
	if (rc != 0) {
		D_ERROR("OBJ_SYNC_RPC failed: rc = "DF_RC"\n", DP_RC(rc));
		D_GOTO(out_req, rc);
	}

	dc_pool_put(pool);
	return 0;

out_req:
	crt_req_decref(req);

out:
	if (pool != NULL)
		dc_pool_put(pool);

	tse_task_complete(task, rc);
	return rc;
}
