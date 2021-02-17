#!/usr/bin/python
"""
  (C) Copyright 2020-2021 Intel Corporation.

  SPDX-License-Identifier: BSD-2-Clause-Patent
"""
import time
from osa_utils import OSAUtils
from test_utils_pool import TestPool


class OSADmgNegativeTest(OSAUtils):
    # pylint: disable=too-many-ancestors
    """
    Test Class Description: This test runs
    passes invalid parameters to dmg OSA
    related commands and check whether
    appropriate errors are returned.

    :avocado: recursive
    """
    def setUp(self):
        """Set up for test case."""
        super(OSADmgNegativeTest, self).setUp()
        self.dmg_command = self.get_dmg_command()
        # Start an additional server.
        self.extra_servers = self.params.get("test_servers",
                                             "/run/extra_servers/*")
        # Dmg command test seeuence
        self.test_seq = self.params.get("dmg_cmd_test",
                                        "/run/dmg_cmds/*")

    def run_osa_dmg_test(self, num_pool):
        """Run the offline extend without data.

        Args:
            num_pool (int) : total pools to create for testing purposes.
        """
        # Create a pool
        self.dmg_command.exit_status_exception = False
        pool = {}
        pool_uuid = []

        for val in range(0, num_pool):
            pool[val] = TestPool(self.context, dmg_command=self.dmg_command)
            pool[val].get_params(self)
            # Split total SCM and NVME size for creating multiple pools.
            pool[val].scm_size.value = int(pool[val].scm_size.value /
                                           num_pool)
            pool[val].nvme_size.value = int(pool[val].nvme_size.value /
                                            num_pool)
            pool[val].create()
            pool_uuid.append(pool[val].uuid)
            self.pool = pool[val]

        # Start the additional servers and extend the pool
        self.log.info("Extra Servers = %s", self.extra_servers)
        self.start_additional_servers(self.extra_servers)
        # Give sometime for the additional server to come up.
        time.sleep(5)

        # Get rank, target from the test_dmg_seeuence
        for val in range(0, num_pool):
            for i in range(len(self.test_seq)):
                self.pool = pool[val]
                scm_size = self.pool.scm_size
                nvme_size = self.pool.nvme_size
                rank = self.test_seq[i][0]
                target = "{}".format(self.test_seq[i][1])
                output = self.dmg_command.pool_extend(self.pool.uuid,
                                                      rank,
                                                      scm_size,
                                                      nvme_size)
                self.log.info(output)
                # Exclude and reintegrate commands presently
                # don't take a list.
                # First perform OSA dmg exclude
                output = self.dmg_command.pool_exclude(self.pool.uuid,
                                                       rank,
                                                       target)
                self.log.info(output)
                time.sleep(10)
                # Now reintegrate the excluded rank.
                output = self.dmg_command.pool_reintegrate(self.pool.uuid,
                                                           rank,
                                                           target)
                self.log.info(output)

    def test_osa_dmg_cmd(self):
        """
        JIRA ID: DAOS-5866

        Test Description: Test

        :avocado: tags=all,daily_regression,hw,medium,ib2
        :avocado: tags=osa,osa_dmg_negative_test,dmg_negative_test
        """
        # Perform extend testing with 1 pool
        self.run_osa_dmg_test(1)
