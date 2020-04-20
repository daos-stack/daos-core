#!/usr/bin/python3

"""Test code for dfuse"""

import os
import sys
import time
import uuid
import json
import signal
import subprocess
import tempfile
import pickle

from collections import OrderedDict

class DFTestFail(Exception):
    """Used to indicate test failure"""
    pass

def umount(path):
    """Umount dfuse from a given path"""
    cmd = ['fusermount', '-u', path]
    ret = subprocess.run(cmd)
    print('rc from umount {}'.format(ret.returncode))
    return ret.returncode

def load_conf():
    """Load the build config file"""
    file_self = os.path.dirname(os.path.abspath(__file__))
    json_file = None
    while True:
        new_file = os.path.join(file_self, '.build_vars.json')
        if os.path.exists(new_file):
            json_file = new_file
            break
        file_self = os.path.dirname(file_self)
        if file_self == '/':
            raise Exception('build file not found')
    ofh = open(json_file, 'r')
    conf = json.load(ofh)
    ofh.close()
    return conf

class DaosServer():
    """Manage a DAOS server instance"""

    def __init__(self, conf):
        self.running = False

        self._sp = None
        self.conf = conf
        self._agent = None
        self.agent_dir = None

        socket_dir = '/tmp/daos_sockets'
        if not os.path.exists(socket_dir):
            os.mkdir(socket_dir)
        if os.path.exists('/tmp/daos_server.log'):
            os.unlink('/tmp/daos_server.log')

        self._agent_dir = tempfile.TemporaryDirectory(prefix='daos_agent_')
        self.agent_dir = self._agent_dir.name

    def __del__(self):
        if self.running:
            self.stop()

    def start(self):
        """Start a DAOS server"""

        daos_server = os.path.join(self.conf['PREFIX'], 'bin', 'daos_server')

        self_dir = os.path.dirname(os.path.abspath(__file__))

        server_config = os.path.join(self_dir, 'daos_server.yml')

        cmd = [daos_server, '--config={}'.format(server_config),
               'start', '-t' '4', '--insecure', '-d', self.agent_dir,
               '--recreate-superblocks']

        server_env = os.environ.copy()
        server_env['CRT_PHY_ADDR_STR'] = 'ofi+sockets'
        server_env['OFI_INTERFACE'] = 'lo'
        server_env['D_LOG_MASK'] = 'INFO'
        server_env['DD_MASK'] = 'all'
        server_env['DAOS_DISABLE_REQ_FWD'] = '1'
        self._sp = subprocess.Popen(cmd, env=server_env)

        agent_config = os.path.join(self_dir, 'daos_agent.yml')

        agent_env = os.environ.copy()
        # DAOS-??? Need to set this for agent
        agent_env['LD_LIBRARY_PATH'] = os.path.join(self.conf['PREFIX'],
                                                    'lib64')
        agent_bin = os.path.join(self.conf['PREFIX'], 'bin', 'daos_agent')

        self._agent = subprocess.Popen([agent_bin,
                                        '--config-path={}'.format(agent_config),
                                        '-i', '--runtime_dir', self.agent_dir],
                                       env=agent_env)
        time.sleep(10)
        self.running = True

    def stop(self):
        """Stop a previously started DAOS server"""
        if self._agent:
            self._agent.send_signal(signal.SIGINT)
            ret = self._agent.wait(timeout=5)
            print('rc from agent is {}'.format(ret))

        if not self._sp:
            return
        self._sp.send_signal(signal.SIGTERM)
        ret = self._sp.wait(timeout=5)
        print('rc from server is {}'.format(ret))
        log_test('/tmp/daos_server.log')
        self.running = False

def il_cmd(dfuse, cmd):
    """Run a command under the interception library"""
    my_env = os.environ.copy()
    my_env['CRT_PHY_ADDR_STR'] = 'ofi+sockets'
    my_env['OFI_INTERFACE'] = 'lo'
    log_file = tempfile.NamedTemporaryFile(prefix='daos_dfuse_il_',
                                           suffix='.log', delete=False)
    symlink_file('/tmp/dfuse_il_latest.log', log_file.name)
    my_env['D_LOG_FILE'] = log_file.name
    my_env['LD_PRELOAD'] = os.path.join(dfuse.conf['PREFIX'],
                                        'lib64', 'libioil.so')
    my_env['D_LOG_MASK'] = 'DEBUG'
    my_env['DD_MASK'] = 'all'
    ret = subprocess.run(cmd, env=my_env)
    print('Logged il to {}'.format(log_file.name))
    print(ret)
    print('Log results for il')
    log_test(log_file.name)
    return ret

