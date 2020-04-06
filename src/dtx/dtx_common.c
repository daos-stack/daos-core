/**
 * (C) Copyright 2019-2020 Intel Corporation.
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
 * dtx: DTX common logic
 */
#define D_LOGFAC	DD_FAC(dtx)

#include <abt.h>
#include <uuid/uuid.h>
#include <daos_srv/pool.h>
#include <daos_srv/container.h>
#include <daos_srv/vos.h>
#include <daos_srv/dtx_srv.h>
#include <daos_srv/daos_server.h>
#include "dtx_internal.h"

struct dtx_batched_commit_args {
	d_list_t		 dbca_link;
	struct ds_cont_child	*dbca_cont;
	void			*dbca_deregistering;
};

void
dtx_aggregate(void *arg)
{
	struct ds_cont_child	*cont = arg;

	while (1) {
		struct dtx_stat		stat = { 0 };
		int			rc;

		rc = vos_dtx_aggregate(cont->sc_hdl);
		if (rc != 0)
			break;

		ABT_thread_yield();

		if (cont->sc_open == 0)
			break;

		vos_dtx_stat(cont->sc_hdl, &stat);

		if (stat.dtx_committed_count <= DTX_AGG_THRESHOLD_CNT_LOWER)
			break;

		if (stat.dtx_committed_count >= DTX_AGG_THRESHOLD_CNT_UPPER)
			continue;

		if (stat.dtx_oldest_committed_time == 0 ||
		    dtx_hlc_age2sec(stat.dtx_oldest_committed_time) <=
		    DTX_AGG_THRESHOLD_AGE_LOWER)
			break;
	}

	cont->sc_dtx_aggregating = 0;
	ds_cont_child_put(cont);
}

static inline void
dtx_free_committable(struct dtx_entry *dtes)
{
	D_FREE(dtes);
}

static inline void
dtx_free_dbca(struct dtx_batched_commit_args *dbca)
{
	d_list_del(&dbca->dbca_link);
	ds_cont_child_put(dbca->dbca_cont);
	D_FREE_PTR(dbca);
}

static void
dtx_flush_on_deregister(struct dss_module_info *dmi,
			struct dtx_batched_commit_args *dbca)
{
	struct ds_cont_child	*cont = dbca->dbca_cont;
	struct ds_pool_child	*pool = cont->sc_pool;
	int			 rc;

	D_ASSERT(dbca->dbca_deregistering != NULL);
	do {
		struct dtx_entry	*dtes = NULL;

		rc = vos_dtx_fetch_committable(cont->sc_hdl,
					       DTX_THRESHOLD_COUNT, NULL,
					       DAOS_EPOCH_MAX, &dtes);
		if (rc <= 0)
			break;

		rc = dtx_commit(pool->spc_uuid, cont->sc_uuid,
				dtes, rc, pool->spc_map_version);
		dtx_free_committable(dtes);
	} while (rc >= 0);

	if (rc < 0)
		D_ERROR(DF_UUID": Fail to flush CoS cache: rc = %d\n",
			DP_UUID(cont->sc_uuid), rc);

	/*
	 * dtx_batched_commit_deregister() set force flush and wait for
	 * flush done, then free the dbca.
	 */
	d_list_del_init(&dbca->dbca_link);
	rc = ABT_future_set(dbca->dbca_deregistering, NULL);
	D_ASSERTF(rc == ABT_SUCCESS, "ABT_future_set failed for DTX "
		  "flush on "DF_UUID": rc = %d\n", DP_UUID(cont->sc_uuid), rc);
}

