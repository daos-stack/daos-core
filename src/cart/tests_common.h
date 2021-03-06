/*
 * (C) Copyright 2019-2021 Intel Corporation.
 *
 * SPDX-License-Identifier: BSD-2-Clause-Patent
 */
/**
 * Common functions to be shared among tests
 */
#ifndef __TESTS_COMMON_H__
#define __TESTS_COMMON_H__
#include <semaphore.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <cart/api.h>

#include "crt_internal.h"

#define DBG_PRINT(x...)							\
	do {								\
		D_INFO(x);						\
		if (opts.is_server)					\
			fprintf(stderr, "SRV [rank=%d pid=%d]\t",       \
			opts.self_rank,					\
			opts.mypid);					\
		else							\
			fprintf(stderr, "CLI [rank=%d pid=%d]\t",       \
			opts.self_rank,					\
			opts.mypid);					\
		fprintf(stderr, x);					\
	} while (0)

struct test_options {
	bool		is_initialized;
	d_rank_t	self_rank;
	int		mypid;
	int		num_attach_retries;
	bool		is_server;
	bool		assert_on_error;
	volatile int	shutdown;
	int		delay_shutdown_sec;
	bool		is_swim_enabled;
};

static struct test_options opts = { .is_initialized = false };

void
tc_test_init(d_rank_t rank, int num_attach_retries, bool is_server,
	     bool assert_on_error)
{
	opts.is_initialized	= true;
	opts.self_rank		= rank;
	opts.mypid		= getpid();
	opts.is_server		= is_server;
	opts.num_attach_retries	= num_attach_retries;
	opts.assert_on_error	= assert_on_error;
	opts.shutdown		= 0;
	opts.is_swim_enabled	= false;

	/* Use 2 second delay as a default for all tests for now */
	opts.delay_shutdown_sec	= 2;
}

static inline void
tc_test_swim_enable(bool is_swim_enabled)
{
	opts.is_swim_enabled	= is_swim_enabled;
}

static inline int
tc_drain_queue(crt_context_t ctx)
{
	int	rc;
	int	i;

	/* TODO: Need better mechanism for tests to drain all queues */
	for (i = 0; i < 1000; i++)
		crt_progress(ctx, 1000);

	/* Drain the queue. Progress until 1 second timeout.  We need
	 * a more robust method
	 */
	do {
		rc = crt_progress(ctx, 1000000);
		if (rc != 0 && rc != -DER_TIMEDOUT) {
			D_ERROR("crt_progress failed rc: %d.\n", rc);
			return rc;
		}

		if (rc == -DER_TIMEDOUT)
			break;
	} while (1);

	D_DEBUG(DB_TEST, "Done draining queue\n");
	return 0;
}

void
tc_set_shutdown_delay(int delay_sec)
{
	opts.delay_shutdown_sec = delay_sec;
}

void
tc_progress_stop(void)
{
	opts.shutdown = 1;
}

/* Write a completion file to signal graceful server shutdown */
void
write_completion_file(void)
{
	FILE	*fptr;
	char	*dir;
	char	*completion_file;
	char	*tmp_str;
	char	pid[6];
	pid_t	_pid;

	_pid = getpid();
	sprintf(pid, "%d", _pid);

	dir = getenv("DAOS_TEST_SHARED_DIR");
	D_ASSERTF(dir != NULL,
		"DAOS_TEST_SHARED_DIR must be set for --write_completion_file "
		"option.\n");
	tmp_str = strcat(dir, "/test-servers-completed.txt.");
	completion_file = strcat(tmp_str, pid);

	unlink(completion_file);
	fptr = fopen(completion_file, "w");
	D_ASSERTF(fptr != NULL, "Error opening completion file for writing.\n");
	DBG_PRINT("Wrote completion file: %s.\n", completion_file);
	fclose(fptr);
}