def symlink_file(a, b):
    """Create a symlink from a to b"""
    if os.path.exists(a):
        os.remove(a)
    os.symlink(b, a)

class DFuse():
    """Manage a dfuse instance"""
    def __init__(self, daos, conf, pool=None, container=None):
        self.dir = '/tmp/dfs_test'
        self.pool = pool
        self.container = container
        self.conf = conf
        self._daos = daos
        self._sp = None

        log_file = tempfile.NamedTemporaryFile(prefix='daos_dfuse_', suffix='.log', delete=False)
        self.log_file = log_file.name

        symlink_file('/tmp/dfuse_latest.log', self.log_file)

        if not os.path.exists(self.dir):
            os.mkdir(self.dir)

    def start(self):
        """Start a dfuse instance"""
        dfuse_bin = os.path.join(self.conf['PREFIX'], 'bin', 'dfuse')
        my_env = os.environ.copy()

        single_threaded = False
        caching = False

        pre_inode = os.stat(self.dir).st_ino

        my_env['CRT_PHY_ADDR_STR'] = 'ofi+sockets'
        my_env['OFI_INTERFACE'] = 'eth0'
        my_env['D_LOG_MASK'] = 'INFO,dfuse=DEBUG,dfs=DEBUG'
        my_env['DD_MASK'] = 'all'
        my_env['DD_SUBSYS'] = 'all'
        my_env['D_LOG_FILE'] = self.log_file
        my_env['DAOS_AGENT_DRPC_DIR'] = self._daos.agent_dir

        cmd = ['valgrind', '--quiet']

        cmd.extend(['--leak-check=full', '--show-leak-kinds=all'])


        cmd.extend(['--suppressions={}'.format(os.path.join('src',
                                                            'cart',
                                                            'utils',
                                                            'memcheck-cart.supp')),
                    '--suppressions={}'.format(os.path.join('utils',
                                                            'memcheck-daos-client.supp'))])

        cmd.extend(['--xml=yes', '--xml-file=dfuse.%p.memcheck'])

        cmd.extend([dfuse_bin, '-s', '0', '-m', self.dir, '-f'])

        if single_threaded:
            cmd.append('-S')

        if caching:
            cmd.append('--enable-caching')

        if self.pool:
            cmd.extend(['--pool', self.pool])
        if self.container:
            cmd.extend(['--container', self.container])
        self._sp = subprocess.Popen(cmd, env=my_env)
        print('Started dfuse at {}'.format(self.dir))
        print('Log file is {}'.format(self.log_file))

        total_time = 0
        while os.stat(self.dir).st_ino == pre_inode:
            print('Dfuse not started, waiting...')
            try:
                ret = self._sp.wait(timeout=1)
                print('dfuse command exited with {}'.format(ret))
                self._sp = None
                raise Exception('dfuse died waiting for start')
            except subprocess.TimeoutExpired:
                pass
            total_time += 1
            if total_time > 10:
                raise Exception('Timeout starting dfuse')

    def _close_files(self):
        work_done = False
        for fname in os.listdir('/proc/self/fd'):
            try:
                tfile = os.readlink(os.path.join('/proc/self/fd', fname))
            except FileNotFoundError:
                continue
            if tfile.startswith(self.dir):
                print('closing file {}'.format(tfile))
                os.close(int(fname))
                work_done = True
        return work_done

    def __del__(self):
        if self._sp:
            self.stop()

    def stop(self):
        """Stop a previously started dfuse instance"""
        if not self._sp:
            return

        print('Stopping fuse')
        ret = umount(self.dir)
        if ret:
            self._close_files()
            umount(self.dir)

        ret = self._sp.wait(timeout=20)
        print('rc from dfuse {}'.format(ret))
        self._sp = None
        log_test(self.log_file)

    def wait_for_exit(self):
        """Wait for dfuse to exit"""
        ret = self._sp.wait()
        print('rc from dfuse {}'.format(ret))
        self._sp = None

def get_pool_list():
    """Return a list of valid pool names"""
    pools = []

    for fname in os.listdir('/mnt/daos'):
        if len(fname) != 36:
            continue
        try:
            uuid.UUID(fname)
        except ValueError:
            continue
        pools.append(fname)
    return pools