void
dtx_batched_commit(void *arg)
{
	struct dss_module_info		*dmi = dss_get_module_info();
	struct dtx_batched_commit_args	*dbca;

	while (1) {
		struct ds_cont_child		*cont;
		struct dtx_entry		*dtes = NULL;
		struct dtx_stat			 stat = { 0 };
		int				 rc;

		if (d_list_empty(&dmi->dmi_dtx_batched_list))
			goto check;

		dbca = d_list_entry(dmi->dmi_dtx_batched_list.next,
				    struct dtx_batched_commit_args, dbca_link);
		cont = dbca->dbca_cont;
		if (dbca->dbca_deregistering != NULL) {
			dtx_flush_on_deregister(dmi, dbca);
			goto check;
		}

		d_list_move_tail(&dbca->dbca_link, &dmi->dmi_dtx_batched_list);
		vos_dtx_stat(cont->sc_hdl, &stat);

		if ((stat.dtx_committable_count > DTX_THRESHOLD_COUNT) ||
		    (stat.dtx_oldest_committable_time != 0 &&
		     dtx_hlc_age2sec(stat.dtx_oldest_committable_time) >
		     DTX_COMMIT_THRESHOLD_AGE)) {
			rc = vos_dtx_fetch_committable(cont->sc_hdl,
						DTX_THRESHOLD_COUNT, NULL,
						DAOS_EPOCH_MAX, &dtes);
			if (rc > 0) {
				rc = dtx_commit(cont->sc_pool->spc_uuid,
					cont->sc_uuid, dtes, rc,
					cont->sc_pool->spc_map_version);
				dtx_free_committable(dtes);

				if (dbca->dbca_deregistering) {
					dtx_flush_on_deregister(dmi, dbca);
					goto check;
				}

				if (!cont->sc_dtx_aggregating)
					vos_dtx_stat(cont->sc_hdl, &stat);
			}
		}

		if (!cont->sc_dtx_aggregating &&
		    (stat.dtx_committed_count >= DTX_AGG_THRESHOLD_CNT_UPPER ||
		     (stat.dtx_committed_count > DTX_AGG_THRESHOLD_CNT_LOWER &&
		      stat.dtx_oldest_committed_time != 0 &&
		      dtx_hlc_age2sec(stat.dtx_oldest_committed_time) >=
				DTX_AGG_THRESHOLD_AGE_UPPER))) {
			ds_cont_child_get(cont);
			cont->sc_dtx_aggregating = 1;
			rc = dss_ult_create(dtx_aggregate, cont,
				DSS_ULT_AGGREGATE, DSS_TGT_SELF, 0, NULL);
			if (rc != 0) {
				cont->sc_dtx_aggregating = 0;
				ds_cont_child_put(cont);
			}
		}

check:
		if (dss_xstream_exiting(dmi->dmi_xstream))
			break;
		ABT_thread_yield();
	}

	while (!d_list_empty(&dmi->dmi_dtx_batched_list)) {
		dbca = d_list_entry(dmi->dmi_dtx_batched_list.next,
				    struct dtx_batched_commit_args, dbca_link);
		dtx_free_dbca(dbca);
	}
}

/**
 * Init local dth handle.
 */
static void
dtx_handle_init(struct dtx_id *dti, daos_handle_t coh, daos_epoch_t epoch,
		uint32_t pm_ver, uint32_t dti_cos_cnt, struct dtx_id *dti_cos,
		bool leader, bool solo, uint16_t rdg_cnt, uint32_t rdg_size,
		struct dtx_rdg_unit *rdgs, struct dtx_handle *dth)
{
	struct dtx_entry	*dte = &dth->dth_entry;

	dte->dte_xid = *dti;
	dte->dte_rdg_size = rdg_size;
	dte->dte_rdg_cnt = rdg_cnt;
	dte->dte_rdgs = rdgs;

	dth->dth_coh = coh;
	dth->dth_epoch = epoch;
	dth->dth_ver = pm_ver;

	dth->dth_sync = 0;
	dth->dth_leader = leader ? 1 : 0;
	dth->dth_solo = solo ? 1 : 0;
	dth->dth_compounded = 0;
	dth->dth_dti_cos_done = 0;
	dth->dth_touch_shared = 0;
	dth->dth_renew = 0;
	dth->dth_actived = 0;

	dth->dth_gen = 0;
	dth->dth_ent = NULL;

	dth->dth_dti_cos = dti_cos;
	dth->dth_dti_cos_cnt = dti_cos_cnt;

	memset(&dth->dth_intent, 0, sizeof(*dth) -
	       offsetof(struct dtx_handle, dth_oid) + sizeof(dth->dth_oid));
}

/**
 * Initialize the DTX handle for per modification based part.
 *
 * \param dth		[IN]	Pointer to the DTX handle.
 * \param oid		[IN]	The target object (shard) ID.
 * \param dkey_hash	[IN]	Hash of the dkey to be modified if applicable.
 * \param intent	[IN]	The intent of related modification.
 */
