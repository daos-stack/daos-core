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

#ifndef __DAOS_DTX_H__
#define __DAOS_DTX_H__

#include <time.h>
#include <uuid/uuid.h>

/**
 * DAOS two-phase commit transaction identifier,
 * generated by client, globally unique.
 */
struct daos_tx_id {
	/** The uuid of the transaction */
	uuid_t			dti_uuid;
	/** The timestamp in second for the transaction */
	uint64_t		dti_sec;
};

static inline void
daos_generate_dti(struct daos_tx_id *dti)
{
	uuid_generate(dti->dti_uuid);
	dti->dti_sec = time(NULL);
}

static inline void
daos_copy_dti(struct daos_tx_id *des, struct daos_tx_id *src)
{
	if (src != NULL)
		*des = *src;
	else
		memset(des, 0, sizeof(*des));
}

#define DF_DTI		DF_UUID
#define DP_DTI(dti)	DP_UUID((dti)->dti_uuid)

enum daos_ops_intent {
	DAOS_INTENT_DEFAULT	= 0,	/* fetch/enumerate/query */
	DAOS_INTENT_PURGE	= 1,	/* purge/aggregation */
	DAOS_INTENT_UPDATE	= 2,	/* write/insert */
	DAOS_INTENT_PUNCH	= 3,	/* punch/delete */
	DAOS_INTENT_REBUILD	= 4,	/* for rebuild related scan */
	DAOS_INTENT_CHECK	= 5,	/* check aborted or not */
};

enum daos_dtx_vbt {
	/* invisible case */
	DTX_VBT_INVISIBLE		= 0,
	/* visible, no (or not care) pending modification */
	DTX_VBT_VISIBLE_CLEAN		= 1,
	/* visible but with dirty modification or garbage */
	DTX_VBT_VISIBLE_DIRTY		= 2,
};

#endif /* __DAOS_DTX_H__ */
