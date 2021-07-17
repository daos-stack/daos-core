//
// (C) Copyright 2020-2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

// This file defines RAS event related protobuf messages communicated over dRPC
// and gRPC.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: shared/event.proto

package shared

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

// RASEvent describes a RAS event in the DAOS system.
type RASEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                              // Unique event identifier, 64-char.
	Msg       string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`                             // Human readable message describing event.
	Timestamp string `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`                 // Fully qualified timestamp (us) incl timezone.
	Type      uint32 `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"`                          // Event type.
	Severity  uint32 `protobuf:"varint,5,opt,name=severity,proto3" json:"severity,omitempty"`                  // Event severity.
	Hostname  string `protobuf:"bytes,6,opt,name=hostname,proto3" json:"hostname,omitempty"`                   // (optional) Hostname of node involved in event.
	Rank      uint32 `protobuf:"varint,7,opt,name=rank,proto3" json:"rank,omitempty"`                          // (optional) DAOS rank involved in event.
	RankInc   uint64 `protobuf:"varint,8,opt,name=rank_inc,json=rankInc,proto3" json:"rank_inc,omitempty"`     // (optional) Rank incarnation number.
	HwId      string `protobuf:"bytes,9,opt,name=hw_id,json=hwId,proto3" json:"hw_id,omitempty"`               // (optional) Hardware component involved in event.
	ProcId    uint64 `protobuf:"varint,10,opt,name=proc_id,json=procId,proto3" json:"proc_id,omitempty"`       // (optional) Process involved in event.
	ThreadId  uint64 `protobuf:"varint,11,opt,name=thread_id,json=threadId,proto3" json:"thread_id,omitempty"` // (optional) Thread involved in event.
	JobId     string `protobuf:"bytes,12,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`           // (optional) Job involved in event.
	PoolUuid  string `protobuf:"bytes,13,opt,name=pool_uuid,json=poolUuid,proto3" json:"pool_uuid,omitempty"`  // (optional) Pool UUID involved in event.
	ContUuid  string `protobuf:"bytes,14,opt,name=cont_uuid,json=contUuid,proto3" json:"cont_uuid,omitempty"`  // (optional) Container UUID involved in event.
	ObjId     string `protobuf:"bytes,15,opt,name=obj_id,json=objId,proto3" json:"obj_id,omitempty"`           // (optional) Object involved in event.
	CtlOp     string `protobuf:"bytes,16,opt,name=ctl_op,json=ctlOp,proto3" json:"ctl_op,omitempty"`           // (optional) Recommended automatic action.
	// Types that are assignable to ExtendedInfo:
	//	*RASEvent_StrInfo
	//	*RASEvent_EngineStateInfo
	//	*RASEvent_PoolSvcInfo
	ExtendedInfo isRASEvent_ExtendedInfo `protobuf_oneof:"extended_info"`
}

func (x *RASEvent) Reset() {
	*x = RASEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RASEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RASEvent) ProtoMessage() {}

func (x *RASEvent) ProtoReflect() protoreflect.Message {
	mi := &file_shared_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RASEvent.ProtoReflect.Descriptor instead.
func (*RASEvent) Descriptor() ([]byte, []int) {
	return file_shared_event_proto_rawDescGZIP(), []int{0}
}

func (x *RASEvent) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RASEvent) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *RASEvent) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *RASEvent) GetType() uint32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *RASEvent) GetSeverity() uint32 {
	if x != nil {
		return x.Severity
	}
	return 0
}

func (x *RASEvent) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *RASEvent) GetRank() uint32 {
	if x != nil {
		return x.Rank
	}
	return 0
}

func (x *RASEvent) GetRankInc() uint64 {
	if x != nil {
		return x.RankInc
	}
	return 0
}

func (x *RASEvent) GetHwId() string {
	if x != nil {
		return x.HwId
	}
	return ""
}

func (x *RASEvent) GetProcId() uint64 {
	if x != nil {
		return x.ProcId
	}
	return 0
}

