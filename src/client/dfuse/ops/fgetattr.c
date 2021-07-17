/**
 * (C) Copyright 2016-2021 Intel Corporation.
 *
 * SPDX-License-Identifier: BSD-2-Clause-Patent
 */

#include "dfuse_common.h"
#include "dfuse.h"

void
dfuse_cb_getattr(fuse_req_t req, struct dfuse_inode_entry *ie)
{
	struct stat	stat = {};
	int		rc;

	rc = dfs_ostat(ie->ie_dfs->dfs_ns, ie->ie_obj, &stat);
	if (rc != 0)
		D_GOTO(err, rc);

	if (ie->ie_dfs->dfs_multi_user) {
		rc = dfuse_get_uid(ie);
		if (rc)
			D_GOTO(err, rc);
		stat.st_uid = ie->ie_stat.st_uid;
		stat.st_gid = ie->ie_stat.st_gid;
	}

	/* Copy the inode number from the inode struct, to avoid having to
	 * recompute it each time.
	 */

	stat.st_ino = ie->ie_stat.st_ino;
	DFUSE_REPLY_ATTR(ie, req, &stat);

	return;
err:
	DFUSE_REPLY_ERR_RAW(ie, req, rc);
}