void *
tc_progress_fn(void *data)
{
	int		rc;
	int		idx = -1;
	crt_context_t	*p_ctx = (crt_context_t *)data;

	D_ASSERTF(opts.is_initialized == true, "tc_test_init not called.\n");

	rc = crt_context_idx(*p_ctx, &idx);
	if (rc != 0) {
		D_ERROR("crt_context_idx() failed; rc=%d\n", rc);
		assert(0);
	}

	while (opts.shutdown == 0)
		crt_progress(*p_ctx, 1000);

	if (opts.is_swim_enabled && idx == 0)
		crt_swim_fini();

	rc = tc_drain_queue(*p_ctx);
	D_ASSERTF(rc == 0, "tc_drain_queue() failed with rc=%d\n", rc);

	if (opts.delay_shutdown_sec > 0)
		sleep(opts.delay_shutdown_sec);

	rc = crt_context_destroy(*p_ctx, 1);
	D_ASSERTF(rc == 0, "Failed to destroy context %p rc=%d\n",
		  p_ctx, rc);

	pthread_exit(rc ? *p_ctx : NULL);

	return NULL;
}

struct wfr_status {
	sem_t	sem;
	int	rc;
	int	num_ctx;
};

static inline void
tc_sync_timedwait(struct wfr_status *wfrs, int sec, int line_number)
{
	struct timespec	deadline;
	int		rc;

	rc = clock_gettime(CLOCK_REALTIME, &deadline);
	D_ASSERTF(rc == 0, "clock_gettime() failed at line %d rc: %d\n",
		  line_number, rc);

	deadline.tv_sec += sec;

	rc = sem_timedwait(&wfrs->sem, &deadline);
	D_ASSERTF(rc == 0, "Sync timed out at line %d rc: %d\n",
		  line_number, rc);
}

static void
ctl_client_cb(const struct crt_cb_info *info)
{
	struct wfr_status		*wfrs;
	struct crt_ctl_ep_ls_out	*out_ls_args;
	char				*addr_str;
	int				 i;

	wfrs = (struct wfr_status *)info->cci_arg;

	if (info->cci_rc == 0) {
		out_ls_args = crt_reply_get(info->cci_rpc);
		wfrs->num_ctx = out_ls_args->cel_ctx_num;
		wfrs->rc = out_ls_args->cel_rc;

		D_DEBUG(DB_TEST, "ctx_num: %d\n",
			out_ls_args->cel_ctx_num);
		addr_str = out_ls_args->cel_addr_str.iov_buf;
		for (i = 0; i < out_ls_args->cel_ctx_num; i++) {
			D_DEBUG(DB_TEST, "    %s\n", addr_str);
				addr_str += (strlen(addr_str) + 1);
		}
	} else {
		wfrs->rc = info->cci_rc;
	}

	sem_post(&wfrs->sem);
}

