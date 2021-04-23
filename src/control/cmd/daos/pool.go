//
// (C) Copyright 2020-2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package main

import (
	"fmt"
	"io"
	"strings"
	"unsafe"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/daos-stack/daos/src/control/build"
	"github.com/daos-stack/daos/src/control/cmd/dmg/pretty"
	"github.com/daos-stack/daos/src/control/common/proto/convert"
	mgmtpb "github.com/daos-stack/daos/src/control/common/proto/mgmt"
	"github.com/daos-stack/daos/src/control/lib/control"
	"github.com/daos-stack/daos/src/control/lib/txtfmt"
)

/*
#cgo CFLAGS: -I${SRCDIR}/../../../utils
#cgo LDFLAGS: -ldaos_cmd_hdlrs

#include "daos.h"
#include "daos_hdlr.h"
*/
import "C"

type poolBaseCmd struct {
	daosCmd
	poolUUID uuid.UUID

	cPoolHandle C.daos_handle_t

	SysName string `long:"sys-name" short:"G" description:"DAOS system name"`
	Args    struct {
		Pool string `positional-arg-name:"<pool name or UUID>"`
	} `positional-args:"yes" required:"yes"`
}

func (cmd *poolBaseCmd) poolUUIDPtr() *C.uchar {
	return (*C.uchar)(unsafe.Pointer(&cmd.poolUUID[0]))
}

func (cmd *poolBaseCmd) resolvePool() (err error) {
	// TODO: Resolve name
	cmd.poolUUID, err = uuid.Parse(cmd.Args.Pool)
	if err != nil {
		return
	}

	return
}

func (cmd *poolBaseCmd) resolveAndConnect() (func(), error) {
	if err := cmd.resolvePool(); err != nil {
		return nil, err
	}

	if err := cmd.connectPool(); err != nil {
		return nil, err
	}

	return func() {
		cmd.disconnectPool()
	}, nil
}

func (cmd *poolBaseCmd) connectPool() error {
	sysName := cmd.SysName
	if sysName == "" {
		sysName = build.DefaultSystemName
	}

	cSysName := C.CString(sysName)
	defer C.free(unsafe.Pointer(cSysName))
	rc := C.daos_pool_connect(cmd.poolUUIDPtr(), cSysName,
		C.DAOS_PC_RW, &cmd.cPoolHandle, nil, nil)
	return daosError(rc)
}

func (cmd *poolBaseCmd) disconnectPool() {
	rc := C.daos_pool_disconnect(cmd.cPoolHandle, nil)
	if err := daosError(rc); err != nil {
		cmd.log.Errorf("pool disconnect failed: %s", err)
	}
}

type poolCmd struct {
	ListContainers poolContainersListCmd `command:"list-containers" alias:"list-cont" alias:"ls" description:"container operations for the specified pool"`
	Query          poolQueryCmd          `command:"query" description:"query pool info"`
	ListAttrs      poolListAttrsCmd      `command:"list-attributes" alias:"list-attrs" description:"list pool attributes"`
	GetAttr        poolGetAttrCmd        `command:"get-atttribute" alias:"get-attr" description:"get pool attribute"`
	SetAttr        poolSetAttrCmd        `command:"set-attribute" alias:"set-attr" description:"set pool attribute"`
	DelAttr        poolDelAttrCmd        `command:"delete-attribute" alias:"del-attr" description:"delete pool attribute"`
	AutoTest       poolAutoTestCmd       `command:"autotest" description:"verify setup with smoke tests"`
}

type poolContainersListCmd struct {
	poolBaseCmd
}

