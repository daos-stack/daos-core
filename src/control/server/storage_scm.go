//
// (C) Copyright 2018-2019 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//

package main

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/daos-stack/daos/src/control/common"
	pb "github.com/daos-stack/daos/src/control/common/proto/mgmt"
	"github.com/daos-stack/daos/src/control/log"
	"github.com/daos-stack/go-ipmctl/ipmctl"
	"github.com/pkg/errors"
)

var (
	msgScmNotInited        = "scm storage not initialized"
	msgScmAlreadyFormatted = "scm storage has already been formatted and " +
		"reformat not implemented"
	msgScmMountEmpty = "scm mount must be specified in config"
	msgScmBadDevList = "expecting one scm dcpm pmem device " +
		"per-server in config"
	msgScmDevEmpty          = "scm dcpm device list must contain path"
	msgScmClassNotSupported = "operation unsupported on scm class"
	msgIpmctlDiscoverFail   = "ipmctl module discovery"
)

// ScmStorage interface specifies basic functionality for subsystem
type ScmStorage interface {
	Setup() error
	Teardown() error
	Format(int) error
	Discover() error
}

// scmStorage gives access to underlying storage interface implementation
// for accessing SCM devices (API) in addition to storage of device
// details.
//
// IpaCtl provides necessary methods to interact with Storage Class
// Memory modules through libipmctl via go-ipmctl bindings.
type scmStorage struct {
	ipmctl      ipmctl.IpmCtl  // ipmctl NVM API interface
	config      *configuration // server configuration structure
	modules     []*pb.ScmModule
	initialized bool
	formatted   bool
}

// TODO: implement remaining methods for scmStorage
// func (s *scmStorage) Update(params interface{}) interface{} {return nil}
// func (s *scmStorage) BurnIn(params interface{}) (fioPath string, cmds []string, env string, err error) {
// return
// }

// Setup placeholder implementation for scmStorage
func (s *scmStorage) Setup() error {
	resp := new(pb.ScanStorageResp)
	s.Discover(resp)

	if resp.Scmstate.Status != pb.ResponseStatus_CTRL_SUCCESS {
		return errors.New("scm scan: " + resp.Scmstate.Error)
	}

	return nil
}

// Teardown placeholder implementation for scmStorage
func (s *scmStorage) Teardown() error {
	s.initialized = false
	return nil
}

func loadModules(mms []ipmctl.DeviceDiscovery) (pbMms []*pb.ScmModule) {
	for _, c := range mms {
		pbMms = append(
			pbMms,
			&pb.ScmModule{
				Loc: &pb.ScmModule_Location{
					Channel:    uint32(c.Channel_id),
					Channelpos: uint32(c.Channel_pos),
					Memctrlr:   uint32(c.Memory_controller_id),
					Socket:     uint32(c.Socket_id),
				},
				Physicalid: uint32(c.Physical_id),
				Capacity:   c.Capacity,
			})
	}
	return
}

// Discover method implementation for scmStorage
func (s *scmStorage) Discover(resp *pb.ScanStorageResp) {
	addStateDiscover := func(
		status pb.ResponseStatus, errMsg string,
		infoMsg string) *pb.ResponseState {

		return addState(
			status, errMsg, infoMsg, common.UtilLogDepth+1,
			"scm storage discover")
	}

	if s.initialized {
		resp.Scmstate = addStateDiscover(
			pb.ResponseStatus_CTRL_SUCCESS, "", "")
		resp.Modules = s.modules
		return
	}

	mms, err := s.ipmctl.Discover()
	if err != nil {
		resp.Scmstate = addStateDiscover(
			pb.ResponseStatus_CTRL_ERR_SCM,
			msgIpmctlDiscoverFail+": "+err.Error(), "")
		return
	}
	s.modules = loadModules(mms)

	resp.Scmstate = addStateDiscover(pb.ResponseStatus_CTRL_SUCCESS, "", "")
	resp.Modules = s.modules

	s.initialized = true
}

// clearMount unmounts then removes mount point.
//
// NOTE: requires elevated privileges
func (s *scmStorage) clearMount(mntPoint string) (err error) {
	if err = s.config.ext.unmount(mntPoint); err != nil {
		return
	}

	if err = s.config.ext.remove(mntPoint); err != nil {
		return
	}

	return
}

// reFormat wipes fs signatures and formats dev with ext4.
//
// NOTE: Requires elevated privileges and is a destructive operation, prompt
//       user for confirmation before running.
func (s *scmStorage) reFormat(devPath string) (err error) {
	log.Debugf("wiping all fs identifiers on device %s", devPath)

	if err = s.config.ext.runCommand(
		fmt.Sprintf("wipefs -a %s", devPath)); err != nil {

		return errors.WithMessage(err, "wipefs")
	}

	if err = s.config.ext.runCommand(
		fmt.Sprintf("mkfs.ext4 %s", devPath)); err != nil {

		return errors.WithMessage(err, "mkfs format")
	}

	return
}

