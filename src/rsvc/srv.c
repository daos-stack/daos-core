/*
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
 * \file
 *
 * ds_rsvc: Replicated Service Server
 */

#define D_LOGFAC DD_FAC(rsvc)

#include <sys/stat.h>
#include <daos_srv/daos_server.h>
#include <daos_srv/rsvc.h>
#include "rpc.h"

static struct ds_rsvc_class *rsvc_classes[DS_RSVC_CLASS_COUNT];

void
ds_rsvc_class_register(enum ds_rsvc_class_id id, struct ds_rsvc_class *class)
{
	D_ASSERT(class != NULL);
	D_ASSERT(rsvc_classes[id] == NULL);
	rsvc_classes[id] = class;
}

void
ds_rsvc_class_unregister(enum ds_rsvc_class_id id)
{
	D_ASSERT(rsvc_classes[id] != NULL);
	rsvc_classes[id] = NULL;
}

static struct ds_rsvc_class *
rsvc_class(enum ds_rsvc_class_id id)
{
	D_ASSERTF(id >= 0 && id < DS_RSVC_CLASS_COUNT, "%d\n", id);
	D_ASSERT(rsvc_classes[id] != NULL);
	return rsvc_classes[id];
}

static char *
state_str(enum ds_rsvc_state state)
{
	switch (state) {
	case DS_RSVC_UP_EMPTY:
		return "UP_EMPTY";
	case DS_RSVC_UP:
		return "UP";
	case DS_RSVC_DRAINING:
		return "DRAINING";
	case DS_RSVC_DOWN:
		return "DOWN";
	default:
		return "UNKNOWN";
	}
}

/* Allocate and initialize a ds_rsvc object. */
static int
alloc_init(enum ds_rsvc_class_id class, d_iov_t *id, uuid_t db_uuid,
	   struct ds_rsvc **svcp)
{
	struct ds_rsvc *svc;
	int		rc;

	rc = rsvc_class(class)->sc_alloc(id, &svc);
	if (rc != 0)
		goto err;

	D_INIT_LIST_HEAD(&svc->s_entry);
	svc->s_class = class;
	D_ASSERT(svc->s_id.iov_buf != NULL);
	D_ASSERT(svc->s_id.iov_len > 0);
	D_ASSERT(svc->s_id.iov_buf_len >= svc->s_id.iov_len);
	uuid_copy(svc->s_db_uuid, db_uuid);
	svc->s_state = DS_RSVC_DOWN;

	rc = rsvc_class(class)->sc_name(&svc->s_id, &svc->s_name);
	if (rc != 0)
		goto err_svc;

	rc = rsvc_class(class)->sc_locate(&svc->s_id, &svc->s_db_path);
	if (rc != 0)
		goto err_name;

	rc = ABT_mutex_create(&svc->s_mutex);
	if (rc != ABT_SUCCESS) {
		D_ERROR("%s: failed to create mutex: %d\n", svc->s_name, rc);
		rc = dss_abterr2der(rc);
		goto err_db_path;
	}

	rc = ABT_cond_create(&svc->s_state_cv);
	if (rc != ABT_SUCCESS) {
		D_ERROR("%s: failed to create state_cv: %d\n", svc->s_name, rc);
		rc = dss_abterr2der(rc);
		goto err_mutex;
	}

	rc = ABT_cond_create(&svc->s_leader_ref_cv);
	if (rc != ABT_SUCCESS) {
		D_ERROR("%s: failed to create leader_ref_cv: %d\n", svc->s_name,
			rc);
		rc = dss_abterr2der(rc);
		goto err_state_cv;
	}

	*svcp = svc;
	return 0;

err_state_cv:
	ABT_cond_free(&svc->s_state_cv);
err_mutex:
	ABT_mutex_free(&svc->s_mutex);
err_db_path:
	D_FREE(svc->s_db_path);
err_name:
	D_FREE(svc->s_name);
err_svc:
	rsvc_class(class)->sc_free(svc);
err:
	return rc;
}

static void
fini_free(struct ds_rsvc *svc)
{
	D_ASSERT(d_list_empty(&svc->s_entry));
	D_ASSERTF(svc->s_ref == 0, "%d\n", svc->s_ref);
	D_ASSERTF(svc->s_leader_ref == 0, "%d\n", svc->s_leader_ref);
	ABT_cond_free(&svc->s_leader_ref_cv);
	ABT_cond_free(&svc->s_state_cv);
	ABT_mutex_free(&svc->s_mutex);
	D_FREE(svc->s_db_path);
	D_FREE(svc->s_name);
	rsvc_class(svc->s_class)->sc_free(svc);
}

static void
ds_rsvc_get(struct ds_rsvc *svc)
{
	svc->s_ref++;
}

