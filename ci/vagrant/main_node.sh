#!/bin/bash

# shellcheck source=ci/common.sh
. "$WORKDIR"/ci/common.sh

DSL_REPO_var="DAOS_STACK_$(toupper "${DISTRO:-}")_LOCAL_REPO"

if ! cd "$WORKDIR"; then
    echo "Failed to chdir $WORKDIR"
    exit 1
fi

case "$1" in
    "start-vagrant")
        mkdir -p ~/.local/share/libvirt/images/ ~/tmp/
        # The following seems to need to be done or else vagrant fails to start
        sudo virsh net-list --all
        vagrant up || exit 1
    ;;
    "vagrant-status")
        vagrant status || exit 1
    ;;
    "config-vagrant-nodes")
        DAOS_STACK_EL_7_LOCAL_REPO="${!DSL_REPO_var}"                 \
        REPOSITORY_URL="$REPOSITORY_URL"                              \
        JENKINS_URL="$JENKINS_URL"                                    \
        INST_REPOS="$INST_REPOS"                                      \
        INST_RPMS="$INST_RPMS"                                        \
        GPG_KEY_URLS=""                                               \
        CONFIG_POWER_ONLY=false                                       \
        DISTRO=EL_7                                                   \
        NODESTRING="$NODESTRING"                                      \
        FOR_DAOS=true                                                 \
        ci/provisioning/post_provision_config.sh || exit 1
    ;;
    "run-tests")
        REMOTE_ACCT=${REMOTE_ACCT:-vagrant}        \
        STAGE_NAME=local_test                      \
        TEST_TAG="${TEST_TAG:-basic}"              \
        FTEST_ARG=auto                             \
        PRAGMA_SUFFIX=-vm                          \
        NODE_COUNT=${NODE_COUNT}                   \
        OPERATIONS_EMAIL=\"$OPERATIONS_EMAIL\"     \
        NODELIST="$NODESTRING"                     \
        ci/functional/test_main.sh || exit 1
    ;;
esac
