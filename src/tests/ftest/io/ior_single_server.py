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
  The Government's rights to use, modify, reproduce, release, perform, display,
  or disclose this software are subject to the terms of the Apache License as
  provided in Contract No. B609815.
  Any reproduction of computer software, computer software documentation, or
  portions thereof marked with this legend must also reproduce the markings.
'''
from __future__ import print_function

import os
import json

from apricot import Test

import agent_utils
import server_utils
import write_host_file
import ior_utils
from daos_api import DaosContext, DaosPool, DaosApiError

class IorSingleServer(Test):
    """
    Tests IOR with Single Server config.
    :avocado: recursive
    """
    def setUp(self):
        # get paths from the build_vars generated by build
        with open('../../../.build_vars.json') as build_file:
            build_paths = json.load(build_file)
        self.basepath = os.path.normpath(build_paths['PREFIX'] + "/../")

        self.server_group = self.params.get("name", '/server_config/',
                                            'daos_server')
        self.daosctl = self.basepath + '/install/bin/daosctl'

        # setup the DAOS python API
        self.context = DaosContext(build_paths['PREFIX'] + '/lib/')
        self.pool = None

        self.hostlist_servers = self.params.get("test_servers",
                                                '/run/hosts/test_machines/*')
        self.hostfile_servers = (
            write_host_file.write_host_file(self.hostlist_servers,
                                            self.workdir))
        print("Host file servers is: {}".format(self.hostfile_servers))

        self.hostlist_clients = self.params.get("test_clients",
                                                '/run/hosts/test_machines/*')
        self.hostfile_clients = (
            write_host_file.write_host_file(self.hostlist_clients,
                                            self.workdir, None))
        print("Host file clientsis: {}".format(self.hostfile_clients))

        self.agent_sessions = agent_utils.run_agent(self.basepath,
                                                    self.hostlist_servers,
                                                    self.hostlist_clients)
        server_utils.run_server(self.hostfile_servers, self.server_group,
                                self.basepath)

        #if int(str(self.name).split("-")[0]) == 1:
        #    ior_utils.build_ior(self.basepath)

    def tearDown(self):
        try:
            if self.hostfile_clients is not None:
                os.remove(self.hostfile_clients)
            if self.hostfile_servers is not None:
                os.remove(self.hostfile_servers)
            if self.pool is not None and self.pool.attached:
                self.pool.destroy(1)
        finally:
            if self.agent_sessions:
                agent_utils.stop_agent(self.agent_sessions,
                                       self.hostlist_clients)
            server_utils.stop_server(hosts=self.hostlist_servers)

    def test_singleserver(self):
        """
        Test IOR with Single Server config.

        :avocado: tags=all,daosio,large,singleserver
        """

        # parameters used in pool create
        createmode = self.params.get("mode", '/run/createtests/createmode/*/')
        createuid = os.geteuid()
        creategid = os.getegid()
        createsetid = self.params.get("setname", '/run/createtests/createset/')
        createsize = self.params.get("size", '/run/createtests/createsize/')
        createsvc = self.params.get("svcn", '/run/createtests/createsvc/')

        # ior parameters
        client_processes = self.params.get("np", '/run/ior/client_processes/*/')
        iteration = self.params.get("iter", '/run/ior/iteration/')
        ior_flags = self.params.get("F", '/run/ior/iorflags/')
        transfer_size = self.params.get("t",
                                        '/run/ior/transfersize_blocksize/*/')
        block_size = self.params.get("b", '/run/ior/transfersize_blocksize/*/')
        object_class = self.params.get("o", '/run/ior/objectclass/')

        try:
            # initialize a python pool object then create the underlying
            # daos storage
            self.pool = DaosPool(self.context)
            self.pool.create(createmode, createuid, creategid,
                             createsize, createsetid, None, None, createsvc)
            pool_uuid = self.pool.get_uuid_str()
            print ("pool_uuid: {}".format(pool_uuid))
            tmp_rank_list = []
            svc_list = ""
            for item in range(createsvc):
                tmp_rank_list.append(int(self.pool.svc.rl_ranks[item]))
                svc_list += str(tmp_rank_list[item]) + ":"
            svc_list = svc_list[:-1]

            ior_utils.run_ior_daos(self.hostfile_clients, ior_flags, iteration,
                                   block_size, transfer_size, pool_uuid,
                                   svc_list, object_class, self.basepath,
                                   client_processes)

        except (DaosApiError, ior_utils.IorFailed) as excep:
            self.fail("<Single Server Test FAILED>\n {}".format(excep))