/** Put a replicated service reference. */
void
ds_rsvc_put(struct ds_rsvc *svc)
{
	D_ASSERTF(svc->s_ref > 0, "%d\n", svc->s_ref);
	svc->s_ref--;
	if (svc->s_ref == 0) {
		rdb_stop(svc->s_db);
		fini_free(svc);
	}
}

static struct d_hash_table rsvc_hash;

static struct ds_rsvc *
rsvc_obj(d_list_t *rlink)
{
	return container_of(rlink, struct ds_rsvc, s_entry);
}

static bool
rsvc_key_cmp(struct d_hash_table *htable, d_list_t *rlink, const void *key,
	     unsigned int ksize)
{
	struct ds_rsvc *svc = rsvc_obj(rlink);

	if (svc->s_id.iov_len != ksize)
		return false;
	return memcmp(svc->s_id.iov_buf, key, svc->s_id.iov_len) == 0;
}

static void
rsvc_rec_addref(struct d_hash_table *htable, d_list_t *rlink)
{
	struct ds_rsvc *svc = rsvc_obj(rlink);

	svc->s_ref++;
}

static bool
rsvc_rec_decref(struct d_hash_table *htable, d_list_t *rlink)
{
	struct ds_rsvc *svc = rsvc_obj(rlink);

	D_ASSERTF(svc->s_ref > 0, "%d\n", svc->s_ref);
	svc->s_ref--;
	return svc->s_ref == 0;
}

static void
rsvc_rec_free(struct d_hash_table *htable, d_list_t *rlink)
{
	struct ds_rsvc *svc = rsvc_obj(rlink);

	rdb_stop(svc->s_db);
	fini_free(svc);
}

static d_hash_table_ops_t rsvc_hash_ops = {
	.hop_key_cmp	= rsvc_key_cmp,
	.hop_rec_addref	= rsvc_rec_addref,
	.hop_rec_decref	= rsvc_rec_decref,
	.hop_rec_free	= rsvc_rec_free
};

static int
rsvc_hash_init(void)
{
	return d_hash_table_create_inplace(D_HASH_FT_NOLOCK, 4 /* bits */,
					   NULL /* priv */, &rsvc_hash_ops,
					   &rsvc_hash);
}

static void
rsvc_hash_fini(void)
{
	d_hash_table_destroy_inplace(&rsvc_hash, true /* force */);
}

/**
 * Look up a replicated service.
 *
 * \param[in]	class	replicated service class
 * \param[in]	id	replicated service ID
 * \param[out]	svc	replicated service
 */
int
ds_rsvc_lookup(enum ds_rsvc_class_id class, d_iov_t *id,
	       struct ds_rsvc **svc)
{
	d_list_t       *entry;
	bool		nonexist = false;

	entry = d_hash_rec_find(&rsvc_hash, id->iov_buf, id->iov_len);
	if (entry == NULL) {
		char	       *path;
		struct stat	buf;
		int		rc;

		/*
		 * See if the DB exists. If an error prevents us from find that
		 * out, return -DER_NOTLEADER so that the client tries other
		 * replicas.
		 */
		rc = rsvc_class(class)->sc_locate(id, &path);
		if (rc != 0)
			goto out;
		rc = stat(path, &buf);
		if (rc != 0) {
			if (errno == ENOENT) {
				nonexist = true;
			} else {
				char *name = NULL;

				rsvc_class(class)->sc_name(id, &name);
				D_ERROR("%s: failed to stat %s: %d\n", name,
					path, errno);
				if (name != NULL)
					D_FREE(name);
			}
		}
		D_FREE(path);
	}
out:
	if (nonexist)
		return -DER_NOTREPLICA;
	if (entry == NULL)
		return -DER_NOTLEADER;
	*svc = rsvc_obj(entry);
	return 0;
}

/*
 * Is svc up (i.e., ready to accept RPCs)? If not, the caller may always report
 * -DER_NOTLEADER, even if svc->s_db is in leader state, in which case the
 * client will retry the RPC.
 */
static bool
up(struct ds_rsvc *svc)
{
	return !svc->s_stop && svc->s_state == DS_RSVC_UP;
}

/**
 * Retrieve the latest leader hint from \a db and fill it into \a hint.
 *
 * \param[in]	svc	replicated service
 * \param[out]	hint	rsvc hint
 */
void
ds_rsvc_set_hint(struct ds_rsvc *svc, struct rsvc_hint *hint)
{
	int rc;

	rc = rdb_get_leader(svc->s_db, &hint->sh_term, &hint->sh_rank);
	if (rc != 0)
		return;
	hint->sh_flags |= RSVC_HINT_VALID;
}

static void
get_leader(struct ds_rsvc *svc)
{
	svc->s_leader_ref++;
}

static void
put_leader(struct ds_rsvc *svc)
{
	D_ASSERTF(svc->s_leader_ref > 0, "%d\n", svc->s_leader_ref);
	svc->s_leader_ref--;
	if (svc->s_leader_ref == 0)
		ABT_cond_broadcast(svc->s_leader_ref_cv);
}