def assert_file_size(ofd, size):
    """Verify the file size is as expected"""
    stat = os.fstat(ofd.fileno())
    print('Checking file size is {} {}'.format(size, stat.st_size))
    assert stat.st_size == size

def import_daos(server, conf):
    """Return a handle to the pydaos module"""

    if sys.version_info.major < 3:
        pydir = 'python{}.{}'.format(sys.version_info.major, sys.version_info.minor)
    else:
        pydir = 'python{}'.format(sys.version_info.major)

    sys.path.append(os.path.join(conf['PREFIX'],
                                 'lib64',
                                 pydir,
                                 'site-packages'))

    os.environ['CRT_PHY_ADDR_STR'] = 'ofi+sockets'
    os.environ['OFI_INTERFACE'] = 'lo'
    os.environ["DAOS_AGENT_DRPC_DIR"] = server.agent_dir

    daos = __import__('pydaos')
    return daos

def show_cont(conf, pool):
    """Create a container and return a container list"""
    daos_bin = os.path.join(conf['PREFIX'], 'bin', 'daos')
    cmd = [daos_bin, 'container', 'create', '--svc', '0', '--pool', pool]
    rc = subprocess.run(cmd, stdout=subprocess.PIPE)
    print('rc is {}'.format(rc))

    cmd = [daos_bin, 'pool', 'list-containers', '--svc', '0', '--pool', pool]
    rc = subprocess.run(cmd, stdout=subprocess.PIPE)
    print('rc is {}'.format(rc))
    return rc.stdout.strip()

def make_pool(daos, conf):
    """Create a DAOS pool"""

    time.sleep(2)

    daos_raw = __import__('pydaos.raw')

    context = daos.raw.DaosContext(os.path.join(conf['PREFIX'], 'lib64'))

    pool_con = daos.raw.DaosPool(context)

    try:
        pool_con.create(511, os.geteuid(), os.getegid(), 1024*1014*128, b'daos_server')
    except daos_raw.daos_api.DaosApiError:
        time.sleep(10)
        pool_con.create(511, os.geteuid(), os.getegid(), 1024*1014*128, b'daos_server')

    return get_pool_list()

def run_tests(dfuse):
    """Run some tests"""
    path = dfuse.dir

    fname = os.path.join(path, 'test_file3')
    ofd = open(fname, 'w')
    ofd.write('hello')
    print(os.fstat(ofd.fileno()))
    ofd.flush()
    print(os.stat(fname))
    assert_file_size(ofd, 5)
    ofd.truncate(0)
    assert_file_size(ofd, 0)
    ofd.truncate(1024*1024)
    assert_file_size(ofd, 1024*1024)
    ofd.truncate(0)
    ofd.seek(0)
    ofd.write('world\n')
    ofd.flush()
    assert_file_size(ofd, 6)
    print(os.fstat(ofd.fileno()))
    ofd.close()
    il_cmd(dfuse, ['cat', fname])

def stat_and_check(dfuse, pre_stat):
    """Check that dfuse started"""
    post_stat = os.stat(dfuse.dir)
    if pre_stat.st_dev == post_stat.st_dev:
        raise DFTestFail('Device # unchanged')
    if post_stat.st_ino != 1:
        raise DFTestFail('Unexpected inode number')

def check_no_file(dfuse):
    """Check that a non-existent file doesn't exist"""
    try:
        os.stat(os.path.join(dfuse.dir, 'no-file'))
        raise DFTestFail('file exists')
    except FileNotFoundError:
        pass

EFILES = ['src/common/misc.c',
          'src/common/prop.c',
          'src/cart/crt_hg_proc.c',
          'src/security/cli_security.c',
          'src/client/dfuse/dfuse_core.c']

lp = None
lt = None

