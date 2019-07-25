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

// CreatePoolReq supplies new pool parameters.
type CreatePoolReq struct {
	Scmbytes             uint64   `protobuf:"varint,1,opt,name=scmbytes,proto3" json:"scmbytes,omitempty"`
	Nvmebytes            uint64   `protobuf:"varint,2,opt,name=nvmebytes,proto3" json:"nvmebytes,omitempty"`
	Ranks                string   `protobuf:"bytes,3,opt,name=ranks,proto3" json:"ranks,omitempty"`
	Numsvcreps           uint32   `protobuf:"varint,4,opt,name=numsvcreps,proto3" json:"numsvcreps,omitempty"`
	User                 string   `protobuf:"bytes,5,opt,name=user,proto3" json:"user,omitempty"`
	Usergroup            string   `protobuf:"bytes,6,opt,name=usergroup,proto3" json:"usergroup,omitempty"`
	Sys                  string   `protobuf:"bytes,7,opt,name=sys,proto3" json:"sys,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreatePoolReq) Reset()         { *m = CreatePoolReq{} }
func (m *CreatePoolReq) String() string { return proto.CompactTextString(m) }
func (*CreatePoolReq) ProtoMessage()    {}
func (*CreatePoolReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_137e83b9de26bb58, []int{0}
}
func (m *CreatePoolReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreatePoolReq.Unmarshal(m, b)
}
func (m *CreatePoolReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreatePoolReq.Marshal(b, m, deterministic)
}
func (dst *CreatePoolReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreatePoolReq.Merge(dst, src)
}
func (m *CreatePoolReq) XXX_Size() int {
	return xxx_messageInfo_CreatePoolReq.Size(m)
}
func (m *CreatePoolReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CreatePoolReq.DiscardUnknown(m)
}

var xxx_messageInfo_CreatePoolReq proto.InternalMessageInfo

func (m *CreatePoolReq) GetScmbytes() uint64 {
	if m != nil {
		return m.Scmbytes
	}
	return 0
}

func (m *CreatePoolReq) GetNvmebytes() uint64 {
	if m != nil {
		return m.Nvmebytes
	}
	return 0
}

func (m *CreatePoolReq) GetRanks() string {
	if m != nil {
		return m.Ranks
	}
	return ""
}

func (m *CreatePoolReq) GetNumsvcreps() uint32 {
	if m != nil {
		return m.Numsvcreps
	}
	return 0
}

func (m *CreatePoolReq) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *CreatePoolReq) GetUsergroup() string {
	if m != nil {
		return m.Usergroup
	}
	return ""
}

func (m *CreatePoolReq) GetSys() string {
	if m != nil {
		return m.Sys
	}
	return ""
}

// CreatePoolResp returns created pool uuid and ranks.
type CreatePoolResp struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Uuid                 string   `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Svcreps              string   `protobuf:"bytes,3,opt,name=svcreps,proto3" json:"svcreps,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreatePoolResp) Reset()         { *m = CreatePoolResp{} }
func (m *CreatePoolResp) String() string { return proto.CompactTextString(m) }
func (*CreatePoolResp) ProtoMessage()    {}
func (*CreatePoolResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_137e83b9de26bb58, []int{1}
}
func (m *CreatePoolResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreatePoolResp.Unmarshal(m, b)
}
func (m *CreatePoolResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreatePoolResp.Marshal(b, m, deterministic)
}
func (dst *CreatePoolResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreatePoolResp.Merge(dst, src)
}
func (m *CreatePoolResp) XXX_Size() int {
	return xxx_messageInfo_CreatePoolResp.Size(m)
}
func (m *CreatePoolResp) XXX_DiscardUnknown() {
	xxx_messageInfo_CreatePoolResp.DiscardUnknown(m)
}

var xxx_messageInfo_CreatePoolResp proto.InternalMessageInfo

func (m *CreatePoolResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *CreatePoolResp) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *CreatePoolResp) GetSvcreps() string {
	if m != nil {
		return m.Svcreps
	}
	return ""
}

// DestroyPoolReq supplies pool identifier and force flag.
type DestroyPoolReq struct {
	Uuid                 string   `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Sys                  string   `protobuf:"bytes,2,opt,name=sys,proto3" json:"sys,omitempty"`
	Force                bool     `protobuf:"varint,3,opt,name=force,proto3" json:"force,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DestroyPoolReq) Reset()         { *m = DestroyPoolReq{} }
func (m *DestroyPoolReq) String() string { return proto.CompactTextString(m) }
func (*DestroyPoolReq) ProtoMessage()    {}
func (*DestroyPoolReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_137e83b9de26bb58, []int{2}
}
func (m *DestroyPoolReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DestroyPoolReq.Unmarshal(m, b)
}
func (m *DestroyPoolReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DestroyPoolReq.Marshal(b, m, deterministic)
}
func (dst *DestroyPoolReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DestroyPoolReq.Merge(dst, src)
}
func (m *DestroyPoolReq) XXX_Size() int {
	return xxx_messageInfo_DestroyPoolReq.Size(m)
}
func (m *DestroyPoolReq) XXX_DiscardUnknown() {
	xxx_messageInfo_DestroyPoolReq.DiscardUnknown(m)
}

var xxx_messageInfo_DestroyPoolReq proto.InternalMessageInfo

func (m *DestroyPoolReq) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *DestroyPoolReq) GetSys() string {
	if m != nil {
		return m.Sys
	}
	return ""
}

func (m *DestroyPoolReq) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

// DestroyPoolResp returns resultant state of destroy operation.
type DestroyPoolResp struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DestroyPoolResp) Reset()         { *m = DestroyPoolResp{} }
func (m *DestroyPoolResp) String() string { return proto.CompactTextString(m) }
func (*DestroyPoolResp) ProtoMessage()    {}
func (*DestroyPoolResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_137e83b9de26bb58, []int{3}
}
func (m *DestroyPoolResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DestroyPoolResp.Unmarshal(m, b)
}
func (m *DestroyPoolResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DestroyPoolResp.Marshal(b, m, deterministic)
}
func (dst *DestroyPoolResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DestroyPoolResp.Merge(dst, src)
}
func (m *DestroyPoolResp) XXX_Size() int {
	return xxx_messageInfo_DestroyPoolResp.Size(m)
}
func (m *DestroyPoolResp) XXX_DiscardUnknown() {
	xxx_messageInfo_DestroyPoolResp.DiscardUnknown(m)
}

var xxx_messageInfo_DestroyPoolResp proto.InternalMessageInfo

func (m *DestroyPoolResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func init() {
	proto.RegisterType((*CreatePoolReq)(nil), "mgmt.CreatePoolReq")
	proto.RegisterType((*CreatePoolResp)(nil), "mgmt.CreatePoolResp")
	proto.RegisterType((*DestroyPoolReq)(nil), "mgmt.DestroyPoolReq")
	proto.RegisterType((*DestroyPoolResp)(nil), "mgmt.DestroyPoolResp")
}

func init() { proto.RegisterFile("pool.proto", fileDescriptor_pool_137e83b9de26bb58) }

var fileDescriptor_pool_137e83b9de26bb58 = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xc1, 0x4e, 0xc3, 0x30,
	0x0c, 0x86, 0x95, 0xad, 0xed, 0x36, 0x4b, 0x1b, 0xc8, 0x42, 0x28, 0x42, 0x08, 0x55, 0x3d, 0x95,
	0x0b, 0x17, 0x1e, 0x01, 0x8e, 0x1c, 0x50, 0x0e, 0xdc, 0xbb, 0x2e, 0x4c, 0x13, 0x4b, 0x13, 0xe2,
	0x64, 0x52, 0x9f, 0x8f, 0x17, 0x43, 0x49, 0xd6, 0xb5, 0x1c, 0x38, 0xc5, 0xff, 0x6f, 0xd9, 0xfe,
	0xec, 0x00, 0x18, 0xad, 0x8f, 0x4f, 0xc6, 0x6a, 0xa7, 0x31, 0x53, 0x7b, 0xe5, 0xaa, 0x1f, 0x06,
	0xeb, 0x17, 0x2b, 0x1b, 0x27, 0xdf, 0xb5, 0x3e, 0x0a, 0xf9, 0x8d, 0x77, 0xb0, 0xa4, 0x56, 0x6d,
	0x7b, 0x27, 0x89, 0xb3, 0x92, 0xd5, 0x99, 0xb8, 0x68, 0xbc, 0x87, 0x55, 0x77, 0x52, 0x32, 0x25,
	0x67, 0x31, 0x39, 0x1a, 0x78, 0x03, 0xb9, 0x6d, 0xba, 0x2f, 0xe2, 0xf3, 0x92, 0xd5, 0x2b, 0x91,
	0x04, 0x3e, 0x00, 0x74, 0x5e, 0xd1, 0xa9, 0xb5, 0xd2, 0x10, 0xcf, 0x4a, 0x56, 0xaf, 0xc5, 0xc4,
	0x41, 0x84, 0xcc, 0x93, 0xb4, 0x3c, 0x8f, 0x45, 0x31, 0x0e, 0x73, 0xc2, 0xbb, 0xb7, 0xda, 0x1b,
	0x5e, 0xc4, 0xc4, 0x68, 0xe0, 0x35, 0xcc, 0xa9, 0x27, 0xbe, 0x88, 0x7e, 0x08, 0xab, 0x0f, 0xd8,
	0x4c, 0x97, 0x20, 0x83, 0xb7, 0x50, 0x90, 0x6b, 0x9c, 0x4f, 0x3b, 0xe4, 0xe2, 0xac, 0xe2, 0x34,
	0x7f, 0xd8, 0x45, 0xf8, 0x30, 0xcd, 0x1f, 0x76, 0xc8, 0x61, 0x31, 0xe0, 0x25, 0xf2, 0x41, 0x56,
	0x6f, 0xb0, 0x79, 0x95, 0xe4, 0xac, 0xee, 0x87, 0xeb, 0x0c, 0xf5, 0x6c, 0x52, 0x7f, 0xe6, 0x99,
	0x5d, 0x78, 0xc2, 0x25, 0x3e, 0xb5, 0x6d, 0x65, 0xec, 0xb7, 0x14, 0x49, 0x54, 0x8f, 0x70, 0xf5,
	0xa7, 0xdb, 0xff, 0x98, 0xdb, 0x22, 0xfe, 0xd1, 0xf3, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x53,
	0xa2, 0x98, 0x8c, 0xb1, 0x01, 0x00, 0x00,
}
