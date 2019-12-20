// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pool.proto

package mgmt

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// PoolRebuildState indicates the pool's rebuild state.
type PoolRebuildState int32

const (
	PoolRebuildState_IDLE PoolRebuildState = 0
	PoolRebuildState_DONE PoolRebuildState = 1
	PoolRebuildState_BUSY PoolRebuildState = 2
)

var PoolRebuildState_name = map[int32]string{
	0: "IDLE",
	1: "DONE",
	2: "BUSY",
}
var PoolRebuildState_value = map[string]int32{
	"IDLE": 0,
	"DONE": 1,
	"BUSY": 2,
}

func (x PoolRebuildState) String() string {
	return proto.EnumName(PoolRebuildState_name, int32(x))
}
func (PoolRebuildState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{0}
}

// PoolCreateReq supplies new pool parameters.
type PoolCreateReq struct {
	Scmbytes             uint64   `protobuf:"varint,1,opt,name=scmbytes,proto3" json:"scmbytes,omitempty"`
	Nvmebytes            uint64   `protobuf:"varint,2,opt,name=nvmebytes,proto3" json:"nvmebytes,omitempty"`
	Ranks                []uint32 `protobuf:"varint,3,rep,packed,name=ranks,proto3" json:"ranks,omitempty"`
	Numsvcreps           uint32   `protobuf:"varint,4,opt,name=numsvcreps,proto3" json:"numsvcreps,omitempty"`
	User                 string   `protobuf:"bytes,5,opt,name=user,proto3" json:"user,omitempty"`
	Usergroup            string   `protobuf:"bytes,6,opt,name=usergroup,proto3" json:"usergroup,omitempty"`
	Uuid                 string   `protobuf:"bytes,7,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Sys                  string   `protobuf:"bytes,8,opt,name=sys,proto3" json:"sys,omitempty"`
	Acl                  []string `protobuf:"bytes,9,rep,name=acl,proto3" json:"acl,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PoolCreateReq) Reset()         { *m = PoolCreateReq{} }
func (m *PoolCreateReq) String() string { return proto.CompactTextString(m) }
func (*PoolCreateReq) ProtoMessage()    {}
func (*PoolCreateReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{0}
}
func (m *PoolCreateReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolCreateReq.Unmarshal(m, b)
}
func (m *PoolCreateReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolCreateReq.Marshal(b, m, deterministic)
}
func (dst *PoolCreateReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolCreateReq.Merge(dst, src)
}
func (m *PoolCreateReq) XXX_Size() int {
	return xxx_messageInfo_PoolCreateReq.Size(m)
}
func (m *PoolCreateReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolCreateReq.DiscardUnknown(m)
}

var xxx_messageInfo_PoolCreateReq proto.InternalMessageInfo

func (m *PoolCreateReq) GetScmbytes() uint64 {
	if m != nil {
		return m.Scmbytes
	}
	return 0
}

func (m *PoolCreateReq) GetNvmebytes() uint64 {
	if m != nil {
		return m.Nvmebytes
	}
	return 0
}

func (m *PoolCreateReq) GetRanks() []uint32 {
	if m != nil {
		return m.Ranks
	}
	return nil
}

func (m *PoolCreateReq) GetNumsvcreps() uint32 {
	if m != nil {
		return m.Numsvcreps
	}
	return 0
}

func (m *PoolCreateReq) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *PoolCreateReq) GetUsergroup() string {
	if m != nil {
		return m.Usergroup
	}
	return ""
}

func (m *PoolCreateReq) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *PoolCreateReq) GetSys() string {
	if m != nil {
		return m.Sys
	}
	return ""
}

func (m *PoolCreateReq) GetAcl() []string {
	if m != nil {
		return m.Acl
	}
	return nil
}