func (cmd *poolContainersListCmd) Execute(_ []string) error {
	extra_cont_margin := C.size_t(16)

	cleanup, err := cmd.resolveAndConnect()
	if err != nil {
		return err
	}
	defer cleanup()

	// First call gets the current number of containers.
	var ncont C.daos_size_t
	rc := C.daos_pool_list_cont(cmd.cPoolHandle, &ncont, nil, nil)
	if err := daosError(rc); err != nil {
		return errors.Wrap(err, "pool list containers failed")
	}

	// No containers.
	if ncont == 0 {
		return nil
	}

	var cConts *C.struct_daos_pool_cont_info
	// Extend ncont with a safety margin to account for containers
	// that might have been created since the first API call.
	ncont += extra_cont_margin
	cConts = (*C.struct_daos_pool_cont_info)(C.calloc(C.sizeof_struct_daos_pool_cont_info, ncont))
	if cConts == nil {
		return errors.New("malloc() failed")
	}
	defer C.free(unsafe.Pointer(cConts))

	rc = C.daos_pool_list_cont(cmd.cPoolHandle, &ncont, cConts, nil)
	if err := daosError(rc); err != nil {
		return err
	}

	conts := (*[1 << 30]C.struct_daos_pool_cont_info)(unsafe.Pointer(cConts))[:ncont:ncont]
	contUUIDs := make([]uuid.UUID, ncont)
	for i, cont := range conts {
		contUUIDs[i] = uuid.Must(uuidFromC(cont.pci_uuid))
	}

	if cmd.jsonOutputEnabled() {
		return cmd.outputJSON(contUUIDs, nil)
	}

	for _, cont := range contUUIDs {
		cmd.log.Infof("%s", cont)
	}

	return nil
}

type poolQueryCmd struct {
	poolBaseCmd
}

func convertPoolSpaceInfo(in *C.struct_daos_pool_space, mt C.uint) *mgmtpb.StorageUsageStats {
	if in == nil {
		return nil
	}

	return &mgmtpb.StorageUsageStats{
		Total: uint64(in.ps_space.s_total[mt]),
		Free:  uint64(in.ps_space.s_free[mt]),
		Min:   uint64(in.ps_free_min[mt]),
		Max:   uint64(in.ps_free_max[mt]),
		Mean:  uint64(in.ps_free_mean[mt]),
	}
}

func convertPoolRebuildStatus(in *C.struct_daos_rebuild_status) *mgmtpb.PoolRebuildStatus {
	if in == nil {
		return nil
	}

	out := &mgmtpb.PoolRebuildStatus{
		Status: int32(in.rs_errno),
	}
	if out.Status == 0 {
		out.Objects = uint64(in.rs_obj_nr)
		out.Records = uint64(in.rs_rec_nr)
		switch {
		case in.rs_version == 0:
			out.State = mgmtpb.PoolRebuildStatus_IDLE
		case in.rs_done == 1:
			out.State = mgmtpb.PoolRebuildStatus_DONE
		default:
			out.State = mgmtpb.PoolRebuildStatus_BUSY
		}
	}

	return out
}

// This is not great... But it allows us to leverage the existing
// pretty printer that dmg uses for this info. Better to find some
// way to unify all of this and remove redundancy/manual conversion.
//
// We're basically doing the same thing as ds_mgmt_drpc_pool_query()
// to stuff the info into a protobuf message and then using the
// automatic conversion from proto to control. Kind of ugly but
// gets the job done. We could potentially create some function
// that's shared between this code and the drpc handlers to deal
// with stuffing the protobuf message but it's probably overkill.
func convertPoolInfo(pinfo *C.daos_pool_info_t) (*control.PoolQueryResp, error) {
	pqp := new(mgmtpb.PoolQueryResp)

	pqp.Uuid = uuid.Must(uuidFromC(pinfo.pi_uuid)).String()
	pqp.TotalTargets = uint32(pinfo.pi_ntargets)
	pqp.DisabledTargets = uint32(pinfo.pi_ndisabled)
	pqp.ActiveTargets = uint32(pinfo.pi_space.ps_ntargets)
	pqp.TotalNodes = uint32(pinfo.pi_nnodes)
	pqp.Leader = uint32(pinfo.pi_leader)
	pqp.Version = uint32(pinfo.pi_map_ver)

	pqp.Scm = convertPoolSpaceInfo(&pinfo.pi_space, C.DAOS_MEDIA_SCM)
	pqp.Nvme = convertPoolSpaceInfo(&pinfo.pi_space, C.DAOS_MEDIA_NVME)
	pqp.Rebuild = convertPoolRebuildStatus(&pinfo.pi_rebuild_st)

	pqr := new(control.PoolQueryResp)
	return pqr, convert.Types(pqp, pqr)
}

const (
	dpiQuerySpace   = C.DPI_SPACE
	dpiQueryRebuild = C.DPI_REBUILD_STATUS
	dpiQueryAll     = C.uint64_t(^uint64(0)) // DPI_ALL is -1
)