def setup_log_test():
    """Setup and import the log tracing code"""
    file_self = os.path.dirname(os.path.abspath(__file__))
    logparse_dir = os.path.join(os.path.dirname(file_self), '../../../src/cart/test/util')
    crt_mod_dir = os.path.realpath(logparse_dir)
    if crt_mod_dir not in sys.path:
        sys.path.append(crt_mod_dir)

    global lp
    global lt

    lp = __import__('cart_logparse')
    lt = __import__('cart_logtest')

    lt.mismatch_alloc_ok['crt_proc_d_rank_list_t'] = ('rank_list',
                                                      'rank_list->rl_ranks')
    lt.mismatch_alloc_ok['path_gen'] = ('*fpath')
    lt.mismatch_alloc_ok['get_attach_info'] = ('reqb')
    lt.mismatch_alloc_ok['iod_fetch'] = ('biovs')
    lt.mismatch_alloc_ok['bio_sgl_init'] = ('sgl->bs_iovs')
    lt.mismatch_alloc_ok['process_credential_response'] = ('bytes')
    lt.mismatch_alloc_ok['pool_map_find_tgts'] = ('*tgt_pp')
    lt.mismatch_alloc_ok['daos_acl_dup'] = ('acl_copy')
    lt.mismatch_alloc_ok['dfuse_pool_lookup'] = ('ie', 'dfs', 'dfp')
    lt.mismatch_alloc_ok['pool_prop_read'] = ('prop->dpp_entries[idx].dpe_str',
                                              'prop->dpp_entries[idx].dpe_val_ptr')
    lt.mismatch_alloc_ok['cont_prop_read'] = ('prop->dpp_entries[idx].dpe_str')
    lt.mismatch_alloc_ok['cont_iv_prop_g2l'] = ('prop_entry->dpe_str')
    lt.mismatch_alloc_ok['notify_ready'] = ('reqb')
    lt.mismatch_alloc_ok['pack_daos_response'] = ('body')
    lt.mismatch_alloc_ok['ds_mgmt_drpc_get_attach_info'] = ('body')
    lt.mismatch_alloc_ok['pool_prop_default_copy'] = ('entry_def->dpe_str')
    lt.mismatch_alloc_ok['daos_prop_dup'] = ('entry_dup->dpe_str')
    lt.mismatch_alloc_ok['auth_cred_to_iov'] = ('packed')

    lt.mismatch_free_ok['d_rank_list_free'] = ('rank_list',
                                               'rank_list->rl_ranks')
    lt.mismatch_free_ok['pool_prop_default_copy'] = ('entry_def->dpe_str')
    lt.mismatch_free_ok['pool_svc_store_uuid_cb'] = ('path')
    lt.mismatch_free_ok['ds_mgmt_svc_start'] = ('uri')
    lt.mismatch_free_ok['daos_acl_free'] = ('acl')
    lt.mismatch_free_ok['drpc_free'] = ('pointer')
    lt.mismatch_free_ok['pool_child_add_one'] = ('path')
    lt.mismatch_free_ok['get_tgt_rank'] = ('tgts')
    lt.mismatch_free_ok['bio_sgl_fini'] = ('sgl->bs_iovs')
    lt.mismatch_free_ok['daos_iov_free'] = ('iov->iov_buf'),
    lt.mismatch_free_ok['daos_prop_free'] = ('entry->dpe_str',
                                             'entry->dpe_val_ptr')
    lt.mismatch_free_ok['main'] = ('dfs')
    lt.mismatch_free_ok['start_one'] = ('path')
    lt.mismatch_free_ok['pool_svc_load_uuid_cb'] = ('path')
    lt.mismatch_free_ok['ie_sclose'] = ('ie', 'dfs', 'dfp')
    lt.mismatch_free_ok['notify_ready'] = ('req.uri')
    lt.mismatch_free_ok['get_tgt_rank'] = ('tgts')

    lt.memleak_ok.append('dfuse_start')
    #lt.memleak_ok.append('expand_vector')
    #lt.memleak_ok.append('d_rank_list_alloc')
    #lt.memleak_ok.append('get_tpv')
    #lt.memleak_ok.append('get_new_entry')
    #lt.memleak_ok.append('get_attach_info')
    #lt.memleak_ok.append('drpc_call_create')

def log_test(filename):
    """Run the log checker on filename, logging to stdout"""

    global lp
    global lt

    lp = __import__('cart_logparse')
    lt = __import__('cart_logtest')

    lt.shown_logs = set()

    log_iter = lp.LogIter(filename)
    lto = lt.LogTest(log_iter)

    for efile in EFILES:
        lto.set_error_ok(efile)

    try:
        lto.check_log_file(abort_on_warning=False)
    except lt.LogCheckError:
        print('Error detected')

def create_and_read_via_il(dfuse, path):
    """Create file in dir, write to and and read
    through the interception library"""

    fname = os.path.join(path, 'test_file')
    ofd = open(fname, 'w')
    ofd.write('hello ')
    ofd.write('world\n')
    ofd.flush()
    assert_file_size(ofd, 12)
    print(os.fstat(ofd.fileno()))
    ofd.close()
    il_cmd(dfuse, ['cat', fname])

