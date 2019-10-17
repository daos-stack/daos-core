// Code generated by protoc-gen-go. DO NOT EDIT.
// source: system.proto

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

// SystemMember refers to a data-plane instance that is a member of DAOS
// system running on host with the control-plane listening at "Addr".
type SystemMember struct {
	Addr                 string   `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Uuid                 string   `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Rank                 uint32   `protobuf:"varint,3,opt,name=rank,proto3" json:"rank,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SystemMember) Reset()         { *m = SystemMember{} }
func (m *SystemMember) String() string { return proto.CompactTextString(m) }
func (*SystemMember) ProtoMessage()    {}
func (*SystemMember) Descriptor() ([]byte, []int) {
	return fileDescriptor_system_15d58b874904f062, []int{0}
}
func (m *SystemMember) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemMember.Unmarshal(m, b)
}
func (m *SystemMember) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemMember.Marshal(b, m, deterministic)
}
func (dst *SystemMember) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemMember.Merge(dst, src)
}
func (m *SystemMember) XXX_Size() int {
	return xxx_messageInfo_SystemMember.Size(m)
}
func (m *SystemMember) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemMember.DiscardUnknown(m)
}

var xxx_messageInfo_SystemMember proto.InternalMessageInfo

func (m *SystemMember) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *SystemMember) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *SystemMember) GetRank() uint32 {
	if m != nil {
		return m.Rank
	}
	return 0
}

// SystemStopReq supplies system shutdown parameters.
type SystemStopReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SystemStopReq) Reset()         { *m = SystemStopReq{} }
func (m *SystemStopReq) String() string { return proto.CompactTextString(m) }
func (*SystemStopReq) ProtoMessage()    {}
func (*SystemStopReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_system_15d58b874904f062, []int{1}
}
func (m *SystemStopReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemStopReq.Unmarshal(m, b)
}
func (m *SystemStopReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemStopReq.Marshal(b, m, deterministic)
}
func (dst *SystemStopReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemStopReq.Merge(dst, src)
}
func (m *SystemStopReq) XXX_Size() int {
	return xxx_messageInfo_SystemStopReq.Size(m)
}
func (m *SystemStopReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemStopReq.DiscardUnknown(m)
}

var xxx_messageInfo_SystemStopReq proto.InternalMessageInfo

// SystemStopResp returns status of shutdown attempt and results
// of attempts to stop system members.
type SystemStopResp struct {
	Results              []*SystemStopResp_Result `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *SystemStopResp) Reset()         { *m = SystemStopResp{} }
func (m *SystemStopResp) String() string { return proto.CompactTextString(m) }
func (*SystemStopResp) ProtoMessage()    {}
func (*SystemStopResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_system_15d58b874904f062, []int{2}
}
func (m *SystemStopResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemStopResp.Unmarshal(m, b)
}
func (m *SystemStopResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemStopResp.Marshal(b, m, deterministic)
}
func (dst *SystemStopResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemStopResp.Merge(dst, src)
}
func (m *SystemStopResp) XXX_Size() int {
	return xxx_messageInfo_SystemStopResp.Size(m)
}
func (m *SystemStopResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemStopResp.DiscardUnknown(m)
}

var xxx_messageInfo_SystemStopResp proto.InternalMessageInfo

func (m *SystemStopResp) GetResults() []*SystemStopResp_Result {
	if m != nil {
		return m.Results
	}
	return nil
}

type SystemStopResp_Result struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Action               string   `protobuf:"bytes,2,opt,name=action,proto3" json:"action,omitempty"`
	Errored              bool     `protobuf:"varint,3,opt,name=errored,proto3" json:"errored,omitempty"`
	Msg                  string   `protobuf:"bytes,4,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SystemStopResp_Result) Reset()         { *m = SystemStopResp_Result{} }
func (m *SystemStopResp_Result) String() string { return proto.CompactTextString(m) }
func (*SystemStopResp_Result) ProtoMessage()    {}
func (*SystemStopResp_Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_system_15d58b874904f062, []int{2, 0}
}
func (m *SystemStopResp_Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemStopResp_Result.Unmarshal(m, b)
}
func (m *SystemStopResp_Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemStopResp_Result.Marshal(b, m, deterministic)
}
func (dst *SystemStopResp_Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemStopResp_Result.Merge(dst, src)
}
func (m *SystemStopResp_Result) XXX_Size() int {
	return xxx_messageInfo_SystemStopResp_Result.Size(m)
}
func (m *SystemStopResp_Result) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemStopResp_Result.DiscardUnknown(m)
}

var xxx_messageInfo_SystemStopResp_Result proto.InternalMessageInfo

func (m *SystemStopResp_Result) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SystemStopResp_Result) GetAction() string {
	if m != nil {
		return m.Action
	}
	return ""
}

func (m *SystemStopResp_Result) GetErrored() bool {
	if m != nil {
		return m.Errored
	}
	return false
}

func (m *SystemStopResp_Result) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

// SystemMemberQueryReq supplies system query parameters.
type SystemMemberQueryReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SystemMemberQueryReq) Reset()         { *m = SystemMemberQueryReq{} }
func (m *SystemMemberQueryReq) String() string { return proto.CompactTextString(m) }
func (*SystemMemberQueryReq) ProtoMessage()    {}
func (*SystemMemberQueryReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_system_15d58b874904f062, []int{3}
}
func (m *SystemMemberQueryReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemMemberQueryReq.Unmarshal(m, b)
}
func (m *SystemMemberQueryReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemMemberQueryReq.Marshal(b, m, deterministic)
}
func (dst *SystemMemberQueryReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemMemberQueryReq.Merge(dst, src)
}
func (m *SystemMemberQueryReq) XXX_Size() int {
	return xxx_messageInfo_SystemMemberQueryReq.Size(m)
}
func (m *SystemMemberQueryReq) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemMemberQueryReq.DiscardUnknown(m)
}

var xxx_messageInfo_SystemMemberQueryReq proto.InternalMessageInfo

// SystemMemberQueryResp returns active system members.
type SystemMemberQueryResp struct {
	Members              []*SystemMember `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *SystemMemberQueryResp) Reset()         { *m = SystemMemberQueryResp{} }
