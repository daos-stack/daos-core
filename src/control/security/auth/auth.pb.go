// Code generated by protoc-gen-go. DO NOT EDIT.
// source: auth.proto

package auth

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

// Authentication token includes a packed structure of the specified flavor
type Flavor int32

const (
	Flavor_AUTH_NONE Flavor = 0
	Flavor_AUTH_SYS  Flavor = 1
)

var Flavor_name = map[int32]string{
	0: "AUTH_NONE",
	1: "AUTH_SYS",
}

var Flavor_value = map[string]int32{
	"AUTH_NONE": 0,
	"AUTH_SYS":  1,
}

func (x Flavor) String() string {
	return proto.EnumName(Flavor_name, int32(x))
}

func (Flavor) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{0}
}

type Token struct {
	Flavor               Flavor   `protobuf:"varint,1,opt,name=flavor,proto3,enum=auth.Flavor" json:"flavor,omitempty"`
	Data                 []byte   `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{0}
}

func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (m *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(m, src)
}
func (m *Token) XXX_Size() int {
	return xxx_messageInfo_Token.Size(m)
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetFlavor() Flavor {
	if m != nil {
		return m.Flavor
	}
	return Flavor_AUTH_NONE
}

func (m *Token) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// Token structure for AUTH_SYS flavor
type Sys struct {
	Stamp                uint64   `protobuf:"varint,1,opt,name=stamp,proto3" json:"stamp,omitempty"`
	Machinename          string   `protobuf:"bytes,2,opt,name=machinename,proto3" json:"machinename,omitempty"`
	User                 string   `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	Group                string   `protobuf:"bytes,4,opt,name=group,proto3" json:"group,omitempty"`
	Groups               []string `protobuf:"bytes,5,rep,name=groups,proto3" json:"groups,omitempty"`
	Secctx               string   `protobuf:"bytes,6,opt,name=secctx,proto3" json:"secctx,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Sys) Reset()         { *m = Sys{} }
func (m *Sys) String() string { return proto.CompactTextString(m) }
func (*Sys) ProtoMessage()    {}
func (*Sys) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{1}
}

func (m *Sys) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Sys.Unmarshal(m, b)
}
func (m *Sys) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Sys.Marshal(b, m, deterministic)
}
func (m *Sys) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Sys.Merge(m, src)
}
func (m *Sys) XXX_Size() int {
	return xxx_messageInfo_Sys.Size(m)
}
func (m *Sys) XXX_DiscardUnknown() {
	xxx_messageInfo_Sys.DiscardUnknown(m)
}

var xxx_messageInfo_Sys proto.InternalMessageInfo

func (m *Sys) GetStamp() uint64 {
	if m != nil {
		return m.Stamp
	}
	return 0
}

func (m *Sys) GetMachinename() string {
	if m != nil {
		return m.Machinename
	}
	return ""
}

func (m *Sys) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *Sys) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *Sys) GetGroups() []string {
	if m != nil {
		return m.Groups
	}
	return nil
}

func (m *Sys) GetSecctx() string {
	if m != nil {
		return m.Secctx
	}
	return ""
}

type SysVerifier struct {
	Signature            []byte   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SysVerifier) Reset()         { *m = SysVerifier{} }
func (m *SysVerifier) String() string { return proto.CompactTextString(m) }
func (*SysVerifier) ProtoMessage()    {}
func (*SysVerifier) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{2}
}

func (m *SysVerifier) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SysVerifier.Unmarshal(m, b)
}
func (m *SysVerifier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SysVerifier.Marshal(b, m, deterministic)
}
func (m *SysVerifier) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SysVerifier.Merge(m, src)
}
func (m *SysVerifier) XXX_Size() int {
	return xxx_messageInfo_SysVerifier.Size(m)
}
func (m *SysVerifier) XXX_DiscardUnknown() {
	xxx_messageInfo_SysVerifier.DiscardUnknown(m)
}

var xxx_messageInfo_SysVerifier proto.InternalMessageInfo

func (m *SysVerifier) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

// SecurityCredential includes the auth token and a verifier that can be used by
// the server to verify the integrity of the token.
// Token and verifier are expected to have the same flavor type.
type Credential struct {
	Token                *Token   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Verifier             *Token   `protobuf:"bytes,2,opt,name=verifier,proto3" json:"verifier,omitempty"`
	Origin               string   `protobuf:"bytes,3,opt,name=origin,proto3" json:"origin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Credential) Reset()         { *m = Credential{} }
func (m *Credential) String() string { return proto.CompactTextString(m) }
func (*Credential) ProtoMessage()    {}
func (*Credential) Descriptor() ([]byte, []int) {
	return fileDescriptor_8bbd6f3875b0e874, []int{3}
}

func (m *Credential) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Credential.Unmarshal(m, b)
}
func (m *Credential) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Credential.Marshal(b, m, deterministic)
}
func (m *Credential) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Credential.Merge(m, src)
}
func (m *Credential) XXX_Size() int {
	return xxx_messageInfo_Credential.Size(m)
}
func (m *Credential) XXX_DiscardUnknown() {
	xxx_messageInfo_Credential.DiscardUnknown(m)
}