void
dtx_sub_init(struct dtx_handle *dth, daos_unit_oid_t *oid,
	     uint64_t dkey_hash, uint16_t intent)
{
	if (dth != NULL) {
		dth->dth_intent = intent;
		dth->dth_op_seq++;
		dth->dth_dkey_hash = dkey_hash;
		dth->dth_oid = *oid;

		D_DEBUG(DB_TRACE, "Sub init DTX "DF_DTI" for object "DF_UOID
			" dkey %lu, intent %s\n",
			DP_DTI(&dth->dth_entry.dte_xid),
			DP_UOID(*oid), (unsigned long)dkey_hash,
			intent == DAOS_INTENT_PUNCH ? "Punch" : "Update");
	}
}

/**
 * Start the leader DTX handle in DRAM.
 *
 * \param dti		[IN]	The DTX identifier.
 * \param coh		[IN]	Container open handle.
 * \param epoch		[IN]	Epoch for the DTX.
 * \param pm_ver	[IN]	Pool map version for the DTX.
 * \param rdg_size	[IN]	The size for the redundancy groups (RDG) array.
 * \param rdg_cnt	[IN]	The count of RDGs to be modified.
 * \param rdgs		[IN]	The array of RDGs to be modified.
 * \param tgts		[IN]	targets for distribute transaction.
 * \param tgt_cnt	[IN]	number of targets.
 * \param dti_cos_cnt	[IN]	The @dti_cos array size.
 * \param dti_cos	[IN]	DTXs array to be committed because of shared.
 * \param dlh		[OUT]	Pointer to the leader DTX handle.
 *
 * \return			Zero on success, negative value if error.
 */
int
dtx_leader_begin(struct dtx_id *dti, daos_handle_t coh, daos_epoch_t epoch,
		 uint32_t pm_ver, uint32_t rdg_size, uint16_t rdg_cnt,
		 struct dtx_rdg_unit *rdgs, struct daos_shard_tgt *tgts,
		 uint32_t tgt_cnt, uint32_t dti_cos_cnt,
		 struct dtx_id *dti_cos, struct dtx_leader_handle *dlh)
{
	struct dtx_handle	*dth = &dlh->dlh_handle;
	int			 i;

	/* Single replica case. */
	if (tgt_cnt == 0) {
		if (!daos_is_zero_dti(dti))
			goto init;

		daos_dti_gen(&dth->dth_entry.dte_xid, true);
		return 0;
	}

	dlh->dlh_future = ABT_FUTURE_NULL;
	D_ALLOC_ARRAY(dlh->dlh_subs, tgt_cnt);
	if (dlh->dlh_subs == NULL)
		return -DER_NOMEM;

	for (i = 0; i < tgt_cnt; i++)
		dlh->dlh_subs[i].dss_tgt = tgts[i];
	dlh->dlh_sub_cnt = tgt_cnt;

	if (daos_is_zero_dti(dti)) {
		daos_dti_gen(&dth->dth_entry.dte_xid, true); /* zero it */
		return 0;
	}

init:
	dtx_handle_init(dti, coh, epoch, pm_ver, dti_cos_cnt, dti_cos,
			true, tgt_cnt == 0 ? true : false,
			rdg_cnt, rdg_size, rdgs, dth);

	D_DEBUG(DB_TRACE,
		"Start leader DTX "DF_DTI" ver %u, cos_count %d\n",
		DP_DTI(dti), dth->dth_ver, dti_cos_cnt);

	return 0;
}

static int
dtx_leader_wait(struct dtx_leader_handle *dlh)
{
	int	rc;

	rc = ABT_future_wait(dlh->dlh_future);
	D_ASSERTF(rc == ABT_SUCCESS, "ABT_future_wait failed %d.\n", rc);

	ABT_future_free(&dlh->dlh_future);
	D_DEBUG(DB_TRACE, "dth "DF_DTI" rc %d\n",
		DP_DTI(&dlh->dlh_handle.dth_entry.dte_xid), rc);

	return dlh->dlh_result;
};

/**
 * Stop the leader DTX thandle.
 *
 * \param dlh		[IN]	The DTX handle on leader node.
 * \param cont		[IN]	Per-thread container cache.
 * \param result	[IN]	Operation result.
 *
 * \return			Zero on success, negative value if error.
 */
