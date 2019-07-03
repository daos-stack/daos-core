#!/usr/bin/python
'''
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
  The Governments rights to use, modify, reproduce, release, perform, display,
  or disclose this software are subject to the terms of the Apache License as
  provided in Contract No. B609815.
  Any reproduction of computer software, computer software documentation, or
  portions thereof marked with this legend must also reproduce the markings.
'''
from __future__ import print_function

import os
import time
import traceback

from apricot import TestWithServers

import agent_utils
import server_utils
import write_host_file
import check_for_pool
import ior_utils
from mpio_utils import MpioUtils, MpioFailed
import io_utilities
from daos_api import DaosPool, DaosServer, DaosContainer, DaosApiError

class RebuildWithIOR(TestWithServers):
    """
    This class contains tests for pool rebuild that feature I/O going on
    during the rebuild using IOR.
    :avocado: recursive
    """

    def setUp(self):
        super(RebuildWithIOR, self).setUp()

        self.hostfile_clients = (
          write_host_file.write_host_file(self.hostlist_clients, self.workdir,
                                          None))
    
#    def tearDown(self):
#        try:
#            if self.pool is not None:
#                self.pool.destroy(1)
#        finally:
#            super(RebuildWithIOR, self).tearDown()

    def test_rebuild_with_ior(self):
        """
        Jira ID: DAOS-951

        Test Description: Trigger a rebuild while I/O is ongoing.
                          I/O performed using IOR.

        Use Cases:
          -- single pool, single client performing continous read/write/verify
             sequence while failure/rebuild is triggered in another process

        :avocado: tags=pool,rebuild,rebuildwithior
        """

        try:
            # initialize MpioUtils
            self.mpio = MpioUtils()
            if self.mpio.mpich_installed(self.hostlist_clients) is False:
                self.fail("Exiting Test: Mpich not installed")

            # use the uid/gid of the user running the test, these should
            # be perfectly valid
            createuid = os.geteuid()
            creategid = os.getegid()

            # parameters used in pool create that are in yaml
            createmode = self.params.get("mode", '/run/testparams/createmode/')
            createsetid = self.params.get("setname",
                                          '/run/testparams/createset/')
            createsize = self.params.get("size", '/run/testparams/createsize/')
            createsvc = self.params.get("svcn", '/run/testparams/createsvc/')

            # ior parameters
            client_processes = self.params.get("np", '/run/ior/client_processes/*/')
            iteration = self.params.get("iter", '/run/ior/iteration/')
            iorflags_write = self.params.get("write", '/run/ior/iorflags/')
            iorflags_read = self.params.get("read", '/run/ior/iorflags/')
            transfer_size = self.params.get("t",
                                            '/run/ior/transfersize_blocksize/*/')
            block_size = self.params.get("b", '/run/ior/transfersize_blocksize/*/')
            oclass = self.params.get("oclass", '/run/ior/object_class/')

            # initialize a python pool object then create the underlying
            # daos storage
            self.pool = DaosPool(self.context)
            self.pool.create(createmode, createuid, creategid,
                        createsize, createsetid, None, None, createsvc)

            pool_uuid = self.pool.get_uuid_str()
            svc_list = ""
            for i in range(createsvc):
                svc_list += str(int(self.pool.svc.rl_ranks[i])) + ":"
            svc_list = svc_list[:-1]

            # connect to the pool
            self.pool.connect(1 << 1)

            # get pool status and make sure it all looks good before we start
            self.pool.pool_query()
            if self.pool.pool_info.pi_ndisabled != 0:
                self.fail("Number of disabled targets reporting incorrectly.\n")
            if self.pool.pool_info.pi_rebuild_st.rs_errno != 0:
                self.fail("Rebuild error but rebuild hasn't run.\n")
            if self.pool.pool_info.pi_rebuild_st.rs_done != 1:
                self.fail("Rebuild is running but device hasn't failed yet.\n")
            if self.pool.pool_info.pi_rebuild_st.rs_obj_nr != 0:
                self.fail("Rebuilt objs not zero.\n")
            if self.pool.pool_info.pi_rebuild_st.rs_rec_nr != 0:
                self.fail("Rebuilt recs not zero.\n")
            dummy_pool_version = self.pool.pool_info.pi_rebuild_st.rs_version

            self.pool.disconnect()
            # perform first set of io using IOR
            ior_utils.run_ior_mpiio(self.basepath, self.mpio.mpichinstall,
                                    pool_uuid, svc_list, client_processes,
                                    self.hostfile_clients, iorflags_write, iteration,
                                    transfer_size, block_size, True, oclass)

            self.pool.connect(1 << 1)

            # trigger the rebuild
            rank = self.params.get("rank", '/run/testparams/ranks/*')
            server = DaosServer(self.context, self.server_group, rank)
            server.kill(1)
            self.pool.exclude([rank])
            #self.pool.connect(1 << 1)

            # wait for the rebuild to finish
            while True:
                self.pool.pool_query()
                if self.pool.pool_info.pi_rebuild_st.rs_done == 1:
                    print(1)
                    break
                else:
                    time.sleep(2)

            # perform second set of io using IOR
#            ior_utils.run_ior_mpiio(self.basepath, self.mpio.mpichinstall,
#                                    pool_uuid, svc_list, client_processes,
#                                    self.hostfile_clients, iorflags_write, iteration,
#                                    transfer_size, block_size, True, oclass, "testFile2")

            # check rebuild statistics
            if self.pool.pool_info.pi_ndisabled != 8:
                self.fail("Number of disabled targets reporting incorrectly: {}"
                          .format(self.pool.pool_info.pi_ndisabled))
            if self.pool.pool_info.pi_rebuild_st.rs_errno != 0:
                self.fail("Rebuild error reported: {}".format(
                    self.pool.pool_info.pi_rebuild_st.rs_errno))
            if self.pool.pool_info.pi_rebuild_st.rs_obj_nr <= 0:
                self.fail("No objects have been rebuilt.")
            if self.pool.pool_info.pi_rebuild_st.rs_rec_nr <= 0:
                self.fail("No records have been rebuilt.")

            self.pool.disconnect()

            # perform second set of io using IOR
            ior_utils.run_ior_mpiio(self.basepath, self.mpio.mpichinstall,
                                    pool_uuid, svc_list, client_processes,
                                    self.hostfile_clients, iorflags_write, iteration,
                                    transfer_size, block_size, True, oclass, "testFile2")

            # check data intergrity using ior for both ior runs
            ior_utils.run_ior_mpiio(self.basepath, self.mpio.mpichinstall,
                                     pool_uuid, svc_list, client_processes,
                                     self.hostfile_clients, iorflags_read, iteration,
                                     transfer_size, block_size, True, oclass)
            ior_utils.run_ior_mpiio(self.basepath, self.mpio.mpichinstall,
                                     pool_uuid, svc_list, client_processes,
                                     self.hostfile_clients, iorflags_read, iteration,
                                     transfer_size, block_size, True, oclass, "testFile2")

        except (ValueError, DaosApiError, MpioFailed) as excep:
            print(excep)
            print(traceback.format_exc())
            self.fail("Expecting to pass but test has failed.\n")

