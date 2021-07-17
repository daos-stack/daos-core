#!/usr/bin/python
"""
  (C) Copyright 2018-2021 Intel Corporation.

  SPDX-License-Identifier: BSD-2-Clause-Patent
"""
from avocado.core.exceptions import TestFail
from dfuse_test_base import DfuseTestBase
from env_modules import load_mpi
from general_utils import DaosTestError
from general_utils import run_command


class PosixSimul(DfuseTestBase):
    # pylint: disable=too-few-public-methods
    # pylint: disable=too-many-ancestors
    """Tests a posix container with simul

    From : https://github.com/LLNL/simul

    "simul" is an MPI coordinated test of parallel filesystem system calls and
    library functions.  It was designed to perform filesystem operations
    simultaneously from many nodes and processes to test the correctness
    and coherence of parallel filesystems.

    List of tests:
        Test #0: open, shared mode.
        Test #1: close, shared mode.
        Test #2: file stat, shared mode.
        Test #3: lseek, shared mode.
        Test #4: read, shared mode.
        Test #5: write, shared mode.
        Test #6: chdir, shared mode.
        Test #7: directory stat, shared mode.
        Test #8: statfs, shared mode.
        Test #9: readdir, shared mode.
        Test #10: mkdir, shared mode.
        Test #11: rmdir, shared mode.
        Test #12: unlink, shared mode.
        Test #13: rename, shared mode.
        Test #14: creat, shared mode.
        Test #15: truncate, shared mode.
        Test #16: symlink, shared mode.
        Test #17: readlink, shared mode.
        Test #18: link to one file, shared mode.
        Test #19: link to a file per process, shared mode.
        Test #20: fcntl locking, shared mode.
        Test #21: open, individual mode.
        Test #22: close, individual mode.
        Test #23: file stat, individual mode.
        Test #24: lseek, individual mode.
        Test #25: read, individual mode.
        Test #26: write, individual mode.
        Test #27: chdir, individual mode.
        Test #28: directory stat, individual mode.
        Test #29: statfs, individual mode.
        Test #30: readdir, individual mode.
        Test #31: mkdir, individual mode.
        Test #32: rmdir, individual mode.
        Test #33: unlink, individual mode.
        Test #34: rename, individual mode.
        Test #35: creat, individual mode.
        Test #36: truncate, individual mode.
        Test #37: symlink, individual mode.
        Test #38: readlink, individual mode.
        Test #39: link to one file, individual mode.
        Test #40: link to a file per process, individual mode.
        Test #41: fcntl locking, individual mode.

    :avocado: recursive
    """

    def run_simul(self, include=None, exclude=None):
        """ Run simul

        include str:  comma-separated list of tests to include
        exclude str:  comma-separated list of tests to exclude

        If include value is set, exclude value is ignored and vice versa.

        """
        simul_openmpi = "/usr/lib64/openmpi3/bin/simul"
        simul_mpich = "/usr/lib64/mpich/bin/simul"
        simul_dict = {"openmpi": simul_openmpi,
                      "mpich": simul_mpich
                      }

        clients = self.params.get("test_clients", '/run/hosts/*')
        dfuse_mount_dir = self.params.get("mount_dir", '/run/dfuse/*')

        if not self.pool:
            self.log.info("Create a pool")
            self.add_pool()
        if not self.container:
            self.log.info("Create container")
            self.add_container(self.pool)

        self.start_dfuse(self.hostlist_clients, self.pool, self.container)
        self.dfuse.check_running()

        # Load openmpi
        for mpi, simul in simul_dict.items():
            load_mpi(mpi)

            for client in clients:
                if include:
                    cmd = "/usr/bin/ssh {0} {1} -vv -d {2} -i {3}".format(
                        client,
                        simul,
                        dfuse_mount_dir,
                        include)
                else:
                    cmd = "/usr/bin/ssh {0} {1} -vv -d {2} -e {3}".format(
                        client,
                        simul,
                        dfuse_mount_dir,
                        exclude)
                try:
                    run_command(cmd)
                except DaosTestError:
                    self.stop_dfuse()
                    self.fail("Simul failed on {}".format(mpi))

                self.log.info("Simul passed on %s", mpi)

        self.stop_dfuse()

    def test_posix_simul(self):
        """Test simul.

        :avocado: tags=all,full_regression
        :avocado: tags=vm
        :avocado: tags=posix,simul
        """
        self.run_simul(exclude="9,18,30,39,40")

    def test_posix_expected_failures(self):
        """Test simul, expected failures

        :avocado: tags=all,full_regression
        :avocado: tags=vm
        :avocado: tags=posix,simul_failure
        """
        expected_failures = ["9", "18", "30", "39", "40"]
        for test in expected_failures:
            try:
                self.run_simul(include=test)
            except TestFail as exception:
                # Only catch expected message, otherwise fail
                if "Simul failed" in exception.args[0]:
                    self.log.info("Expected failure")
                else:
                    raise exception
