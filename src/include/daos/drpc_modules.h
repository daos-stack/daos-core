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
 * provided in Contract No. 8F-30005.
 * Any reproduction of computer software, computer software documentation, or
 * portions thereof marked with this legend must also reproduce the markings.
 */

#ifndef __DAOS_DRPC_MODULES_H__
#define __DAOS_DRPC_MODULES_H__

/**
 * DAOS dRPC Modules
 *
 * dRPC modules are used to multiplex communications over the Unix Domain Socket
 * to appropriate handlers. They are populated in the Drpc__Call structure.
 *
 * dRPC module IDs must be unique. This is a list of all DAOS dRPC modules.
 */

enum drpc_module {
	DRPC_MODULE_TEST		= 0,	/* Reserved for testing */
	DRPC_MODULE_SECURITY_AGENT	= 1,
	DRPC_MODULE_SRV			= 2,	/* daos_server */

	NUM_DRPC_MODULES			/* Must be last */
};

/** Methods of DRPC_MODULE_SRV */
enum drpc_srv_method {
	DRPC_SRV_NOTIFY_READY = 0
};

#endif /* __DAOS_DRPC_MODULES_H__ */
