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
import threading

from itertools import product
from avocado import fail_on
from test_utils_pool import TestPool
from write_host_file import write_host_file
from command_utils import CommandFailure
from osa_utils import OSAUtils

try:
    # python 3.x
    import queue as queue
except ImportError:
    # python 2.7
    import Queue as queue


class NvmePoolExtend(OSAUtils):
    # pylint: disable=too-many-ancestors
    """
    Test Class Description: This test runs
    NVME Pool Extend test cases.
    - Start the few daos servers.
    - Create a pool
    - Run IOR with write mode
    - Start a new server and extend the pool
    - Verify IOR written data after extending the pool.

    :avocado: recursive
    """
    def setUp(self):
        """Set up for test case."""
        super(NvmePoolExtend, self).setUp()
        self.dmg_command = self.get_dmg_command()
        self.ior_w_flags = self.params.get("write_flags", '/run/ior/iorflags/*')
        self.ior_r_flags = self.params.get("read_flags", '/run/ior/iorflags/*')
        self.ior_apis = self.params.get("ior_api", '/run/ior/iorflags/*')
        self.ior_test_sequence = self.params.get("ior_test_sequence",
                                                 '/run/ior/iorflags/*')
        self.ior_daos_oclass = self.params.get("obj_class",
                                               '/run/ior/iorflags/*')
        self.ior_dfs_oclass = self.params.get(
            "obj_class", '/run/ior/iorflags/*')
        # Start an additional server.
        self.extra_servers = self.params.get("test_servers",
                                             "/run/extra_servers/*")
        # Recreate the client hostfile without slots defined
        self.hostfile_clients = write_host_file(
            self.hostlist_clients, self.workdir, None)
        self.pool = None
        self.out_queue = queue.Queue()
        self.container_info = {}

    @fail_on(CommandFailure)
    def get_pool_version(self):
        """Get the pool version.

        Returns:
            int: pool_version_value

        """
        data = self.dmg_command.pool_query(self.pool.uuid)
        return int(data["version"])

    def run_ior_thread(self, action, oclass, api, test):
        """Start the IOR thread for either writing or
        reading data to/from a container.
        Args:
            action (str): Start the IOR thread with Read or
                          Write
            oclass (str): IOR object class
            API (str): IOR API
            test (list): IOR test sequence
            flags (str): IOR flags
        """
        if action == "Write":
            flags = self.ior_w_flags
        else:
            flags = self.ior_r_flags

        # Add a thread for these IOR arguments
        process = threading.Thread(target=self.ior_thread,
                                   kwargs={"pool": self.pool,
                                           "oclass": oclass,
                                           "api": api,
                                           "test": test,
                                           "flags": flags,
                                           "results":
                                           self.out_queue})
        # Launch the IOR thread
        process.start()
        # Wait for the thread to finish
        process.join()

    def run_nvme_pool_extend(self, num_pool):
        """Run Pool Extend
        Args:
            int : total pools to create for testing purposes.
        """
        pool = {}
        total_servers = len(self.hostlist_servers) * 2
        self.log.info("Total Daos Servers (Initial): %d", total_servers)

        for val in range(0, num_pool):
            # Create a pool
            pool[val] = TestPool(self.context, dmg_command=self.dmg_command)
            pool[val].get_params(self)
            # Split total SCM and NVME size for creating multiple pools.
            pool[val].scm_size.value = int(pool[val].scm_size.value /
                                           num_pool)
            pool[val].nvme_size.value = int(pool[val].nvme_size.value /
                                            num_pool)
            pool[val].create()

        for val in range(0, num_pool):
            self.pool = pool[val]
            for oclass, api, test in product(self.ior_dfs_oclass,
                                             self.ior_apis,
                                             self.ior_test_sequence):
                self.run_ior_thread("Write", oclass, api, test)

                scm_pool_size = self.pool.scm_size
                nvme_pool_size = self.pool.nvme_size
                self.pool.display_pool_daos_space("Pool space: Beginning")
                pver_begin = self.get_pool_version()

                # Start the additional servers and extend the pool
                self.log.info("Extra Servers = %s", self.extra_servers)
                self.start_additional_servers(self.extra_servers)
                # Give some time for the additional server to come up.
                # Extending two ranks (two servers per node)
                time.sleep(25)
                self.log.info("Pool Version at the beginning %s", pver_begin)
                output = self.dmg_command.pool_extend(self.pool.uuid,
                                                      "6,7", scm_pool_size,
                                                      nvme_pool_size)
                self.log.info(output)
                self.is_rebuild_done(3)
                self.assert_on_rebuild_failure()

                pver_extend = self.get_pool_version()
                self.log.info("Pool Version after extend %s", pver_extend)
                # Check pool version incremented after pool extend
                self.assertTrue(pver_extend > pver_begin,
                                "Pool Version Error:  After extend")

                # Verify the data after pool extend
                self.run_ior_thread("Read", oclass, api, test)
                # Get the pool space at the end of the test
                display_string = "Pool{} space at the End".format(val)
                self.pool = pool[val]
                self.pool.display_pool_daos_space(display_string)

                # Stop the extra node servers (rank 6 and 7)
                output = self.dmg_command.system_stop(self.pool.uuid,
                                                      "6,7")
                self.log.info(output)
                self.is_rebuild_done(3)
                self.assert_on_rebuild_failure()
        pool[val].destroy()

    def test_nvme_pool_extend(self):
        """Test ID: DAOS-2086
        Test Description: NVME Pool Extend

        :avocado: tags=all,full_regression,hw,large,nvme
        :avocado: tags=nvme_pool_extend
        """
        self.run_nvme_pool_extend(1)
