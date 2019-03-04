#!/usr/bin/python
'''
    (C) Copyright 2018Copyright 2018-2019 Intel Corporation.

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
    '''

import os
import time
import traceback
import sys
import json
import uuid
from apricot       import TestWithServers

import ServerUtils
import WriteHostFile
import daos_api
from daos_api import DaosContext, DaosPool, DaosContainer, DaosApiError

class OpenClose(TestWithServers):
    """
    Tests DAOS container open/close function with handle parameter.
    :avocado: recursive
    """

    def tearDown(self):
        try:
            if self.pool is not None and self.pool.attached:
                self.pool.destroy(1)
        finally:
            try:
                super(OpenClose, self).tearDown()
            except ServerFailed as e:
                pass

    def test_closehandle(self):
        """
        Test container close function with container handle paramter.

        :avocado: tags=container,openclose,closehandle
        """
        saved_coh = None

        # parameters used in pool create
        createmode = self.params.get("mode",'/run/pool/createmode/')
        createuid  = os.geteuid()
        creategid  = os.getegid()
        createsetid = self.params.get("setname",'/run/pool/createset/')
        createsize  = self.params.get("size",'/run/pool/createsize/')
        coh_params = self.params.get("coh",'/run/container/container_handle/*/')

        expected_result = coh_params[1]

        try:
            # initialize a python pool object then create the underlying
            # daos storage
            self.pool = DaosPool(self.Context)
            self.pool.create(createmode, createuid, creategid,
                            createsize, createsetid, None)

            poh = self.pool.handle
            self.pool.connect(1 << 1)

            # Container initialization and creation
            self.Container1 = DaosContainer(self.Context)
            self.Container1.create(poh)
            str_cuuid = self.Container1.get_uuid_str()
            cuuid = uuid.UUID(str_cuuid)
            self.Container1.open(poh, cuuid, 2, None)

            # Defining 'good' and 'bad' container handles
            saved_coh = self.Container1.coh
            if coh_params[0] == 'GOOD':
                coh = self.Container1.coh
            else:
                # create a second container, open to get a handle
                # then close & destroy so handle is invalid
                self.Container2 = DaosContainer(self.Context)
                self.Container2.create(poh)
                self.Container2.open(poh, cuuid, 2, None)
                coh = self.Container2.coh
                self.Container2.close()
                self.Container2.destroy()

            # close container with either good or bad handle
            self.Container1.close(coh)

            if expected_result in ['FAIL']:
                self.fail("Test was expected to fail but it passed.\n")

        except DaosApiError as e:
            if expected_result == 'PASS':
                print(e)
                print(traceback.format_exc())
                self.fail("Test was expected to pass but it failed.\n")

            # close above failed so close for real with the right coh
            if saved_coh is not None:
                self.Container1.close(saved_coh)

        finally:
            self.Container1.destroy(1)
            self.pool.disconnect()
            self.pool.destroy(1)
            self.pool = None
