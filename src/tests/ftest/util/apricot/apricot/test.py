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

# Some useful test classes inherited from avocado.Test

from __future__ import print_function

import os
import json
import re

from avocado import Test as avocadoTest
from avocado import skip, TestFail
from avocado.utils import process
from ClusterShell.NodeSet import NodeSet, NodeSetParseError

import fault_config_utils
import agent_utils
import write_host_file
from daos_api import DaosContext, DaosLog, DaosApiError
from server_utils import ServerManager, ServerFailed


# pylint: disable=invalid-name
def skipForTicket(ticket):
    """Skip a test with a comment about a ticket."""
    return skip("Skipping until {} is fixed.".format(ticket))
# pylint: enable=invalid-name


class Test(avocadoTest):
    """Basic Test class.

    :avocado: recursive
    """

    def __init__(self, *args, **kwargs):
        """Initialize a Test object."""
        super(Test, self).__init__(*args, **kwargs)
        # set a default timeout of 1 minute
        # tests that want longer should set a timeout in their .yaml file
        # all tests should set a timeout and 60 seconds will enforce that
        if not self.timeout:
            self.timeout = 60

        item_list = self.logdir.split('/')
        for index, item in enumerate(item_list):
            if item == 'job-results':
                self.job_id = item_list[index + 1]
                break

        self.log.info("Job-ID: %s", self.job_id)
        self.log.info("Test PID: %s", os.getpid())

        self.basepath = None
        self.orterun = None
        self.prefix = None
        self.ompi_prefix = None
        self.tmp = None
        self.server_group = None
        self.daosctl = None
        self.context = None
        self.pool = None
        self.container = None
        self.hostlist_servers = None
        self.hostfile_servers = None
        self.hostfile_servers_slots = 1
        self.partition_servers = None
        self.hostlist_clients = None
        self.hostfile_clients = None
        self.hostfile_clients_slots = 1
        self.partition_clients = None
        self.d_log = None
        self.uri_file = None
        self.fault_file = None

    # pylint: disable=invalid-name
    def cancelForTicket(self, ticket):
        """Skip a test due to a ticket needing to be completed."""
        return self.cancel("Skipping until {} is fixed.".format(ticket))
    # pylint: enable=invalid-name


class TestWithoutServers(Test):
    """Run tests without DAOS servers.

    :avocado: recursive
    """

    def setUp(self):
        """Set up run before each test."""
        super(TestWithoutServers, self).setUp()
        # get paths from the build_vars generated by build
        with open('../../../.build_vars.json') as build_vars:
            build_paths = json.load(build_vars)
        self.basepath = os.path.normpath(os.path.join(build_paths['PREFIX'],
                                                      '..') + os.path.sep)
        self.prefix = build_paths['PREFIX']
        self.ompi_prefix = build_paths["OMPI_PREFIX"]
        self.tmp = os.path.join(self.prefix, 'tmp')
        self.bin = os.path.join(self.prefix, 'bin')
        self.daos_test = os.path.join(self.bin, 'daos_test')
        self.orterun = os.path.join(self.ompi_prefix, "bin", "orterun")
        self.daosctl = os.path.join(self.bin, 'daosctl')

        # setup fault injection, this MUST be before API setup
        fault_list = self.params.get("fault_list", '/run/faults/*/')
        if fault_list:
            # not using workdir because the huge path was messing up
            # orterun or something, could re-evaluate this later
            tmp = os.path.join(self.basepath, 'install', 'tmp')
            self.fault_file = fault_config_utils.write_fault_file(tmp,
                                                                  fault_list,
                                                                  None)
            os.environ["D_FI_CONFIG"] = self.fault_file

        self.context = DaosContext(self.prefix + '/lib/')
        self.d_log = DaosLog(self.context)

    def tearDown(self):
        """Tear down after each test case."""
        super(TestWithoutServers, self).tearDown()

        if self.fault_file:
            os.remove(self.fault_file)

    def multi_log(self, msg, log_type="info"):
        """Log the provided message to the daos log and the test log.

        Args:
            msg (str): message to log
            log_type (str, optional): logging method name to call with the
                message.  Defaults to "info".
        """
        for log_object in (self.d_log, self.log):
            try:
                getattr(log_object, log_type)(msg)
            except AttributeError as error:
                self.fail("Error logging '{}': {}".format(msg, error))