// PoolCreateResp returns created pool uuid and ranks.
type PoolCreateResp struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Svcreps              []uint32 `protobuf:"varint,2,rep,packed,name=svcreps,proto3" json:"svcreps,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PoolCreateResp) Reset()         { *m = PoolCreateResp{} }
func (m *PoolCreateResp) String() string { return proto.CompactTextString(m) }
func (*PoolCreateResp) ProtoMessage()    {}
func (*PoolCreateResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{1}
}
func (m *PoolCreateResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolCreateResp.Unmarshal(m, b)
}
func (m *PoolCreateResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolCreateResp.Marshal(b, m, deterministic)
}
func (dst *PoolCreateResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolCreateResp.Merge(dst, src)
}
func (m *PoolCreateResp) XXX_Size() int {
	return xxx_messageInfo_PoolCreateResp.Size(m)
}
func (m *PoolCreateResp) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolCreateResp.DiscardUnknown(m)
}

var xxx_messageInfo_PoolCreateResp proto.InternalMessageInfo

func (m *PoolCreateResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *PoolCreateResp) GetSvcreps() []uint32 {
	if m != nil {
		return m.Svcreps
	}
	return nil
}

// PoolDestroyReq supplies pool identifier and force flag.
type PoolDestroyReq struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Sys                  string   `protobuf:"bytes,2,opt,name=sys,proto3" json:"sys,omitempty"`
	Force                bool     `protobuf:"varint,3,opt,name=force,proto3" json:"force,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PoolDestroyReq) Reset()         { *m = PoolDestroyReq{} }
func (m *PoolDestroyReq) String() string { return proto.CompactTextString(m) }
func (*PoolDestroyReq) ProtoMessage()    {}
func (*PoolDestroyReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{2}
}
func (m *PoolDestroyReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolDestroyReq.Unmarshal(m, b)
}
func (m *PoolDestroyReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolDestroyReq.Marshal(b, m, deterministic)
}
func (dst *PoolDestroyReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolDestroyReq.Merge(dst, src)
}
func (m *PoolDestroyReq) XXX_Size() int {
	return xxx_messageInfo_PoolDestroyReq.Size(m)
}
func (m *PoolDestroyReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolDestroyReq.DiscardUnknown(m)
}

var xxx_messageInfo_PoolDestroyReq proto.InternalMessageInfo

func (m *PoolDestroyReq) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *PoolDestroyReq) GetSys() string {
	if m != nil {
		return m.Sys
	}
	return ""
}

func (m *PoolDestroyReq) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

// PoolDestroyResp returns resultant state of destroy operation.
type PoolDestroyResp struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PoolDestroyResp) Reset()         { *m = PoolDestroyResp{} }
func (m *PoolDestroyResp) String() string { return proto.CompactTextString(m) }
func (*PoolDestroyResp) ProtoMessage()    {}
func (*PoolDestroyResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{3}
}
func (m *PoolDestroyResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolDestroyResp.Unmarshal(m, b)
}
func (m *PoolDestroyResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolDestroyResp.Marshal(b, m, deterministic)
}
func (dst *PoolDestroyResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolDestroyResp.Merge(dst, src)
}
func (m *PoolDestroyResp) XXX_Size() int {
	return xxx_messageInfo_PoolDestroyResp.Size(m)
}
func (m *PoolDestroyResp) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolDestroyResp.DiscardUnknown(m)
}

var xxx_messageInfo_PoolDestroyResp proto.InternalMessageInfo