/**
 * As a convenience for general replicated service RPC handlers, this function
 * looks up the replicated service by id, checks that it is up, and takes a
 * reference to the leader fields. svcp is filled only if zero is returned. If
 * the replicated service is not up, hint is filled.
 */
int
ds_rsvc_lookup_leader(enum ds_rsvc_class_id class, d_iov_t *id,
		      struct ds_rsvc **svcp, struct rsvc_hint *hint)
{
	struct ds_rsvc *svc;
	int		rc;

	rc = ds_rsvc_lookup(class, id, &svc);
	if (rc != 0)
		return rc;
	if (!up(svc)) {
		if (hint != NULL)
			ds_rsvc_set_hint(svc, hint);
		ds_rsvc_put(svc);
		return -DER_NOTLEADER;
	}
	get_leader(svc);
	*svcp = svc;
	return 0;
}

/**
 * As a convenience for general replicated service RPC handlers, this function
 * puts svc returned by ds_rsvc_lookup_leader.
 */
void
ds_rsvc_put_leader(struct ds_rsvc *svc)
{
	put_leader(svc);
	ds_rsvc_put(svc);
}

static void
change_state(struct ds_rsvc *svc, enum ds_rsvc_state state)
{
	D_DEBUG(DB_MD, "%s: term "DF_U64" state %s to %s\n", svc->s_name,
		svc->s_term, state_str(svc->s_state), state_str(state));
	svc->s_state = state;
	ABT_cond_broadcast(svc->s_state_cv);
}

static int
rsvc_step_up_cb(struct rdb *db, uint64_t term, void *arg)
{
	struct ds_rsvc *svc = arg;
	int		rc;

	ABT_mutex_lock(svc->s_mutex);
	if (svc->s_stop) {
		D_DEBUG(DB_MD, "%s: skip term "DF_U64" due to stopping\n",
			svc->s_name, term);
		rc = 0;
		goto out_mutex;
	}
	D_ASSERTF(svc->s_state == DS_RSVC_DOWN, "%d\n", svc->s_state);
	svc->s_term = term;
	D_DEBUG(DB_MD, "%s: stepping up to "DF_U64"\n", svc->s_name,
		svc->s_term);

	rc = rsvc_class(svc->s_class)->sc_step_up(svc);
	if (rc == DER_UNINIT) {
		change_state(svc, DS_RSVC_UP_EMPTY);
		rc = 0;
		goto out_mutex;
	} else if (rc != 0) {
		D_ERROR("%s: failed to step up as leader "DF_U64": %d\n",
			svc->s_name, term, rc);
		goto out_mutex;
	}

	change_state(svc, DS_RSVC_UP);
out_mutex:
	ABT_mutex_unlock(svc->s_mutex);
	return rc;
}

/* Bootstrap a self-only, single-replica DB created in start. */
static int
bootstrap_self(struct ds_rsvc *svc, void *arg)
{
	int rc;

	D_DEBUG(DB_MD, "%s: bootstrapping\n", svc->s_name);
	ABT_mutex_lock(svc->s_mutex);

	/*
	 * This single-replica DB shall change from DS_RSVC_DOWN to
	 * DS_RSVC_UP_EMPTY state promptly.
	 */
	while (svc->s_state == DS_RSVC_DOWN)
		ABT_cond_wait(svc->s_state_cv, svc->s_mutex);
	D_ASSERTF(svc->s_state == DS_RSVC_UP_EMPTY, "%d\n", svc->s_state);

	rc = rsvc_class(svc->s_class)->sc_bootstrap(svc, arg);
	if (rc != 0)
		goto out_mutex;

	/* Try stepping up again. */
	rc = rsvc_class(svc->s_class)->sc_step_up(svc);
	if (rc != 0) {
		D_ASSERT(rc != DER_UNINIT);
		goto out_mutex;
	}

	change_state(svc, DS_RSVC_UP);
out_mutex:
	ABT_mutex_unlock(svc->s_mutex);
	D_DEBUG(DB_MD, "%s: bootstrapped: %d\n", svc->s_name, rc);
	return rc;
}

