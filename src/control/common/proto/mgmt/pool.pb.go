// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pool.proto

package mgmt

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

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
	return fileDescriptor_8a14d8612184524f, []int{0}
}

func (m *PoolCreateReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolCreateReq.Unmarshal(m, b)
}
func (m *PoolCreateReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolCreateReq.Marshal(b, m, deterministic)
}
func (m *PoolCreateReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolCreateReq.Merge(m, src)
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
	return fileDescriptor_8a14d8612184524f, []int{1}
}

func (m *PoolCreateResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolCreateResp.Unmarshal(m, b)
}
func (m *PoolCreateResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolCreateResp.Marshal(b, m, deterministic)
}
func (m *PoolCreateResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolCreateResp.Merge(m, src)
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
	return fileDescriptor_8a14d8612184524f, []int{2}
}

func (m *PoolDestroyReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolDestroyReq.Unmarshal(m, b)
}
func (m *PoolDestroyReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolDestroyReq.Marshal(b, m, deterministic)
}
func (m *PoolDestroyReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolDestroyReq.Merge(m, src)
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
	return fileDescriptor_8a14d8612184524f, []int{3}
}

func (m *PoolDestroyResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PoolDestroyResp.Unmarshal(m, b)
}
func (m *PoolDestroyResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PoolDestroyResp.Marshal(b, m, deterministic)
}
func (m *PoolDestroyResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolDestroyResp.Merge(m, src)
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
	return fileDescriptor_8a14d8612184524f, []int{4}
}

func (m *ListPoolsReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListPoolsReq.Unmarshal(m, b)
}
func (m *ListPoolsReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListPoolsReq.Marshal(b, m, deterministic)
}
func (m *ListPoolsReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListPoolsReq.Merge(m, src)
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
	return fileDescriptor_8a14d8612184524f, []int{5}
}

func (m *ListPoolsResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListPoolsResp.Unmarshal(m, b)
}
func (m *ListPoolsResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListPoolsResp.Marshal(b, m, deterministic)
}
func (m *ListPoolsResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListPoolsResp.Merge(m, src)
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
	return fileDescriptor_8a14d8612184524f, []int{5, 0}
}

func (m *ListPoolsResp_Pool) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListPoolsResp_Pool.Unmarshal(m, b)
}
func (m *ListPoolsResp_Pool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListPoolsResp_Pool.Marshal(b, m, deterministic)
}
func (m *ListPoolsResp_Pool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListPoolsResp_Pool.Merge(m, src)
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
	return fileDescriptor_8a14d8612184524f, []int{6}
}