int
dtx_leader_end(struct dtx_leader_handle *dlh, struct ds_cont_child *cont,
	       int result)
{
	struct dtx_handle		*dth = &dlh->dlh_handle;
	daos_epoch_t			 epoch = dth->dth_epoch;
	int				 flags = 0;
	int				 rc = 0;

	if (dlh == NULL)
		return result;

	if (dlh->dlh_sub_cnt == 0) {
		if (result == -DER_AGAIN)
			dth->dth_renew = 1;

		if (daos_is_zero_dti(&dth->dth_entry.dte_xid))
			return result;

		goto out;
	}

	D_ASSERT(cont != NULL);

	rc = dtx_leader_wait(dlh);
	if (daos_is_zero_dti(&dth->dth_entry.dte_xid)) {
		D_FREE(dlh->dlh_subs);

		return result < 0 ? result : rc;
	}

	if (rc == -DER_AGAIN)
		dth->dth_renew = 1;

	if (result < 0 || rc < 0 || !dth->dth_actived)
		D_GOTO(out, result = result < 0 ? result : rc);

	/* XXX: Synchronously commit compounded DTX to avoid recovery issues. */
	if (dth->dth_compounded) {
		dth->dth_sync = 1;
		goto sync;
	}

	D_ASSERT(daos_oid_cmp(dth->dth_oid.id_pub,
			      dth->dth_entry.dte_rdgs[0].dru_oid) == 0);

	if (dth->dth_touch_shared)
		flags |= DCF_SHARED;

	rc = vos_dtx_check_sync(dth->dth_coh, dth->dth_oid, &epoch);
	if (rc == 0)
		rc = vos_dtx_add_cos(dth->dth_coh, &dth->dth_entry,
				     dth->dth_dkey_hash, dth->dth_epoch,
				     dth->dth_gen, flags);
	if (rc == -DER_AGAIN) {
		D_WARN(DF_UUID": Fail to add DTX "DF_DTI" to CoS "
		       "because of using old epoch "DF_U64"\n",
		       DP_UUID(cont->sc_uuid), DP_DTI(&dth->dth_entry.dte_xid),
		       dth->dth_epoch);
		dth->dth_renew = 1;
		D_GOTO(out, result = rc);
	}

	/* The DTX has been aborted by resync ULT, ask the client to retry. */
	if (rc == -DER_NONEXIST) {
		D_WARN(DF_UUID": Fail to add DTX "DF_DTI" with eph "
		       DF_U64" to CoS because it is aborted by resync.\n",
		       DP_UUID(cont->sc_uuid), DP_DTI(&dth->dth_entry.dte_xid),
		       dth->dth_epoch);
		D_GOTO(out, result = -DER_INPROGRESS);
	}

	if (rc != 0 && epoch < dth->dth_epoch) {
		D_WARN(DF_UUID": Fail to add DTX "DF_DTI" to CoS cache: rc = "
		       DF_RC ". Try to commit it sychronously.\n",
		       DP_UUID(cont->sc_uuid), DP_DTI(&dth->dth_entry.dte_xid),
		       DP_RC(rc));
		dth->dth_sync = 1;
	}

sync:
	if (dth->dth_sync) {
		rc = dtx_commit(cont->sc_pool->spc_uuid, cont->sc_uuid,
				&dth->dth_entry, 1,
				cont->sc_pool->spc_map_version);
		if (rc != 0) {
			D_ERROR(DF_UUID": Fail to sync commit DTX "DF_DTI
				": rc = "DF_RC"\n", DP_UUID(cont->sc_uuid),
				DP_DTI(&dth->dth_entry.dte_xid), DP_RC(rc));
			D_GOTO(out, result = rc);
		}
	}

out:
	if (result < 0 && dlh->dlh_sub_cnt > 0)
		dtx_abort(cont->sc_pool->spc_uuid, cont->sc_uuid,
			  dth->dth_epoch, &dth->dth_entry, 1,
			  cont->sc_pool->spc_map_version);

	D_DEBUG(DB_TRACE, "Stop leader DTX "DF_DTI" ver %u, %s, %s "
		"participator(s): rc = "DF_RC"\n",
		DP_DTI(&dth->dth_entry.dte_xid), dth->dth_ver,
		dth->dth_sync ? "sync" : "async",
		dth->dth_solo ? "single" : "multiple", DP_RC(result));

	D_ASSERTF(result <= 0, "unexpected return value %d\n", result);

	D_FREE(dlh->dlh_subs);

	return result;
}