static void
rsvc_step_down_cb(struct rdb *db, uint64_t term, void *arg)
{
	struct ds_rsvc *svc = arg;

	D_DEBUG(DB_MD, "%s: stepping down from "DF_U64"\n", svc->s_name, term);
	ABT_mutex_lock(svc->s_mutex);
	D_ASSERTF(svc->s_term == term, DF_U64" == "DF_U64"\n", svc->s_term,
		  term);
	D_ASSERT(svc->s_state == DS_RSVC_UP_EMPTY ||
		 svc->s_state == DS_RSVC_UP);

	if (svc->s_state == DS_RSVC_UP) {
		/* Stop accepting new leader references. */
		change_state(svc, DS_RSVC_DRAINING);

		rsvc_class(svc->s_class)->sc_drain(svc);

		/* TODO: Abort all in-flight RPCs we sent. */

		/* Wait for all leader references to be released. */
		for (;;) {
			if (svc->s_leader_ref == 0)
				break;
			D_DEBUG(DB_MD, "%s: waiting for %d leader refs\n",
				svc->s_name, svc->s_leader_ref);
			ABT_cond_wait(svc->s_leader_ref_cv, svc->s_mutex);
		}

		rsvc_class(svc->s_class)->sc_step_down(svc);
	}

	change_state(svc, DS_RSVC_DOWN);
	ABT_mutex_unlock(svc->s_mutex);
	D_DEBUG(DB_MD, "%s: stepped down from "DF_U64"\n", svc->s_name, term);
}

static int stop(struct ds_rsvc *svc, bool destroy);

static void
rsvc_stopper(void *arg)
{
	struct ds_rsvc *svc = arg;

	d_hash_rec_delete_at(&rsvc_hash, &svc->s_entry);
	stop(svc, false /* destroy */);
}

static void
rsvc_stop_cb(struct rdb *db, int err, void *arg)
{
	struct ds_rsvc *svc = arg;
	int		rc;

	ds_rsvc_get(svc);
	rc = dss_ult_create(rsvc_stopper, svc, DSS_ULT_MISC, DSS_TGT_SELF,
			    0, NULL);
	if (rc != 0) {
		D_ERROR("%s: failed to create service stopper: %d\n",
			svc->s_name, rc);
		ds_rsvc_put(svc);
	}
}

static struct rdb_cbs rsvc_rdb_cbs = {
	.dc_step_up	= rsvc_step_up_cb,
	.dc_step_down	= rsvc_step_down_cb,
	.dc_stop	= rsvc_stop_cb
};

static bool
self_only(d_rank_list_t *replicas)
{
	d_rank_t	self;
	int		rc;

	rc = crt_group_rank(NULL /* grp */, &self);
	D_ASSERTF(rc == 0, "%d\n", rc);
	return replicas != NULL && replicas->rl_nr == 1 &&
	       replicas->rl_ranks[0] == self;
}

static int
start(enum ds_rsvc_class_id class, d_iov_t *id, uuid_t db_uuid, bool create,
      size_t size, d_rank_list_t *replicas, void *arg, struct ds_rsvc **svcp)
{
	struct ds_rsvc *svc;
	int		rc;

	rc = alloc_init(class, id, db_uuid, &svc);
	if (rc != 0)
		goto err;

	if (create) {
		rc = rdb_create(svc->s_db_path, svc->s_db_uuid, size, replicas);
		if (rc != 0)
			goto err_svc;
	}

	rc = rdb_start(svc->s_db_path, svc->s_db_uuid, &rsvc_rdb_cbs, svc,
		       &svc->s_db);
	if (rc != 0)
		goto err_creation;

	if (create && self_only(replicas) &&
	    rsvc_class(class)->sc_bootstrap != NULL) {
		rc = bootstrap_self(svc, arg);
		if (rc != 0)
			goto err_db;
	}

	svc->s_ref = 1;
	*svcp = svc;
	return 0;

err_db:
	rdb_stop(svc->s_db);
err_creation:
	if (create)
		rdb_destroy(svc->s_db_path, svc->s_db_uuid);
err_svc:
	fini_free(svc);
err:
	return rc;
}

/**
 * Start a replicated service. If \a create is false, all remaining input
 * parameters are ignored; otherwise, create the replica first. If \a replicas
 * is NULL, all remaining input parameters are ignored; otherwise, bootstrap
 * the replicated service.
 *
 * \param[in]	class		replicated service class
 * \param[in]	id		replicated service ID
 * \param[in]	db_uuid		DB UUID
 * \param[in]	create		whether to create the replica before starting
 * \param[in]	size		replica size in bytes
 * \param[in]	replicas	optional initial membership
 * \param[in]	arg		argument for cbs.sc_bootstrap
 *
 * \retval -DER_ALREADY		replicated service already started
 * \retval -DER_CANCELED	replicated service stopping
 */
