/**
 * (C) Copyright 2015-2021 Intel Corporation.
 *
 * SPDX-License-Identifier: BSD-2-Clause-Patent
 */
#define D_LOGFAC	DD_FAC(client)

#include <daos/pool.h>
#include <daos/pool_map.h>
#include <daos/task.h>
#include "client_internal.h"
#include "task_internal.h"

int
daos_pool_connect(const uuid_t uuid, const char *grp,
		  unsigned int flags,
		  daos_handle_t *poh, daos_pool_info_t *info, daos_event_t *ev)
{
	daos_pool_connect_t	*args;
	tse_task_t		*task;
	int			 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_CONNECT);
	if (!daos_uuid_valid(uuid))
		return -DER_INVAL;

	rc = dc_task_create(dc_pool_connect, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->grp		= grp;
	args->flags		= flags;
	args->poh		= poh;
	args->info		= info;
	uuid_copy((unsigned char *)args->uuid, uuid);
	args->label		= NULL;

	return dc_task_schedule(task, true);
}

int
daos_pool_connect_by_label(const char *label, const char *grp,
			   unsigned int flags, daos_handle_t *poh,
			   daos_pool_info_t *info, daos_event_t *ev)
{
	daos_pool_connect_t	*args;
	tse_task_t		*task;
	size_t			 label_len = 0;
	int			 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_CONNECT);
	if (label)
		label_len = strnlen(label, DAOS_PROP_LABEL_MAX_LEN+1);
	if (!label || (label_len == 0) ||
	    (label_len > DAOS_PROP_LABEL_MAX_LEN)) {
		D_ERROR("invalid label parameter\n");
		return -DER_INVAL;
	}

	rc = dc_task_create(dc_pool_connect_lbl, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->grp		= grp;
	args->flags		= flags;
	args->poh		= poh;
	args->info		= info;
	uuid_clear(args->uuid);
	args->label		= label;

	return dc_task_schedule(task, true);
}

int
daos_pool_disconnect(daos_handle_t poh, daos_event_t *ev)
{
	daos_pool_disconnect_t	*args;
	tse_task_t		*task;
	int			 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_DISCONNECT);
	rc = dc_task_create(dc_pool_disconnect, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->poh = poh;

	return dc_task_schedule(task, true);
}

int
daos_pool_local2global(daos_handle_t poh, d_iov_t *glob)
{
	return dc_pool_local2global(poh, glob);
}

int
daos_pool_global2local(d_iov_t glob, daos_handle_t *poh)
{
	return dc_pool_global2local(glob, poh);
}

int
daos_pool_query(daos_handle_t poh, d_rank_list_t *tgts, daos_pool_info_t *info,
		daos_prop_t *pool_prop, daos_event_t *ev)
{
	daos_pool_query_t	*args;
	tse_task_t		*task;
	int			 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_QUERY);

	if (pool_prop != NULL && !daos_prop_valid(pool_prop, true, false)) {
		D_ERROR("invalid pool_prop parameter.\n");
		return -DER_INVAL;
	}

	rc = dc_task_create(dc_pool_query, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->poh	= poh;
	args->tgts	= tgts;
	args->info	= info;
	args->prop	= pool_prop;

	return dc_task_schedule(task, true);
}

int
daos_pool_query_target(daos_handle_t poh, uint32_t tgt_idx, d_rank_t rank,
		       daos_target_info_t *info, daos_event_t *ev)
{
	daos_pool_query_target_t	*args;
	tse_task_t			*task;
	int				 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_QUERY_INFO);

	rc = dc_task_create(dc_pool_query_target, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->poh	= poh;
	args->tgt_idx	= tgt_idx;
	args->rank	= rank;
	args->info	= info;

	return dc_task_schedule(task, true);

}

