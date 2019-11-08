// Code generated by protoc-gen-go. DO NOT EDIT.
// source: features.proto

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

type FeatureName struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FeatureName) Reset()         { *m = FeatureName{} }
func (m *FeatureName) String() string { return proto.CompactTextString(m) }
func (*FeatureName) ProtoMessage()    {}
func (*FeatureName) Descriptor() ([]byte, []int) {
	return fileDescriptor_2216f05915163cdf, []int{0}
}

func (m *FeatureName) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureName.Unmarshal(m, b)
}
func (m *FeatureName) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureName.Marshal(b, m, deterministic)
}
func (m *FeatureName) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureName.Merge(m, src)
}
func (m *FeatureName) XXX_Size() int {
	return xxx_messageInfo_FeatureName.Size(m)
}
func (m *FeatureName) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureName.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureName proto.InternalMessageInfo

func (m *FeatureName) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Category struct {
	Category             string   `protobuf:"bytes,1,opt,name=category,proto3" json:"category,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Category) Reset()         { *m = Category{} }
func (m *Category) String() string { return proto.CompactTextString(m) }
func (*Category) ProtoMessage()    {}
func (*Category) Descriptor() ([]byte, []int) {
	return fileDescriptor_2216f05915163cdf, []int{1}
}

func (m *Category) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Category.Unmarshal(m, b)
}
func (m *Category) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Category.Marshal(b, m, deterministic)
}
func (m *Category) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Category.Merge(m, src)
}
func (m *Category) XXX_Size() int {
	return xxx_messageInfo_Category.Size(m)
}
func (m *Category) XXX_DiscardUnknown() {
	xxx_messageInfo_Category.DiscardUnknown(m)
}

var xxx_messageInfo_Category proto.InternalMessageInfo

func (m *Category) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

// Feature represents a management task that can be performed by server.
type Feature struct {
	Category             *Category    `protobuf:"bytes,1,opt,name=category,proto3" json:"category,omitempty"`
	Fname                *FeatureName `protobuf:"bytes,2,opt,name=fname,proto3" json:"fname,omitempty"`
	Description          string       `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Feature) Reset()         { *m = Feature{} }
func (m *Feature) String() string { return proto.CompactTextString(m) }
func (*Feature) ProtoMessage()    {}
func (*Feature) Descriptor() ([]byte, []int) {
	return fileDescriptor_2216f05915163cdf, []int{2}
}

func (m *Feature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Feature.Unmarshal(m, b)
}
func (m *Feature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Feature.Marshal(b, m, deterministic)
}
func (m *Feature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Feature.Merge(m, src)
}
func (m *Feature) XXX_Size() int {
	return xxx_messageInfo_Feature.Size(m)
}
func (m *Feature) XXX_DiscardUnknown() {
	xxx_messageInfo_Feature.DiscardUnknown(m)
}

var xxx_messageInfo_Feature proto.InternalMessageInfo

func (m *Feature) GetCategory() *Category {
	if m != nil {
		return m.Category
	}
	return nil
}

func (m *Feature) GetFname() *FeatureName {
	if m != nil {
		return m.Fname
	}
	return nil
}

func (m *Feature) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterType((*FeatureName)(nil), "ctl.FeatureName")
	proto.RegisterType((*Category)(nil), "ctl.Category")
	proto.RegisterType((*Feature)(nil), "ctl.Feature")
}

func init() { proto.RegisterFile("features.proto", fileDescriptor_2216f05915163cdf) }

var fileDescriptor_2216f05915163cdf = []byte{
	// 162 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0x4b, 0x4d, 0x2c,
	0x29, 0x2d, 0x4a, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x2e, 0xc9, 0x51,
	0x52, 0xe4, 0xe2, 0x76, 0x83, 0x08, 0xfb, 0x25, 0xe6, 0xa6, 0x0a, 0x09, 0x71, 0xb1, 0xe4, 0x25,
	0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x4a, 0x6a, 0x5c, 0x1c, 0xce,
	0x89, 0x25, 0xa9, 0xe9, 0xf9, 0x45, 0x95, 0x42, 0x52, 0x5c, 0x1c, 0xc9, 0x50, 0x36, 0x54, 0x0d,
	0x9c, 0xaf, 0x54, 0xc7, 0xc5, 0x0e, 0x35, 0x4a, 0x48, 0x13, 0x4d, 0x19, 0xb7, 0x11, 0xaf, 0x5e,
	0x72, 0x49, 0x8e, 0x1e, 0xcc, 0x1c, 0x84, 0x2e, 0x21, 0x35, 0x2e, 0xd6, 0x34, 0xb0, 0x95, 0x4c,
	0x60, 0x75, 0x02, 0x60, 0x75, 0x48, 0x4e, 0x0a, 0x82, 0x48, 0x0b, 0x29, 0x70, 0x71, 0xa7, 0xa4,
	0x16, 0x27, 0x17, 0x65, 0x16, 0x94, 0x64, 0xe6, 0xe7, 0x49, 0x30, 0x83, 0x2d, 0x47, 0x16, 0x4a,
	0x62, 0x03, 0x7b, 0xcb, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x67, 0xa6, 0x50, 0xa1, 0xe8, 0x00,
	0x00, 0x00,
}