int
ds_rsvc_start(enum ds_rsvc_class_id class, d_iov_t *id, uuid_t db_uuid,
	      bool create, size_t size, d_rank_list_t *replicas, void *arg)
{
	uuid_t			 db_uuid_buf;
	struct ds_rsvc		*svc;
	struct ds_rsvc_class	*impl;
	d_list_t		*entry;
	int			 rc;

	impl = rsvc_class(class);
	if (!create) {
		rc = impl->sc_load_uuid(id, db_uuid_buf);
		if (rc != 0)
			goto out;
		db_uuid = db_uuid_buf;
	}

	entry = d_hash_rec_find(&rsvc_hash, id->iov_buf, id->iov_len);
	if (entry != NULL) {
		svc = rsvc_obj(entry);
		D_DEBUG(DB_MD, "%s: found: stop=%d\n", svc->s_name,
			svc->s_stop);
		if (svc->s_stop)
			rc = -DER_CANCELED;
		else
			rc = -DER_ALREADY;
		ds_rsvc_put(svc);
		goto out;
	}

	rc = start(class, id, db_uuid, create, size, replicas, arg, &svc);
	if (rc != 0)
		goto out;

	rc = d_hash_rec_insert(&rsvc_hash, svc->s_id.iov_buf, svc->s_id.iov_len,
			       &svc->s_entry, true /* exclusive */);
	if (rc != 0) {
		D_DEBUG(DB_MD, "%s: insert: %d\n", svc->s_name, rc);
		stop(svc, create /* destroy */);
		goto out;
	}

	if (create) {
		rc = impl->sc_store_uuid(id, db_uuid);
		if (rc != 0) {
			stop(svc, create /* destroy */);
			goto out;
		}
	}
	D_DEBUG(DB_MD, "%s: started replicated service\n", svc->s_name);
	ds_rsvc_put(svc);
out:
	if (rc != 0 && rc != -DER_ALREADY && !(create && rc == -DER_EXIST))
		D_ERROR("Failed to start replicated service: %d\n", rc);
	return rc;
}

static int
stop(struct ds_rsvc *svc, bool destroy)
{
	int rc = 0;

	ABT_mutex_lock(svc->s_mutex);

	if (svc->s_stop) {
		ABT_mutex_unlock(svc->s_mutex);
		D_DEBUG(DB_MD, "%s: stopping already\n", svc->s_name);
		return -DER_CANCELED;
	}
	svc->s_stop = true;
	D_DEBUG(DB_MD, "%s: stopping\n", svc->s_name);

	if (svc->s_state == DS_RSVC_UP || svc->s_state == DS_RSVC_UP_EMPTY)
		/*
		 * The service has stepped up. If it is still the leader of
		 * svc->s_term, the following rdb_resign() call will trigger
		 * the matching rsvc_step_down_cb() callback in svc->s_term;
		 * otherwise, the callback must already be pending. Either way,
		 * the service shall eventually enter the DS_RSVC_DOWN state.
		 */
		rdb_resign(svc->s_db, svc->s_term);
	while (svc->s_state != DS_RSVC_DOWN)
		ABT_cond_wait(svc->s_state_cv, svc->s_mutex);

	if (destroy)
		rc = remove(svc->s_db_path);

	ABT_mutex_unlock(svc->s_mutex);
	ds_rsvc_put(svc);
	return rc;
}

/**
 * Stop a replicated service. If destroy is false, all remaining parameters are
 * ignored; otherwise, destroy the service afterward.
 *
 * \param[in]	class		replicated service class
 * \param[in]	id		replicated service ID
 * \param[in]	destroy		whether to destroy the replica after stopping
 *
 * \retval -DER_ALREADY		replicated service already stopped
 * \retval -DER_CANCELED	replicated service stopping
 */
int
ds_rsvc_stop(enum ds_rsvc_class_id class, d_iov_t *id, bool destroy)
{
	struct ds_rsvc		*svc;
	int			 rc;

	rc = ds_rsvc_lookup(class, id, &svc);
	if (rc != 0)
		return -DER_ALREADY;
	d_hash_rec_delete_at(&rsvc_hash, &svc->s_entry);
	rc = stop(svc, destroy);
	if (!rc && destroy)
		rc = rsvc_class(class)->sc_delete_uuid(id);
	return rc;
}

struct stop_ult {
	d_list_t	su_entry;
	ABT_thread	su_thread;
};

struct stop_all_arg {
	d_list_t		saa_list;	/* of stop_ult objects */
	enum ds_rsvc_class_id	saa_class;
};

static int
stop_all_cb(d_list_t *entry, void *varg)
{
	struct ds_rsvc	       *svc = rsvc_obj(entry);
	struct stop_all_arg    *arg = varg;
	struct stop_ult	       *ult;
	int			rc;

	if (svc->s_class != arg->saa_class)
		return 0;

	D_ALLOC_PTR(ult);
	if (ult == NULL)
		return -DER_NOMEM;

	d_hash_rec_addref(&rsvc_hash, &svc->s_entry);
	rc = dss_ult_create(rsvc_stopper, svc, DSS_ULT_POOL_SRV, 0, 0,
			    &ult->su_thread);
	if (rc != 0) {
		d_hash_rec_decref(&rsvc_hash, &svc->s_entry);
		D_FREE(ult);
		return rc;
	}

	d_list_add(&ult->su_entry, &arg->saa_list);
	return 0;
}

/**
 * Stop all replicated services of \a class.
 *
 * \param[in]	class	replicated service class
 */
