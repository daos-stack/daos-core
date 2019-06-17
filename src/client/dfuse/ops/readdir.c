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

#include "dfuse_common.h"
#include "dfuse.h"

#define LOOP_COUNT 512

struct iterate_data {
	fuse_req_t			req;
	struct dfuse_inode_entry	*inode;
	struct dfuse_obj_hdl		*oh;
	size_t				size;
	size_t				fuse_size;
	size_t				b_off;
	uint8_t				stop;
};

int
filler_cb(dfs_t *dfs, dfs_obj_t *dir, const char name[], void *_udata)
{
	struct iterate_data	*udata = (struct iterate_data *)_udata;
	struct dfuse_projection_info *fs_handle = fuse_req_userdata(udata->req);
	struct dfuse_obj_hdl	*oh = udata->oh;
	dfs_obj_t		*obj;
	daos_obj_id_t		oid;
	struct stat		stbuf = {0};
	int			ns = 0;
	int			rc;

	/*
	 * MSC - from fuse fuse_add_direntry: "From the 'stbuf' argument the
	 * st_ino field and bits 12-15 of the st_mode field are used. The other
	 * fields are ignored." So we only need to lookup the entry for the
	 * mode.
	 */
	DFUSE_TRA_DEBUG(udata->inode, "Adding entry name '%s'", name);

	rc = dfs_lookup_rel(dfs, dir, name, O_RDONLY, &obj, &stbuf.st_mode);
	if (rc)
		return rc;

	rc = dfs_obj2id(obj, &oid);
	if (rc)
		D_GOTO(out, rc);

	rc = dfuse_lookup_inode(fs_handle, udata->inode->ie_dfs, &oid,
				&stbuf.st_ino);
	if (rc)
		D_GOTO(out, rc);

	/** if we have not run out of space on the fuse buffer */
	if (oh->doh_cur_off == 0) {
		/** add the entry to the buffer to return to fuse first */
		ns = fuse_add_direntry(udata->req, oh->doh_buf + udata->b_off,
				       udata->fuse_size - udata->b_off, name,
				       &stbuf, (off_t)(&oh->doh_anchor));

		/*
		 * If entry does not fit within the readdir buff size, save
		 * state in oh and re-add it to the buf. this will not be
		 * returned for current readdir call.
		 */
		if (ns > udata->fuse_size - udata->b_off) {
			oh->doh_start_off[oh->doh_idx] = udata->b_off;
			oh->doh_cur_off = udata->b_off;

			ns = fuse_add_direntry(udata->req,
					       oh->doh_buf + udata->b_off,
					       udata->size - udata->b_off,
					       name, &stbuf,
					       (off_t)(&oh->doh_anchor));

			D_ASSERT(ns <= udata->size - udata->b_off);
			oh->doh_cur_off += ns;

			/** no need to issue futher dfs_iterate() calls */
			udata->stop = 1;
		} else {
			/** entry did fit, just increment the buf offset */
			udata->b_off += ns;
		}
	} else {
insert:
		/** fuse size exceeded so add entry & update the oh counters */
		ns = fuse_add_direntry(udata->req,
				       oh->doh_buf + oh->doh_cur_off,
				       udata->size - oh->doh_cur_off, name,
				       &stbuf, (off_t)(&oh->doh_anchor));

		/*
		 * If entry does not fit, realloc to fit the entries that were
		 * already enumerated.
		 */
		if (ns > udata->size - oh->doh_cur_off) {
			udata->size = udata->size * 2;
			oh->doh_buf = realloc(oh->doh_buf, udata->size);
			if (oh->doh_buf == NULL)
				D_GOTO(out, rc = -ENOMEM);
			goto insert;
		}

		oh->doh_cur_off += ns;
		/*
		 * we need to keep track of offsets where we exceed 4k in
		 * entries for further calls to readdir
		 */
		if (oh->doh_cur_off - oh->doh_start_off[oh->doh_idx] >
		    udata->fuse_size) {
			oh->doh_idx++;
			oh->doh_start_off[oh->doh_idx] = oh->doh_cur_off - ns;
		}
	}

out:
	dfs_release(obj);
	/* we return the negative errno back to DFS */
	return rc;
}

void
dfuse_cb_readdir(fuse_req_t req, struct dfuse_inode_entry *inode,
		 size_t size, off_t offset, struct fuse_file_info *fi)
{
	struct dfuse_obj_hdl	*oh = (struct dfuse_obj_hdl *)fi->fh;
	uint32_t		nr = LOOP_COUNT;
	size_t			buf_size;
	struct iterate_data	udata;
	int			rc;

	if (offset < 0)
		D_GOTO(err, rc = EINVAL);

	D_ASSERT(oh);

	if (offset == 0) {
		memset(&oh->doh_anchor, 0, sizeof(oh->doh_anchor));
	} else {
		if (offset != (off_t)&oh->doh_anchor)
			D_GOTO(err, rc = EIO);
	}

	/** if there was anything left in the old buffer */
	if (offset && oh->doh_cur_off) {
		/** if remaining does not fit in the fuse buf, return a block */
		if (size < oh->doh_cur_off - oh->doh_start_off[oh->doh_idx]) {
			fuse_reply_buf(req, oh->doh_buf +
				       oh->doh_start_off[oh->doh_idx],
				       oh->doh_start_off[oh->doh_idx + 1] -
				       oh->doh_start_off[oh->doh_idx]);
			oh->doh_idx++;
		} else {
			fuse_reply_buf(req, oh->doh_buf +
				       oh->doh_start_off[oh->doh_idx],
				       oh->doh_cur_off -
				       oh->doh_start_off[oh->doh_idx]);
			oh->doh_cur_off = 0;
			oh->doh_idx = 0;
		}
		return;
	}

	/*
	 * the DFS size should be less than what we want in fuse to account for
	 * the fuse metadata for each entry, so just use 1/2 for now
	 */
	buf_size = size * READDIR_BLOCKS / 2;

	/** buffer will be freed when this oh is closed */
	D_ALLOC(oh->doh_buf, buf_size);
	if (!oh->doh_buf)
		D_GOTO(err, rc = ENOMEM);

	udata.req = req;
	udata.size = buf_size;
	udata.fuse_size = size;
	udata.b_off = 0;
	udata.inode = inode;
	udata.oh = oh;
	udata.stop = 0;

	while (!daos_anchor_is_eof(&oh->doh_anchor)) {
		/** should not be here if extra buf space is used */
		D_ASSERT(oh->doh_cur_off == 0);

		rc = dfs_iterate(oh->doh_dfs, oh->doh_obj, &oh->doh_anchor, &nr,
				 buf_size - udata.b_off, filler_cb, &udata);

		/** if entry does not fit in buffer, just return */
		if (rc == -E2BIG)
			break;
		/** otherwise a different error occured */
		if (rc)
			D_GOTO(err, rc = -rc);

		/** if buffer is full, break enumeration */
		if (udata.stop)
			break;
	}

	oh->doh_idx = 0;
	fuse_reply_buf(req, oh->doh_buf, udata.b_off);
	return;

err:
	DFUSE_FUSE_REPLY_ERR(req, rc);
}
