# vi: set ft=ruby :
#
# To prepare a node to be able to use this file:
#
# $ sudo dnf -y install vagrant-libvirt libvirt-client
# $ sudo systemctl start libvirtd.socket

# This can only reasonably be done on a newish distro such as Fedora (i.e. 33)
# as it requires a fairly recent QEMU implementation.  EL7 is not new enough.
# EL8 might be but is untested.

# Only needed if user doesn't have blanket sudo permissions:
# # cat <<EOF > /etc/sudoers.d/vagrant
## Allow Vagrant to manage /etc/exports
## https://www.vagrantup.com/docs/synced-folders/nfs.html
#Cmnd_Alias VAGRANT_EXPORTS_CHOWN = /usr/bin/chown 0\:0 /tmp/vagrant*
#Cmnd_Alias VAGRANT_EXPORTS_MV = /usr/bin/mv -f /tmp/* /etc/exports
#Cmnd_Alias VAGRANT_EXPORTS_CAT = /usr/bin/cat /etc/exports
#Cmnd_Alias VAGRANT_NFSD_CHECK = /usr/bin/systemctl status --no-pager nfs-server.service
#Cmnd_Alias VAGRANT_NFSD_START = /usr/bin/systemctl start nfs-server.service
#Cmnd_Alias VAGRANT_NFSD_APPLY = /usr/sbin/exportfs -ar
##%vagrant ALL=(root) NOPASSWD: VAGRANT_EXPORTS_CHOWN, VAGRANT_EXPORTS_MV, VAGRANT_EXPORTS_CAT, VAGRANT_NFSD_CHECK, VAGRANT_NFSD_START, VAGRANT_NFSD_APPLY
#$(id -un) ALL=(root) NOPASSWD: VAGRANT_EXPORTS_CHOWN, VAGRANT_EXPORTS_MV, VAGRANT_EXPORTS_CAT, VAGRANT_NFSD_CHECK, VAGRANT_NFSD_START, VAGRANT_NFSD_APPLY
#

require 'etc'
ENV['VAGRANT_EXPERIMENTAL'] = 'typed_triggers'

Vagrant.configure("2") do |config|

        # disable the default shared folder
        config.vm.synced_folder ".", "/vagrant", disabled: true
        # but share the user's daoshome
        # can't do this if user's daoshome is NFS mounted
        #config.vm.synced_folder ".", File.dirname(__FILE__), type: "nfs",
        #          nfs_version: 4,
        #          nfs_udp: false

        config.vm.provider :libvirt do |libvirt, override|
                override.vm.box = "centos/7"
                libvirt.cpus = 4
                libvirt.numa_nodes = [
                        {:cpus => "0-1", :memory => 1024 * 8},
                        {:cpus => "2-3", :memory => 1024 * 8}
                ]
        end

        # Increase yum timeout for slow mirrors
        config.vm.provision "Extend yum mirror timeout",
                type: "shell",
                inline: "sed -i -e '/^distroverpkg/atimeout=300' /etc/yum.conf"

        # Vagrant (dumbly) adds the host's name to the loopback /etc/hosts
        # entry.  https://github.com/hashicorp/vagrant/issues/7263
        # They seem to be preferring simplicity-in-simple scenarios and
        # do really bad things to that end.  Things that blow up in real-
        # world complexity levels
        config.vm.provision "Fix /etc/hosts",
                type: "shell",
                inline: "sed -i -e \"/^127.0.0.1/s/\$HOSTNAME//g\" -e '/^127.0.0.1[         ]*$/d' /etc/hosts"

        # Verbose booting
        config.vm.provision "Fix grub", type: "shell",
                            inline: "sed -ie 's/ rhgb quiet//' /boot/grub2/grub.cfg /etc/sysconfig/grub"

        # The VMs will have IPv6 but no IPv6 connectivity so alter
        # their gai.conf to prefer IPv4 addresses over IPv6
        config.vm.provision "Fix gai.conf", type: "shell",
                            inline: "echo \"precedence ::ffff:0:0/96  100\" > /etc/gai.conf
for i in all default; do
    echo 1 > /proc/sys/net/ipv6/conf/$i/disable_ipv6
done
if ! grep ip_resolve= /etc/yum.conf; then
    sed -i -e '/^\\[main\\]$/aip_resolve=4' /etc/yum.conf
fi"

        # Install needed packages for daos
        config.vm.provision "Install epel-release", type: "shell",
                            inline: "yum -y install epel-release"
        config.vm.provision "Install basic packages 1", \
                            type: "shell",              \
                            inline: "yum -y install librdmacm libcmocka ed \
                                     python-clustershell python3-pip strace"

        # Allow cluster hosts to ssh to each other
        if not(File.exist?("id_rsa"))
                system("ssh-keygen -t rsa -N '' -f id_rsa -C \"Vagrant cluster\"")
        end

        # Add the generated SSH public key to each host's
        # authorized_keys file.
        config.vm.provision "Copy id_rsa.pub", type: "file",
                                               source: "id_rsa.pub",
                                               destination: "/tmp/id_rsa.pub"
        config.vm.provision "Fix authorized_keys", type: "shell",
                                                   inline: "mkdir -m 0700 -p /root/.ssh