int
ds_rsvc_stop_all(enum ds_rsvc_class_id class)
{
	struct stop_all_arg	arg;
	struct stop_ult	       *ult;
	struct stop_ult	       *ult_tmp;
	int			rc;

	D_INIT_LIST_HEAD(&arg.saa_list);
	arg.saa_class = class;
	rc = d_hash_table_traverse(&rsvc_hash, stop_all_cb, &arg);

	/* Wait for the stopper ULTs to return. */
	d_list_for_each_entry_safe(ult, ult_tmp, &arg.saa_list, su_entry) {
		d_list_del_init(&ult->su_entry);
		ABT_thread_join(ult->su_thread);
		ABT_thread_free(&ult->su_thread);
		D_FREE(ult);
	}

	if (rc != 0)
		D_ERROR("failed to stop all replicated services: %d\n", rc);
	return rc;
}

/**
 * Stop a replicated service if it is in leader state. Currently, this is used
 * only for testing.
 *
 * \param[in]	class	replicated service class
 * \param[in]	id	replicated service ID
 * \param[out]	hint	rsvc hint
 */
int
ds_rsvc_stop_leader(enum ds_rsvc_class_id class, d_iov_t *id,
		    struct rsvc_hint *hint)
{
	struct ds_rsvc *svc;
	int		rc;

	rc = ds_rsvc_lookup_leader(class, id, &svc, hint);
	if (rc != 0)
		return rc;
	/* Drop our leader reference to allow the service to step down. */
	put_leader(svc);

	d_hash_rec_delete_at(&rsvc_hash, &svc->s_entry);

	return stop(svc, false /* destroy */);
}

int
ds_rsvc_add_replicas_s(struct ds_rsvc *svc, d_rank_list_t *ranks, size_t size)
{
	int	rc;

	rc = ds_rsvc_dist_start(svc->s_class, &svc->s_id, svc->s_db_uuid, ranks,
				true /* create */, false /* bootstrap */, size);

	/* TODO: Attempt to only add replicas that were successfully started */
	if (rc != 0)
		goto out_stop;
	rc = rdb_add_replicas(svc->s_db, ranks);
out_stop:
	/* Clean up ranks that were not added */
	if (ranks->rl_nr > 0) {
		D_ASSERT(rc != 0);
		ds_rsvc_dist_stop(svc->s_class, &svc->s_id, ranks,
				  true /* destroy */);
	}
	return rc;
}

int
ds_rsvc_add_replicas(enum ds_rsvc_class_id class, d_iov_t *id,
		     d_rank_list_t *ranks, size_t size, struct rsvc_hint *hint)
{
	struct ds_rsvc	*svc;
	int		 rc;

	rc = ds_rsvc_lookup_leader(class, id, &svc, hint);
	if (rc != 0)
		return rc;
	rc = ds_rsvc_add_replicas_s(svc, ranks, size);
	ds_rsvc_set_hint(svc, hint);
	put_leader(svc);
	return rc;
}

int
ds_rsvc_remove_replicas_s(struct ds_rsvc *svc, d_rank_list_t *ranks)
{
	d_rank_list_t	*stop_ranks;
	int		 rc;

	rc = daos_rank_list_dup(&stop_ranks, ranks);
	if (rc != 0)
		return rc;
	rc = rdb_remove_replicas(svc->s_db, ranks);

	/* filter out failed ranks */
	daos_rank_list_filter(ranks, stop_ranks, true /* exclude */);
	if (stop_ranks->rl_nr > 0)
		ds_rsvc_dist_stop(svc->s_class, &svc->s_id, stop_ranks,
				  true /* destroy */);
	d_rank_list_free(stop_ranks);
	return rc;
}

int
ds_rsvc_remove_replicas(enum ds_rsvc_class_id class, d_iov_t *id,
			d_rank_list_t *ranks, struct rsvc_hint *hint)
{
	struct ds_rsvc	*svc;
	int		 rc;

	rc = ds_rsvc_lookup_leader(class, id, &svc, hint);
	if (rc != 0)
		return rc;
	rc = ds_rsvc_remove_replicas_s(svc, ranks);
	ds_rsvc_set_hint(svc, hint);
	put_leader(svc);
	return rc;
}

/*************************** Distributed Operations ***************************/

enum rdb_start_flag {
	RDB_AF_CREATE		= 0x1,
	RDB_AF_BOOTSTRAP	= 0x2
};

enum rdb_stop_flag {
	RDB_OF_DESTROY		= 0x1
};

static inline void
rsvc_corpc_incr_cb(d_rank_list_t *ranks, d_rank_t target, int rpc_rc,
		   int *nsuccess)
{
	if (nsuccess != NULL && rpc_rc == 0 &&
	    (ranks == NULL || d_rank_in_rank_list(ranks, target)))
		++(*nsuccess);
}