/**
 * Start the non-leader DTX handle in DRAM.
 *
 * \param dti		[IN]	The DTX identifier.
 * \param coh		[IN]	Container open handle.
 * \param epoch		[IN]	Epoch for the DTX.
 * \param pm_ver	[IN]	Pool map version for the DTX.
 * \param rdg_size	[IN]	The size for the redundancy groups (RDG) array.
 * \param rdg_cnt	[IN]	The count of RDGs to be modified.
 * \param rdgs		[IN]	The array of RDGs to be modified.
 * \param dti_cos_cnt	[IN]	The @dti_cos array size.
 * \param dti_cos	[IN]	DTXs array to be committed because of shared.
 * \param dth		[OUT]	Pointer to the DTX handle.
 *
 * \return			Zero on success, negative value if error.
 */
int
dtx_begin(struct dtx_id *dti, daos_handle_t coh, daos_epoch_t epoch,
	  uint32_t pm_ver, uint32_t rdg_size, uint16_t rdg_cnt,
	  struct dtx_rdg_unit *rdgs, uint32_t dti_cos_cnt,
	  struct dtx_id *dti_cos, struct dtx_handle *dth)
{
	if (dth == NULL || daos_is_zero_dti(dti))
		return 0;

	dtx_handle_init(dti, coh, epoch, pm_ver, dti_cos_cnt, dti_cos, false,
			false, rdg_cnt, rdg_size, rdgs, dth);

	D_DEBUG(DB_TRACE,
		"Start non-leader DTX "DF_DTI" ver %u, cos_count %d\n",
		DP_DTI(&dth->dth_entry.dte_xid), dth->dth_ver, dti_cos_cnt);

	return 0;
}

int
dtx_end(struct dtx_handle *dth, struct ds_cont_hdl *cont_hdl,
	struct ds_cont_child *cont, int result)
{
	int rc = 0;

	if (dth == NULL || daos_is_zero_dti(&dth->dth_entry.dte_xid))
		goto out;

	if (result < 0) {
		if (dth->dth_dti_cos_cnt > 0) {
			/* XXX: For non-leader replica, even if we fail to
			 *	make related modification for some reason,
			 *	we still need to commit the DTXs for CoS.
			 *	Because other replica may have already
			 *	committed them. For leader case, it is
			 *	not important even if we fail to commit
			 *	the CoS DTXs, because they are still in
			 *	CoS cache, and can be committed next time.
			 */
			rc = vos_dtx_commit(cont->sc_hdl, dth->dth_dti_cos,
					    dth->dth_dti_cos_cnt);
			if (rc != 0)
				D_ERROR(DF_UUID": Fail to DTX CoS commit: %d\n",
					DP_UUID(cont->sc_uuid), rc);
		}
	}

	D_DEBUG(DB_TRACE, "Stop non-leader DTX "DF_DTI" ver %u, rc = "DF_RC"\n",
		DP_DTI(&dth->dth_entry.dte_xid), dth->dth_ver, DP_RC(result));

	D_ASSERTF(result <= 0, "unexpected return value %d\n", result);

out:
	return result;
}


int
dtx_batched_commit_register(struct ds_cont_child *cont)
{
	struct dtx_batched_commit_args	*dbca;
	d_list_t			*head;

	D_ASSERT(cont != NULL);

	head = &dss_get_module_info()->dmi_dtx_batched_list;
	d_list_for_each_entry(dbca, head, dbca_link) {
		if (dbca->dbca_deregistering != NULL)
			continue;

		if (uuid_compare(dbca->dbca_cont->sc_uuid,
				 cont->sc_uuid) == 0)
			return 0;
	}

	D_ALLOC_PTR(dbca);
	if (dbca == NULL)
		return -DER_NOMEM;

	ds_cont_child_get(cont);
	dbca->dbca_cont = cont;
	d_list_add_tail(&dbca->dbca_link, head);
	return 0;
}