def run_dfuse(server, conf):
    """Run several dfuse instances"""

    daos = import_daos(server, conf)

    pools = get_pool_list()
    while len(pools) < 1:
        pools = make_pool(daos, conf)

    dfuse = DFuse(server, conf)
    try:
        pre_stat = os.stat(dfuse.dir)
    except OSError:
        umount(dfuse.dir)
        raise
    container = str(uuid.uuid4())
    dfuse.start()
    print(os.statvfs(dfuse.dir))
    subprocess.run(['df', '-h'])
    subprocess.run(['df', '-i', dfuse.dir])
    print('Running dfuse with nothing')
    stat_and_check(dfuse, pre_stat)
    check_no_file(dfuse)
    for pool in pools:
        pool_stat = os.stat(os.path.join(dfuse.dir, pool))
        print('stat for {}'.format(pool))
        print(pool_stat)
        cdir = os.path.join(dfuse.dir, pool, container)
        os.mkdir(cdir)
        #create_and_read_via_il(dfuse, cdir)
    dfuse.stop()

    uns_container = container

    container2 = str(uuid.uuid4())
    dfuse = DFuse(server, conf, pool=pools[0])
    pre_stat = os.stat(dfuse.dir)
    dfuse.start()
    print('Running dfuse with pool only')
    stat_and_check(dfuse, pre_stat)
    check_no_file(dfuse)
    cpath = os.path.join(dfuse.dir, container2)
    os.mkdir(cpath)
    cdir = os.path.join(dfuse.dir, container)
    create_and_read_via_il(dfuse, cdir)

    dfuse.stop()

    dfuse = DFuse(server, conf, pool=pools[0], container=container)
    pre_stat = os.stat(dfuse.dir)
    dfuse.start()
    print('Running fuse with both')

    stat_and_check(dfuse, pre_stat)

    create_and_read_via_il(dfuse, dfuse.dir)

    run_tests(dfuse)

    daos_bin = os.path.join(conf['PREFIX'], 'bin', 'daos')

    dfuse.stop()

    dfuse = DFuse(server, conf, pool=pools[0], container=container2)
    dfuse.start()

    uns_path = os.path.join(dfuse.dir, 'ep0')

    uns_container = str(uuid.uuid4())

    cmd = [daos_bin, 'container', 'create', '--svc', '0',
           '--pool', pools[0], '--cont', uns_container, '--path', uns_path,
           '--type', 'POSIX']

    print('Inserting entry point')
    rc = subprocess.run(cmd)
    print('rc is {}'.format(rc))
    print(os.stat(uns_path))
    print(os.stat(uns_path))
    print(os.listdir(dfuse.dir))

    dfuse.stop()

    print('Trying UNS')
    dfuse = DFuse(server, conf)
    dfuse.start()

    # List the root container.
    print(os.listdir(os.path.join(dfuse.dir, pools[0], container2)))

    uns_path = os.path.join(dfuse.dir, pools[0], container2, 'ep0', 'ep')
    direct_path = os.path.join(dfuse.dir, pools[0], uns_container)

    uns_container = str(uuid.uuid4())

    # Make a link within the new container.
    cmd = [daos_bin, 'container', 'create', '--svc', '0',
           '--pool', pools[0], '--cont', uns_container,
           '--path', uns_path, '--type', 'POSIX']

    print('Inserting entry point')
    rc = subprocess.run(cmd)
    print('rc is {}'.format(rc))

    # List the root container again.
    print(os.listdir(os.path.join(dfuse.dir, pools[0], container2)))

    # List the target container.
    files = os.listdir(direct_path)
    print(files)
    # List the target container through UNS.
    print(os.listdir(uns_path))
    direct_stat = os.stat(os.path.join(direct_path, files[0]))
    uns_stat = os.stat(uns_path)
    print(direct_stat)
    print(uns_stat)
    assert uns_stat.st_ino == direct_stat.st_ino

    dfuse.stop()
    print('Trying UNS with previous cont')
    dfuse = DFuse(server, conf)
    dfuse.start()

    files = os.listdir(direct_path)
    print(files)
    print(os.listdir(uns_path))

    direct_stat = os.stat(os.path.join(direct_path, files[0]))
    uns_stat = os.stat(uns_path)
    print(direct_stat)
    print(uns_stat)
    assert uns_stat.st_ino == direct_stat.st_ino

    print('Reached the end, no errors')