func (m *PoolDestroyResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

// ListPoolsReq represents a request to list pools on a given DAOS system.
type ListPoolsReq struct {
	Sys                  string   `protobuf:"bytes,1,opt,name=sys,proto3" json:"sys,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListPoolsReq) Reset()         { *m = ListPoolsReq{} }
func (m *ListPoolsReq) String() string { return proto.CompactTextString(m) }
func (*ListPoolsReq) ProtoMessage()    {}
func (*ListPoolsReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{4}
}
func (m *ListPoolsReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListPoolsReq.Unmarshal(m, b)
}
func (m *ListPoolsReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListPoolsReq.Marshal(b, m, deterministic)
}
func (dst *ListPoolsReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListPoolsReq.Merge(dst, src)
}
func (m *ListPoolsReq) XXX_Size() int {
	return xxx_messageInfo_ListPoolsReq.Size(m)
}
func (m *ListPoolsReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ListPoolsReq.DiscardUnknown(m)
}

var xxx_messageInfo_ListPoolsReq proto.InternalMessageInfo

func (m *ListPoolsReq) GetSys() string {
	if m != nil {
		return m.Sys
	}
	return ""
}

// ListPoolsResp returns the list of pools in the system.
type ListPoolsResp struct {
	Status               int32                 `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Pools                []*ListPoolsResp_Pool `protobuf:"bytes,2,rep,name=pools,proto3" json:"pools,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ListPoolsResp) Reset()         { *m = ListPoolsResp{} }
func (m *ListPoolsResp) String() string { return proto.CompactTextString(m) }
func (*ListPoolsResp) ProtoMessage()    {}
func (*ListPoolsResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{5}
}
func (m *ListPoolsResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListPoolsResp.Unmarshal(m, b)
}
func (m *ListPoolsResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListPoolsResp.Marshal(b, m, deterministic)
}
func (dst *ListPoolsResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListPoolsResp.Merge(dst, src)
}
func (m *ListPoolsResp) XXX_Size() int {
	return xxx_messageInfo_ListPoolsResp.Size(m)
}
func (m *ListPoolsResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ListPoolsResp.DiscardUnknown(m)
}

var xxx_messageInfo_ListPoolsResp proto.InternalMessageInfo

func (m *ListPoolsResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *ListPoolsResp) GetPools() []*ListPoolsResp_Pool {
	if m != nil {
		return m.Pools
	}
	return nil
}

type ListPoolsResp_Pool struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Svcreps              []uint32 `protobuf:"varint,2,rep,packed,name=svcreps,proto3" json:"svcreps,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListPoolsResp_Pool) Reset()         { *m = ListPoolsResp_Pool{} }
func (m *ListPoolsResp_Pool) String() string { return proto.CompactTextString(m) }
func (*ListPoolsResp_Pool) ProtoMessage()    {}
func (*ListPoolsResp_Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{5, 0}
}
func (m *ListPoolsResp_Pool) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListPoolsResp_Pool.Unmarshal(m, b)
}
func (m *ListPoolsResp_Pool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListPoolsResp_Pool.Marshal(b, m, deterministic)
}
func (dst *ListPoolsResp_Pool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListPoolsResp_Pool.Merge(dst, src)
}
func (m *ListPoolsResp_Pool) XXX_Size() int {
	return xxx_messageInfo_ListPoolsResp_Pool.Size(m)
}
func (m *ListPoolsResp_Pool) XXX_DiscardUnknown() {
	xxx_messageInfo_ListPoolsResp_Pool.DiscardUnknown(m)
}

var xxx_messageInfo_ListPoolsResp_Pool proto.InternalMessageInfo

func (m *ListPoolsResp_Pool) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *ListPoolsResp_Pool) GetSvcreps() []uint32 {
	if m != nil {
		return m.Svcreps
	}
	return nil
}

// ListContainers
// Initial implementation differs from C API
// (numContainers not provided in request - get whole list)
type ListContReq struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListContReq) Reset()         { *m = ListContReq{} }
func (m *ListContReq) String() string { return proto.CompactTextString(m) }
func (*ListContReq) ProtoMessage()    {}
func (*ListContReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{6}
}
func (m *ListContReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListContReq.Unmarshal(m, b)
}
func (m *ListContReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListContReq.Marshal(b, m, deterministic)
}
func (dst *ListContReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListContReq.Merge(dst, src)
}
func (m *ListContReq) XXX_Size() int {
	return xxx_messageInfo_ListContReq.Size(m)
}
func (m *ListContReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ListContReq.DiscardUnknown(m)
}

var xxx_messageInfo_ListContReq proto.InternalMessageInfo

