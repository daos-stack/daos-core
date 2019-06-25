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
	return fileDescriptor_pool_537e58f36e0e7747, []int{0}
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
	return fileDescriptor_pool_537e58f36e0e7747, []int{1}
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
	return fileDescriptor_pool_537e58f36e0e7747, []int{2}
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

func init() {
	proto.RegisterType((*CreatePoolReq)(nil), "mgmt.CreatePoolReq")
	proto.RegisterType((*CreatePoolResp)(nil), "mgmt.CreatePoolResp")
	proto.RegisterType((*DestroyPoolReq)(nil), "mgmt.DestroyPoolReq")
}

func init() { proto.RegisterFile("pool.proto", fileDescriptor_pool_537e58f36e0e7747) }

var fileDescriptor_pool_537e58f36e0e7747 = []byte{
	// 272 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x51, 0xcb, 0x6a, 0xc3, 0x30,
	0x10, 0x44, 0x89, 0xf3, 0x5a, 0x48, 0x28, 0xa2, 0x50, 0x11, 0x4a, 0x31, 0x3e, 0xf9, 0xe4, 0x42,
	0xfb, 0x09, 0xcd, 0xb1, 0x87, 0xa2, 0x7e, 0x81, 0xe3, 0x6c, 0x43, 0x68, 0xec, 0x75, 0xb4, 0x92,
	0xc1, 0xdf, 0xd7, 0x1f, 0x2b, 0x92, 0xe2, 0xc4, 0x27, 0xef, 0xcc, 0x78, 0xa4, 0x99, 0x15, 0x40,
	0x4b, 0x74, 0x2e, 0x5a, 0x43, 0x96, 0x64, 0x52, 0x1f, 0x6b, 0xbb, 0x5d, 0xb1, 0xe9, 0x22, 0x91,
	0xfd, 0x09, 0x58, 0x7f, 0x18, 0x2c, 0x2d, 0x7e, 0x11, 0x9d, 0x35, 0x5e, 0xe4, 0x16, 0x96, 0x5c,
	0xd5, 0xfb, 0xde, 0x22, 0x2b, 0x91, 0x8a, 0x3c, 0xd1, 0x37, 0x2c, 0x9f, 0x61, 0xd5, 0x74, 0x35,
	0x46, 0x71, 0x12, 0xc4, 0x3b, 0x21, 0x1f, 0x61, 0x66, 0xca, 0xe6, 0x97, 0xd5, 0x34, 0x15, 0xf9,
	0x4a, 0x47, 0x20, 0x5f, 0x00, 0x1a, 0x57, 0x73, 0x57, 0x19, 0x6c, 0x59, 0x25, 0xa9, 0xc8, 0xd7,
	0x7a, 0xc4, 0x48, 0x09, 0x89, 0x63, 0x34, 0x6a, 0x16, 0x4c, 0x61, 0xf6, 0xf7, 0xf8, 0xef, 0xd1,
	0x90, 0x6b, 0xd5, 0x3c, 0x08, 0x77, 0x42, 0x3e, 0xc0, 0x94, 0x7b, 0x56, 0x8b, 0xc0, 0xfb, 0x31,
	0x23, 0xd8, 0x8c, 0x4b, 0x70, 0x2b, 0x5f, 0x61, 0xce, 0xb6, 0xb4, 0x2e, 0x76, 0xd8, 0xbc, 0x3d,
	0x15, 0xbe, 0x79, 0xb1, 0x2b, 0x89, 0x35, 0x5e, 0x1c, 0xb2, 0xfd, 0x0e, 0xb2, 0xbe, 0xfe, 0x16,
	0x62, 0xb8, 0xd3, 0x21, 0xb4, 0xf2, 0x31, 0xdc, 0xe9, 0x20, 0x15, 0x2c, 0x86, 0xdc, 0xb1, 0xd2,
	0x00, 0xb3, 0x4f, 0xd8, 0xec, 0x90, 0xad, 0xa1, 0x7e, 0x58, 0xdb, 0xe0, 0x17, 0x23, 0xff, 0x35,
	0xe8, 0xe4, 0x16, 0xd4, 0xaf, 0xe8, 0x87, 0x4c, 0x85, 0xe1, 0xbc, 0xa5, 0x8e, 0x60, 0x3f, 0x0f,
	0x6f, 0xf1, 0xfe, 0x1f, 0x00, 0x00, 0xff, 0xff, 0xcf, 0x5f, 0x5f, 0xba, 0xaa, 0x01, 0x00, 0x00,
}