def run_il_test(server, conf):
    """Run a basic interception library test"""
    daos = import_daos(server, conf)

    pools = get_pool_list()

    # TODO: This doesn't work with two pools, there appears to be a bug
    # relating to re-using container uuids across pools.
    while len(pools) < 1:
        pools = make_pool(daos, conf)

    print('pools are ', ','.join(pools))

    containers = ['62176a51-8229-4e4c-ad1b-43aaace8a97a',
                  '4ef12a58-c544-406c-8acf-56a2c0589cd6']

    dfuse = DFuse(server, conf)
    dfuse.start()

    dirs = []

    for p in pools:
        for c in containers:
            d = os.path.join(dfuse.dir, p, c)
            try:
                print('Making directory {}'.format(d))
                os.mkdir(d)
            except FileExistsError:
                pass
            dirs.append(d)

    f = os.path.join(dirs[0], 'file')
    fd = open(f, 'w')
    fd.write('Hello')
    fd.close()
    il_cmd(dfuse, ['cp', f, dirs[-1]])
    il_cmd(dfuse, ['cp', '/bin/bash', dirs[-1]])
    il_cmd(dfuse, ['md5sum', os.path.join(dirs[-1], 'bash')])
    dfuse.stop()

def run_in_fg(server, conf):
    """Run dfuse in the foreground.

    Block until ctrl-c is pressed.
    """
    daos = import_daos(server, conf)

    pools = get_pool_list()

    while len(pools) < 1:
        pools = make_pool(daos, conf)

    dfuse = DFuse(server, conf, pool=pools[0])
    dfuse.start()
    container = str(uuid.uuid4())
    t_dir = os.path.join(dfuse.dir, container)
    os.mkdir(t_dir)
    print('Running at {}'.format(t_dir))
    print('daos container create --svc 0 --type POSIX --pool {} --path {}/uns-link'.format(
        pools[0], t_dir))
    print('cd {}/uns-link'.format(t_dir))
    print('daos container destroy --svc 0 --path {}/uns-link'.format(t_dir))
    print('daos pool list-containers --svc 0 --pool {}'.format(pools[0]))
    try:
        dfuse.wait_for_exit()
    except KeyboardInterrupt:
        pass
    dfuse = None

def test_pydaos_kv(server, conf):
    """Test the KV interface"""

    daos = import_daos(server, conf)

    file_self = os.path.dirname(os.path.abspath(__file__))
    mod_dir = os.path.join(os.path.dirname(file_self), '../../../src/client/pydaos')
    if mod_dir not in sys.path:
        sys.path.append(mod_dir)

    dbm = __import__('daosdbm')

    pools = get_pool_list()

    while len(pools) < 1:
        pools = make_pool(daos, conf)

    pool = pools[0]

    container = show_cont(conf, pool)

    print(container)
    c_uuid = container.decode().split(' ')[-1]
    kvg = dbm.daos_named_kv('ofi+sockets', 'lo', pool, c_uuid)

    kv = kvg.get_kv_by_name('Dave')
    kv['a'] = 'a'
    kv['b'] = 'b'
    kv['list'] = pickle.dumps(list(range(1, 100000)))
    for k in range(1, 100):
        kv[str(k)] = pickle.dumps(list(range(1, 10)))
    print(type(kv))
    print(kv)
    print(kv['a'])

    print("First iteration")
    data = OrderedDict()
    for key in kv:
        print('key is {}, len {}'.format(key, len(kv[key])))
        print(type(kv[key]))
        data[key] = None

    print("Bulk loading")

    data['no-key'] = None

    kv.bget(data, value_size=16)
    print("Second iteration")
    failed = False
    for key in data:
        if data[key]:
            print('key is {}, len {}'.format(key, len(data[key])))
        elif key == 'no-key':
            pass
        else:
            failed = True
            print('Key is None {}'.format(key))

    if failed:
        print("That's not good")

def main():
    """Main entry point"""

    setup_log_test()
    conf = load_conf()
    server = DaosServer(conf)
    server.start()

    if len(sys.argv) == 2 and sys.argv[1] == 'launch':
        run_in_fg(server, conf)
    elif len(sys.argv) == 2 and sys.argv[1] == 'kv':
        test_pydaos_kv(server, conf)
    elif len(sys.argv) == 2 and sys.argv[1] == 'all':
        run_il_test(server, conf)
        run_dfuse(server, conf)
        test_pydaos_kv(server, conf)
    else:
        run_il_test(server, conf)
        run_dfuse(server, conf)

if __name__ == '__main__':
    main()