// makeMount creates a mount target directory and mounts device there.
//
// NOTE: requires elevated privileges
func (s *scmStorage) makeMount(
	devPath string, mntPoint string, devType string, mntOpts string) (err error) {

	var flags uintptr
	flags = syscall.MS_NOATIME | syscall.MS_SILENT
	flags |= syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID

	if err = s.config.ext.mkdir(mntPoint); err != nil {
		return
	}

	if err = s.config.ext.mount(
		devPath, mntPoint, devType, flags, mntOpts); err != nil {

		return
	}

	return
}

// addMret populates and adds to response SCM mount results list in addition
// to logging any err.
func addMret(
	resp *pb.FormatStorageResp, op string, mntPoint string,
	status pb.ResponseStatus, errMsg string, logDepth int) {

	resp.Mrets = append(
		resp.Mrets,
		&pb.ScmMountResult{
			Mntpoint: mntPoint,
			State: &pb.ResponseState{
				Status: status,
				Error:  errMsg,
			},
		})

	if errMsg != "" {
		log.Errordf(logDepth, "scm storage "+op+": "+errMsg)
	}
}

// Format attempts to format (forcefully) SCM mounts on a given server
// as specified in config file and populates resp ScmMountResult.
func (s *scmStorage) Format(i int, resp *pb.FormatStorageResp) {
	var devType, devPath, mntOpts string

	srv := s.config.Servers[i]
	mntPoint := srv.ScmMount

	// wraps around addMret to provide format specific function
	addMretFormat := func(status pb.ResponseStatus, errMsg string) {
		// log depth should be stack layer registering result
		addMret(
			resp, "format", mntPoint, status, errMsg,
			common.UtilLogDepth+1)
	}

	if !s.initialized {
		addMretFormat(pb.ResponseStatus_CTRL_ERR_APP, msgScmNotInited)
		return
	}

	if s.formatted {
		addMretFormat(pb.ResponseStatus_CTRL_ERR_APP, msgScmAlreadyFormatted)
		return
	}

	if mntPoint == "" {
		addMretFormat(pb.ResponseStatus_CTRL_ERR_CONF, msgScmMountEmpty)
		return
	}

	switch srv.ScmClass {
	case scmDCPM:
		devType = "ext4"
		mntOpts = "dax"

		if len(srv.ScmList) != 1 {
			addMretFormat(
				pb.ResponseStatus_CTRL_ERR_CONF, msgScmBadDevList)
			return
		}

		devPath = srv.ScmList[0]
		if devPath == "" {
			addMretFormat(
				pb.ResponseStatus_CTRL_ERR_CONF, msgScmDevEmpty)
			return
		}

		if err := s.clearMount(mntPoint); err != nil {
			addMretFormat(pb.ResponseStatus_CTRL_ERR_APP, err.Error())
			return
		}

		log.Debugf(
			"formatting scm device %s, should be quick!...", devPath)

		if err := s.reFormat(devPath); err != nil {
			addMretFormat(
				pb.ResponseStatus_CTRL_ERR_APP, err.Error())
			return
		}

		log.Debugf("format complete.\n")
	case scmRAM:
		devPath = "tmpfs"
		devType = "tmpfs"

		if err := s.clearMount(mntPoint); err != nil {
			addMretFormat(pb.ResponseStatus_CTRL_ERR_APP, err.Error())
			return
		}

		if srv.ScmSize >= 0 {
			mntOpts = "size=" + strconv.Itoa(srv.ScmSize) + "g"
			break
		}
		log.Debugf("no scm_size specified in config for ram tmpfs")
	default:
		addMretFormat(
			pb.ResponseStatus_CTRL_ERR_CONF,
			string(srv.ScmClass)+": "+msgScmClassNotSupported)
		return
	}

	log.Debugf(
		"mounting scm device %s at %s (%s)...",
		devPath, mntPoint, devType)

	if err := s.makeMount(devPath, mntPoint, devType, mntOpts); err != nil {
		addMretFormat(pb.ResponseStatus_CTRL_ERR_APP, err.Error())
		return
	}

	log.Debugf("mount complete.\n")
	addMretFormat(pb.ResponseStatus_CTRL_SUCCESS, "")

	s.formatted = true
	return
}

// newScmStorage creates a new instance of ScmStorage struct.
//
// NvmMgmt is the implementation of ipmctl interface in go-ipmctl
func newScmStorage(config *configuration) *scmStorage {
	return &scmStorage{
		ipmctl: &ipmctl.NvmMgmt{},
		config: config,
	}
}