func (m *ListContReq) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

type ListContResp struct {
	Status               int32                `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Containers           []*ListContResp_Cont `protobuf:"bytes,2,rep,name=containers,proto3" json:"containers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ListContResp) Reset()         { *m = ListContResp{} }
func (m *ListContResp) String() string { return proto.CompactTextString(m) }
func (*ListContResp) ProtoMessage()    {}
func (*ListContResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{7}
}
func (m *ListContResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListContResp.Unmarshal(m, b)
}
func (m *ListContResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListContResp.Marshal(b, m, deterministic)
}
func (dst *ListContResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListContResp.Merge(dst, src)
}
func (m *ListContResp) XXX_Size() int {
	return xxx_messageInfo_ListContResp.Size(m)
}
func (m *ListContResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ListContResp.DiscardUnknown(m)
}

var xxx_messageInfo_ListContResp proto.InternalMessageInfo

func (m *ListContResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *ListContResp) GetContainers() []*ListContResp_Cont {
	if m != nil {
		return m.Containers
	}
	return nil
}

type ListContResp_Cont struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListContResp_Cont) Reset()         { *m = ListContResp_Cont{} }
func (m *ListContResp_Cont) String() string { return proto.CompactTextString(m) }
func (*ListContResp_Cont) ProtoMessage()    {}
func (*ListContResp_Cont) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{7, 0}
}
func (m *ListContResp_Cont) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListContResp_Cont.Unmarshal(m, b)
}
func (m *ListContResp_Cont) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListContResp_Cont.Marshal(b, m, deterministic)
}
func (dst *ListContResp_Cont) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListContResp_Cont.Merge(dst, src)
}
func (m *ListContResp_Cont) XXX_Size() int {
	return xxx_messageInfo_ListContResp_Cont.Size(m)
}
func (m *ListContResp_Cont) XXX_DiscardUnknown() {
	xxx_messageInfo_ListContResp_Cont.DiscardUnknown(m)
}

var xxx_messageInfo_ListContResp_Cont proto.InternalMessageInfo

func (m *ListContResp_Cont) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

// PoolQueryReq represents a pool query request.
type PoolQueryReq struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PoolQueryReq) Reset()         { *m = PoolQueryReq{} }
func (m *PoolQueryReq) String() string { return proto.CompactTextString(m) }
func (*PoolQueryReq) ProtoMessage()    {}
func (*PoolQueryReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{8}
}
func (m *PoolQueryReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolQueryReq.Unmarshal(m, b)
}
func (m *PoolQueryReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolQueryReq.Marshal(b, m, deterministic)
}
func (dst *PoolQueryReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolQueryReq.Merge(dst, src)
}
func (m *PoolQueryReq) XXX_Size() int {
	return xxx_messageInfo_PoolQueryReq.Size(m)
}
func (m *PoolQueryReq) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolQueryReq.DiscardUnknown(m)
}

var xxx_messageInfo_PoolQueryReq proto.InternalMessageInfo

func (m *PoolQueryReq) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

// StorageUsageStats represents usage statistics for a storage subsystem.
type StorageUsageStats struct {
	Total                uint64   `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Free                 uint64   `protobuf:"varint,2,opt,name=free,proto3" json:"free,omitempty"`
	Min                  uint64   `protobuf:"varint,3,opt,name=min,proto3" json:"min,omitempty"`
	Max                  uint64   `protobuf:"varint,4,opt,name=max,proto3" json:"max,omitempty"`
	Mean                 uint64   `protobuf:"varint,5,opt,name=mean,proto3" json:"mean,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StorageUsageStats) Reset()         { *m = StorageUsageStats{} }
func (m *StorageUsageStats) String() string { return proto.CompactTextString(m) }
func (*StorageUsageStats) ProtoMessage()    {}
func (*StorageUsageStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{9}
}
func (m *StorageUsageStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StorageUsageStats.Unmarshal(m, b)
}
func (m *StorageUsageStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StorageUsageStats.Marshal(b, m, deterministic)
}
func (dst *StorageUsageStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StorageUsageStats.Merge(dst, src)
}
func (m *StorageUsageStats) XXX_Size() int {
	return xxx_messageInfo_StorageUsageStats.Size(m)
}
func (m *StorageUsageStats) XXX_DiscardUnknown() {
	xxx_messageInfo_StorageUsageStats.DiscardUnknown(m)
}

var xxx_messageInfo_StorageUsageStats proto.InternalMessageInfo

func (m *StorageUsageStats) GetTotal() uint64 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *StorageUsageStats) GetFree() uint64 {
	if m != nil {
		return m.Free
	}
	return 0
}

