#!/usr/bin/python
'''
  (C) Copyright 2017Copyright 2017-2019 Intel Corporation.

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
from apricot import TestWithoutServers
from avocado import main


import ServerUtils
import WriteHostFile
from daos_api import DaosContext, DaosPool, DaosContainer, DaosApiError
from conversion import c_uuid_to_str

class SimpleCreateDeleteTest(TestWithoutServers):
    """
    Tests DAOS container basics including create, destroy, open, query
    and close.

    :avocado: recursive
    """
    def test_container_basics(self):
        """
        Test basic container create/destroy/open/close/query.  Nothing fancy
        just making sure they work at a rudimentary level

        :avocado: tags=container,containercreate,containerdestroy,basecont
        """

        pool = None
        hostlist = None

        try:
            hostlist = self.params.get("test_machines",'/run/hosts/*')
            hostfile = WriteHostFile.WriteHostFile(hostlist, self.workdir)

            ServerUtils.runServer(hostfile, self.server_group, self.basepath)

            # give it time to start
            time.sleep(2)

            # parameters used in pool create
            createmode = self.params.get("mode",'/run/conttests/createmode/')
            createuid  = self.params.get("uid",'/run/conttests/createuid/')
            creategid  = self.params.get("gid",'/run/conttests/creategid/')
            createsetid = self.params.get("setname",'/run/conttests/createset/')
            createsize  = self.params.get("size",'/run/conttests/createsize/')

            # initialize a python pool object then create the underlying
            # daos storage
            pool = DaosPool(self.Context)
            pool.create(createmode, createuid, creategid,
                        createsize, createsetid, None)

            # need a connection to create container
            pool.connect(1 << 1)

            # create a container
            container = DaosContainer(self.Context)
            container.create(pool.handle)

            # now open it
            container.open()

            # do a query and compare the UUID returned from create with
            # that returned by query
            container.query()

            if container.get_uuid_str() != c_uuid_to_str(
                    container.info.ci_uuid):
                self.fail("Container UUID did not match the one in info'n")

            container.close()

            # wait a few seconds and then destroy
            time.sleep(5)
            container.destroy()

        except DaosApiError as e:
            print(e)
            print(traceback.format_exc())
            self.fail("Test was expected to pass but it failed.\n")
        except Exception as e:
            self.fail("Daos code segfaulted most likely, error: %s" % e)
        finally:
            # cleanup the pool
            if pool is not None:
                pool.disconnect()
                pool.destroy(1)

            ServerUtils.stopServer(hosts=hostlist)

if __name__ == "__main__":
    main()
