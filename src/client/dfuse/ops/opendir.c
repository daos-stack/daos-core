/**
 * (C) Copyright 2016-2021 Intel Corporation.
 *
 * SPDX-License-Identifier: BSD-2-Clause-Patent
 */

#include "dfuse_common.h"
#include "dfuse.h"

void
dfuse_cb_opendir(fuse_req_t req, struct dfuse_inode_entry *ie,
		 struct fuse_file_info *fi)
{
	struct dfuse_obj_hdl	*oh = NULL;
	struct fuse_file_info	fi_out = {0};
	int			rc;

	D_ALLOC_PTR(oh);
	if (!oh)
		D_GOTO(err, rc = ENOMEM);

	DFUSE_TRA_UP(oh, ie, "open handle");

	/** duplicate the file handle for the fuse handle */
	rc = dfs_dup(ie->ie_dfs->dfs_ns, ie->ie_obj, fi->flags,
		     &oh->doh_obj);
	if (rc)
		D_GOTO(err, rc);

	oh->doh_dfs = ie->ie_dfs->dfs_ns;
	oh->doh_ie = ie;

	fi_out.fh = (uint64_t)oh;

#if HAVE_CACHE_READDIR
	if (ie->ie_dfs->dfc_dentry_timeout > 0)
		fi_out.cache_readdir = 1;
#endif

#if 0
	/* These are fuse version dependent so need compile-time checking */
	if (ie->ie_dfs->dfs_attr_timeout)
		fi->flags = FOPEN_KEEP_CACHE | FOPEN_CACHE_DIR;
#endif

	DFUSE_REPLY_OPEN(oh, req, &fi_out);
	return;
err:
	D_FREE(oh);
	DFUSE_REPLY_ERR_RAW(ie, req, rc);
}

void
dfuse_cb_releasedir(fuse_req_t req, struct dfuse_inode_entry *ino,
		    struct fuse_file_info *fi)
{
	struct dfuse_obj_hdl	*oh = (struct dfuse_obj_hdl *)fi->fh;
	int			rc;

	rc = dfs_release(oh->doh_obj);
	if (rc == 0)
		DFUSE_REPLY_ZERO(oh, req);
	else
		DFUSE_REPLY_ERR_RAW(oh, req, rc);
	D_FREE(oh->doh_dre);
	D_FREE(oh);
};