class TestWithServers(TestWithoutServers):
    """Run tests with DAOS servers and at least one client.

    Optionally run DAOS clients on specified hosts.  By default run a single
    DAOS client on the host executing the test.

    :avocado: recursive
    """

    def __init__(self, *args, **kwargs):
        """Initialize a TestWithServers object."""
        super(TestWithServers, self).__init__(*args, **kwargs)

        self.agent_sessions = None
        self.setup_start_servers = True
        self.setup_start_agents = True
        self.server_log = None
        self.client_log = None
        self.log_dir = os.path.split(
            os.getenv("D_LOG_FILE", "/tmp/server.log"))[0]
        self.test_id = "{}-{}".format(
            os.path.split(self.filename)[1], self.name.str_uid)

    def setUp(self):
        """Set up each test case."""
        super(TestWithServers, self).setUp()

        self.server_group = self.params.get(
            "name", "/server_config/", "daos_server")

        # Determine which hosts to use as servers and optionally clients.
        # Support the use of a host type count to test with subsets of the
        # specified hosts lists
        test_machines = self.params.get("test_machines", "/run/hosts/*")
        test_servers = self.params.get("test_servers", "/run/hosts/*")
        test_clients = self.params.get("test_clients", "/run/hosts/*")
        server_count = self.params.get("server_count", "/run/hosts/*")
        client_count = self.params.get("client_count", "/run/hosts/*")
        # If server or client host list are defined through valid slurm
        # partition names override any hosts specified through lists.
        test_servers, self.partition_servers = self.get_partition_hosts(
            "server_partition", test_servers)
        test_clients, self.partition_clients = self.get_partition_hosts(
            "client_partition", test_clients)

        # Supported combinations of yaml hosts arguments:
        #   - test_machines [+ server_count]
        #   - test_servers [+ server_count]
        #   - test_servers [+ server_count] + test_clients [+ client_count]
        if test_machines:
            self.hostlist_servers = test_machines[:server_count]
        elif test_servers and test_clients:
            self.hostlist_servers = test_servers[:server_count]
            self.hostlist_clients = test_clients[:client_count]
        elif test_servers:
            self.hostlist_servers = test_servers[:server_count]
        self.log.info("hostlist_servers:  %s", self.hostlist_servers)
        self.log.info("hostlist_clients:  %s", self.hostlist_clients)

        # If a specific count is specified, verify enough servers/clients are
        # specified to satisy the count
        host_count_checks = (
            ("server", server_count,
             len(self.hostlist_servers) if self.hostlist_servers else 0),
            ("client", client_count,
             len(self.hostlist_clients) if self.hostlist_clients else 0)
        )
        for host_type, expected_count, actual_count in host_count_checks:
            if expected_count:
                self.assertEqual(
                    expected_count, actual_count,
                    "Test requires {} {}; {} specified".format(
                        expected_count, host_type, actual_count))

        # Create host files
        if self.hostlist_clients:
            self.hostfile_clients = write_host_file.write_host_file(
                self.hostlist_clients, self.workdir,
                self.hostfile_clients_slots)

        # Start the clients (agents)
        if self.setup_start_agents:
            self.agent_sessions = agent_utils.run_agent(
                self.basepath, self.hostlist_servers, self.hostlist_clients)

        # Start the servers
        if self.setup_start_servers:
            self.start_servers()

    def tearDown(self):
        """Tear down after each test case."""
        # Destroy any containers first
        errors = self.destroy_containers(self.container)

        # Destroy any pools next
        errors.extend(self.destroy_pools(self.pool))

        # Stop the agents
        errors.extend(self.stop_agents())

        # Stop the servers
        errors.extend(self.stop_servers())

        # Disable fault injection, if enabled
        try:
            super(TestWithServers, self).tearDown()
        except OSError as error:
            errors.append(
                "Error running inheritted teardown(): {}".format(error))

        # Fail the test if any errors occurred during tear down
        if errors:
            self.fail(
                "Errors detected during teardown:\n  - {}".format(
                    "\n  - ".join(errors)))

    def destroy_containers(self, containers):
        """Close and destroy one or more containers.

        Args:
            containers (object): a list of or single DaosContainer or
                TestContainer object(s) to destroy

        Returns:
            list: a list of exceptions raised destroying the containers

        """
        error_list = []
        if containers:
            if not isinstance(containers, (list, tuple)):
                containers = [containers]
            self.multi_log("Destroying containers")
            for container in containers:
                # Only close a container that has been openned by the test
                if not hasattr(container, "opened") or container.opened:
                    try:
                        container.close()
                    except (DaosApiError, TestFail) as error:
                        self.multi_log("  {}".format(error))
                        error_list.append(
                            "Error closing the container: {}".format(error))

                # Only destroy a container that has been created by the test
                if not hasattr(container, "attached") or container.attached:
                    try:
                        container.destroy()
                    except (DaosApiError, TestFail) as error:
                        self.multi_log("  {}".format(error))
                        error_list.append(
                            "Error destroying container: {}".format(error))
        return error_list

    def destroy_pools(self, pools):
        """Disconnect and destroy one or more pools.

        Args:
            pools (object): a list of or single DaosPool or TestPool object(s)
                to destroy

        Returns:
            list: a list of exceptions raised destroying the pools

        """
        error_list = []
        if pools:
            if not isinstance(pools, (list, tuple)):
                pools = [pools]
            self.multi_log("Destroying pools")
            for pool in pools:
                # Only disconnect a pool that has been connected by the test
                if not hasattr(pool, "connected") or pool.connected:
                    try:
                        pool.disconnect()
                    except (DaosApiError, TestFail) as error:
                        self.multi_log("  {}".format(error))
                        error_list.append(
                            "Error disconnecting pool: {}".format(error))

                # Only destroy a pool that has been created by the test
                if not hasattr(pool, "attached") or pool.attached:
                    try:
                        pool.destroy(1)
                    except (DaosApiError, TestFail) as error:
                        self.multi_log("  {}".format(error))
                        error_list.append(
                            "Error destroying pool: {}".format(error))
        return error_list

    def stop_agents(self):
        """Stop the daos agents.

        Returns:
            list: a list of exceptions raised stopping the agents

        """
        error_list = []
        if self.agent_sessions:
            self.multi_log("Stopping agents")
            try:
                agent_utils.stop_agent(
                    self.agent_sessions, self.hostlist_clients)
            except agent_utils.AgentFailed as error:
                self.multi_log("  {}".format(error))
                error_list.append("Error stopping agents: {}".format(error))
        return error_list

    def stop_servers(self):
        """Stop the daos server and I/O servers.

        Returns:
            list: a list of exceptions raised stopping the servers

        """
        error_list = []
        if self.server_managers:
            for server_manager in self.server_managers:
                error = server_manager.stop()
                error_list.append(error)

        return error_list

    def start_servers(self, server_groups=None):
        """Start the servers and clients.

        Args:
            server_groups (dict, optional): [description]. Defaults to None.
        """
        if server_groups is None:
            server_groups = {"daos_server": self.hostlist_servers}
        if isinstance(server_groups, dict):
            # Optionally start servers on a different subset of hosts with a
            # different server group
            self.server_managers = []
            for group, hosts in server_groups.items():
                self.server_manager = ServerManager(
                    self.bin,
                    os.path.join(self.ompi_prefix, "bin"),
                    enable_path=True if self.prefix != "/usr" else False)
                self.server_managers.append(self.server_manager)
                self.server_manager.get_params(self)
                self.server_manager.runner.job.yaml_params.name = group
                self.log.info(
                    "Starting servers: group=%s, hosts=%s", group, hosts)
                self.server_manager.hosts = (
                    hosts, self.workdir, self.hostfile_servers_slots)
                try:
                    avocado_file = "src/tests/ftest/data/daos_avocado_test.yaml"
                    yamlfile = "{}/{}".format(self.basepath, avocado_file)
                    self.server_manager.start(yamlfile)
                except ServerFailed as error:
                    self.multi_log("  {}".format(error))
                    self.fail("Error starting server: {}".format(error))

    def get_partition_hosts(self, partition_key, host_list):
        """[summary].

        Args:
            partition_key ([type]): [description]
            host_list ([type]): [description]

        Returns:
            tuple: [description]

        """
        hosts = []
        partiton_name = self.params.get(partition_key, "/run/hosts/*")
        if partiton_name is not None:
            cmd = "scontrol show partition {}".format(partiton_name)

            try:
                result = process.run(cmd, shell=True, timeout=10)
            except process.CmdError as error:
                self.log.warning(
                    "Unable to obtain hosts from the {} slurm "
                    "partition: {}".format(partiton_name, error))
                result = None
            if result:
                output = result.stdout
                try:
                    hosts = list(
                        NodeSet(re.findall(r"\s+Nodes=(.*)", output)[0]))
                except (NodeSetParseError, IndexError):
                    self.log.warning(
                        "Unable to obtain hosts from the {} slurm partition "
                        "output: {}".format(partiton_name, output))

        if hosts:
            return hosts, partiton_name
        else:
            return host_list, None

    def update_log_file_names(self, test_name=None):
        """Get separate logs for both servers and clients.

        Args:
            test_name (str, optional): name of test variant
        """
        # Determine the path and name of the daos server log using the
        # D_LOG_FILE env or, if not set, the value used in the doas server yaml
        if test_name:
            self.test_id = test_name

        self.server_log = os.path.join(
            self.log_dir, "{}_server_daos.log".format(self.test_id))
        self.client_log = os.path.join(
            self.log_dir, "{}_client_daos.log".format(self.test_id))