if [ -f /tmp/id_rsa.pub ]; then
    awk -v pk=\"$(cat /tmp/id_rsa.pub)\" 'BEGIN{split(pk,s,\" \")} $2 == s[2] {m=1;exit}END{if (m==0)print pk}' /root/.ssh/authorized_keys >> /root/.ssh/authorized_keys
    awk -v pk=\"$(cat /tmp/id_rsa.pub)\" 'BEGIN{split(pk,s,\" \")} $2 == s[2] {m=1;exit}END{if (m==0)print pk}' /home/vagrant/.ssh/authorized_keys >> /home/vagrant/.ssh/authorized_keys
fi
cat /home/vagrant/.ssh/authorized_keys >> /root/.ssh/authorized_keys
chmod 0600 /root/.ssh/authorized_keys"

        # And make the private key available
        config.vm.provision "Copy id_rsa", type: "file",
                                           source: "id_rsa",
                                           destination: "/tmp/id_rsa"
        config.vm.provision "Configure SSH", type: "shell",
                                             inline: "mkdir -m 0700 -p /root/.ssh
cp /tmp/id_rsa /home/vagrant/.ssh/.
cat <<EOF > /home/vagrant/.ssh/config
Host vm*
  #UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
EOF
chown vagrant.vagrant /home/vagrant/.ssh/{id_rsa,config}
chmod 600 /home/vagrant/.ssh/config
cp /tmp/id_rsa /root/.ssh/.
chmod 0600 /root/.ssh/id_rsa"

        config.vm.provision "Configure NTP", type: "shell",
                                             inline: "set -x; gw=$(ip -o route get 1.1.1.1 | cut -d\" \" -f3)
                                                      ed <<EOF /etc/chrony.conf
1i
# This server was configured by Vagrant
server $gw iburst
.
/^server/,/^$/c

.
wq
EOF"
        #
        # Create the cluster
        #
        (1..3).each do |ss_idx|
                config.vm.define "vm#{ss_idx}", autostart: true do |ss|
                        ss.vm.host_name = "vm#{ss_idx}"
                        ss.trigger.before :up do |trigger|
                                trigger.ruby do |env,machine|
                                        (1..2).each do |nvme_idx|
                                                if File.exist?(ENV['HOME'] + "/.local/share/libvirt/images/nvme_disk#{ss_idx}-#{nvme_idx}.img")
                                                        File.delete(ENV['HOME'] + "/.local/share/libvirt/images/nvme_disk#{ss_idx}-#{nvme_idx}.img")
                                                end
                                                system("qemu-img create -f raw " + ENV['HOME'] + "/.local/share/libvirt/images/nvme_disk#{ss_idx}-#{nvme_idx}.img 32G")
                                                system("restorecon " + ENV['HOME'] + "/.local/share/libvirt/images/nvme_disk#{ss_idx}-#{nvme_idx}.img")
                                        end
                                end
                        end
                        ss.vm.provider :libvirt do |lv|
                                # An NVMe drive:
                                (1..2).each do |nvme_idx|
                                        lv.qemuargs :value => "-drive"
                                        lv.qemuargs :value => "format=raw,file=" + ENV['HOME'] + "/.local/share/libvirt/images/nvme_disk#{ss_idx}-#{nvme_idx}.img,if=none,id=NVME#{ss_idx}-#{nvme_idx}"
                                        lv.qemuargs :value => "-device"
                                        lv.qemuargs :value => "nvme,drive=NVME#{ss_idx}-#{nvme_idx},serial=nvme-1234#{ss_idx}#{nvme_idx}"
                                end
                                # PMEM
                                lv.qemuargs :value => "-machine"
                                lv.qemuargs :value => "pc,accel=kvm,nvdimm=on"
                                lv.qemuargs :value => "-m"
                                lv.qemuargs :value => "16G,slots=2,maxmem=48G"
                                lv.qemuargs :value => "-object"
                                lv.qemuargs :value => "memory-backend-file,id=mem#{ss_idx},share=on,mem-path=" + ENV['HOME'] + "/tmp/nvdimm#{ss_idx},size=32768M"
                                lv.qemuargs :value => "-device"
                                lv.qemuargs :value => "nvdimm,id=nvdimm#{ss_idx},memdev=mem#{ss_idx},label-size=2097152"
                        end
                        config.vm.provision "Configure selinux", type: "shell", inline: "selinuxenabled && setenforce 0; cat >/etc/selinux/config<<__EOF
SELINUX=disabled
SELINUXTYPE=targeted
__EOF"
                        config.vm.provision "Allow ssh passwords", type: "shell", inline: "sed -i -e '/PasswordAuthentication no/s/no/yes/' /etc/ssh/sshd_config"
                        config.vm.provision "Install basic tools", type: "shell", inline: "yum -y install time"
                end
        end

        # Update ~/.ssh/config so that "ssh $vm" works without having to use vagrant ssh
        config.trigger.after :up, type: :command do |trigger|
                trigger.run = {inline: 'bash -c "set -x; ' +
                                                'echo \"# Added by ' + __FILE__ + '\" >> ' + ENV['HOME'] + '/.ssh/config; ' +
                                                'vagrant ssh-config >> ' + ENV['HOME'] + '/.ssh/config; ' +
                                                'echo \"# Added by ' + __FILE__ + '\" >> ' + ENV['HOME'] + '/.ssh/config"'}
        end
        config.trigger.before :destroy, type: :command do |trigger|
                trigger.run = {inline: 'bash -c "set -x; ed <<EOF ' + ENV['HOME'] + '/.ssh/config || true' + '
/^# Added by ' + __FILE__.gsub("/", "\\/") + '$/;/^# Added by ' + __FILE__.gsub("/", "\\/") + '$/d
.
wq
EOF"'}
        end
end
