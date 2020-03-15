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
import re

from daos_core_base import DaosCoreBase
from dmg_utils import DmgCommand
from avocado.utils import process


class CSumErrorLog(DaosCoreBase):
    """
    Test Class Description: This test is runs
    daos_test -z (Checksum tests) and verifies
    whether Checksum Error Counters are incremented
    in the NVME device due to checksum fault injection.
    :avocado: recursive
    """
    # pylint: disable=too-many-instance-attributes
    def setUp(self):
        super(CSumErrorLog, self).setUp()
        self.dmg = DmgCommand(os.path.join(self.prefix, "bin"))
        self.dmg.get_params(self)
        self.dmg.hostlist.value = self.hostlist_servers[0]
        self.dmg.set_sub_command("storage")
        self.dmg.sub_command_class.set_sub_command("query")

    def get_nvme_device_id(self):
        self.dmg.sub_command_class.sub_command_class.set_sub_command("smd")
        self.dmg.sub_command_class. \
            sub_command_class.sub_command_class.devices.value = True
        self.dmg.sub_command_class. \
            sub_command_class.sub_command_class.pools.value = True
        try:
            result = self.dmg.run()
        except process.CmdError as details:
            self.fail("dmg command failed: {}".format(details))
        uid = None
        for line in result.stdout.splitlines():
            line = line.strip()
            if re.search("^UUID:",line):
                uid = line.split()
        return uid

    def get_checksum_error_value(self, device_id=None):
        if device_id is None:
            self.fail("No device id provided")
            return
        self.dmg.sub_command_class. \
            sub_command_class.set_sub_command("blobstore-health")
        self.dmg.sub_command_class. \
            sub_command_class. \
            sub_command_class.devuuid.value = "{}".format(device_id)
        try:
            result = self.dmg.run()
        except process.CmdError as details:
            self.fail("dmg command failed: {}".format(details))
        csum_count = None
        for line in result.stdout.splitlines():
            line = line.strip()
            if re.search("^Checksum", line):
                csum_count = line.split()
        return csum_count

    def test_csum_error_logging(self):
        """
        Test ID: DAOS-3927
        Test Description: Write Avocado Test to verify single data after
                          pool/container disconnect/reconnect.
        :avocado: tags=all,full_regression,pr,hw,small,csum_error_log
        """
        dev_id = self.get_nvme_device_id()
        self.log.info("%s", dev_id)
        csum = self.get_checksum_error_value(dev_id[1])
        self.log.info("Checksum Errors : %s", csum[2])
        DaosCoreBase.run_subtest(self)
        csum_latest = self.get_checksum_error_value(dev_id[1])
        self.log.info("Checksum Errors : %s", csum_latest[2])
        if int(csum_latest[2]) > int(csum[2]):
            self.log.info("Checksum Error Logging Test Passed")
        else:
            self.fail("Checksum Error Log not incremented")
