//
// (C) Copyright 2019 Intel Corporation.
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
	"github.com/daos-stack/daos/src/control/client"
	"github.com/daos-stack/daos/src/control/common"
	pb "github.com/daos-stack/daos/src/control/common/proto/mgmt"
	types "github.com/daos-stack/daos/src/control/common/storage"
	log "github.com/daos-stack/daos/src/control/logging"
)

// StorCmd is the struct representing the top-level storage subcommand.
type StorCmd struct {
	Prepare StoragePrepareCmd `command:"prepare" alias:"p" description:"Prepare SCM and NVMe storage attached to remote servers."`
	Scan    StorageScanCmd    `command:"scan" alias:"s" description:"Scan SCM and NVMe storage attached to remote servers."`
	Format  StorageFormatCmd  `command:"format" alias:"f" description:"Format SCM and NVMe storage attached to remote servers."`
	Update  StorageUpdateCmd  `command:"fwupdate" alias:"u" description:"Update firmware on NVMe storage attached to remote servers."`
	Query  QueryStorCmd `command:"query" alias:"q" description:"Query storage commands, including raw NVMe SSD device health stats and internal blobstore health info."`
}

func storagePrepare(conns client.Connect, req *pb.StoragePrepareReq, force bool) {
	log.Info(
		"This could be a destructive operation and storage devices " +
			"may have data erased. Please be patient as it may take several minutes " +
			"and a subsequent reboot maybe required.\n")

	if force || common.GetConsent() {
		log.Info("")
		log.Infof("NVMe & SCM preparation:\n%s", conns.StoragePrepare(req))
	}
}

// StoragePrepareCmd is the struct representing the prep storage subcommand.
type StoragePrepareCmd struct {
	broadcastCmd
	connectedCmd
	types.StoragePrepareNvmeCmd
	types.StoragePrepareScmCmd
	Force bool `short:"f" long:"force" description:"Perform format without prompting for confirmation"`
}

// Execute is run when StoragePrepareCmd activates
func (cmd *StoragePrepareCmd) Execute(args []string) error {
	storagePrepare(
		cmd.conns,
		&pb.StoragePrepareReq{
			Nvme: &pb.PrepareNvmeReq{
				Pciwhitelist: cmd.PCIWhiteList,
				Nrhugepages:  int32(cmd.NrHugepages),
				Targetuser:   cmd.TargetUser,
				Reset_:       cmd.ResetNvme,
			},
			Scm: &pb.PrepareScmReq{
				Reset_: cmd.ResetScm,
			},
		}, cmd.Force)

	return nil
}

// StorageScanCmd is the struct representing the scan storage subcommand.
type StorageScanCmd struct {
	broadcastCmd
	connectedCmd
}

// run NVMe and SCM storage query on all connected servers
func storageScan(conns client.Connect) {
	cCtrlrs, cModules := conns.StorageScan()
	log.Infof("NVMe SSD controller and constituent namespaces:\n%s", cCtrlrs)
	log.Infof("SCM modules:\n%s", cModules)
}

// Execute is run when StorageScanCmd activates
func (s *StorageScanCmd) Execute(args []string) error {
	storageScan(s.conns)
	return nil
}

// StorageFormatCmd is the struct representing the format storage subcommand.
type StorageFormatCmd struct {
	broadcastCmd
	connectedCmd
	Force bool `short:"f" long:"force" description:"Perform format without prompting for confirmation"`
}

// run NVMe and SCM storage format on all connected servers
func storageFormat(conns client.Connect, force bool) {
	log.Info(
		"This is a destructive operation and storage devices " +
			"specified in the server config file will be erased.\n" +
			"Please be patient as it may take several minutes.\n")

	if force || common.GetConsent() {
		log.Info("")
		cCtrlrResults, cMountResults := conns.StorageFormat()
		log.Infof("NVMe storage format results:\n%s", cCtrlrResults)
		log.Infof("SCM storage format results:\n%s", cMountResults)
	}
}

// Execute is run when StorageFormatCmd activates
func (s *StorageFormatCmd) Execute(args []string) error {
	storageFormat(s.conns, s.Force)
	return nil
}

// StorageUpdateCmd is the struct representing the update storage subcommand.
type StorageUpdateCmd struct {
	broadcastCmd
	connectedCmd
	Force        bool   `short:"f" long:"force" description:"Perform update without prompting for confirmation"`
	NVMeModel    string `short:"m" long:"nvme-model" description:"Only update firmware on NVMe SSDs with this model name/number." required:"1"`
	NVMeStartRev string `short:"r" long:"nvme-fw-rev" description:"Only update firmware on NVMe SSDs currently running this firmware revision." required:"1"`
	NVMeFwPath   string `short:"p" long:"nvme-fw-path" description:"Update firmware on NVMe SSDs with image file at this path (path must be accessible on all servers)." required:"1"`
	NVMeFwSlot   int    `short:"s" default:"0" long:"nvme-fw-slot" description:"Update firmware on NVMe SSDs to this firmware register."`
}

