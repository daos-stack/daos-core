#!/usr/bin/python
'''
  (C) Copyright 2019 Intel Corporation.

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
import sys
import json
from avocado import Test

sys.path.append('./util')
sys.path.append('../util')
sys.path.append('../../../utils/py')
sys.path.append('./../../utils/py')
import server_utils
import write_host_file
import ior_utils
from mpio_utils import MpioUtils, MpioFailed
from daos_api import DaosContext, DaosPool, DaosApiError

class EightServers(Test):
    """
    Test class Description: Runs IOR with 8 servers.

    """

    def __init__(self, *args, **kwargs):

        super(EightServers, self).__init__(*args, **kwargs)

        self.basepath = None
        self.server_group = None
        self.context = None
        self.pool = None
        self.num_procs = None
        self.hostlist_servers = None
        self.hostfile_servers = None
        self.hostlist_clients = None
        self.hostfile_clients = None
        self.mpio = None

    def setUp(self):
        # get paths from the build_vars generated by build
        with open('../../../.build_vars.json') as build_file:
            build_paths = json.load(build_file)
        self.basepath = os.path.normpath(build_paths['PREFIX'] + "/../")
        print("<<{}>>".format(self.basepath))
        self.server_group = self.params.get("name", '/server_config/',
                                            'daos_server')

        # setup the DAOS python API
        self.context = DaosContext(build_paths['PREFIX'] + '/lib/')

        self.hostlist_servers = self.params.get("test_servers",
                                                '/run/hosts/test_machines/*')
        self.hostfile_servers = (
            write_host_file.write_host_file(self.hostlist_servers,
                                            self.workdir))
        print("Host file servers is: {}".format(self.hostfile_servers))

        self.hostlist_clients = self.params.get("test_clients",
                                                '/run/hosts/test_machines/*')
        self.num_procs = self.params.get("np", '/run/ior/client_processes/*')
        self.hostfile_clients = (
            write_host_file.write_host_file(self.hostlist_clients, self.workdir,
                                            None))
        print("Host file clients is: {}".format(self.hostfile_clients))

        server_utils.run_server(self.hostfile_servers, self.server_group,
                                self.basepath)

    def tearDown(self):
        try:
            if self.pool is not None and self.pool.attached:
                self.pool.destroy(1)
        finally:
            server_utils.stop_server(hosts=self.hostlist_servers)

    def executable(self, iorflags=None):
        """
        Executable function to run ior for ssf and fpp
        """

        # parameters used in pool create
        createmode = self.params.get("mode", '/run/pool/createmode/*/')
        createuid = os.geteuid()
        creategid = os.getegid()
        createsetid = self.params.get("setname", '/run/pool/createset/')
        createscm_size = self.params.get("scm_size", '/run/pool/createsize/')
        createnvme_size = self.params.get("nvme_size", '/run/pool/createsize/')
        createsvc = self.params.get("svcn", '/run/pool/createsvc/')
        iteration = self.params.get("iter", '/run/ior/iteration/')
        block_size = self.params.get("b", '/run/ior/transfersize_blocksize/*/')
        transfer_size = self.params.get("t",
                                        '/run/ior/transfersize_blocksize/*/')

        try:
            # initialize MpioUtils
            self.mpio = MpioUtils()
            if self.mpio.mpich_installed(self.hostlist_clients) is False:
                self.fail("Exiting Test: Mpich not installed")

            #print self.mpio.mpichinstall
            # initialize a python pool object then create the underlying
            # daos storage
            self.pool = DaosPool(self.context)
            self.pool.create(createmode, createuid, creategid,
                             createscm_size, createsetid, None, None, createsvc,
                             createnvme_size)

            pool_uuid = self.pool.get_uuid_str()
            svc_list = ""
            for i in range(createsvc):
                svc_list += str(int(self.pool.svc.rl_ranks[i])) + ":"
            svc_list = svc_list[:-1]

            print ("svc_list: {}".format(svc_list))

            ior_utils.run_ior_mpiio(self.basepath, self.mpio.mpichinstall,
                                    pool_uuid, svc_list, self.num_procs,
                                    self.hostfile_clients, iorflags, iteration,
                                    transfer_size, block_size, True)

        except (DaosApiError, MpioFailed) as excep:
            print(excep)

    def test_ssf(self):
        """
        Test ID: DAOS-2121
        Test Description: Run IOR with 1,64 and 128 clients config in ssf mode.
        Use Cases: Different combinations of 1/64/128 Clients,
                   1K/4K/32K/128K/512K/1M transfersize and block size of 32M
                   for 1K transfer size and 128M for rest.
        :avocado: tags=ior,mpiio,eightservers,ior_ssf
        """
        ior_flags = self.params.get("F", '/run/ior/iorflags/ssf/')
        self.executable(ior_flags)

    def test_fpp(self):
        """
        Test ID: DAOS-2121
        Test Description: Run IOR with 1,64 and 128 clients config in fpp mode.
        Use Cases: Different combinations of 1/64/128 Clients,
                   1K/4K/32K/128K/512K/1M transfersize and block size of 32M
                   for 1K transfer size and 128M for rest.
        :avocado: tags=ior,mpiio,eightservers,ior_fpp
        """
        ior_flags = self.params.get("F", '/run/ior/iorflags/fpp/')
        self.executable(ior_flags)
