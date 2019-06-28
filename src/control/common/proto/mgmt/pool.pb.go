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
	return fileDescriptor_pool_b3b79726342cd324, []int{0}
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
	Status               DaosRequestStatus `protobuf:"varint,1,opt,name=status,proto3,enum=mgmt.DaosRequestStatus" json:"status,omitempty"`
	Uuid                 string            `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Svcreps              string            `protobuf:"bytes,3,opt,name=svcreps,proto3" json:"svcreps,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *CreatePoolResp) Reset()         { *m = CreatePoolResp{} }
func (m *CreatePoolResp) String() string { return proto.CompactTextString(m) }
func (*CreatePoolResp) ProtoMessage()    {}
func (*CreatePoolResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_b3b79726342cd324, []int{1}
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

func (m *CreatePoolResp) GetStatus() DaosRequestStatus {
	if m != nil {
		return m.Status
	}
	return DaosRequestStatus_SUCCESS
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
	Force                bool     `protobuf:"varint,2,opt,name=force,proto3" json:"force,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DestroyPoolReq) Reset()         { *m = DestroyPoolReq{} }
func (m *DestroyPoolReq) String() string { return proto.CompactTextString(m) }
func (*DestroyPoolReq) ProtoMessage()    {}
func (*DestroyPoolReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_b3b79726342cd324, []int{2}
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

func (m *DestroyPoolReq) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

// DestroyPoolResp returns status of processed request.
type DestroyPoolResp struct {
	Status               DaosRequestStatus `protobuf:"varint,1,opt,name=status,proto3,enum=mgmt.DaosRequestStatus" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *DestroyPoolResp) Reset()         { *m = DestroyPoolResp{} }
func (m *DestroyPoolResp) String() string { return proto.CompactTextString(m) }
func (*DestroyPoolResp) ProtoMessage()    {}
func (*DestroyPoolResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_pool_b3b79726342cd324, []int{3}
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

func (m *DestroyPoolResp) GetStatus() DaosRequestStatus {
	if m != nil {
		return m.Status
	}
	return DaosRequestStatus_SUCCESS
}

func init() {
	proto.RegisterType((*CreatePoolReq)(nil), "mgmt.CreatePoolReq")
	proto.RegisterType((*CreatePoolResp)(nil), "mgmt.CreatePoolResp")
	proto.RegisterType((*DestroyPoolReq)(nil), "mgmt.DestroyPoolReq")
	proto.RegisterType((*DestroyPoolResp)(nil), "mgmt.DestroyPoolResp")
}

func init() { proto.RegisterFile("pool.proto", fileDescriptor_pool_b3b79726342cd324) }

var fileDescriptor_pool_b3b79726342cd324 = []byte{
	// 278 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x51, 0xcb, 0x4e, 0xc3, 0x30,
	0x10, 0x94, 0x69, 0xfa, 0x5a, 0xa9, 0x01, 0x59, 0x95, 0xb0, 0x2a, 0x84, 0xa2, 0x9c, 0x72, 0x0a,
	0x12, 0xdc, 0x38, 0x42, 0x3f, 0x00, 0x99, 0x2f, 0x48, 0xd3, 0xa5, 0xaa, 0x68, 0xb2, 0xa9, 0xd7,
	0x8e, 0x94, 0xef, 0xe3, 0xc7, 0x50, 0x6c, 0x42, 0xc3, 0x95, 0x93, 0x77, 0x66, 0x3c, 0xf6, 0x8c,
	0x0d, 0xd0, 0x10, 0x9d, 0xf2, 0xc6, 0x90, 0x25, 0x19, 0x55, 0x87, 0xca, 0x6e, 0x96, 0x6c, 0xda,
	0x40, 0xa4, 0x5f, 0x02, 0x56, 0xaf, 0x06, 0x0b, 0x8b, 0x6f, 0x44, 0x27, 0x8d, 0x67, 0xb9, 0x81,
	0x05, 0x97, 0xd5, 0xae, 0xb3, 0xc8, 0x4a, 0x24, 0x22, 0x8b, 0xf4, 0x2f, 0x96, 0x77, 0xb0, 0xac,
	0xdb, 0x0a, 0x83, 0x78, 0xe5, 0xc5, 0x0b, 0x21, 0xd7, 0x30, 0x35, 0x45, 0xfd, 0xc9, 0x6a, 0x92,
	0x88, 0x6c, 0xa9, 0x03, 0x90, 0xf7, 0x00, 0xb5, 0xab, 0xb8, 0x2d, 0x0d, 0x36, 0xac, 0xa2, 0x44,
	0x64, 0x2b, 0x3d, 0x62, 0xa4, 0x84, 0xc8, 0x31, 0x1a, 0x35, 0xf5, 0x26, 0x3f, 0xf7, 0xf7, 0xf4,
	0xeb, 0xc1, 0x90, 0x6b, 0xd4, 0xcc, 0x0b, 0x17, 0x42, 0xde, 0xc0, 0x84, 0x3b, 0x56, 0x73, 0xcf,
	0xf7, 0x63, 0x4a, 0x10, 0x8f, 0x4b, 0x70, 0x23, 0x1f, 0x60, 0xc6, 0xb6, 0xb0, 0x2e, 0x74, 0x88,
	0x1f, 0x6f, 0xf3, 0xbe, 0x79, 0xbe, 0x2d, 0x88, 0x35, 0x9e, 0x1d, 0xb2, 0x7d, 0xf7, 0xb2, 0xfe,
	0xd9, 0xe6, 0x63, 0xb8, 0xe3, 0xde, 0xb7, 0xea, 0x63, 0xb8, 0xe3, 0x5e, 0x2a, 0x98, 0x0f, 0xb9,
	0x43, 0xa5, 0x01, 0xa6, 0xcf, 0x10, 0x6f, 0x91, 0xad, 0xa1, 0x6e, 0x78, 0xb6, 0xc1, 0x2f, 0x46,
	0xfe, 0x35, 0x4c, 0x3f, 0xc8, 0x94, 0xe8, 0x0f, 0x5d, 0xe8, 0x00, 0xd2, 0x17, 0xb8, 0xfe, 0xe3,
	0xfd, 0x47, 0xda, 0xdd, 0xcc, 0xff, 0xde, 0xd3, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x50, 0xc7,
	0x70, 0xd8, 0xdc, 0x01, 0x00, 0x00,
}
