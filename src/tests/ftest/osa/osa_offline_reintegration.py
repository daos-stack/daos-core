#!/usr/bin/python
"""
  (C) Copyright 2020 Intel Corporation.

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
import time
import random
import ctypes
from apricot import TestWithServers
from test_utils_pool import TestPool
from pydaos.raw import (DaosContainer, IORequest,
                        DaosObj)


class OSAOfflineReintegration(TestWithServers):
    # pylint: disable=too-many-ancestors
    """
    Test Class Description: This test runs
    daos_server offline reintegration test cases.

    :avocado: recursive
    """
    def setUp(self):
        """Set up for test case."""
        super(OSAOfflineReintegration, self).setUp()
        self.dmg_command = self.get_dmg_command()
        self.no_of_dkeys = self.params.get("no_of_dkeys", '/run/dkeys/*')[0]
        self.no_of_akeys = self.params.get("no_of_akeys", '/run/akeys/*')[0]
        self.record_length = self.params.get("length", '/run/record/*')

    def get_pool_leader(self):
        """Get the pool leader
           Args: None
           Returns : pool_leader (int)
        """
        out = []
        kwargs = {"pool_uuid": self.pool.uuid}
        out = self.dmg_command.get_output("pool_query", **kwargs)
        pool_leader = out[0][3]
        pool_leader = pool_leader.split('=')
        return int(pool_leader[1])

    def get_pool_version(self):
        """Get the pool version
           Args: None
           Returns : pool_version (int)
        """
        out = []
        kwargs = {"pool_uuid": self.pool.uuid}
        out = self.dmg_command.get_output("pool_query", **kwargs)
        pool_version = out[0][4]
        pool_version = pool_version.split('=')
        return int(pool_version[1])

    def write_single_object(self):
        self.pool.connect(2)
        csum = self.params.get("enable_checksum", '/run/container/*')
        container = DaosContainer(self.context)
        input_param = container.cont_input_values
        input_param.enable_chksum = csum
        container.create(poh=self.pool.pool.handle,
                         con_prop=input_param)
        container.open()
        obj = DaosObj(self.context, container)
        obj.create(objcls=1)
        obj.open()
        ioreq = IORequest(self.context,
                          container,
                          obj, objtype=4)
        self.log.info("Writing the Single Dataset")
        record_index = 0
        for dkey in range(self.no_of_dkeys):
            for akey in range(self.no_of_akeys):
                indata = ("{0}".format(str(akey)[0])
                          * self.record_length[record_index])
                d_key_value = "dkey {0}".format(dkey)
                c_dkey = ctypes.create_string_buffer(d_key_value)
                a_key_value = "akey {0}".format(akey)
                c_akey = ctypes.create_string_buffer(a_key_value)
                c_value = ctypes.create_string_buffer(indata)
                c_size = ctypes.c_size_t(ctypes.sizeof(c_value))
                ioreq.single_insert(c_dkey, c_akey, c_value, c_size)
                record_index = record_index + 1
                if record_index == len(self.record_length):
                    record_index = 0

    def run_offline_reintegration_test(self, num_pool, data=False):
        """Run the offline reintegration without data.
           Args: num_pool (int) : Total pools to create
            for testing purpose.
           Returns : None
        """
        # Create a pool
        pool = {}
        pool_uuid = []
        target_list = []
        exclude_servers = len(self.hostlist_servers) - 1
        nvme_size = 54000000000
        scm_size = 6000000000

        # Exclude target : random two targets
        n = random.randint(1, 7)
        target_list.append(n)
        target_list.append(n+1)
        t_string = "{},{}".format(target_list[0], target_list[1])

        # Exclude rank : two ranks other than rank 0.
        n = random.randint(1, exclude_servers)
        rank = n

        for val in range(0, num_pool):
            pool[val] = TestPool(self.context,
                                 dmg_command=self.get_dmg_command())
            pool[val].get_params(self)
            # Split total SCM and NVME size for creating multiple pools.
            pool[val].scm_size.value = int(scm_size / num_pool)
            pool[val].nvme_size.value = int(nvme_size / num_pool)
            pool[val].create()
            pool_uuid.append(pool[val].uuid)
            self.pool = pool[val]
            if data is True:
                self.write_single_object()

        # Exclude and reintegrate the pool_uuid, rank and targets
        for val in range(0, num_pool):
            self.pool = pool[val]
            self.pool.display_pool_daos_space("Pool space: Beginning")
            pver_begin = self.get_pool_version()
            self.log.info("Pool Version at the beginning %d", pver_begin)
            output = self.dmg_command.pool_exclude(self.pool.uuid,
                                                   rank, t_string)
            self.log.info(output)
            time.sleep(10)
            pver_exclude = self.get_pool_version()
            self.log.info("Pool Version after exclude %d", pver_exclude)
            if pver_exclude <= pver_begin:
                self.fail("Pool Version Error:  After exclude")
            output = self.dmg_command.pool_reintegrate(self.pool.uuid,
                                                       rank,
                                                       t_string)
            self.log.info(output)
            time.sleep(10)
            pver_reint = self.get_pool_version()
            self.log.info("Pool Version after reintegrate %d", pver_reint)
            if pver_reint <= pver_exclude:
                self.fail("Pool Version Error: After reintegrate")

        for val in range(0, num_pool):
            display_string = "Pool{} space at the End".format(val)
            self.pool = pool[val]
            self.pool.display_pool_daos_space(display_string)
            pool[val].destroy()

    def test_osa_offline_reintegration(self):
        """Test ID: DAOS-4749
        Test Description: Validate Offline Reintegration

        :avocado: tags=all,pr,hw,large,osa,offline_reintegration
        """
        for x in range(1, 4):
            self.run_offline_reintegration_test(x)
        self.run_offline_reintegration_test(1, True)
