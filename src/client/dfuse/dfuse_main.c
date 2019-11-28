/**
 * (C) Copyright 2016-2019 Intel Corporation.
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

#include <errno.h>
#include <getopt.h>
#include <dlfcn.h>
#include <fuse3/fuse.h>
#include <fuse3/fuse_lowlevel.h>

#define D_LOGFAC DD_FAC(dfuse)

#include "dfuse.h"

#include "daos_fs.h"
#include "daos_api.h"

#include <gurt/common.h>

static int
ll_loop_fn(struct dfuse_info *dfuse_info)
{
	int			ret;

	/*Blocking*/
	if (dfuse_info->di_threaded) {
		struct fuse_loop_config config = {.max_idle_threads = 10};

		ret = fuse_session_loop_mt(dfuse_info->di_session,
					   &config);
	} else {
		ret = fuse_session_loop(dfuse_info->di_session);
	}
	if (ret != 0)
		DFUSE_LOG_ERROR("Fuse loop exited with return code: %d", ret);

	return ret;
}

/*
 * Creates a fuse filesystem for any plugin that needs one.
 *
 * Should be called from the post_start plugin callback and creates
 * a filesystem.
 * Returns 0 on success, or non-zero on error.
 */
bool
dfuse_launch_fuse(struct dfuse_info *dfuse_info,
		  struct fuse_lowlevel_ops *flo,
		  struct fuse_args *args,
		  struct dfuse_projection_info *fs_handle)
{
	int rc;

	dfuse_info->di_handle = fs_handle;

	dfuse_info->di_session = fuse_session_new(args,
						   flo,
						   sizeof(*flo),
						   fs_handle);
	if (!dfuse_info->di_session)
		goto cleanup;

	rc = fuse_session_mount(dfuse_info->di_session,
				dfuse_info->di_mountpoint);
	if (rc != 0) {
		goto cleanup;
	}

	fuse_opt_free_args(args);

	rc = ll_loop_fn(dfuse_info);
	fuse_session_unmount(dfuse_info->di_session);
	if (rc) {
		goto cleanup;
	}

	return true;
cleanup:
	return false;
}

static void
show_help(char *name)
{
	printf("usage: %s -m=PATHSTR -s=RANKS\n"
		"\n"
		"	-m --mountpoint=PATHSTR	Mount point to use\n"
		"	-s --svc=RANKS		pool service replicas like 1,2,3\n"
		"	   --pool=UUID		pool UUID\n"
		"	   --container=UUID	container UUID\n"
		"	   --sys-name=STR	DAOS system name context for servers\n"
		"	-S --singlethreaded	Single threaded\n"
		"	-f --foreground		Run in foreground\n",
		name);
}