func (cmd *poolQueryCmd) Execute(_ []string) error {
	cleanup, err := cmd.resolveAndConnect()
	if err != nil {
		return err
	}
	defer cleanup()

	pinfo := C.daos_pool_info_t{
		pi_bits: dpiQueryAll,
	}
	rc := C.daos_pool_query(cmd.cPoolHandle, nil, &pinfo, nil, nil)
	if err := daosError(rc); err != nil {
		return errors.Wrap(err, "pool query failed")
	}

	pqr, err := convertPoolInfo(&pinfo)
	if err != nil {
		return err
	}

	if cmd.jsonOutputEnabled() {
		return cmd.outputJSON(pqr, nil)
	}

	var bld strings.Builder
	if err := pretty.PrintPoolQueryResponse(pqr, &bld); err != nil {
		return err
	}

	cmd.log.Info(bld.String())

	return nil
}

type attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func printAttributes(out io.Writer, header string, attrs ...*attribute) error {
	fmt.Fprintf(out, "%s\n", header)

	if len(attrs) == 0 {
		fmt.Fprintln(out, "  No attributes found.")
		return nil
	}

	nameTitle := "Name"
	valueTitle := "Value"
	titles := []string{nameTitle}

	table := []txtfmt.TableRow{}
	for _, attr := range attrs {
		row := txtfmt.TableRow{}
		row[nameTitle] = attr.Name
		if attr.Value != "" {
			row[valueTitle] = attr.Value
			if len(titles) == 1 {
				titles = append(titles, valueTitle)
			}
		}
		table = append(table, row)
	}

	tf := txtfmt.NewTableFormatter(titles...)
	tf.InitWriter(out)
	tf.Format(table)
	return nil
}

type poolListAttrsCmd struct {
	poolBaseCmd
}

func (cmd *poolListAttrsCmd) Execute(_ []string) error {
	cleanup, err := cmd.resolveAndConnect()
	if err != nil {
		return err
	}
	defer cleanup()

	expectedSize, totalSize := C.size_t(0), C.size_t(0)
	rc := C.daos_pool_list_attr(cmd.cPoolHandle, nil, &totalSize, nil)
	if err := daosError(rc); err != nil {
		return errors.Wrapf(err, "failed to list attributes for pool %s", cmd.poolUUID)
	}

	attrs := []*attribute{}

	if totalSize > 0 {
		expectedSize = totalSize
		buf := C.malloc(totalSize)
		if buf == nil {
			return errors.New("failed to malloc buf")
		}
		defer C.free(buf)

		rc = C.daos_pool_list_attr(cmd.cPoolHandle, (*C.char)(buf), &totalSize, nil)
		if err := daosError(rc); err != nil {
			return errors.Wrapf(err, "failed to list attributes for pool %s", cmd.poolUUID)
		}

		if expectedSize < totalSize {
			cmd.log.Errorf("attribute size changed; value truncated (%d < %d)", expectedSize, totalSize)
		}

		bufSlice := (*[1 << 30]C.char)(unsafe.Pointer(buf))[:expectedSize:expectedSize]
		len := C.size_t(0)
		for cur := C.size_t(0); cur < expectedSize; cur += len + 1 {
			chunk := bufSlice[cur:]
			len = C.strnlen((*C.char)(&chunk[0]), totalSize-cur)
			if len == totalSize-cur {
				return errors.New("corrupt buffer")
			}
			chunk = bufSlice[cur:len]
			attrs = append(attrs, &attribute{
				Name: C.GoString((*C.char)(&chunk[0])),
			})
		}
	}

	if cmd.jsonOutputEnabled() {
		return cmd.outputJSON(attrs, nil)
	}

	var bld strings.Builder
	title := fmt.Sprintf("Attributes for pool %s:", cmd.poolUUID)
	if err := printAttributes(&bld, title, attrs...); err != nil {
		return err
	}

	cmd.log.Info(bld.String())

	return nil
}

type poolGetAttrCmd struct {
	poolBaseCmd

	Args struct {
		Name string `positional-arg-name:"<attribute name>"`
	} `positional-args:"yes" required:"yes"`
}

