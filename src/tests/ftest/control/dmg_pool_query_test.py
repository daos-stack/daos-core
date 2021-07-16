#!/usr/bin/python
"""
  (C) Copyright 2020-2021 Intel Corporation.

  SPDX-License-Identifier: BSD-2-Clause-Patent
"""
from copy import deepcopy

from ior_test_base import IorTestBase
from control_test_base import ControlTestBase
from avocado.core.exceptions import TestFail


class DmgPoolQueryTest(ControlTestBase, IorTestBase):
    # pylint: disable=too-many-ancestors
    """Test dmg query command.

    Test Class Description:
        Simple test to verify the pool query command of dmg tool.

    :avocado: recursive
    """

    def setUp(self):
        """Set up for dmg pool query."""
        super().setUp()

        # Init the pool
        self.add_pool(connect=False)

    def test_pool_query_basic(self):
        """
        JIRA ID: DAOS-2976

        Test Description: Test basic dmg functionality to query pool info on
        the system. Provided a valid pool UUID, verify the output received from
        pool query command.

        :avocado: tags=all,daily_regression
        :avocado: tags=small,hw
        :avocado: tags=dmg,pool_query,basic
        :avocado: tags=pool_query_basic
        """
        self.log.info("==>   Verify dmg output against expected output:")
        self.pool.set_query_data()

        # We won't be testing free, min, max, and mean because the values
        # fluctuate across test runs. In addition, they're related to object
        # placement and testing them wouldn't be straightforward, so we'll need
        # some separate test cases.
        del self.pool.query_data["response"]["tier_stats"][0]["free"]
        del self.pool.query_data["response"]["tier_stats"][0]["min"]
        del self.pool.query_data["response"]["tier_stats"][0]["max"]
        del self.pool.query_data["response"]["tier_stats"][0]["mean"]
        del self.pool.query_data["response"]["tier_stats"][1]["free"]
        del self.pool.query_data["response"]["tier_stats"][1]["min"]
        del self.pool.query_data["response"]["tier_stats"][1]["max"]
        del self.pool.query_data["response"]["tier_stats"][1]["mean"]

        # Get the expected pool query values from the test yaml.  This should be
        # as simple as:
        #   exp_info = self.params.get("exp_vals", path="/run/*", default={})
        # but this yields an empty dictionary (the default), so it needs to be
        # defined manually:
        exp_info = {
            "status": self.params.get("pool_status", path="/run/exp_vals/*"),
            "uuid": self.pool.uuid.lower(),
            "total_targets": self.params.get(
                "total_targets", path="/run/exp_vals/*"),
            "active_targets": self.params.get(
                "active_targets", path="/run/exp_vals/*"),
            "total_nodes": self.params.get(
                "total_nodes", path="/run/exp_vals/*"),
            "disabled_targets": self.params.get(
                "disabled_targets", path="/run/exp_vals/*"),
            "version": self.params.get("version", path="/run/exp_vals/*"),
            "leader": self.params.get("leader", path="/run/exp_vals/*"),
            "tier_stats": [
                {
                    "total": self.params.get(
                        "total", path="/run/exp_vals/scm/*")
                },
                {
                    "total": self.params.get(
                        "total", path="/run/exp_vals/nvme/*")
                }
            ],
            "rebuild": {
                "status": self.params.get(
                    "rebuild_status", path="/run/exp_vals/rebuild/*"),
                "state": self.params.get(
                    "state", path="/run/exp_vals/rebuild/*"),
                "objects": self.params.get(
                    "objects", path="/run/exp_vals/rebuild/*"),
                "records": self.params.get(
                    "records", path="/run/exp_vals/rebuild/*")
            }
        }

        self.assertDictEqual(
            self.pool.query_data["response"], exp_info,
            "Found difference in dmg pool query output and the expected values")

        self.log.info("All expect values found in dmg pool query output.")

    def test_pool_query_inputs(self):
        """
        JIRA ID: DAOS-2976

        Test Description: Test basic dmg functionality to query pool info on
        the system. Verify the inputs that can be provided to 'query --pool'
        argument of the dmg pool subcommand.

        :avocado: tags=all,daily_regression
        :avocado: tags=small,hw
        :avocado: tags=dmg,pool_query,basic
        :avocado: tags=pool_query_inputs
        """
        # Get test UUIDs
        errors_list = []
        uuids = self.params.get("uuids", '/run/pool_uuids/*')

        # Add a pass case to verify test is working
        uuids.append([self.pool.uuid, "PASS"])

        for uuid in uuids:
            msg = "Call dmg pool query {} that's expected to {}".format(
                uuid[0], uuid[1])
            self.log.info(msg)
            self.pool.uuid = uuid[0]

            try:
                # Call dmg pool query.
                self.pool.set_query_data()
                if uuid[1] == "FAIL":
                    msg = "Query expected to fail, but worked! {}".format(
                        uuid[0])
                    errors_list.append(msg)
            except TestFail:
                if uuid[1] == "PASS":
                    msg = "Query expected to work, but failed! {}".format(
                        uuid[0])
                    errors_list.append(msg)

        # Restore the original UUID.
        self.pool.uuid = uuids[-1][0]

        # Report errors and fail test if needed.
        if errors_list:
            for err in errors_list:
                self.log.error("==>   Failure: %s", err)
            self.fail("Failed dmg pool query input test.")

    def test_pool_query_ior(self):
        """
        JIRA ID: DAOS-2976

        Test Description: Test that pool query command will properly and
        accurately show the size changes once there is content in the pool.

        :avocado: tags=all,daily_regression
        :avocado: tags=small,hw
        :avocado: tags=dmg,pool_query,basic
        :avocado: tags=pool_query_write
        """
        # Store original pool info
        self.pool.set_query_data()
        out_b = deepcopy(self.pool.query_data)
        self.log.info("==>   Pool info before write: \n%s", out_b)

        #  Run ior
        self.log.info("==>   Write data to pool.")
        self.run_ior_with_pool()

        # Check pool written data
        self.pool.set_query_data()
        out_a = deepcopy(self.pool.query_data)
        self.log.info("==>   Pool info after write: \n%s", out_a)

        # The file should have been written into nvme, compare info
        bytes_orig_val = int(out_b["response"]["tier_stats"][1]["free"])
        bytes_curr_val = int(out_a["response"]["tier_stats"][1]["free"])
        if bytes_orig_val <= bytes_curr_val:
            self.fail(
                "Current NVMe free space should be smaller than {}".format(
                    out_b["response"]["tier_stats"][1]["free"]))