func (m *ListContReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListContReq.Unmarshal(m, b)
}
func (m *ListContReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListContReq.Marshal(b, m, deterministic)
}
func (m *ListContReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListContReq.Merge(m, src)
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
	NumContainers        uint64               `protobuf:"varint,3,opt,name=numContainers,proto3" json:"numContainers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ListContResp) Reset()         { *m = ListContResp{} }
func (m *ListContResp) String() string { return proto.CompactTextString(m) }
func (*ListContResp) ProtoMessage()    {}
func (*ListContResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a14d8612184524f, []int{7}
}

func (m *ListContResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListContResp.Unmarshal(m, b)
}
func (m *ListContResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListContResp.Marshal(b, m, deterministic)
}
func (m *ListContResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListContResp.Merge(m, src)
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

func (m *ListContResp) GetNumContainers() uint64 {
	if m != nil {
		return m.NumContainers
	}
	return 0
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
	return fileDescriptor_8a14d8612184524f, []int{7, 0}
}

func (m *ListContResp_Cont) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListContResp_Cont.Unmarshal(m, b)
}
func (m *ListContResp_Cont) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListContResp_Cont.Marshal(b, m, deterministic)
}
func (m *ListContResp_Cont) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListContResp_Cont.Merge(m, src)
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
}

func init() { proto.RegisterFile("pool.proto", fileDescriptor_8a14d8612184524f) }

var fileDescriptor_8a14d8612184524f = []byte{
	// 393 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0x4b, 0x4e, 0xeb, 0x30,
	0x14, 0x95, 0x9b, 0xa4, 0x9f, 0xdb, 0x97, 0xf7, 0x9e, 0xac, 0x0a, 0xac, 0x08, 0xa1, 0x10, 0x31,
	0x08, 0x93, 0x0c, 0x00, 0x89, 0x39, 0x65, 0xd8, 0x01, 0xf2, 0x0e, 0xd2, 0xd4, 0x54, 0x11, 0x49,
	0x1c, 0x6c, 0xa7, 0x52, 0x17, 0xc1, 0x22, 0x58, 0x1e, 0xbb, 0x40, 0xb6, 0x9b, 0x26, 0x15, 0x6d,
	0x47, 0xb9, 0xe7, 0xf8, 0x7e, 0xce, 0xb9, 0xba, 0x01, 0xa8, 0x39, 0x2f, 0x92, 0x5a, 0x70, 0xc5,
	0xb1, 0x5b, 0xae, 0x4b, 0x15, 0x7d, 0x23, 0xf0, 0x5f, 0x39, 0x2f, 0xe6, 0x82, 0xa5, 0x8a, 0x51,
	0xf6, 0x81, 0x03, 0x18, 0xcb, 0xac, 0x5c, 0x6e, 0x15, 0x93, 0x04, 0x85, 0x28, 0x76, 0xe9, 0x1e,
	0xe3, 0x2b, 0x98, 0x54, 0x9b, 0x92, 0xd9, 0xc7, 0x81, 0x79, 0xec, 0x08, 0x3c, 0x03, 0x4f, 0xa4,
	0xd5, 0xbb, 0x24, 0x4e, 0xe8, 0xc4, 0x3e, 0xb5, 0x00, 0x5f, 0x03, 0x54, 0x4d, 0x29, 0x37, 0x99,
	0x60, 0xb5, 0x24, 0x6e, 0x88, 0x62, 0x9f, 0xf6, 0x18, 0x8c, 0xc1, 0x6d, 0x24, 0x13, 0xc4, 0x0b,
	0x51, 0x3c, 0xa1, 0x26, 0xd6, 0x73, 0xf4, 0x77, 0x2d, 0x78, 0x53, 0x93, 0xa1, 0x79, 0xe8, 0x08,
	0x53, 0xd1, 0xe4, 0x2b, 0x32, 0xda, 0x55, 0x34, 0xf9, 0x0a, 0xff, 0x07, 0x47, 0x6e, 0x25, 0x19,
	0x1b, 0x4a, 0x87, 0x9a, 0x49, 0xb3, 0x82, 0x4c, 0x42, 0x47, 0x33, 0x69, 0x56, 0x44, 0xcf, 0xf0,
	0xb7, 0x6f, 0x55, 0xd6, 0xf8, 0x02, 0x86, 0x52, 0xa5, 0xaa, 0xb1, 0x4e, 0x3d, 0xba, 0x43, 0x98,
	0xc0, 0xa8, 0x15, 0x3c, 0x30, 0x5e, 0x5a, 0x18, 0x2d, 0x6c, 0x8f, 0x17, 0x26, 0x95, 0xe0, 0x5b,
	0xbd, 0xaf, 0x56, 0x0d, 0xfa, 0xad, 0x66, 0xd0, 0xa9, 0x99, 0x81, 0xf7, 0xc6, 0x45, 0xc6, 0x88,
	0x13, 0xa2, 0x78, 0x4c, 0x2d, 0x88, 0xee, 0xe0, 0xdf, 0x41, 0xb7, 0xd3, 0x92, 0xa2, 0x10, 0xfe,
	0x2c, 0x72, 0xa9, 0x74, 0xba, 0xd4, 0x63, 0x77, 0x23, 0xd0, 0x7e, 0x44, 0xf4, 0x89, 0xc0, 0xef,
	0xa5, 0x9c, 0xb1, 0x97, 0x80, 0xa7, 0x0f, 0xc1, 0x9a, 0x9b, 0xde, 0x93, 0x44, 0x9f, 0x42, 0x72,
	0x50, 0x9b, 0xe8, 0x88, 0xda, 0xb4, 0xe0, 0x11, 0x5c, 0x0d, 0x8f, 0x5a, 0x3d, 0xbd, 0xaa, 0x1b,
	0x98, 0xea, 0x96, 0x73, 0x5e, 0xa9, 0x13, 0x7b, 0x8a, 0xbe, 0x90, 0x75, 0x65, 0x73, 0xce, 0x28,
	0x7e, 0x02, 0xc8, 0x78, 0xa5, 0xd2, 0xbc, 0x62, 0xa2, 0x95, 0x7d, 0xd9, 0xc9, 0x6e, 0xeb, 0x13,
	0x13, 0xf4, 0x52, 0xf1, 0x2d, 0xf8, 0x55, 0x53, 0xce, 0xbb, 0x5a, 0xc7, 0x5c, 0xed, 0x21, 0x19,
	0x04, 0xe0, 0x6a, 0x74, 0x4c, 0xe3, 0x72, 0x68, 0x7e, 0x97, 0x87, 0x9f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x9d, 0x83, 0xa5, 0xf5, 0x3c, 0x03, 0x00, 0x00,
}