int
tc_wait_for_ranks(crt_context_t ctx, crt_group_t *grp, d_rank_list_t *rank_list,
		  int tag, int total_ctx, double ping_timeout,
		  double total_timeout)
{
	struct wfr_status		ws;
	struct timespec			t1, t2;
	struct crt_ctl_ep_ls_in		*in_args;
	d_rank_t			rank;
	crt_rpc_t			*rpc = NULL;
	crt_endpoint_t			server_ep;
	double				time_s = 0;
	int				i = 0;
	int				rc = 0;

	D_ASSERTF(opts.is_initialized == true, "tc_test_init not called.\n");

	rc = d_gettime(&t1);
	D_ASSERTF(rc == 0, "d_gettime() failed; rc=%d\n", rc);

	rc = sem_init(&ws.sem, 0, 0);
	D_ASSERTF(rc == 0, "sem_init() failed; rc=%d\n", rc);

	server_ep.ep_tag = tag;
	server_ep.ep_grp = grp;

	for (i = 0; i < rank_list->rl_nr; i++) {
		rank = rank_list->rl_ranks[i];
		server_ep.ep_rank = rank;

		rc = crt_req_create(ctx, &server_ep, CRT_OPC_CTL_LS, &rpc);
		D_ASSERTF(rc == 0, "crt_req_create failed; rc=%d\n", rc);

		in_args = crt_req_get(rpc);
		in_args->cel_grp_id = grp->cg_grpid;
		in_args->cel_rank = rank;

		rc = crt_req_set_timeout(rpc, ping_timeout);
		D_ASSERTF(rc == 0, "crt_req_set_timeout failed; rc=%d\n", rc);

		ws.rc = 0;
		ws.num_ctx = 0;

		rc = crt_req_send(rpc, ctl_client_cb, &ws);

		if (rc == 0)
			tc_sync_timedwait(&ws, 120, __LINE__);
		else
			ws.rc = rc;

		while (ws.rc != 0 && time_s < total_timeout) {
			rc = crt_req_create(ctx, &server_ep,
					    CRT_OPC_CTL_LS, &rpc);
			D_ASSERTF(rc == 0,
				  "crt_req_create failed; rc=%d\n", rc);

			in_args = crt_req_get(rpc);
			in_args->cel_grp_id = grp->cg_grpid;
			in_args->cel_rank = rank;

			rc = crt_req_set_timeout(rpc, ping_timeout);
			D_ASSERTF(rc == 0,
				  "crt_req_set_timeout failed; rc=%d\n", rc);

			ws.rc = 0;
			ws.num_ctx = 0;

			rc = crt_req_send(rpc, ctl_client_cb, &ws);

			if (rc == 0)
				tc_sync_timedwait(&ws, 120, __LINE__);
			else
				ws.rc = rc;

			rc = d_gettime(&t2);
			D_ASSERTF(rc == 0, "d_gettime() failed; rc=%d\n", rc);
			time_s = d_time2s(d_timediff(t1, t2));
		}

		if (ws.rc != 0) {
			rc = ws.rc;
			break;
		}

		if (ws.num_ctx < total_ctx) {
			rc = -1;
			break;
		}
	}

	sem_destroy(&ws.sem);

	return rc;
}

int
tc_load_group_from_file(const char *grp_cfg_file, crt_context_t ctx,
			crt_group_t *grp, d_rank_t my_rank, bool delete_file)
{
	FILE		*f;
	int		parsed_rank;
	char		parsed_addr[255];
	int		rc = 0;

	D_ASSERTF(opts.is_initialized == true, "tc_test_init not called.\n");

	f = fopen(grp_cfg_file, "r");
	if (!f) {
		D_ERROR("Failed to open %s for reading\n", grp_cfg_file);
		D_GOTO(out, rc = -DER_NONEXIST);
	}

	while (1) {
		rc = fscanf(f, "%d %s", &parsed_rank, parsed_addr);
		if (rc == EOF) {
			rc = 0;
			break;
		}

		if (parsed_rank == my_rank)
			continue;

		rc = crt_group_primary_rank_add(ctx, grp,
						parsed_rank, parsed_addr);

		if (rc != 0) {
			D_ERROR("Failed to add %d %s; rc=%d\n",
				parsed_rank, parsed_addr, rc);
			break;
		}
	}

	fclose(f);

	if (delete_file)
		unlink(grp_cfg_file);

out:
	return rc;
}

static inline int
tc_sem_timedwait(sem_t *sem, int sec, int line_number)
{
	struct timespec	deadline;
	int		rc;

	rc = clock_gettime(CLOCK_REALTIME, &deadline);
	if (rc != 0) {
		if (opts.assert_on_error)
			D_ASSERTF(rc == 0, "clock_gettime() failed at "
				  "line %d rc: %d\n", line_number, rc);
		D_ERROR("clock_gettime() failed, rc = %d\n", rc);
		D_GOTO(out, rc);
	}

	deadline.tv_sec += sec;
	rc = sem_timedwait(sem, &deadline);
	if (rc != 0) {
		if (opts.assert_on_error)
			D_ASSERTF(rc == 0, "sem_timedwait() failed at "
				  "line %d rc: %d\n", line_number, rc);
		D_ERROR("sem_timedwait() failed, rc = %d\n", rc);
		D_GOTO(out, rc);
	}

out:
	return rc;
}

