//
// (C) Copyright 2019-2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.5.0
// source: ctl/storage_scm.proto

package ctl

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// ScmModule represent Storage Class Memory modules installed.
type ScmModule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Channelid        uint32 `protobuf:"varint,1,opt,name=channelid,proto3" json:"channelid,omitempty"`              // The channel id where module is installed.
	Channelposition  uint32 `protobuf:"varint,2,opt,name=channelposition,proto3" json:"channelposition,omitempty"`  // The channel position where module is installed.
	Controllerid     uint32 `protobuf:"varint,3,opt,name=controllerid,proto3" json:"controllerid,omitempty"`        // The memory controller id attached to module.
	Socketid         uint32 `protobuf:"varint,4,opt,name=socketid,proto3" json:"socketid,omitempty"`                // The socket id attached to module.
	Physicalid       uint32 `protobuf:"varint,5,opt,name=physicalid,proto3" json:"physicalid,omitempty"`            // The physical id of the module.
	Capacity         uint64 `protobuf:"varint,6,opt,name=capacity,proto3" json:"capacity,omitempty"`                // The capacity of the module.
	Uid              string `protobuf:"bytes,7,opt,name=uid,proto3" json:"uid,omitempty"`                           // The uid of the module.
	PartNumber       string `protobuf:"bytes,8,opt,name=partNumber,proto3" json:"partNumber,omitempty"`             // The part number of the module.
	FirmwareRevision string `protobuf:"bytes,9,opt,name=firmwareRevision,proto3" json:"firmwareRevision,omitempty"` // Module's active firmware revision
}

func (x *ScmModule) Reset() {
	*x = ScmModule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScmModule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScmModule) ProtoMessage() {}

func (x *ScmModule) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScmModule.ProtoReflect.Descriptor instead.
func (*ScmModule) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{0}
}

func (x *ScmModule) GetChannelid() uint32 {
	if x != nil {
		return x.Channelid
	}
	return 0
}

func (x *ScmModule) GetChannelposition() uint32 {
	if x != nil {
		return x.Channelposition
	}
	return 0
}

func (x *ScmModule) GetControllerid() uint32 {
	if x != nil {
		return x.Controllerid
	}
	return 0
}

func (x *ScmModule) GetSocketid() uint32 {
	if x != nil {
		return x.Socketid
	}
	return 0
}

func (x *ScmModule) GetPhysicalid() uint32 {
	if x != nil {
		return x.Physicalid
	}
	return 0
}

func (x *ScmModule) GetCapacity() uint64 {
	if x != nil {
		return x.Capacity
	}
	return 0
}

func (x *ScmModule) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *ScmModule) GetPartNumber() string {
	if x != nil {
		return x.PartNumber
	}
	return ""
}

func (x *ScmModule) GetFirmwareRevision() string {
	if x != nil {
		return x.FirmwareRevision
	}
	return ""
}