static int
bcast_create(crt_opcode_t opc, crt_group_t *group, crt_rpc_t **rpc, void *priv)
{
	struct dss_module_info *info = dss_get_module_info();
	crt_opcode_t		opc_full;

	opc_full = DAOS_RPC_OPCODE(opc, DAOS_RSVC_MODULE, DAOS_RSVC_VERSION);
	return crt_corpc_req_create(info->dmi_ctx, group,
				    NULL /* excluded_ranks */, opc_full,
				    NULL /* co_bulk_hdl */, priv, 0 /* flags */,
				    crt_tree_topo(CRT_TREE_FLAT, 0), rpc);
}

/**
 * Perform a distributed create, if \a create is true, and start operation on
 * all replicas of a database with \a dbid spanning \a ranks. This method can
 * be called on any rank. If \a create is false, \a ranks may be NULL.
 *
 * \param[in]	class		replicated service class
 * \param[in]	id		replicated service ID
 * \param[in]	dbid		database UUID
 * \param[in]	ranks		list of replica ranks
 * \param[in]	create		create replicas first
 * \param[in]	bootstrap	start with an initial list of replicas
 * \param[in]	size		size of each replica in bytes if \a create
 */
int
ds_rsvc_dist_start(enum ds_rsvc_class_id class, d_iov_t *id,
		   const uuid_t dbid, const d_rank_list_t *ranks, bool create,
		   bool bootstrap, size_t size)
{
	crt_rpc_t		*rpc;
	struct rsvc_start_in	*in;
	struct rsvc_start_out	*out;
	int			 nstarted = 0;
	int			 rc;

	D_ASSERT(!create || ranks != NULL);
	D_DEBUG(DB_MD, DF_UUID": %s DB\n",
		DP_UUID(dbid), create ? "creating" : "starting");

	/*
	 * If ranks doesn't include myself, creating a group with ranks will
	 * fail; bcast to the primary group instead.
	 */
	rc = bcast_create(RSVC_START, NULL /* group */, &rpc, &nstarted);
	if (rc != 0)
		goto out;
	in = crt_req_get(rpc);
	in->sai_class = class;
	rc = daos_iov_copy(&in->sai_svc_id, id);
	if (rc != 0)
		goto out_rpc;
	uuid_copy(in->sai_db_uuid, dbid);
	if (create)
		in->sai_flags |= RDB_AF_CREATE;
	if (bootstrap)
		in->sai_flags |= RDB_AF_BOOTSTRAP;
	in->sai_size = size;
	in->sai_ranks = (d_rank_list_t *)ranks;
	dss_rpc_send(rpc); /* Ignore return code; errors are detected below. */
	out = crt_reply_get(rpc);
	rsvc_corpc_incr_cb(in->sai_ranks, rpc->cr_ep.ep_rank, out->sao_rc,
			   &nstarted);

	if (!nstarted || (ranks && nstarted < ranks->rl_nr)) {
		D_ERROR(DF_UUID": failed to start%s some replicas\n",
			DP_UUID(dbid), create ? "/create" : "");
		rc = -DER_IO;
	}

	daos_iov_free(&in->sai_svc_id);
out_rpc:
	crt_req_decref(rpc);
out:
	return rc;
}

static void
ds_rsvc_start_handler(crt_rpc_t *rpc)
{
	struct rsvc_start_in	*in = crt_req_get(rpc);
	struct rsvc_start_out	*out = crt_reply_get(rpc);
	bool			 create = in->sai_flags & RDB_AF_CREATE;
	int			 rc;

	if (create && in->sai_ranks == NULL) {
		rc = -DER_PROTO;
		goto out;
	}

	if (in->sai_ranks != NULL) {
		d_rank_t	rank;
		int		i;

		/* Do nothing if I'm not one of the replicas. */
		rc = crt_group_rank(NULL /* grp */, &rank);
		D_ASSERTF(rc == 0, "%d\n", rc);
		if (!daos_rank_list_find(in->sai_ranks, rank, &i))
			goto out;
	}

	rc = ds_rsvc_start(in->sai_class, &in->sai_svc_id, in->sai_db_uuid,
			   create, in->sai_size,
			   (in->sai_flags & RDB_AF_BOOTSTRAP) ?
			   in->sai_ranks : NULL, NULL /* arg */);
out:
	out->sao_rc = rc;
	crt_reply_send(rpc);
}

static int
ds_rsvc_start_aggregator(crt_rpc_t *source, crt_rpc_t *result, void *priv)
{
	struct rsvc_start_in	*in_source	= crt_req_get(source);
	struct rsvc_start_out   *out_source	= crt_reply_get(source);

	rsvc_corpc_incr_cb(in_source->sai_ranks, source->cr_ep.ep_rank,
			   out_source->sao_rc, (int *)priv);
	return 0;
}

/**
 * Perform a distributed stop, and if \a destroy is true, destroy operation on
 * all replicas of a database spanning \a ranks. This method can be called on
 * any rank. \a ranks may be NULL.
 *
 * \param[in]	class		replicated service class
 * \param[in]	id		replicated service ID
 * \param[in]	ranks		list of \a ranks->rl_nr replica ranks
 * \param[in]	destroy		destroy after close
 */