void
tc_cli_start_basic(char *local_group_name, char *srv_group_name,
		   crt_group_t **grp, d_rank_list_t **rank_list,
		   crt_context_t *crt_ctx, pthread_t *progress_thread,
		   unsigned int total_srv_ctx, bool use_cfg,
		   crt_init_options_t *init_opt)
{
	char		*grp_cfg_file;
	uint32_t	 grp_size;
	int		 attach_retries = opts.num_attach_retries;
	int		 rc = 0;

	D_ASSERTF(opts.is_initialized == true, "tc_test_init not called.\n");

	rc = d_log_init();
	D_ASSERTF(rc == 0, "d_log_init failed, rc=%d\n", rc);

	if (init_opt)
		rc = crt_init_opt(local_group_name, 0, init_opt);
	else
		rc = crt_init(local_group_name, 0);
	D_ASSERTF(rc == 0, "crt_init() failed; rc=%d\n", rc);

	rc = crt_context_create(crt_ctx);
	D_ASSERTF(rc == 0, "crt_context_create() failed; rc=%d\n", rc);

	rc = pthread_create(progress_thread, NULL, tc_progress_fn, crt_ctx);
	D_ASSERTF(rc == 0, "pthread_create() failed; rc=%d\n", rc);

	if (use_cfg) {
		while (attach_retries-- > 0) {
			rc = crt_group_attach(srv_group_name, grp);
			if (rc == 0)
				break;
			sleep(1);
		}
		D_ASSERTF(rc == 0, "crt_group_attach failed, rc: %d\n", rc);
		D_ASSERTF(*grp != NULL, "NULL attached remote grp\n");
	} else {
		rc = crt_group_view_create(srv_group_name, grp);
		if (!*grp || rc != 0) {
			D_ERROR("Failed to create group view; rc=%d\n", rc);
			assert(0);
		}

		grp_cfg_file = getenv("CRT_L_GRP_CFG");

		/* load group info from a config file and
		 * delete file upon return
		 */
		rc = tc_load_group_from_file(grp_cfg_file, *crt_ctx, *grp, -1,
					     true);
		D_ASSERTF(rc == 0, "tc_load_group_from_file() failed; rc=%d\n",
			  rc);
	}

	rc = crt_group_size(*grp, &grp_size);
	D_ASSERTF(rc == 0, "crt_group_size() failed; rc=%d\n", rc);

	rc = crt_group_ranks_get(*grp, rank_list);
	D_ASSERTF(rc == 0, "crt_group_ranks_get() failed; rc=%d\n", rc);

	if (!*rank_list) {
		D_ERROR("Rank list is NULL\n");
		assert(0);
	}

	if ((*rank_list)->rl_nr != grp_size) {
		D_ERROR("rank_list differs in size. expected %d got %d\n",
			grp_size, (*rank_list)->rl_nr);
		assert(0);
	}

	rc = crt_group_psr_set(*grp, (*rank_list)->rl_ranks[0]);
	D_ASSERTF(rc == 0, "crt_group_psr_set() failed; rc=%d\n", rc);
}