// ScmNamespace represents SCM namespace as pmem device files created on a ScmRegion.
type ScmNamespace struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid     string              `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Blockdev string              `protobuf:"bytes,2,opt,name=blockdev,proto3" json:"blockdev,omitempty"`
	Dev      string              `protobuf:"bytes,3,opt,name=dev,proto3" json:"dev,omitempty"` // ndctl specific device identifier
	NumaNode uint32              `protobuf:"varint,4,opt,name=numa_node,json=numaNode,proto3" json:"numa_node,omitempty"`
	Size     uint64              `protobuf:"varint,5,opt,name=size,proto3" json:"size,omitempty"`  // pmem block device capacity in bytes
	Mount    *ScmNamespace_Mount `protobuf:"bytes,6,opt,name=mount,proto3" json:"mount,omitempty"` // mount OS info
}

func (x *ScmNamespace) Reset() {
	*x = ScmNamespace{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScmNamespace) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScmNamespace) ProtoMessage() {}

func (x *ScmNamespace) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScmNamespace.ProtoReflect.Descriptor instead.
func (*ScmNamespace) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{1}
}

func (x *ScmNamespace) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *ScmNamespace) GetBlockdev() string {
	if x != nil {
		return x.Blockdev
	}
	return ""
}

func (x *ScmNamespace) GetDev() string {
	if x != nil {
		return x.Dev
	}
	return ""
}

func (x *ScmNamespace) GetNumaNode() uint32 {
	if x != nil {
		return x.NumaNode
	}
	return 0
}

func (x *ScmNamespace) GetSize() uint64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ScmNamespace) GetMount() *ScmNamespace_Mount {
	if x != nil {
		return x.Mount
	}
	return nil
}

// ScmModuleResult represents operation state for specific SCM/PM module.
//
// TODO: replace identifier with serial when returned in scan
type ScmModuleResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Physicalid uint32         `protobuf:"varint,1,opt,name=physicalid,proto3" json:"physicalid,omitempty"` // SCM module identifier
	State      *ResponseState `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`            // state of current operation
}

func (x *ScmModuleResult) Reset() {
	*x = ScmModuleResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScmModuleResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScmModuleResult) ProtoMessage() {}

func (x *ScmModuleResult) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScmModuleResult.ProtoReflect.Descriptor instead.
func (*ScmModuleResult) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{2}
}

func (x *ScmModuleResult) GetPhysicalid() uint32 {
	if x != nil {
		return x.Physicalid
	}
	return 0
}

func (x *ScmModuleResult) GetState() *ResponseState {
	if x != nil {
		return x.State
	}
	return nil
}

// ScmMountResult represents operation state for specific SCM mount point.
type ScmMountResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mntpoint    string         `protobuf:"bytes,1,opt,name=mntpoint,proto3" json:"mntpoint,omitempty"`        // Path where device or tmpfs is mounted
	State       *ResponseState `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`              // state of current operation
	Instanceidx uint32         `protobuf:"varint,3,opt,name=instanceidx,proto3" json:"instanceidx,omitempty"` // Index of I/O Engine instance
}

func (x *ScmMountResult) Reset() {
	*x = ScmMountResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScmMountResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScmMountResult) ProtoMessage() {}

func (x *ScmMountResult) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScmMountResult.ProtoReflect.Descriptor instead.
func (*ScmMountResult) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{3}
}

func (x *ScmMountResult) GetMntpoint() string {
	if x != nil {
		return x.Mntpoint
	}
	return ""
}

func (x *ScmMountResult) GetState() *ResponseState {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *ScmMountResult) GetInstanceidx() uint32 {
	if x != nil {
		return x.Instanceidx
	}
	return 0
}

type PrepareScmReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reset_ bool `protobuf:"varint,1,opt,name=reset,proto3" json:"reset,omitempty"` // Reset SCM devices to memory mode
}

func (x *PrepareScmReq) Reset() {
	*x = PrepareScmReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareScmReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareScmReq) ProtoMessage() {}

func (x *PrepareScmReq) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareScmReq.ProtoReflect.Descriptor instead.
func (*PrepareScmReq) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{4}
}

func (x *PrepareScmReq) GetReset_() bool {
	if x != nil {
		return x.Reset_
	}
	return false
}

type PrepareScmResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespaces     []*ScmNamespace `protobuf:"bytes,1,rep,name=namespaces,proto3" json:"namespaces,omitempty"` // Existing namespace devices (new and old)
	State          *ResponseState  `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
	Rebootrequired bool            `protobuf:"varint,3,opt,name=rebootrequired,proto3" json:"rebootrequired,omitempty"`
}

func (x *PrepareScmResp) Reset() {
	*x = PrepareScmResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareScmResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareScmResp) ProtoMessage() {}

func (x *PrepareScmResp) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareScmResp.ProtoReflect.Descriptor instead.
func (*PrepareScmResp) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{5}
}

func (x *PrepareScmResp) GetNamespaces() []*ScmNamespace {
	if x != nil {
		return x.Namespaces
	}
	return nil
}

func (x *PrepareScmResp) GetState() *ResponseState {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *PrepareScmResp) GetRebootrequired() bool {
	if x != nil {
		return x.Rebootrequired
	}
	return false
}

type ScanScmReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Usage bool `protobuf:"varint,1,opt,name=usage,proto3" json:"usage,omitempty"` // Populate usage statistics in scan
}

func (x *ScanScmReq) Reset() {
	*x = ScanScmReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanScmReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanScmReq) ProtoMessage() {}

func (x *ScanScmReq) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanScmReq.ProtoReflect.Descriptor instead.
func (*ScanScmReq) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{6}
}

func (x *ScanScmReq) GetUsage() bool {
	if x != nil {
		return x.Usage
	}
	return false
}

type ScanScmResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Modules    []*ScmModule    `protobuf:"bytes,1,rep,name=modules,proto3" json:"modules,omitempty"`
	Namespaces []*ScmNamespace `protobuf:"bytes,2,rep,name=namespaces,proto3" json:"namespaces,omitempty"`
	State      *ResponseState  `protobuf:"bytes,3,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *ScanScmResp) Reset() {
	*x = ScanScmResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanScmResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanScmResp) ProtoMessage() {}

func (x *ScanScmResp) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanScmResp.ProtoReflect.Descriptor instead.
func (*ScanScmResp) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{7}
}

func (x *ScanScmResp) GetModules() []*ScmModule {
	if x != nil {
		return x.Modules
	}
	return nil
}

func (x *ScanScmResp) GetNamespaces() []*ScmNamespace {
	if x != nil {
		return x.Namespaces
	}
	return nil
}

func (x *ScanScmResp) GetState() *ResponseState {
	if x != nil {
		return x.State
	}
	return nil
}

type FormatScmReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FormatScmReq) Reset() {
	*x = FormatScmReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FormatScmReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FormatScmReq) ProtoMessage() {}

func (x *FormatScmReq) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FormatScmReq.ProtoReflect.Descriptor instead.
func (*FormatScmReq) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{8}
}

// Mount represents a mounted pmem block device.
type ScmNamespace_Mount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path       string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	TotalBytes uint64 `protobuf:"varint,2,opt,name=total_bytes,json=totalBytes,proto3" json:"total_bytes,omitempty"`
	AvailBytes uint64 `protobuf:"varint,3,opt,name=avail_bytes,json=availBytes,proto3" json:"avail_bytes,omitempty"`
}

func (x *ScmNamespace_Mount) Reset() {
	*x = ScmNamespace_Mount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ctl_storage_scm_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScmNamespace_Mount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScmNamespace_Mount) ProtoMessage() {}

func (x *ScmNamespace_Mount) ProtoReflect() protoreflect.Message {
	mi := &file_ctl_storage_scm_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScmNamespace_Mount.ProtoReflect.Descriptor instead.
func (*ScmNamespace_Mount) Descriptor() ([]byte, []int) {
	return file_ctl_storage_scm_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ScmNamespace_Mount) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *ScmNamespace_Mount) GetTotalBytes() uint64 {
	if x != nil {
		return x.TotalBytes
	}
	return 0
}

func (x *ScmNamespace_Mount) GetAvailBytes() uint64 {
	if x != nil {
		return x.AvailBytes
	}
	return 0
}

var File_ctl_storage_scm_proto protoreflect.FileDescriptor

var file_ctl_storage_scm_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x74, 0x6c, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x63,
	0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x63, 0x74, 0x6c, 0x1a, 0x10, 0x63, 0x74,
	0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xad,
	0x02, 0x0a, 0x09, 0x53, 0x63, 0x6d, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x69, 0x64, 0x12, 0x28, 0x0a, 0x0f, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c,
	0x65, 0x72, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x6f, 0x63, 0x6b,
	0x65, 0x74, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x73, 0x6f, 0x63, 0x6b,
	0x65, 0x74, 0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c,
	0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63,
	0x61, 0x6c, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x63, 0x61, 0x70, 0x61, 0x63, 0x69, 0x74, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75,
	0x69, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x12, 0x2a, 0x0a, 0x10, 0x66, 0x69, 0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x52, 0x65,
	0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x66, 0x69,
	0x72, 0x6d, 0x77, 0x61, 0x72, 0x65, 0x52, 0x65, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x8f,
	0x02, 0x0a, 0x0c, 0x53, 0x63, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75,
	0x75, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x64, 0x65, 0x76, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x64, 0x65, 0x76, 0x12,
	0x10, 0x0a, 0x03, 0x64, 0x65, 0x76, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x64, 0x65,
	0x76, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x75, 0x6d, 0x61, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x6e, 0x75, 0x6d, 0x61, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x12, 0x2d, 0x0a, 0x05, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x63, 0x74, 0x6c, 0x2e, 0x53, 0x63, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x2e, 0x4d, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x05, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x1a, 0x5d, 0x0a, 0x05, 0x4d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x1f,
	0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12,
	0x1f, 0x0a, 0x0b, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73,
	0x22, 0x5b, 0x0a, 0x0f, 0x53, 0x63, 0x6d, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61, 0x6c, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x70, 0x68, 0x79, 0x73, 0x69, 0x63, 0x61,
	0x6c, 0x69, 0x64, 0x12, 0x28, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x74, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x78, 0x0a,
	0x0e, 0x53, 0x63, 0x6d, 0x4d, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x6d, 0x6e, 0x74, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6d, 0x6e, 0x74, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x74, 0x6c,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x69, 0x64, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x69, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x69, 0x64, 0x78, 0x22, 0x25, 0x0a, 0x0d, 0x50, 0x72, 0x65, 0x70, 0x61,
	0x72, 0x65, 0x53, 0x63, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x73, 0x65,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x72, 0x65, 0x73, 0x65, 0x74, 0x22, 0x95,
	0x01, 0x0a, 0x0e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x53, 0x63, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x31, 0x0a, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x74, 0x6c, 0x2e, 0x53, 0x63, 0x6d, 0x4e,
	0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x74, 0x6c, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x26,
	0x0a, 0x0e, 0x72, 0x65, 0x62, 0x6f, 0x6f, 0x74, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x72, 0x65, 0x62, 0x6f, 0x6f, 0x74, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x22, 0x0a, 0x0a, 0x53, 0x63, 0x61, 0x6e, 0x53, 0x63,
	0x6d, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x22, 0x94, 0x01, 0x0a, 0x0b, 0x53,
	0x63, 0x61, 0x6e, 0x53, 0x63, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x12, 0x28, 0x0a, 0x07, 0x6d, 0x6f,
	0x64, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x74,
	0x6c, 0x2e, 0x53, 0x63, 0x6d, 0x4d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x07, 0x6d, 0x6f, 0x64,
	0x75, 0x6c, 0x65, 0x73, 0x12, 0x31, 0x0a, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x74, 0x6c, 0x2e, 0x53,
	0x63, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x0a, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x12, 0x28, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x74, 0x6c, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x22, 0x0e, 0x0a, 0x0c, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x53, 0x63, 0x6d, 0x52, 0x65,
	0x71, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x64, 0x61, 0x6f, 0x73, 0x2d, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2f, 0x64, 0x61, 0x6f, 0x73, 0x2f,
	0x73, 0x72, 0x63, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x74, 0x6c, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ctl_storage_scm_proto_rawDescOnce sync.Once
	file_ctl_storage_scm_proto_rawDescData = file_ctl_storage_scm_proto_rawDesc
)

func file_ctl_storage_scm_proto_rawDescGZIP() []byte {
	file_ctl_storage_scm_proto_rawDescOnce.Do(func() {
		file_ctl_storage_scm_proto_rawDescData = protoimpl.X.CompressGZIP(file_ctl_storage_scm_proto_rawDescData)
	})
	return file_ctl_storage_scm_proto_rawDescData
}

var file_ctl_storage_scm_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_ctl_storage_scm_proto_goTypes = []interface{}{
	(*ScmModule)(nil),          // 0: ctl.ScmModule
	(*ScmNamespace)(nil),       // 1: ctl.ScmNamespace
	(*ScmModuleResult)(nil),    // 2: ctl.ScmModuleResult
	(*ScmMountResult)(nil),     // 3: ctl.ScmMountResult
	(*PrepareScmReq)(nil),      // 4: ctl.PrepareScmReq
	(*PrepareScmResp)(nil),     // 5: ctl.PrepareScmResp
	(*ScanScmReq)(nil),         // 6: ctl.ScanScmReq
	(*ScanScmResp)(nil),        // 7: ctl.ScanScmResp
	(*FormatScmReq)(nil),       // 8: ctl.FormatScmReq
	(*ScmNamespace_Mount)(nil), // 9: ctl.ScmNamespace.Mount
	(*ResponseState)(nil),      // 10: ctl.ResponseState
}
var file_ctl_storage_scm_proto_depIdxs = []int32{
	9,  // 0: ctl.ScmNamespace.mount:type_name -> ctl.ScmNamespace.Mount
	10, // 1: ctl.ScmModuleResult.state:type_name -> ctl.ResponseState
	10, // 2: ctl.ScmMountResult.state:type_name -> ctl.ResponseState
	1,  // 3: ctl.PrepareScmResp.namespaces:type_name -> ctl.ScmNamespace
	10, // 4: ctl.PrepareScmResp.state:type_name -> ctl.ResponseState
	0,  // 5: ctl.ScanScmResp.modules:type_name -> ctl.ScmModule
	1,  // 6: ctl.ScanScmResp.namespaces:type_name -> ctl.ScmNamespace
	10, // 7: ctl.ScanScmResp.state:type_name -> ctl.ResponseState
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_ctl_storage_scm_proto_init() }
func file_ctl_storage_scm_proto_init() {
	if File_ctl_storage_scm_proto != nil {
		return
	}
	file_ctl_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ctl_storage_scm_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScmModule); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ctl_storage_scm_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScmNamespace); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ctl_storage_scm_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScmModuleResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ctl_storage_scm_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScmMountResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ctl_storage_scm_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareScmReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ctl_storage_scm_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareScmResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ctl_storage_scm_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanScmReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ctl_storage_scm_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanScmResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ctl_storage_scm_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FormatScmReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ctl_storage_scm_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScmNamespace_Mount); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ctl_storage_scm_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ctl_storage_scm_proto_goTypes,
		DependencyIndexes: file_ctl_storage_scm_proto_depIdxs,
		MessageInfos:      file_ctl_storage_scm_proto_msgTypes,
	}.Build()
	File_ctl_storage_scm_proto = out.File
	file_ctl_storage_scm_proto_rawDesc = nil
	file_ctl_storage_scm_proto_goTypes = nil
	file_ctl_storage_scm_proto_depIdxs = nil
}