int
main(int argc, char **argv)
{
	struct dfuse_info	*dfuse_info = NULL;
	char			*svcl = NULL;
	struct dfuse_dfs	*dfs = NULL;
	char			c;
	int			ret = -DER_SUCCESS;
	int			rc;

	/* The 'daos' command uses -m as an alias for --scv however
	 * dfuse uses -m for --mountpoint so this is inconsistent
	 * but probably better than changing the meaning of the -m
	 * option here.h
	 */
	struct option long_options[] = {
		{"pool",		required_argument, 0, 'p'},
		{"container",		required_argument, 0, 'c'},
		{"svc",			required_argument, 0, 's'},
		{"sys-name",		required_argument, 0, 'G'},
		{"mountpoint",		required_argument, 0, 'm'},
		{"singlethread",	no_argument,	   0, 'S'},
		{"foreground",		no_argument,	   0, 'f'},
		{"help",		no_argument,	   0, 'h'},
		{0, 0, 0, 0}
	};

	D_ALLOC_PTR(dfuse_info);
	if (!dfuse_info)
		D_GOTO(out_fini, ret = -DER_NOMEM);

	D_INIT_LIST_HEAD(&dfuse_info->di_dfs_list);

	dfuse_info->di_threaded = true;

	while (1) {
		c = getopt_long(argc, argv, "s:m:Sfh",
				long_options, NULL);

		if (c == -1)
			break;

		switch (c) {
		case 'p':
			dfuse_info->di_pool = optarg;
			break;
		case 'c':
			dfuse_info->di_cont = optarg;
			break;
		case 's':
			svcl = optarg;
			break;
		case 'G':
			dfuse_info->di_group = optarg;
			break;
		case 'm':
			dfuse_info->di_mountpoint = optarg;
			break;
		case 'S':
			dfuse_info->di_threaded = false;
			break;
		case 'f':
			dfuse_info->di_foreground = true;
			break;
		case 'h':
			show_help(argv[0]);
			exit(0);
			break;
		case '?':
			show_help(argv[0]);
			exit(1);
			break;
		}
	}

	if (!dfuse_info->di_foreground && getenv("PMIX_RANK")) {
		DFUSE_LOG_WARNING("Not running in background under orterun");
		dfuse_info->di_foreground = true;
	}

	if (!dfuse_info->di_mountpoint) {
		printf("Mountpoint is required\n");
		show_help(argv[0]);
		exit(1);
	}

	/* Is this required, or can we assume some kind of default for
	 * this.
	 */
	if (!svcl) {
		printf("Svcl is required\n");
		show_help(argv[0]);
		exit(1);
	}

	DFUSE_TRA_ROOT(dfuse_info, "dfuse_info");

	if (!dfuse_info->di_foreground) {
		rc = daemon(0, 0);
		if (rc)
			return daos_errno2der(rc);
	}

	rc = daos_init();
	if (rc != -DER_SUCCESS)
		D_GOTO(out, ret = rc);

	dfuse_info->di_svcl = daos_rank_list_parse(svcl, ":");
	if (dfuse_info->di_svcl == NULL) {
		DFUSE_LOG_ERROR("Invalid pool service rank list");
		D_GOTO(out_dfuse, ret = -DER_INVAL);
	}

	D_ALLOC_PTR(dfs);
	if (!dfs) {
		D_GOTO(out_svcl, 0);
	}

	d_list_add(&dfs->dfs_list, &dfuse_info->di_dfs_list);

	if (dfuse_info->di_pool) {
		if (uuid_parse(dfuse_info->di_pool, dfs->dfs_pool) < 0) {
			DFUSE_LOG_ERROR("Invalid pool uuid");
			D_GOTO(out_dfs, ret = -DER_INVAL);
		}

		/** Connect to DAOS pool */
		rc = daos_pool_connect(dfs->dfs_pool, dfuse_info->di_group,
				       dfuse_info->di_svcl, DAOS_PC_RW,
				       &dfs->dfs_poh, &dfs->dfs_pool_info,
				       NULL);
		if (rc != -DER_SUCCESS) {
			DFUSE_LOG_ERROR("Failed to connect to pool (%d)", rc);
			D_GOTO(out_dfs, 0);
		}

		if (dfuse_info->di_cont) {

			if (uuid_parse(dfuse_info->di_cont, dfs->dfs_cont) < 0) {
				DFUSE_LOG_ERROR("Invalid container uuid");
				D_GOTO(out_pool, ret = -DER_INVAL);
			}

			/** Try to open the DAOS container (the mountpoint) */
			rc = daos_cont_open(dfs->dfs_poh, dfs->dfs_cont,
					    DAOS_COO_RW, &dfs->dfs_coh,
					    &dfs->dfs_co_info, NULL);
			if (rc) {
				DFUSE_LOG_ERROR("Failed container open (%d)",
						rc);
				D_GOTO(out_pool, 0);
			}

			rc = dfs_mount(dfs->dfs_poh, dfs->dfs_coh, O_RDWR,
				       &dfs->dfs_ns);
			if (rc) {
				daos_cont_close(dfs->dfs_coh, NULL);
				DFUSE_LOG_ERROR("dfs_mount failed (%d)", rc);
				D_GOTO(out_pool, 0);
			}
			dfs->dfs_ops = &dfuse_dfs_ops;
		} else {
			dfs->dfs_ops = &dfuse_cont_ops;
		}
	} else {
		dfs->dfs_ops = &dfuse_pool_ops;
	}

	rc = dfuse_start(dfuse_info, dfs);
	if (rc != -DER_SUCCESS)
		D_GOTO(out_cont, ret = rc);

	ret = dfuse_destroy_fuse(dfuse_info->di_handle);

	fuse_session_destroy(dfuse_info->di_session);

	D_GOTO(out_dfs, 0);

out_cont:
	if (dfuse_info->di_cont) {
		dfs_umount(dfs->dfs_ns);
		daos_cont_close(dfs->dfs_coh, NULL);
	}
out_pool:
	if (dfuse_info->di_pool)
		daos_pool_disconnect(dfs->dfs_poh, NULL);
out_dfs:
	while ((dfs = d_list_pop_entry(&dfuse_info->di_dfs_list,
				       struct dfuse_dfs, dfs_list)))
	{
		/* Try and close/disconnect all container/pool handles and free
		 * the dfs struct.
		 *
		 * dfuse_destroy_fuse() has already been called here which will
		 * have iterated the inode table and should have dropped all
		 * references to the dfs entries, however depending on the order
		 * some pools/containers may be left open here so check for this
		 * and try them again.
		 */
		if (!daos_handle_is_inval(dfs->dfs_coh)) {

			if (dfs->dfs_ns) {
				rc = dfs_umount(dfs->dfs_ns);
				if (rc != 0)
					DFUSE_TRA_ERROR(dfs,
							"dfs_umount() failed (%d)",
							rc);
			}

			rc = daos_cont_close(dfs->dfs_coh, NULL);
			if (rc != -DER_SUCCESS) {
				DFUSE_TRA_ERROR(dfs,
						"daos_cont_close() failed: (%d)",
						rc);
			}

		} else if (!daos_handle_is_inval(dfs->dfs_poh)) {
			rc = daos_pool_disconnect(dfs->dfs_poh, NULL);
			if (rc != -DER_SUCCESS) {
				DFUSE_TRA_ERROR(dfs,
						"daos_pool_disconnect() failed: (%d)",
						rc);
			}
		}

		D_FREE(dfs);
	}
out_svcl:
	d_rank_list_free(dfuse_info->di_svcl);
out_dfuse:
	DFUSE_TRA_DOWN(dfuse_info);
	D_FREE(dfuse_info);
out_fini:
	daos_fini();
out:
	/* Convert CaRT error numbers to something that can be returned to the
	 * user.  This needs to be less than 256 so only works for CaRT, not
	 * DAOS error numbers.
	 */
	DFUSE_LOG_INFO("Exiting with status %d", ret);
	if (ret)
		return -(ret + DER_ERR_GURT_BASE);
	else
		return 0;
}