int
daos_pool_list_cont(daos_handle_t poh, daos_size_t *ncont,
		    struct daos_pool_cont_info *cbuf, daos_event_t *ev)
{
	daos_pool_list_cont_t	*args;
	tse_task_t		*task;
	int			 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_LIST_CONT);

	if (ncont == NULL) {
		D_ERROR("ncont must be non-NULL\n");
		return -DER_INVAL;
	}

	rc = dc_task_create(dc_pool_list_cont, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->poh	= poh;
	args->ncont	= ncont;
	args->cont_buf	= cbuf;

	return dc_task_schedule(task, true);
}

int
daos_pool_list_attr(daos_handle_t poh, char *buf, size_t *size,
		    daos_event_t *ev)
{
	daos_pool_list_attr_t	*args;
	tse_task_t		*task;
	int			 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_LIST_ATTR);

	rc = dc_task_create(dc_pool_list_attr, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->poh	= poh;
	args->buf	= buf;
	args->size	= size;

	return dc_task_schedule(task, true);
}

static int
free_name(tse_task_t *task, void *args)
{
	char *name = *(char **)args;

	D_FREE(name);
	return 0;
}


int
daos_pool_get_attr(daos_handle_t poh, int n, char const *const names[],
		   void *const values[], size_t sizes[], daos_event_t *ev)
{
	daos_pool_get_attr_t	*args;
	tse_task_t		*task;
	int			 i, rc;
	char			**new_names;

	DAOS_API_ARG_ASSERT(*args, POOL_GET_ATTR);

	D_ALLOC_ARRAY(new_names, n);
	if (!new_names)
		return -DER_NOMEM;

	rc = dc_task_create(dc_pool_get_attr, NULL, ev, &task);
	if (rc) {
		D_FREE(new_names);
		return rc;
	}

	for (i = 0 ; i < n ; i++) {
		/* no easy way to determine if a name storage address is likely
		 * to cause an EFAULT during memory registration, so duplicate
		 * name in heap
		 */
		D_STRNDUP(new_names[i], names[i], DAOS_ATTR_NAME_MAX);
		if (new_names[i] == NULL)
			D_GOTO(error, rc = -DER_NOMEM);
		rc = tse_task_register_comp_cb(task, free_name, &new_names[i],
					      sizeof(char *));
		if (rc)
			D_GOTO(error, rc);
	}

	args = dc_task_get_args(task);
	args->poh	= poh;
	args->n		= n;
	args->names	= (char const *const *)new_names;
	args->values	= values;
	args->sizes	= sizes;

	return dc_task_schedule(task, true);

error:
	/* this will cause new_names[] to be freed by callback */
	tse_task_complete(task, rc);
	D_FREE(new_names);
	return rc;
}

int
daos_pool_set_attr(daos_handle_t poh, int n, char const *const names[],
		   void const *const values[], size_t const sizes[],
		   daos_event_t *ev)
{
	daos_pool_set_attr_t	*args;
	tse_task_t		*task;
	int			 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_SET_ATTR);

	rc = dc_task_create(dc_pool_set_attr, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->poh	= poh;
	args->n		= n;
	args->names	= names;
	args->values	= values;
	args->sizes	= sizes;

	return dc_task_schedule(task, true);
}

int
daos_pool_del_attr(daos_handle_t poh, int n, char const *const names[],
		   daos_event_t *ev)
{
	daos_pool_del_attr_t	*args;
	tse_task_t		*task;
	int			 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_DEL_ATTR);

	rc = dc_task_create(dc_pool_del_attr, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->poh	= poh;
	args->n		= n;
	args->names	= names;

	return dc_task_schedule(task, true);
}

int
daos_pool_stop_svc(daos_handle_t poh, daos_event_t *ev)
{
	daos_pool_stop_svc_t	*args;
	tse_task_t		*task;
	int			 rc;

	DAOS_API_ARG_ASSERT(*args, POOL_STOP_SVC);
	rc = dc_task_create(dc_pool_stop_svc, NULL, ev, &task);
	if (rc)
		return rc;

	args = dc_task_get_args(task);
	args->poh	= poh;

	return dc_task_schedule(task, true);
}