func (x *RASEvent) GetThreadId() uint64 {
	if x != nil {
		return x.ThreadId
	}
	return 0
}

func (x *RASEvent) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

func (x *RASEvent) GetPoolUuid() string {
	if x != nil {
		return x.PoolUuid
	}
	return ""
}

func (x *RASEvent) GetContUuid() string {
	if x != nil {
		return x.ContUuid
	}
	return ""
}

func (x *RASEvent) GetObjId() string {
	if x != nil {
		return x.ObjId
	}
	return ""
}

func (x *RASEvent) GetCtlOp() string {
	if x != nil {
		return x.CtlOp
	}
	return ""
}

func (m *RASEvent) GetExtendedInfo() isRASEvent_ExtendedInfo {
	if m != nil {
		return m.ExtendedInfo
	}
	return nil
}

func (x *RASEvent) GetStrInfo() string {
	if x, ok := x.GetExtendedInfo().(*RASEvent_StrInfo); ok {
		return x.StrInfo
	}
	return ""
}

func (x *RASEvent) GetEngineStateInfo() *RASEvent_EngineStateEventInfo {
	if x, ok := x.GetExtendedInfo().(*RASEvent_EngineStateInfo); ok {
		return x.EngineStateInfo
	}
	return nil
}

func (x *RASEvent) GetPoolSvcInfo() *RASEvent_PoolSvcEventInfo {
	if x, ok := x.GetExtendedInfo().(*RASEvent_PoolSvcInfo); ok {
		return x.PoolSvcInfo
	}
	return nil
}

type isRASEvent_ExtendedInfo interface {
	isRASEvent_ExtendedInfo()
}

type RASEvent_StrInfo struct {
	StrInfo string `protobuf:"bytes,17,opt,name=str_info,json=strInfo,proto3,oneof"` // Opaque data blob.
}

type RASEvent_EngineStateInfo struct {
	EngineStateInfo *RASEvent_EngineStateEventInfo `protobuf:"bytes,18,opt,name=engine_state_info,json=engineStateInfo,proto3,oneof"`
}

type RASEvent_PoolSvcInfo struct {
	PoolSvcInfo *RASEvent_PoolSvcEventInfo `protobuf:"bytes,19,opt,name=pool_svc_info,json=poolSvcInfo,proto3,oneof"`
}

func (*RASEvent_StrInfo) isRASEvent_ExtendedInfo() {}

func (*RASEvent_EngineStateInfo) isRASEvent_ExtendedInfo() {}

func (*RASEvent_PoolSvcInfo) isRASEvent_ExtendedInfo() {}

// ClusterEventReq communicates occurrence of a RAS event in the DAOS system.
type ClusterEventReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sequence uint64    `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"` // Sequence identifier for RAS events.
	Event    *RASEvent `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`        // RAS event.
}

func (x *ClusterEventReq) Reset() {
	*x = ClusterEventReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterEventReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterEventReq) ProtoMessage() {}

func (x *ClusterEventReq) ProtoReflect() protoreflect.Message {
	mi := &file_shared_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterEventReq.ProtoReflect.Descriptor instead.
func (*ClusterEventReq) Descriptor() ([]byte, []int) {
	return file_shared_event_proto_rawDescGZIP(), []int{1}
}

func (x *ClusterEventReq) GetSequence() uint64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *ClusterEventReq) GetEvent() *RASEvent {
	if x != nil {
		return x.Event
	}
	return nil
}

// ClusterEventResp acknowledges receipt of an event notification.
type ClusterEventResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sequence uint64 `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"` // Sequence identifier for RAS events.
	Status   int32  `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`     // DAOS error code.
}

func (x *ClusterEventResp) Reset() {
	*x = ClusterEventResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterEventResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterEventResp) ProtoMessage() {}

func (x *ClusterEventResp) ProtoReflect() protoreflect.Message {
	mi := &file_shared_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterEventResp.ProtoReflect.Descriptor instead.
func (*ClusterEventResp) Descriptor() ([]byte, []int) {
	return file_shared_event_proto_rawDescGZIP(), []int{2}
}

func (x *ClusterEventResp) GetSequence() uint64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *ClusterEventResp) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