void
dtx_batched_commit_deregister(struct ds_cont_child *cont)
{
	struct dtx_batched_commit_args	*dbca;
	d_list_t			*head;
	ABT_future			 future;
	int				 rc;

	D_ASSERT(cont != NULL);
	D_ASSERT(cont->sc_open == 0);

	head = &dss_get_module_info()->dmi_dtx_batched_list;
	d_list_for_each_entry(dbca, head, dbca_link) {
		if (uuid_compare(dbca->dbca_cont->sc_uuid,
				 cont->sc_uuid) != 0)
			continue;

		/*
		 * Notify the dtx_batched_commit ULT to flush the
		 * committable DTXs.
		 *
		 * Then current ULT will wait here until the DTXs
		 * have been committed by dtx_batched_commit ULT
		 * that will wakeup current ULT.
		 */
		D_ASSERT(dbca->dbca_deregistering == NULL);
		rc = ABT_future_create(1, NULL, &future);
		if (rc != ABT_SUCCESS) {
			D_ERROR("ABT_future_create failed for DTX flush on "
				DF_UUID" %d\n", DP_UUID(cont->sc_uuid), rc);
			return;
		}

		dbca->dbca_deregistering = future;
		rc = ABT_future_wait(future);
		D_ASSERTF(rc == ABT_SUCCESS, "ABT_future_wait failed "
			  "for DTX flush (2) on "DF_UUID": rc = %d\n",
			  DP_UUID(cont->sc_uuid), rc);

		D_ASSERT(d_list_empty(&dbca->dbca_link));
		dtx_free_dbca(dbca);
		ABT_future_free(&future);
		break;
	}
}

int
dtx_handle_resend(daos_handle_t coh, daos_obj_id_t *oid, struct dtx_id *dti,
		  uint64_t dkey_hash, daos_epoch_t *epoch)
{
	int	rc;

	if (daos_is_zero_dti(dti))
		/* If DTX is disabled, then means that the appplication does
		 * not care about the replicas consistency. Under such case,
		 * if client resends some modification RPC, then just handle
		 * it as non-resent case, return -DER_NONEXIST.
		 *
		 * It will cause trouble if related modification has ever
		 * been handled before the resending. But since we cannot
		 * trace (if without DTX) whether it has ever been handled
		 * or not, then just handle it as original without DTX case.
		 */
		return -DER_NONEXIST;

again:
	rc = vos_dtx_check_resend(coh, oid, dti, dkey_hash, epoch);
	switch (rc) {
	case DTX_ST_PREPARED:
		return 0;
	case DTX_ST_COMMITTED:
		return -DER_ALREADY;
	case -DER_NONEXIST:
		if (dtx_hlc_age2sec(dti->dti_hlc) >
		    DTX_AGG_THRESHOLD_AGE_LOWER ||
		    DAOS_FAIL_CHECK(DAOS_DTX_LONG_TIME_RESEND)) {
			D_DEBUG(DB_IO, "Not sure about whether the old RPC "
				DF_DTI" is resent or not.\n", DP_DTI(dti));
			rc = -DER_EP_OLD;
		}
		return rc;
	case -DER_AGAIN:
		/* Re-index committed DTX table is not completed yet,
		 * let's wait and retry.
		 */
		ABT_thread_yield();
		goto again;
	default:
		return rc >= 0 ? -DER_INVAL : rc;
	}
}

static void
dtx_comp_cb(void **arg)
{
	struct dtx_leader_handle	*dlh;
	uint32_t			i;

	dlh = arg[0];
	for (i = 0; i < dlh->dlh_sub_cnt; i++) {
		struct dtx_sub_status	*sub = &dlh->dlh_subs[i];

		if (sub->dss_result == 0)
			continue;

		/* Ignore DER_INPROGRESS if there are other failures */
		if (dlh->dlh_result == 0 || dlh->dlh_result == -DER_INPROGRESS)
			dlh->dlh_result = sub->dss_result;
	}
}

static void
dtx_sub_comp_cb(struct dtx_leader_handle *dlh, int idx, int rc)
{
	struct dtx_sub_status	*sub = &dlh->dlh_subs[idx];
	ABT_future		future = dlh->dlh_future;

	sub->dss_result = rc;
	rc = ABT_future_set(future, dlh);
	D_ASSERTF(rc == ABT_SUCCESS, "ABT_future_set failed %d.\n", rc);

	D_DEBUG(DB_TRACE, "execute from rank %d tag %d, rc %d.\n",
		sub->dss_tgt.st_rank, sub->dss_tgt.st_tgt_idx,
		sub->dss_result);
}