func (m *SystemMemberQueryResp) String() string { return proto.CompactTextString(m) }
func (*SystemMemberQueryResp) ProtoMessage()    {}
func (*SystemMemberQueryResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_system_15d58b874904f062, []int{4}
}
func (m *SystemMemberQueryResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SystemMemberQueryResp.Unmarshal(m, b)
}
func (m *SystemMemberQueryResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SystemMemberQueryResp.Marshal(b, m, deterministic)
}
func (dst *SystemMemberQueryResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SystemMemberQueryResp.Merge(dst, src)
}
func (m *SystemMemberQueryResp) XXX_Size() int {
	return xxx_messageInfo_SystemMemberQueryResp.Size(m)
}
func (m *SystemMemberQueryResp) XXX_DiscardUnknown() {
	xxx_messageInfo_SystemMemberQueryResp.DiscardUnknown(m)
}

var xxx_messageInfo_SystemMemberQueryResp proto.InternalMessageInfo

func (m *SystemMemberQueryResp) GetMembers() []*SystemMember {
	if m != nil {
		return m.Members
	}
	return nil
}

func init() {
	proto.RegisterType((*SystemMember)(nil), "mgmt.SystemMember")
	proto.RegisterType((*SystemStopReq)(nil), "mgmt.SystemStopReq")
	proto.RegisterType((*SystemStopResp)(nil), "mgmt.SystemStopResp")
	proto.RegisterType((*SystemStopResp_Result)(nil), "mgmt.SystemStopResp.Result")
	proto.RegisterType((*SystemMemberQueryReq)(nil), "mgmt.SystemMemberQueryReq")
	proto.RegisterType((*SystemMemberQueryResp)(nil), "mgmt.SystemMemberQueryResp")
}

func init() { proto.RegisterFile("system.proto", fileDescriptor_system_15d58b874904f062) }

var fileDescriptor_system_15d58b874904f062 = []byte{
	// 249 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xb1, 0x4a, 0xc4, 0x40,
	0x10, 0x86, 0xd9, 0x24, 0x24, 0x3a, 0xde, 0x9d, 0x32, 0xe8, 0xb1, 0x68, 0x13, 0x52, 0xa5, 0x90,
	0x14, 0x8a, 0x8f, 0x60, 0x23, 0x58, 0xb8, 0xd7, 0xda, 0xe4, 0xdc, 0xe5, 0x08, 0xba, 0xb7, 0x71,
	0x76, 0x53, 0xdc, 0x03, 0xf9, 0x9e, 0xb2, 0x93, 0x04, 0x22, 0xd8, 0xfd, 0xf3, 0xf1, 0x31, 0xfc,
	0x33, 0xb0, 0xf2, 0x27, 0x1f, 0x8c, 0x6d, 0x7a, 0x72, 0xc1, 0x61, 0x66, 0x0f, 0x36, 0x54, 0x2f,
	0xb0, 0xda, 0x31, 0x7d, 0x35, 0x76, 0x6f, 0x08, 0x11, 0xb2, 0x56, 0x6b, 0x92, 0xa2, 0x14, 0xf5,
	0xb9, 0xe2, 0x1c, 0xd9, 0x30, 0x74, 0x5a, 0x26, 0x23, 0x8b, 0x39, 0x32, 0x6a, 0x8f, 0x9f, 0x32,
	0x2d, 0x45, 0xbd, 0x56, 0x9c, 0xab, 0x4b, 0x58, 0x8f, 0xbb, 0x76, 0xc1, 0xf5, 0xca, 0x7c, 0x57,
	0x3f, 0x02, 0x36, 0x4b, 0xe2, 0x7b, 0x7c, 0x82, 0x82, 0x8c, 0x1f, 0xbe, 0x82, 0x97, 0xa2, 0x4c,
	0xeb, 0x8b, 0x87, 0xbb, 0x26, 0xf6, 0x68, 0xfe, 0x6a, 0x8d, 0x62, 0x47, 0xcd, 0xee, 0xed, 0x3b,
	0xe4, 0x23, 0xc2, 0x0d, 0x24, 0x9d, 0x9e, 0xea, 0x25, 0x9d, 0xc6, 0x2d, 0xe4, 0xed, 0x47, 0xe8,
	0xdc, 0x71, 0xaa, 0x37, 0x4d, 0x28, 0xa1, 0x30, 0x44, 0x8e, 0x8c, 0xe6, 0x8e, 0x67, 0x6a, 0x1e,
	0xf1, 0x0a, 0x52, 0xeb, 0x0f, 0x32, 0x63, 0x3d, 0xc6, 0x6a, 0x0b, 0xd7, 0xcb, 0x27, 0xbc, 0x0d,
	0x86, 0x4e, 0xb1, 0xff, 0x33, 0xdc, 0xfc, 0xc3, 0x7d, 0x8f, 0xf7, 0x50, 0x58, 0x46, 0xf3, 0x15,
	0xb8, 0xbc, 0x62, 0xb4, 0xd5, 0xac, 0xec, 0x73, 0x7e, 0xf8, 0xe3, 0x6f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x83, 0x32, 0x5b, 0xa3, 0x80, 0x01, 0x00, 0x00,
}
