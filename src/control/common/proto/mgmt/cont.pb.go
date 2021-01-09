// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cont.proto

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

// ContSetOwnerReq supplies new pool parameters.
type ContSetOwnerReq struct {
	Sys                  string   `protobuf:"bytes,1,opt,name=sys,proto3" json:"sys,omitempty"`
	ContUUID             string   `protobuf:"bytes,2,opt,name=contUUID,proto3" json:"contUUID,omitempty"`
	PoolUUID             string   `protobuf:"bytes,3,opt,name=poolUUID,proto3" json:"poolUUID,omitempty"`
	Owneruser            string   `protobuf:"bytes,4,opt,name=owneruser,proto3" json:"owneruser,omitempty"`
	Ownergroup           string   `protobuf:"bytes,5,opt,name=ownergroup,proto3" json:"ownergroup,omitempty"`
	SvcRanks             []uint32 `protobuf:"varint,6,rep,packed,name=svc_ranks,json=svcRanks,proto3" json:"svc_ranks,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContSetOwnerReq) Reset()         { *m = ContSetOwnerReq{} }
func (m *ContSetOwnerReq) String() string { return proto.CompactTextString(m) }
func (*ContSetOwnerReq) ProtoMessage()    {}
func (*ContSetOwnerReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e578f20ac30c5ed, []int{0}
}

func (m *ContSetOwnerReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContSetOwnerReq.Unmarshal(m, b)
}
func (m *ContSetOwnerReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContSetOwnerReq.Marshal(b, m, deterministic)
}
func (m *ContSetOwnerReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContSetOwnerReq.Merge(m, src)
}
func (m *ContSetOwnerReq) XXX_Size() int {
	return xxx_messageInfo_ContSetOwnerReq.Size(m)
}
func (m *ContSetOwnerReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ContSetOwnerReq.DiscardUnknown(m)
}

var xxx_messageInfo_ContSetOwnerReq proto.InternalMessageInfo

func (m *ContSetOwnerReq) GetSys() string {
	if m != nil {
		return m.Sys
	}
	return ""
}

func (m *ContSetOwnerReq) GetContUUID() string {
	if m != nil {
		return m.ContUUID
	}
	return ""
}

func (m *ContSetOwnerReq) GetPoolUUID() string {
	if m != nil {
		return m.PoolUUID
	}
	return ""
}

func (m *ContSetOwnerReq) GetOwneruser() string {
	if m != nil {
		return m.Owneruser
	}
	return ""
}

func (m *ContSetOwnerReq) GetOwnergroup() string {
	if m != nil {
		return m.Ownergroup
	}
	return ""
}

func (m *ContSetOwnerReq) GetSvcRanks() []uint32 {
	if m != nil {
		return m.SvcRanks
	}
	return nil
}

// ContSetOwnerResp returns created pool uuid and ranks.
type ContSetOwnerResp struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContSetOwnerResp) Reset()         { *m = ContSetOwnerResp{} }
func (m *ContSetOwnerResp) String() string { return proto.CompactTextString(m) }
func (*ContSetOwnerResp) ProtoMessage()    {}
func (*ContSetOwnerResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_7e578f20ac30c5ed, []int{1}
}

func (m *ContSetOwnerResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContSetOwnerResp.Unmarshal(m, b)
}
func (m *ContSetOwnerResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContSetOwnerResp.Marshal(b, m, deterministic)
}
func (m *ContSetOwnerResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContSetOwnerResp.Merge(m, src)
}
func (m *ContSetOwnerResp) XXX_Size() int {
	return xxx_messageInfo_ContSetOwnerResp.Size(m)
}
func (m *ContSetOwnerResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ContSetOwnerResp.DiscardUnknown(m)
}

var xxx_messageInfo_ContSetOwnerResp proto.InternalMessageInfo

func (m *ContSetOwnerResp) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func init() {
	proto.RegisterType((*ContSetOwnerReq)(nil), "mgmt.ContSetOwnerReq")
	proto.RegisterType((*ContSetOwnerResp)(nil), "mgmt.ContSetOwnerResp")
}

func init() {
	proto.RegisterFile("cont.proto", fileDescriptor_7e578f20ac30c5ed)
}

var fileDescriptor_7e578f20ac30c5ed = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0xce, 0xcf, 0x2b,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xc9, 0x4d, 0xcf, 0x2d, 0x51, 0xda, 0xc6, 0xc8,
	0xc5, 0xef, 0x9c, 0x9f, 0x57, 0x12, 0x9c, 0x5a, 0xe2, 0x5f, 0x9e, 0x97, 0x5a, 0x14, 0x94, 0x5a,
	0x28, 0x24, 0xc0, 0xc5, 0x5c, 0x5c, 0x59, 0x2c, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0x62,
	0x0a, 0x49, 0x71, 0x71, 0x80, 0x74, 0x86, 0x86, 0x7a, 0xba, 0x48, 0x30, 0x81, 0x85, 0xe1, 0x7c,
	0x90, 0x5c, 0x41, 0x7e, 0x7e, 0x0e, 0x58, 0x8e, 0x19, 0x22, 0x07, 0xe3, 0x0b, 0xc9, 0x70, 0x71,
	0xe6, 0x83, 0x4c, 0x2d, 0x2d, 0x4e, 0x2d, 0x92, 0x60, 0x01, 0x4b, 0x22, 0x04, 0x84, 0xe4, 0xb8,
	0xb8, 0xc0, 0x9c, 0xf4, 0xa2, 0xfc, 0xd2, 0x02, 0x09, 0x56, 0xb0, 0x34, 0x92, 0x88, 0x90, 0x34,
	0x17, 0x67, 0x71, 0x59, 0x72, 0x7c, 0x51, 0x62, 0x5e, 0x76, 0xb1, 0x04, 0x9b, 0x02, 0xb3, 0x06,
	0x6f, 0x10, 0x47, 0x71, 0x59, 0x72, 0x10, 0x88, 0xaf, 0xa4, 0xc5, 0x25, 0x80, 0xea, 0xee, 0xe2,
	0x02, 0x21, 0x31, 0x2e, 0xb6, 0xe2, 0x92, 0xc4, 0x92, 0x52, 0x88, 0xdb, 0x59, 0x83, 0xa0, 0xbc,
	0x24, 0x36, 0xb0, 0x8f, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc6, 0xbe, 0x8a, 0xbd, 0xff,
	0x00, 0x00, 0x00,
}