func (cmd *poolGetAttrCmd) Execute(_ []string) error {
	cleanup, err := cmd.resolveAndConnect()
	if err != nil {
		return err
	}
	defer cleanup()

	attrName := C.CString(cmd.Args.Name)
	defer C.free(unsafe.Pointer(attrName))
	expectedSize, attrSize := C.size_t(0), C.size_t(0)
	rc := C.daos_pool_get_attr(cmd.cPoolHandle, 1, &attrName, nil, &attrSize, nil)
	if err := daosError(rc); err != nil {
		return errors.Wrapf(err, "failed to get attribute %q for pool %s", cmd.Args.Name, cmd.poolUUID)
	}

	attr := &attribute{
		Name: cmd.Args.Name,
	}

	if attrSize > 0 {
		expectedSize = attrSize
		buf := C.malloc(attrSize)
		if buf == nil {
			return errors.New("failed to malloc buf")
		}
		defer C.free(buf)

		rc = C.daos_pool_get_attr(cmd.cPoolHandle, 1, &attrName, &buf, &attrSize, nil)
		if err := daosError(rc); err != nil {
			return errors.Wrapf(err, "failed to get attribute %q for pool %s", cmd.Args.Name, cmd.poolUUID)
		}

		if expectedSize < attrSize {
			cmd.log.Errorf("attribute size changed; value truncated (%d < %d)", expectedSize, attrSize)
		}

		attr.Value = C.GoString((*C.char)(buf))
	}

	if cmd.jsonOutputEnabled() {
		return cmd.outputJSON(attr, nil)
	}

	var bld strings.Builder
	title := fmt.Sprintf("Attributes for pool %s:", cmd.poolUUID)
	if err := printAttributes(&bld, title, attr); err != nil {
		return err
	}

	cmd.log.Info(bld.String())

	return nil
}

type poolSetAttrCmd struct {
	poolBaseCmd

	Args struct {
		Name  string `positional-arg-name:"<attribute name>"`
		Value string `positional-arg-name:"<attribute value>"`
	} `positional-args:"yes" required:"yes"`
}

func (cmd *poolSetAttrCmd) Execute(_ []string) error {
	cleanup, err := cmd.resolveAndConnect()
	if err != nil {
		return err
	}
	defer cleanup()

	attrName := C.CString(cmd.Args.Name)
	defer C.free(unsafe.Pointer(attrName))
	attrValue := C.CString(cmd.Args.Value)
	defer C.free(unsafe.Pointer(attrValue))
	valueLen := C.uint64_t(len(cmd.Args.Value) + 1)
	rc := C.daos_pool_set_attr(cmd.cPoolHandle, 1, &attrName, (*unsafe.Pointer)(unsafe.Pointer(&attrValue)), &valueLen, nil)
	if err := daosError(rc); err != nil {
		return errors.Wrapf(err, "failed to set attribute %q for pool %s", cmd.Args.Name, cmd.poolUUID)
	}

	if cmd.jsonOutputEnabled() {
		return cmd.outputJSON(nil, nil)
	}

	return nil
}

type poolDelAttrCmd struct {
	poolBaseCmd

	Args struct {
		Name string `positional-arg-name:"<attribute name>"`
	} `positional-args:"yes" required:"yes"`
}

func (cmd *poolDelAttrCmd) Execute(_ []string) error {
	cleanup, err := cmd.resolveAndConnect()
	if err != nil {
		return err
	}
	defer cleanup()

	attrName := C.CString(cmd.Args.Name)
	defer C.free(unsafe.Pointer(attrName))
	rc := C.daos_pool_del_attr(cmd.cPoolHandle, 1, &attrName, nil)
	if err := daosError(rc); err != nil {
		return errors.Wrapf(err, "failed to delete attribute %q for pool %s", cmd.Args.Name, cmd.poolUUID)
	}

	if cmd.jsonOutputEnabled() {
		return cmd.outputJSON(nil, nil)
	}

	return nil
}

type poolAutoTestCmd struct {
	poolBaseCmd
}

func (cmd *poolAutoTestCmd) Execute(_ []string) error {
	cleanup, err := cmd.resolveAndConnect()
	if err != nil {
		return err
	}
	defer cleanup()

	ap, deallocCmdArgs, err := allocCmdArgs(cmd.log)
	if err != nil {
		return err
	}
	defer deallocCmdArgs()

	ap.pool = cmd.cPoolHandle
	if err := copyUUID(&ap.p_uuid, cmd.poolUUID); err != nil {
		return err
	}
	ap.p_op = C.POOL_AUTOTEST

	rc := C.pool_autotest_hdlr(ap)
	if err := daosError(rc); err != nil {
		return errors.Wrapf(err, "failed to run autotest for pool %s", cmd.poolUUID)
	}

	return nil
}