var xxx_messageInfo_Credential proto.InternalMessageInfo

func (m *Credential) GetToken() *Token {
	if m != nil {
		return m.Token
	}
	return nil
}

func (m *Credential) GetVerifier() *Token {
	if m != nil {
		return m.Verifier
	}
	return nil
}

func (m *Credential) GetOrigin() string {
	if m != nil {
		return m.Origin
	}
	return ""
}

func init() {
	proto.RegisterEnum("auth.Flavor", Flavor_name, Flavor_value)
	proto.RegisterType((*Token)(nil), "auth.Token")
	proto.RegisterType((*Sys)(nil), "auth.Sys")
	proto.RegisterType((*SysVerifier)(nil), "auth.SysVerifier")
	proto.RegisterType((*Credential)(nil), "auth.Credential")
}

func init() { proto.RegisterFile("auth.proto", fileDescriptor_8bbd6f3875b0e874) }

var fileDescriptor_8bbd6f3875b0e874 = []byte{
	// 300 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0xdd, 0x4a, 0xc3, 0x30,
	0x14, 0xb6, 0xae, 0x2d, 0xeb, 0x69, 0x95, 0x71, 0x10, 0xe9, 0x85, 0x17, 0xb5, 0x28, 0x0e, 0x85,
	0x5d, 0xcc, 0x27, 0x18, 0xa2, 0x78, 0x35, 0x21, 0x9d, 0x82, 0x57, 0x12, 0xb7, 0x6c, 0x0b, 0x6e,
	0x4d, 0x49, 0xd2, 0xe1, 0x9e, 0xc4, 0xd7, 0x95, 0x9c, 0x06, 0x15, 0xef, 0xbe, 0xbf, 0x7e, 0xe4,
	0x7c, 0x05, 0xe0, 0xad, 0x5d, 0x8f, 0x1a, 0xad, 0xac, 0xc2, 0xd0, 0xe1, 0x72, 0x02, 0xd1, 0x4c,
	0x7d, 0x88, 0x1a, 0x2f, 0x20, 0x5e, 0x6e, 0xf8, 0x4e, 0xe9, 0x3c, 0x28, 0x82, 0xe1, 0xf1, 0x38,
	0x1b, 0x51, 0xf6, 0x81, 0x34, 0xe6, 0x3d, 0x44, 0x08, 0x17, 0xdc, 0xf2, 0xfc, 0xb0, 0x08, 0x86,
	0x19, 0x23, 0x5c, 0x7e, 0x05, 0xd0, 0xab, 0xf6, 0x06, 0x4f, 0x20, 0x32, 0x96, 0x6f, 0x1b, 0x2a,
	0x08, 0x59, 0x47, 0xb0, 0x80, 0x74, 0xcb, 0xe7, 0x6b, 0x59, 0x8b, 0x9a, 0x6f, 0x05, 0x7d, 0x98,
	0xb0, 0xbf, 0x92, 0xeb, 0x6c, 0x8d, 0xd0, 0x79, 0x8f, 0x2c, 0xc2, 0xae, 0x6b, 0xa5, 0x55, 0xdb,
	0xe4, 0x21, 0x89, 0x1d, 0xc1, 0x53, 0x88, 0x09, 0x98, 0x3c, 0x2a, 0x7a, 0xc3, 0x84, 0x79, 0xe6,
	0x74, 0x23, 0xe6, 0x73, 0xfb, 0x99, 0xc7, 0x14, 0xf7, 0xac, 0xbc, 0x81, 0xb4, 0xda, 0x9b, 0x17,
	0xa1, 0xe5, 0x52, 0x0a, 0x8d, 0x67, 0x90, 0x18, 0xb9, 0xaa, 0xb9, 0x6d, 0xb5, 0xa0, 0x47, 0x66,
	0xec, 0x57, 0x28, 0x1b, 0x80, 0x3b, 0x2d, 0x16, 0xa2, 0xb6, 0x92, 0x6f, 0xf0, 0x1c, 0x22, 0xeb,
	0x76, 0xa1, 0x5c, 0x3a, 0x4e, 0xbb, 0x35, 0x68, 0x2a, 0xd6, 0x39, 0x78, 0x05, 0xfd, 0x9d, 0xaf,
	0xa6, 0xb3, 0xfe, 0xa5, 0x7e, 0x4c, 0xf7, 0x3c, 0xa5, 0xe5, 0x4a, 0xd6, 0xfe, 0x44, 0xcf, 0xae,
	0x2f, 0x21, 0xee, 0xe6, 0xc5, 0x23, 0x48, 0x26, 0xcf, 0xb3, 0xc7, 0xb7, 0xe9, 0xd3, 0xf4, 0x7e,
	0x70, 0x80, 0x19, 0xf4, 0x89, 0x56, 0xaf, 0xd5, 0x20, 0x78, 0x8f, 0xe9, 0x7f, 0xdd, 0x7e, 0x07,
	0x00, 0x00, 0xff, 0xff, 0x40, 0x3e, 0xe9, 0x5a, 0xbd, 0x01, 0x00, 0x00,
}
