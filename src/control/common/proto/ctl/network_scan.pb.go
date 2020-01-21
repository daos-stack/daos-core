// Code generated by protoc-gen-go. DO NOT EDIT.
// source: network_scan.proto

package ctl

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

type DeviceScanRequest struct {
	Provider             string   `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeviceScanRequest) Reset()         { *m = DeviceScanRequest{} }
func (m *DeviceScanRequest) String() string { return proto.CompactTextString(m) }
func (*DeviceScanRequest) ProtoMessage()    {}
func (*DeviceScanRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cbceb864bca85ee, []int{0}
}

func (m *DeviceScanRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeviceScanRequest.Unmarshal(m, b)
}
func (m *DeviceScanRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeviceScanRequest.Marshal(b, m, deterministic)
}
func (m *DeviceScanRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeviceScanRequest.Merge(m, src)
}
func (m *DeviceScanRequest) XXX_Size() int {
	return xxx_messageInfo_DeviceScanRequest.Size(m)
}
func (m *DeviceScanRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeviceScanRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeviceScanRequest proto.InternalMessageInfo

func (m *DeviceScanRequest) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

type DeviceScanReply struct {
	Provider             string   `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	Device               string   `protobuf:"bytes,2,opt,name=device,proto3" json:"device,omitempty"`
	Numanode             uint32   `protobuf:"varint,3,opt,name=numanode,proto3" json:"numanode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeviceScanReply) Reset()         { *m = DeviceScanReply{} }
func (m *DeviceScanReply) String() string { return proto.CompactTextString(m) }
func (*DeviceScanReply) ProtoMessage()    {}
func (*DeviceScanReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cbceb864bca85ee, []int{1}
}

func (m *DeviceScanReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeviceScanReply.Unmarshal(m, b)
}
func (m *DeviceScanReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeviceScanReply.Marshal(b, m, deterministic)
}
func (m *DeviceScanReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeviceScanReply.Merge(m, src)
}
func (m *DeviceScanReply) XXX_Size() int {
	return xxx_messageInfo_DeviceScanReply.Size(m)
}
func (m *DeviceScanReply) XXX_DiscardUnknown() {
	xxx_messageInfo_DeviceScanReply.DiscardUnknown(m)
}

var xxx_messageInfo_DeviceScanReply proto.InternalMessageInfo

func (m *DeviceScanReply) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *DeviceScanReply) GetDevice() string {
	if m != nil {
		return m.Device
	}
	return ""
}

func (m *DeviceScanReply) GetNumanode() uint32 {
	if m != nil {
		return m.Numanode
	}
	return 0
}

type ProviderListRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProviderListRequest) Reset()         { *m = ProviderListRequest{} }
func (m *ProviderListRequest) String() string { return proto.CompactTextString(m) }
func (*ProviderListRequest) ProtoMessage()    {}
func (*ProviderListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cbceb864bca85ee, []int{2}
}

func (m *ProviderListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProviderListRequest.Unmarshal(m, b)
}
func (m *ProviderListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProviderListRequest.Marshal(b, m, deterministic)
}
func (m *ProviderListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProviderListRequest.Merge(m, src)
}
func (m *ProviderListRequest) XXX_Size() int {
	return xxx_messageInfo_ProviderListRequest.Size(m)
}
func (m *ProviderListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProviderListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProviderListRequest proto.InternalMessageInfo

type ProviderListReply struct {
	Provider             string   `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProviderListReply) Reset()         { *m = ProviderListReply{} }
func (m *ProviderListReply) String() string { return proto.CompactTextString(m) }
func (*ProviderListReply) ProtoMessage()    {}
func (*ProviderListReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cbceb864bca85ee, []int{3}
}

func (m *ProviderListReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProviderListReply.Unmarshal(m, b)
}
func (m *ProviderListReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProviderListReply.Marshal(b, m, deterministic)
}
func (m *ProviderListReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProviderListReply.Merge(m, src)
}
func (m *ProviderListReply) XXX_Size() int {
	return xxx_messageInfo_ProviderListReply.Size(m)
}
func (m *ProviderListReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ProviderListReply.DiscardUnknown(m)
}

var xxx_messageInfo_ProviderListReply proto.InternalMessageInfo

func (m *ProviderListReply) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func init() {
	proto.RegisterType((*DeviceScanRequest)(nil), "ctl.DeviceScanRequest")
	proto.RegisterType((*DeviceScanReply)(nil), "ctl.DeviceScanReply")
	proto.RegisterType((*ProviderListRequest)(nil), "ctl.ProviderListRequest")
	proto.RegisterType((*ProviderListReply)(nil), "ctl.ProviderListReply")
}

func init() { proto.RegisterFile("network_scan.proto", fileDescriptor_2cbceb864bca85ee) }

var fileDescriptor_2cbceb864bca85ee = []byte{
	// 159 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xca, 0x4b, 0x2d, 0x29,
	0xcf, 0x2f, 0xca, 0x8e, 0x2f, 0x4e, 0x4e, 0xcc, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x4e, 0x2e, 0xc9, 0x51, 0xd2, 0xe7, 0x12, 0x74, 0x49, 0x2d, 0xcb, 0x4c, 0x4e, 0x0d, 0x4e, 0x4e,
	0xcc, 0x0b, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x92, 0xe2, 0xe2, 0x28, 0x28, 0xca, 0x2f,
	0xcb, 0x4c, 0x49, 0x2d, 0x92, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0xf3, 0x95, 0x12, 0xb9,
	0xf8, 0x91, 0x35, 0x14, 0xe4, 0x54, 0xe2, 0x53, 0x2e, 0x24, 0xc6, 0xc5, 0x96, 0x02, 0x56, 0x2e,
	0xc1, 0x04, 0x96, 0x81, 0xf2, 0x40, 0x7a, 0xf2, 0x4a, 0x73, 0x13, 0xf3, 0xf2, 0x53, 0x52, 0x25,
	0x98, 0x15, 0x18, 0x35, 0x78, 0x83, 0xe0, 0x7c, 0x25, 0x51, 0x2e, 0xe1, 0x00, 0xa8, 0x7e, 0x9f,
	0xcc, 0xe2, 0x12, 0xa8, 0xab, 0x40, 0x4e, 0x45, 0x15, 0x26, 0x60, 0x77, 0x12, 0x1b, 0xd8, 0x9f,
	0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xad, 0x8f, 0x1f, 0x65, 0xfd, 0x00, 0x00, 0x00,
}