func (m *StorageUsageStats) GetMin() uint64 {
	if m != nil {
		return m.Min
	}
	return 0
}

func (m *StorageUsageStats) GetMax() uint64 {
	if m != nil {
		return m.Max
	}
	return 0
}

func (m *StorageUsageStats) GetMean() uint64 {
	if m != nil {
		return m.Mean
	}
	return 0
}

// PoolRebuildStatus represents a pool's rebuild status.
type PoolRebuildStatus struct {
	Status               int32            `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	State                PoolRebuildState `protobuf:"varint,2,opt,name=state,proto3,enum=mgmt.PoolRebuildState" json:"state,omitempty"`
	Objects              uint64           `protobuf:"varint,3,opt,name=objects,proto3" json:"objects,omitempty"`
	Records              uint64           `protobuf:"varint,4,opt,name=records,proto3" json:"records,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *PoolRebuildStatus) Reset()         { *m = PoolRebuildStatus{} }
func (m *PoolRebuildStatus) String() string { return proto.CompactTextString(m) }
func (*PoolRebuildStatus) ProtoMessage()    {}
func (*PoolRebuildStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{10}
}
func (m *PoolRebuildStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolRebuildStatus.Unmarshal(m, b)
}
func (m *PoolRebuildStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolRebuildStatus.Marshal(b, m, deterministic)
}
func (dst *PoolRebuildStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolRebuildStatus.Merge(dst, src)
}
func (m *PoolRebuildStatus) XXX_Size() int {
	return xxx_messageInfo_PoolRebuildStatus.Size(m)
}
func (m *PoolRebuildStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolRebuildStatus.DiscardUnknown(m)
}

var xxx_messageInfo_PoolRebuildStatus proto.InternalMessageInfo

func (m *PoolRebuildStatus) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *PoolRebuildStatus) GetState() PoolRebuildState {
	if m != nil {
		return m.State
	}
	return PoolRebuildState_IDLE
}

func (m *PoolRebuildStatus) GetObjects() uint64 {
	if m != nil {
		return m.Objects
	}
	return 0
}

func (m *PoolRebuildStatus) GetRecords() uint64 {
	if m != nil {
		return m.Records
	}
	return 0
}

// PoolQueryResp represents a pool query response.
type PoolQueryResp struct {
	Status               int32              `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Uuid                 string             `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Totaltargets         uint32             `protobuf:"varint,3,opt,name=totaltargets,proto3" json:"totaltargets,omitempty"`
	Activetargets        uint32             `protobuf:"varint,4,opt,name=activetargets,proto3" json:"activetargets,omitempty"`
	Disabled             bool               `protobuf:"varint,5,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Rebuild              *PoolRebuildStatus `protobuf:"bytes,6,opt,name=rebuild,proto3" json:"rebuild,omitempty"`
	Scm                  *StorageUsageStats `protobuf:"bytes,7,opt,name=scm,proto3" json:"scm,omitempty"`
	Nvme                 *StorageUsageStats `protobuf:"bytes,8,opt,name=nvme,proto3" json:"nvme,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *PoolQueryResp) Reset()         { *m = PoolQueryResp{} }
func (m *PoolQueryResp) String() string { return proto.CompactTextString(m) }
func (*PoolQueryResp) ProtoMessage()    {}
func (*PoolQueryResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_726a694af50615e4, []int{11}
}
func (m *PoolQueryResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolQueryResp.Unmarshal(m, b)
}
func (m *PoolQueryResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolQueryResp.Marshal(b, m, deterministic)
}
func (dst *PoolQueryResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolQueryResp.Merge(dst, src)
}
func (m *PoolQueryResp) XXX_Size() int {
	return xxx_messageInfo_PoolQueryResp.Size(m)
}
func (m *PoolQueryResp) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolQueryResp.DiscardUnknown(m)
}