// run NVMe and SCM storage update on all connected servers
func storageUpdate(conns client.Connect, req *pb.StorageUpdateReq, force bool) {
	log.Info(
		"This could be a destructive operation and storage devices " +
			"specified in the server config file will have firmware " +
			"updated. Please check this is a supported upgrade path " +
			"and be patient as it may take several minutes.\n")

	if force || common.GetConsent() {
		log.Info("")
		cCtrlrResults, cModuleResults := conns.StorageUpdate(req)
		log.Infof("NVMe storage update results:\n%s", cCtrlrResults)
		log.Infof("SCM storage update results:\n%s", cModuleResults)
	}
}

// Execute is run when StorageUpdateCmd activates
func (u *StorageUpdateCmd) Execute(args []string) error {
	// only populate nvme fwupdate params for the moment
	storageUpdate(
		u.conns,
		&pb.StorageUpdateReq{
			Nvme: &pb.UpdateNvmeReq{
				Model: u.NVMeModel, Startrev: u.NVMeStartRev,
				Path: u.NVMeFwPath, Slot: int32(u.NVMeFwSlot),
			},
		}, u.Force)

	return nil
}

// TODO: implement burn-in subcommand

//func getFioConfig(c *ishell.Context) (configPath string, err error) {
//	// fetch existing configuration files
//	paths, err := mgmtClient.FetchFioConfigPaths()
//	if err != nil {
//		return
//	}
//	// cut prefix to display filenames not full path
//	configChoices := functional.Map(
//		paths, func(p string) string { return filepath.Base(p) })
//	// add custom path option
//	configChoices = append(configChoices, "custom path")
//
//	choiceIdx := c.MultiChoice(
//		configChoices,
//		"Select the fio configuration you would like to run on the selected controllers.")
//
//	// if custom path selected (last index), process input
//	if choiceIdx == len(configChoices)-1 {
//		// disable the '>>>' for cleaner same line input.
//		c.ShowPrompt(false)
//		defer c.ShowPrompt(true) // revert after user input.
//
//		c.Print("Please enter fio configuration file-path [has file extension .fio]: ")
//		path := c.ReadLine()
//
//		if path == "" {
//			return "", fmt.Errorf("no filepath provided")
//		}
//		if filepath.Ext(path) != ".fio" {
//			return "", fmt.Errorf("provided filepath does not end in .fio")
//		}
//		if !filepath.IsAbs(path) {
//			return "", fmt.Errorf("provided filepath is not absolute")
//		}
//
//		configPath = path
//	} else {
//		configPath = paths[choiceIdx]
//	}
//
//	return
//}

//func nvmeTaskLookup(
//	c *ishell.Context, ctrlrs []*pb.NvmeController, feature string) error {
//
//	switch feature {
//	case "nvme-fw-update":
//		params, err := getUpdateParams(c)
//		if err != nil {
//			c.Println("Problem reading user inputs: ", err)
//			return err
//		}
//
//		for _, ctrlr := range ctrlrs {
//			c.Printf("\nController: %+v\n", ctrlr)
//			c.Printf(
//				"\t- Updating firmware on slot %d with image %s.\n",
//				params.Slot, params.Path)
//
//			params.Ctrlr = ctrlr
//
//			newFwrev, err := mgmtClient.UpdateNvmeCtrlr(params)
//			if err != nil {
//				c.Println("\nProblem updating firmware: ", err)
//				return err
//			}
//			c.Printf(
//				"\nSuccessfully updated firmware from revision %s to %s!\n",
//				params.Ctrlr.Fwrev, newFwrev)
//		}
//	case "nvme-burn-in":
//		configPath, err := getFioConfig(c)
//		if err != nil {
//			c.Println("Problem reading user inputs: ", err)
//			return err
//		}
//
//		for _, ctrlr := range ctrlrs {
//			c.Printf("\nController: %+v\n", ctrlr)
//			c.Printf(
//				"\t- Running burn-in validation with spdk fio plugin using job file %s.\n\n",
//				filepath.Base(configPath))
//			_, err := mgmtClient.BurnInNvme(ctrlr.Id, configPath)
//			if err != nil {
//				return err
//			}
//		}
//	default:
//		c.Printf("Sorry, task '%s' has not been implemented.\n", feature)
//	}
//
//	return nil
//}
