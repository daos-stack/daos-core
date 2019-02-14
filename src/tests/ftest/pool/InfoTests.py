#!/usr/bin/python
"""
  (C) Copyright 2018-2019 Intel Corporation.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

  GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
  The Government's rights to use, modify, reproduce, release, perform, display,
  or disclose this software are subject to the terms of the Apache License as
  provided in Contract No. B609815.
  Any reproduction of computer software, computer software documentation, or
  portions thereof marked with this legend must also reproduce the markings.
"""
import os
import sys
import json
import logging
from avocado import Test, main

sys.path.append('./util')
sys.path.append('../util')
sys.path.append('./../../utils/py')
sys.path.append('../../../utils/py')

import AgentUtils
import ServerUtils
import WriteHostFile
from daos_api import DaosContext, DaosPool, DaosLog, DaosApiError
from conversion import c_uuid_to_str

class InfoTests(Test):
    """
    Tests DAOS pool query.
    """
    def setUp(self):
        # get paths from the build_vars generated by build
        with open(os.path.join(os.path.dirname(os.path.realpath(__file__)),
                               "../../../../.build_vars.json")) as f:
            build_paths = json.load(f)
        self.basepath = os.path.normpath(build_paths['PREFIX'] + "/../")
        self.tmp = build_paths['PREFIX'] + '/tmp'
        self.server_group = self.params.get("server_group", '/server/',
                                            'daos_server')

        context = DaosContext(build_paths['PREFIX'] + '/lib/')

        self.pool = DaosPool(context)
        self.d_log = DaosLog(context)
        self.hostlist = self.params.get("test_machines1", '/run/hosts/')
        self.hostfile = WriteHostFile.WriteHostFile(self.hostlist, self.tmp)

        AgentUtils.run_agent(self.basepath, self.hostlist)
        ServerUtils.runServer(self.hostfile, self.server_group, self.basepath)

    def tearDown(self):
        # shut 'er down
        try:
            if self.pool:
                self.pool.destroy(1)
            os.remove(self.hostfile)
        finally:
            AgentUtils.stop_agent(self.hostlist)
            ServerUtils.stopServer(hosts=self.hostlist)

    def test_simple_query(self):
        """
        Test querying a pool created on a single server.

        :avocado: tags=pool,poolquery,infotest
        """
        # create pool
        mode = self.params.get("mode", '/run/testparams/modes/*', 0731)
        if mode == 73:
             self.cancel('Cancel the mode test 73 because of DAOS-1877')

        uid = os.geteuid()
        gid = os.getegid()
        size = self.params.get("size", '/run/testparams/sizes/*', 0)
        group = self.server_group

        self.pool.create(mode, uid, gid, size, group, None)

        # connect to the pool
        flags = self.params.get("perms", '/run/testparams/connectperms/*', '')
        connect_flags = 1 << flags
        self.pool.connect(connect_flags)

        # query the pool
        pool_info = self.pool.pool_query()

        # check uuid
        uuid_str = c_uuid_to_str(pool_info.pi_uuid)
        if uuid_str != self.pool.get_uuid_str():
            self.d_log.error("UUID str does not match expected string")
            self.fail("UUID str does not match expected string")

        '''
        # validate size of pool is what we expect
        This check is currently disabled, as space is not implemented in
        DAOS C API yet.
        if size != pool_info.pi_space:
            self.d_log.error("expected size {0} did not match actual size {1}"
                      .format(size, pool_info.pi_space))
            self.fail("expected size {0} did not match actual size {1}"
                      .format(size, pool_info.pi_space))
        '''

        # number of targets
        if pool_info.pi_ntargets != len(self.hostlist):
            self.d_log.error("found number of targets in pool did not match "
                      	     "expected number, 1. num targets: {0}"
                      	     .format(pool_info.pi_ntargets))
            self.fail("found number of targets in pool did not match "
                      "expected number, 1. num targets: {0}"
                      .format(pool_info.pi_ntargets))

        # number of disabled targets
        if pool_info.pi_ndisabled > 0:
            self.d_log.error("found disabled targets, none expected to be")
            self.fail("found disabled targets, none expected to be disabled")

        # mode
        if pool_info.pi_mode != mode:
            self.d_log.error("found different mode than expected. expected {0}, "
                      	     "found {1}.".format(mode, pool_info.pi_mode))
            self.fail("found different mode than expected. expected {0}, "
                      "found {1}.".format(mode, pool_info.pi_mode))

        # uid
        if pool_info.pi_uid != uid:
            self.d_log.error("found actual pool uid {0} does not match expected "
                      	     "uid {1}".format(pool_info.pi_uid, uid))
            self.fail("found actual pool uid {0} does not match expected uid "
                      "{1}".format(pool_info.pi_uid, uid))

        # gid
        if pool_info.pi_gid != gid:
            self.d_log.error("found actual pool gid {0} does not match expected "
                      	     "gid {1}".format(pool_info.pi_gid, gid))
            self.fail("found actual pool gid {0} does not match expected gid "
                      "{1}".format(pool_info.pi_gid, gid))