int
ds_rsvc_dist_stop(enum ds_rsvc_class_id class, d_iov_t *id,
		  const d_rank_list_t *ranks, bool destroy)
{
	crt_rpc_t		*rpc;
	struct rsvc_stop_in	*in;
	struct rsvc_stop_out	*out;
	int			 nstopped = 0;
	int			 rc;

	/*
	 * If ranks doesn't include myself, creating a group with ranks will
	 * fail; bcast to the primary group instead.
	 */
	rc = bcast_create(RSVC_STOP, NULL /* group */, &rpc, &nstopped);
	if (rc != 0)
		goto out;
	in = crt_req_get(rpc);
	in->soi_class = class;
	rc = daos_iov_copy(&in->soi_svc_id, id);
	if (rc != 0)
		goto out_rpc;
	if (destroy)
		in->soi_flags |= RDB_OF_DESTROY;
	in->soi_ranks = (d_rank_list_t *)ranks;

	dss_rpc_send(rpc); /* Ignore return code; errors are detected below. */
	out = crt_reply_get(rpc);
	rsvc_corpc_incr_cb(in->soi_ranks, rpc->cr_ep.ep_rank, out->soo_rc,
			   &nstopped);

	if (!nstopped || (ranks && nstopped < ranks->rl_nr)) {
		D_ERROR("failed to stop%s some replicas\n",
			destroy ? "/destroy" : "");
		rc = -DER_IO;
	}

	daos_iov_free(&in->soi_svc_id);
out_rpc:
	crt_req_decref(rpc);
out:
	return rc;
}

static void
ds_rsvc_stop_handler(crt_rpc_t *rpc)
{
	struct rsvc_stop_in	*in = crt_req_get(rpc);
	struct rsvc_stop_out	*out = crt_reply_get(rpc);
	int			 rc = 0;

	if (in->soi_ranks != NULL) {
		d_rank_t	rank;
		int		i;

		/* Do nothing if I'm not one of the replicas. */
		rc = crt_group_rank(NULL /* grp */, &rank);
		D_ASSERTF(rc == 0, "%d\n", rc);
		if (!daos_rank_list_find(in->soi_ranks, rank, &i))
			goto out;
	}

	rc = ds_rsvc_stop(in->soi_class, &in->soi_svc_id,
			  in->soi_flags & RDB_OF_DESTROY);
out:
	out->soo_rc = rc;
	crt_reply_send(rpc);
}

static int
ds_rsvc_stop_aggregator(crt_rpc_t *source, crt_rpc_t *result, void *priv)
{
	struct rsvc_stop_in	*in_source	= crt_req_get(source);
	struct rsvc_stop_out	*out_source	= crt_reply_get(source);

	rsvc_corpc_incr_cb(in_source->soi_ranks, source->cr_ep.ep_rank,
			   out_source->soo_rc, (int *)priv);
	return 0;
}

static struct crt_corpc_ops ds_rsvc_start_co_ops = {
	.co_aggregate	= ds_rsvc_start_aggregator,
	.co_pre_forward	= NULL,
};

static struct crt_corpc_ops ds_rsvc_stop_co_ops = {
	.co_aggregate	= ds_rsvc_stop_aggregator,
	.co_pre_forward	= NULL,
};

#define X(a, b, c, d, e)	\
{				\
	.dr_opc       = a,	\
	.dr_hdlr      = d,	\
	.dr_corpc_ops = e,	\
}

static struct daos_rpc_handler rsvc_handlers[] = {
	RSVC_PROTO_SRV_RPC_LIST,
};

#undef X

size_t
ds_rsvc_get_md_cap(void)
{
	const size_t	size_default = 1 << 27 /* 128 MB */;
	char	       *v;
	int		n;

	v = getenv("DAOS_MD_CAP"); /* in MB */
	if (v == NULL)
		return size_default;
	n = atoi(v);
	if (n < size_default >> 20) {
		D_ERROR("metadata capacity too low; using %zu MB\n",
			size_default >> 20);
		return size_default;
	}
	return (size_t)n << 20;
}

static int
rsvc_module_init(void)
{
	return rsvc_hash_init();
}

static int
rsvc_module_fini(void)
{
	rsvc_hash_fini();
	return 0;
}
struct dss_module rsvc_module = {
	.sm_name	= "rsvc",
	.sm_mod_id	= DAOS_RSVC_MODULE,
	.sm_ver		= DAOS_RSVC_VERSION,
	.sm_init	= rsvc_module_init,
	.sm_fini	= rsvc_module_fini,
	.sm_proto_fmt	= &rsvc_proto_fmt,
	.sm_cli_count	= 0,
	.sm_handlers	= rsvc_handlers,
	.sm_key		= NULL,
};