var xxx_messageInfo_PoolQueryResp proto.InternalMessageInfo

func (m *PoolQueryResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *PoolQueryResp) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *PoolQueryResp) GetTotaltargets() uint32 {
	if m != nil {
		return m.Totaltargets
	}
	return 0
}

func (m *PoolQueryResp) GetActivetargets() uint32 {
	if m != nil {
		return m.Activetargets
	}
	return 0
}

func (m *PoolQueryResp) GetDisabled() bool {
	if m != nil {
		return m.Disabled
	}
	return false
}

func (m *PoolQueryResp) GetRebuild() *PoolRebuildStatus {
	if m != nil {
		return m.Rebuild
	}
	return nil
}

func (m *PoolQueryResp) GetScm() *StorageUsageStats {
	if m != nil {
		return m.Scm
	}
	return nil
}

func (m *PoolQueryResp) GetNvme() *StorageUsageStats {
	if m != nil {
		return m.Nvme
	}
	return nil
}

func init() {
	proto.RegisterType((*PoolCreateReq)(nil), "mgmt.PoolCreateReq")
	proto.RegisterType((*PoolCreateResp)(nil), "mgmt.PoolCreateResp")
	proto.RegisterType((*PoolDestroyReq)(nil), "mgmt.PoolDestroyReq")
	proto.RegisterType((*PoolDestroyResp)(nil), "mgmt.PoolDestroyResp")
	proto.RegisterType((*ListPoolsReq)(nil), "mgmt.ListPoolsReq")
	proto.RegisterType((*ListPoolsResp)(nil), "mgmt.ListPoolsResp")
	proto.RegisterType((*ListPoolsResp_Pool)(nil), "mgmt.ListPoolsResp.Pool")
	proto.RegisterType((*ListContReq)(nil), "mgmt.ListContReq")
	proto.RegisterType((*ListContResp)(nil), "mgmt.ListContResp")
	proto.RegisterType((*ListContResp_Cont)(nil), "mgmt.ListContResp.Cont")
	proto.RegisterType((*PoolQueryReq)(nil), "mgmt.PoolQueryReq")
	proto.RegisterType((*StorageUsageStats)(nil), "mgmt.StorageUsageStats")
	proto.RegisterType((*PoolRebuildStatus)(nil), "mgmt.PoolRebuildStatus")
	proto.RegisterType((*PoolQueryResp)(nil), "mgmt.PoolQueryResp")
	proto.RegisterEnum("mgmt.PoolRebuildState", PoolRebuildState_name, PoolRebuildState_value)
}

func init() { proto.RegisterFile("pool.proto", fileDescriptor_pool_726a694af50615e4) }