struct dtx_ult_arg {
	dtx_sub_func_t			func;
	void				*func_arg;
	struct dtx_leader_handle	*dlh;
};

static void
dtx_leader_exec_ops_ult(void *arg)
{
	struct dtx_ult_arg	  *ult_arg = arg;
	struct dtx_leader_handle  *dlh = ult_arg->dlh;
	ABT_future		  future = dlh->dlh_future;
	uint32_t		  i;
	int			  rc = 0;

	D_ASSERT(future != ABT_FUTURE_NULL);
	for (i = 0; i < dlh->dlh_sub_cnt; i++) {
		struct dtx_sub_status *sub = &dlh->dlh_subs[i];

		sub->dss_result = 0;

		if (sub->dss_tgt.st_rank == TGTS_IGNORE) {
			int ret;

			ret = ABT_future_set(future, dlh);
			D_ASSERTF(ret == ABT_SUCCESS,
				  "ABT_future_set failed %d.\n", ret);
			continue;
		}

		rc = ult_arg->func(dlh, ult_arg->func_arg, i,
				   dtx_sub_comp_cb);
		if (rc) {
			dlh->dlh_subs[i].dss_result = rc;
			break;
		}
	}

	if (rc != 0) {
		for (i++; i < dlh->dlh_sub_cnt; i++) {
			int ret;

			ret = ABT_future_set(future, dlh);
			D_ASSERTF(ret == ABT_SUCCESS,
				  "ABT_future_set failed %d.\n", ret);
		}
	}
}

/**
 * Execute the operations on all targets.
 */
int
dtx_leader_exec_ops(struct dtx_leader_handle *dlh, dtx_sub_func_t func,
		    void *func_arg)
{
	struct dtx_ult_arg	*ult_arg;
	int			rc;

	if (dlh->dlh_sub_cnt == 0)
		goto exec;

	D_ALLOC_PTR(ult_arg);
	if (ult_arg == NULL)
		return -DER_NOMEM;
	ult_arg->func	= func;
	ult_arg->func_arg = func_arg;
	ult_arg->dlh	= dlh;

	/* the future should already be freed */
	D_ASSERT(dlh->dlh_future == ABT_FUTURE_NULL);
	rc = ABT_future_create(dlh->dlh_sub_cnt, dtx_comp_cb, &dlh->dlh_future);
	if (rc != ABT_SUCCESS) {
		D_ERROR("ABT_future_create failed %d.\n", rc);
		D_FREE_PTR(ult_arg);
		return dss_abterr2der(rc);
	}

	/*
	 * XXX ideally, we probably should create ULT for each shard, but
	 * for performance reasons, let's only create one for all remote
	 * targets for now.
	 */
	dlh->dlh_result = 0;
	rc = dss_ult_create(dtx_leader_exec_ops_ult, ult_arg, DSS_ULT_IOFW,
			    dss_get_module_info()->dmi_tgt_id, 0, NULL);
	if (rc != 0) {
		D_ERROR("ult create failed "DF_RC"\n", DP_RC(rc));
		D_FREE_PTR(ult_arg);
		ABT_future_free(&dlh->dlh_future);
		D_GOTO(out, rc);
	}

exec:
	/* Then execute the local operation */
	rc = func(dlh, func_arg, -1, NULL);
out:
	return rc;
}

int
dtx_obj_sync(uuid_t po_uuid, uuid_t co_uuid, daos_handle_t coh,
	     daos_unit_oid_t oid, daos_epoch_t epoch, uint32_t map_ver)
{
	int	rc = 0;

	while (1) {
		struct dtx_entry	*dtes = NULL;

		rc = vos_dtx_fetch_committable(coh, DTX_THRESHOLD_COUNT,
					       &oid.id_pub, epoch, &dtes);
		if (rc < 0) {
			D_ERROR(DF_UOID" fail to fetch dtx: rc = %d\n",
				DP_UOID(oid), rc);
			break;
		}

		if (rc == 0)
			break;

		rc = dtx_commit(po_uuid, co_uuid, dtes, rc, map_ver);
		dtx_free_committable(dtes);
		if (rc < 0) {
			D_ERROR(DF_UOID" fail to commit dtx: rc = %d\n",
				DP_UOID(oid), rc);
			break;
		}
	}

	if (rc == 0)
		rc = vos_dtx_mark_sync(coh, oid, epoch);

	return rc;
}