// EngineStateEventInfo defines extended fields for state change events.
type RASEvent_EngineStateEventInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Instance uint32 `protobuf:"varint,1,opt,name=instance,proto3" json:"instance,omitempty"` // Control-plane harness instance index.
	Errored  bool   `protobuf:"varint,2,opt,name=errored,proto3" json:"errored,omitempty"`   // Rank in error state.
	Error    string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`        // Message associated with error.
}

func (x *RASEvent_EngineStateEventInfo) Reset() {
	*x = RASEvent_EngineStateEventInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RASEvent_EngineStateEventInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RASEvent_EngineStateEventInfo) ProtoMessage() {}

func (x *RASEvent_EngineStateEventInfo) ProtoReflect() protoreflect.Message {
	mi := &file_shared_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RASEvent_EngineStateEventInfo.ProtoReflect.Descriptor instead.
func (*RASEvent_EngineStateEventInfo) Descriptor() ([]byte, []int) {
	return file_shared_event_proto_rawDescGZIP(), []int{0, 0}
}

func (x *RASEvent_EngineStateEventInfo) GetInstance() uint32 {
	if x != nil {
		return x.Instance
	}
	return 0
}

func (x *RASEvent_EngineStateEventInfo) GetErrored() bool {
	if x != nil {
		return x.Errored
	}
	return false
}

func (x *RASEvent_EngineStateEventInfo) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// PoolSvcEventInfo defines extended fields for pool service change events.
type RASEvent_PoolSvcEventInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SvcReps []uint32 `protobuf:"varint,1,rep,packed,name=svc_reps,json=svcReps,proto3" json:"svc_reps,omitempty"` // Pool service replica ranks.
	Version uint64   `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`                       // Raft leadership term.
}

func (x *RASEvent_PoolSvcEventInfo) Reset() {
	*x = RASEvent_PoolSvcEventInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_event_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RASEvent_PoolSvcEventInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RASEvent_PoolSvcEventInfo) ProtoMessage() {}

func (x *RASEvent_PoolSvcEventInfo) ProtoReflect() protoreflect.Message {
	mi := &file_shared_event_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RASEvent_PoolSvcEventInfo.ProtoReflect.Descriptor instead.
func (*RASEvent_PoolSvcEventInfo) Descriptor() ([]byte, []int) {
	return file_shared_event_proto_rawDescGZIP(), []int{0, 1}
}

func (x *RASEvent_PoolSvcEventInfo) GetSvcReps() []uint32 {
	if x != nil {
		return x.SvcReps
	}
	return nil
}

func (x *RASEvent_PoolSvcEventInfo) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

var File_shared_event_proto protoreflect.FileDescriptor

var file_shared_event_proto_rawDesc = []byte{
	0x0a, 0x12, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x22, 0x88, 0x06, 0x0a,
	0x08, 0x52, 0x41, 0x53, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x08, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73,
	0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73,
	0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x61, 0x6e, 0x6b, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x04, 0x72, 0x61, 0x6e, 0x6b, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x61, 0x6e,
	0x6b, 0x5f, 0x69, 0x6e, 0x63, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72, 0x61, 0x6e,
	0x6b, 0x49, 0x6e, 0x63, 0x12, 0x13, 0x0a, 0x05, 0x68, 0x77, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x77, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x72, 0x6f,
	0x63, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x70, 0x72, 0x6f, 0x63,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x68, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x69, 0x64, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x74, 0x68, 0x72, 0x65, 0x61, 0x64, 0x49, 0x64, 0x12,
	0x15, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6a, 0x6f, 0x62, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6f, 0x6f, 0x6c, 0x5f, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6f, 0x6f, 0x6c, 0x55,
	0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x5f, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x55, 0x75, 0x69, 0x64,
	0x12, 0x15, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x5f, 0x69, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x6f, 0x62, 0x6a, 0x49, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x63, 0x74, 0x6c, 0x5f, 0x6f,
	0x70, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x74, 0x6c, 0x4f, 0x70, 0x12, 0x1b,
	0x0a, 0x08, 0x73, 0x74, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x11, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x00, 0x52, 0x07, 0x73, 0x74, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x53, 0x0a, 0x11, 0x65,
	0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x6e, 0x66, 0x6f,
	0x18, 0x12, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e,
	0x52, 0x41, 0x53, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52,
	0x0f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x47, 0x0a, 0x0d, 0x70, 0x6f, 0x6f, 0x6c, 0x5f, 0x73, 0x76, 0x63, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64,
	0x2e, 0x52, 0x41, 0x53, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x6f, 0x6f, 0x6c, 0x53, 0x76,
	0x63, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x48, 0x00, 0x52, 0x0b, 0x70, 0x6f,
	0x6f, 0x6c, 0x53, 0x76, 0x63, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x62, 0x0a, 0x14, 0x45, 0x6e, 0x67,
	0x69, 0x6e, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x1a, 0x47, 0x0a,
	0x10, 0x50, 0x6f, 0x6f, 0x6c, 0x53, 0x76, 0x63, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x76, 0x63, 0x5f, 0x72, 0x65, 0x70, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0d, 0x52, 0x07, 0x73, 0x76, 0x63, 0x52, 0x65, 0x70, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x0f, 0x0a, 0x0d, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64,
	0x65, 0x64, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x55, 0x0a, 0x0f, 0x43, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x52,
	0x41, 0x53, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x46,
	0x0a, 0x10, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x3c, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x61, 0x6f, 0x73, 0x2d, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2f,
	0x64, 0x61, 0x6f, 0x73, 0x2f, 0x73, 0x72, 0x63, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shared_event_proto_rawDescOnce sync.Once
	file_shared_event_proto_rawDescData = file_shared_event_proto_rawDesc
)

func file_shared_event_proto_rawDescGZIP() []byte {
	file_shared_event_proto_rawDescOnce.Do(func() {
		file_shared_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_event_proto_rawDescData)
	})
	return file_shared_event_proto_rawDescData
}

var file_shared_event_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_shared_event_proto_goTypes = []interface{}{
	(*RASEvent)(nil),                      // 0: shared.RASEvent
	(*ClusterEventReq)(nil),               // 1: shared.ClusterEventReq
	(*ClusterEventResp)(nil),              // 2: shared.ClusterEventResp
	(*RASEvent_EngineStateEventInfo)(nil), // 3: shared.RASEvent.EngineStateEventInfo
	(*RASEvent_PoolSvcEventInfo)(nil),     // 4: shared.RASEvent.PoolSvcEventInfo
}
var file_shared_event_proto_depIdxs = []int32{
	3, // 0: shared.RASEvent.engine_state_info:type_name -> shared.RASEvent.EngineStateEventInfo
	4, // 1: shared.RASEvent.pool_svc_info:type_name -> shared.RASEvent.PoolSvcEventInfo
	0, // 2: shared.ClusterEventReq.event:type_name -> shared.RASEvent
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_shared_event_proto_init() }
func file_shared_event_proto_init() {
	if File_shared_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RASEvent); i {
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
		file_shared_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterEventReq); i {
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
		file_shared_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterEventResp); i {
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
		file_shared_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RASEvent_EngineStateEventInfo); i {
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
		file_shared_event_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RASEvent_PoolSvcEventInfo); i {
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
	file_shared_event_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*RASEvent_StrInfo)(nil),
		(*RASEvent_EngineStateInfo)(nil),
		(*RASEvent_PoolSvcInfo)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_shared_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_event_proto_goTypes,
		DependencyIndexes: file_shared_event_proto_depIdxs,
		MessageInfos:      file_shared_event_proto_msgTypes,
	}.Build()
	File_shared_event_proto = out.File
	file_shared_event_proto_rawDesc = nil
	file_shared_event_proto_goTypes = nil
	file_shared_event_proto_depIdxs = nil
}