var fileDescriptor_pool_726a694af50615e4 = []byte{
	// 633 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0xdb, 0x6e, 0xd3, 0x4c,
	0x10, 0xfe, 0x7d, 0x6a, 0xd3, 0x69, 0xdd, 0x3f, 0x5d, 0x55, 0xc5, 0x8a, 0x10, 0x32, 0x16, 0x17,
	0x29, 0x20, 0x0b, 0x0a, 0x12, 0xf7, 0x3d, 0x5c, 0x20, 0x55, 0x1c, 0x36, 0xea, 0x05, 0x97, 0x1b,
	0x7b, 0x1b, 0x19, 0x6c, 0x6f, 0xd8, 0x5d, 0x57, 0x8d, 0x78, 0x06, 0xde, 0x84, 0x17, 0xe3, 0x2d,
	0xd0, 0xec, 0xda, 0x39, 0x34, 0x4d, 0xae, 0x3c, 0x33, 0xfb, 0x8d, 0xe7, 0x9b, 0x6f, 0x76, 0x07,
	0x60, 0x2a, 0x44, 0x99, 0x4e, 0xa5, 0xd0, 0x82, 0xf8, 0xd5, 0xa4, 0xd2, 0xc9, 0x5f, 0x07, 0xc2,
	0x2f, 0x42, 0x94, 0x17, 0x92, 0x33, 0xcd, 0x29, 0xff, 0x49, 0x06, 0xd0, 0x53, 0x59, 0x35, 0x9e,
	0x69, 0xae, 0x22, 0x27, 0x76, 0x86, 0x3e, 0x9d, 0xfb, 0xe4, 0x29, 0xec, 0xd5, 0x77, 0x15, 0xb7,
	0x87, 0xae, 0x39, 0x5c, 0x04, 0xc8, 0x31, 0x04, 0x92, 0xd5, 0x3f, 0x54, 0xe4, 0xc5, 0xde, 0x30,
	0xa4, 0xd6, 0x21, 0xcf, 0x00, 0xea, 0xa6, 0x52, 0x77, 0x99, 0xe4, 0x53, 0x15, 0xf9, 0xb1, 0x33,
	0x0c, 0xe9, 0x52, 0x84, 0x10, 0xf0, 0x1b, 0xc5, 0x65, 0x14, 0xc4, 0xce, 0x70, 0x8f, 0x1a, 0x1b,
	0xeb, 0xe0, 0x77, 0x22, 0x45, 0x33, 0x8d, 0x76, 0xcc, 0xc1, 0x22, 0x60, 0x32, 0x9a, 0x22, 0x8f,
	0x76, 0xdb, 0x8c, 0xa6, 0xc8, 0x49, 0x1f, 0x3c, 0x35, 0x53, 0x51, 0xcf, 0x84, 0xd0, 0xc4, 0x08,
	0xcb, 0xca, 0x68, 0x2f, 0xf6, 0x30, 0xc2, 0xb2, 0x32, 0x39, 0x87, 0xc3, 0xe5, 0x56, 0xd5, 0x94,
	0x9c, 0xc0, 0x8e, 0xd2, 0x4c, 0x37, 0xb6, 0xd3, 0x80, 0xb6, 0x1e, 0x89, 0x60, 0xb7, 0x23, 0xec,
	0x9a, 0x5e, 0x3a, 0x37, 0xb9, 0xb6, 0xff, 0xb8, 0xe4, 0x4a, 0x4b, 0x31, 0x43, 0xbd, 0x3a, 0x36,
	0xce, 0x3a, 0x1b, 0x77, 0xc1, 0xe6, 0x18, 0x82, 0x5b, 0x21, 0x33, 0x1e, 0x79, 0xb1, 0x33, 0xec,
	0x51, 0xeb, 0x24, 0xa7, 0xf0, 0xff, 0xca, 0xdf, 0x36, 0x53, 0x4a, 0x62, 0x38, 0xb8, 0x2e, 0x94,
	0x46, 0xb8, 0xc2, 0xb2, 0x6d, 0x09, 0x67, 0x5e, 0x22, 0xf9, 0xed, 0x40, 0xb8, 0x04, 0xd9, 0xd2,
	0x5e, 0x0a, 0x01, 0x5e, 0x04, 0xdb, 0xdc, 0xfe, 0x59, 0x94, 0xe2, 0x55, 0x48, 0x57, 0x72, 0x53,
	0xb4, 0xa8, 0x85, 0x0d, 0xde, 0x83, 0x8f, 0xee, 0xa3, 0xad, 0x6e, 0x96, 0xea, 0x39, 0xec, 0xe3,
	0x2f, 0x2f, 0x44, 0xad, 0x37, 0xe8, 0x94, 0xfc, 0xb2, 0x4d, 0x59, 0xc8, 0x16, 0xc2, 0x1f, 0x00,
	0x32, 0x51, 0x6b, 0x56, 0xd4, 0x5c, 0x76, 0xac, 0x9f, 0x2c, 0x58, 0x77, 0xf9, 0xa9, 0x31, 0x96,
	0xa0, 0x83, 0x01, 0xf8, 0x18, 0x7b, 0xb4, 0x78, 0x02, 0x07, 0xd8, 0xd5, 0xd7, 0x86, 0xcb, 0x4d,
	0x83, 0x4c, 0x1a, 0x38, 0x1a, 0x69, 0x21, 0xd9, 0x84, 0xdf, 0x28, 0x36, 0xe1, 0x23, 0xcd, 0xb4,
	0x99, 0xa5, 0x16, 0x9a, 0x95, 0xed, 0xf3, 0xb0, 0x0e, 0xa6, 0xdf, 0x4a, 0xce, 0xdb, 0x67, 0x61,
	0x6c, 0x1c, 0x52, 0x55, 0xd4, 0x66, 0xe6, 0x3e, 0x45, 0xd3, 0x44, 0xd8, 0xbd, 0x79, 0x06, 0x18,
	0x61, 0xf7, 0x98, 0x57, 0x71, 0x56, 0x9b, 0xfb, 0xef, 0x53, 0x63, 0xe3, 0x28, 0x8f, 0xcc, 0x00,
	0xf8, 0xb8, 0x29, 0xca, 0x7c, 0x64, 0x55, 0xd8, 0xa4, 0xce, 0x6b, 0x08, 0xd0, 0xb2, 0xa5, 0x0f,
	0xcf, 0x4e, 0xac, 0x30, 0x0f, 0xf2, 0x39, 0xb5, 0x20, 0x1c, 0x98, 0x18, 0x7f, 0xe7, 0x99, 0x56,
	0x2d, 0xaf, 0xce, 0xc5, 0x13, 0xc9, 0x33, 0x21, 0x73, 0xd5, 0xf2, 0xeb, 0xdc, 0xe4, 0x8f, 0x6b,
	0xb7, 0x44, 0xab, 0xd5, 0x96, 0x49, 0x75, 0x22, 0xba, 0x4b, 0x57, 0x24, 0x81, 0x03, 0x23, 0x91,
	0x66, 0x72, 0xc2, 0xdb, 0xb2, 0x21, 0x5d, 0x89, 0x91, 0x17, 0x10, 0xb2, 0x4c, 0x17, 0x77, 0xbc,
	0x03, 0xd9, 0x45, 0xb1, 0x1a, 0xc4, 0xdd, 0x94, 0x17, 0x8a, 0x8d, 0x4b, 0x9e, 0x1b, 0xbd, 0x7a,
	0x74, 0xee, 0x93, 0xb7, 0xc8, 0xde, 0xb4, 0x6b, 0x36, 0xc6, 0xfc, 0x82, 0xac, 0xe9, 0x48, 0x3b,
	0x1c, 0x39, 0x05, 0x4f, 0x65, 0x95, 0xd9, 0x23, 0x73, 0xf8, 0xda, 0xb8, 0x29, 0x62, 0xc8, 0x2b,
	0xf0, 0x71, 0xd1, 0x99, 0x05, 0xb3, 0x05, 0x6b, 0x40, 0x2f, 0xdf, 0x40, 0xff, 0xa1, 0xfa, 0xa4,
	0x07, 0xfe, 0xc7, 0xcb, 0xeb, 0xab, 0xfe, 0x7f, 0x68, 0x5d, 0x7e, 0xfe, 0x74, 0xd5, 0x77, 0xd0,
	0x3a, 0xbf, 0x19, 0x7d, 0xeb, 0xbb, 0xe3, 0x1d, 0xb3, 0x93, 0xdf, 0xfd, 0x0b, 0x00, 0x00, 0xff,
	0xff, 0x62, 0xb0, 0x04, 0xaa, 0xa1, 0x05, 0x00, 0x00,
}