void
tc_srv_start_basic(char *srv_group_name, crt_context_t *crt_ctx,
		   pthread_t *progress_thread, crt_group_t **grp,
		   uint32_t *grp_size, crt_init_options_t *init_opt)
{
	char		*env_self_rank;
	char		*grp_cfg_file;
	char		*my_uri;
	d_rank_t	 my_rank;
	int		 rc = 0;

	D_ASSERTF(opts.is_initialized == true, "tc_test_init not called.\n");

	env_self_rank = getenv("CRT_L_RANK");
	my_rank = atoi(env_self_rank);

	rc = d_log_init();
	D_ASSERT(rc == 0);

	if (init_opt) {
		rc = crt_init_opt(srv_group_name, CRT_FLAG_BIT_SERVER |
				  CRT_FLAG_BIT_AUTO_SWIM_DISABLE, init_opt);
	} else {
		rc = crt_init(srv_group_name, CRT_FLAG_BIT_SERVER |
			      CRT_FLAG_BIT_AUTO_SWIM_DISABLE);
	}
	D_ASSERTF(rc == 0, "crt_init() failed, rc: %d\n", rc);

	*grp = crt_group_lookup(NULL);
	if (!(*grp)) {
		D_ERROR("Failed to lookup group\n");
		assert(0);
	}

	rc = crt_rank_self_set(my_rank);
	D_ASSERTF(rc == 0, "crt_rank_self_set(%d) failed; rc=%d\n",
		  my_rank, rc);

	rc = crt_context_create(crt_ctx);
	D_ASSERTF(rc == 0, "crt_context_create() failed; rc=%d\n", rc);

	rc = pthread_create(progress_thread, NULL, tc_progress_fn, crt_ctx);
	D_ASSERTF(rc == 0, "pthread_create() failed; rc=%d\n", rc);

	grp_cfg_file = getenv("CRT_L_GRP_CFG");

	rc = crt_rank_uri_get(*grp, my_rank, 0, &my_uri);
	D_ASSERTF(rc == 0, "crt_rank_uri_get() failed; rc=%d\n", rc);

	/* load group info from a config file and delete file upon return */
	rc = tc_load_group_from_file(grp_cfg_file, crt_ctx[0], *grp, my_rank,
				     true);
	D_ASSERTF(rc == 0, "tc_load_group_from_file() failed; rc=%d\n", rc);

	D_FREE(my_uri);

	if (opts.is_swim_enabled) {
		rc = crt_swim_init(0);
		D_ASSERTF(rc == 0, "crt_swim_init() failed; rc=%d\n", rc);
	}

	rc = crt_group_size(NULL, grp_size);
	D_ASSERTF(rc == 0, "crt_group_size() failed; rc=%d\n", rc);
}

struct tc_log_msg_cb_resp {
	sem_t	sem;
};

static void
tc_log_msg_cb(const struct crt_cb_info *info)
{
	struct tc_log_msg_cb_resp	*resp;

	if (info->cci_rc != 0) {
		D_WARN("Add Log message CB failed\n");
		D_ASSERTF(info->cci_rc == 0,
			  "Send Log RPC did not respond\n");
	}
	resp = (struct tc_log_msg_cb_resp *)info->cci_arg;
	sem_post(&resp->sem);
}

int
tc_log_msg(crt_context_t ctx, crt_group_t *grp, d_rank_t rank,
	   char *msg) {
	int32_t				 rc = 0;
	struct crt_ctl_log_add_msg_in	*send_args;
	crt_rpc_t			*rpc_req = NULL;
	crt_endpoint_t			 ep;
	crt_opcode_t			 opcode = CRT_OPC_CTL_LOG_ADD_MSG;
	struct tc_log_msg_cb_resp	 resp;

	/* Initialize response structure */
	rc = sem_init(&resp.sem, 0, 0);
	D_ASSERTF(rc == 0, "sem_init() failed\n");

	/* Fill in the endpoint info */
	ep.ep_grp = grp;
	ep.ep_rank = rank;
	ep.ep_tag = 0;

	rc = crt_req_create(ctx, &ep, opcode, &rpc_req);
	if (rc != 0) {
		D_ERROR("crt_req_create() failed. rc %d.\n", rc);
		D_GOTO(exit, rc);
	}

	crt_req_addref(rpc_req);
	send_args =  crt_req_get(rpc_req);
	send_args->log_msg = msg;

	/* send the request */
	rc = crt_req_send(rpc_req, tc_log_msg_cb, &resp);
	if (rc < 0) {
		D_WARN("rpc failed, message: \"%s \"not sent\n", msg);
		goto cleanup;
	}

	/* Wait for response */
	rc = tc_sem_timedwait(&resp.sem, 30, __LINE__);
	if (rc < 0) {
		D_WARN("Messaage logged timed out: %s\n", msg);
		crt_req_abort(rpc_req);
		goto cleanup;
	}

	/* Decrement reference */
cleanup:
	crt_req_decref(rpc_req);
exit:
	return rc;
}

#endif /* __TESTS_COMMON_H__ */
