// Code generated by protoc-gen-go. DO NOT EDIT.
// source: storage_scm.proto

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

// ScmModule represent Storage Class Memory modules installed.
type ScmModule struct {
	// string uid = 1; // The uid of the module.
	Physicalid uint32 `protobuf:"varint,1,opt,name=physicalid,proto3" json:"physicalid,omitempty"`
	// string handle = 3; // The device handle of the module.
	// string serial = 8; // The serial number of the module.
	Capacity uint64 `protobuf:"varint,2,opt,name=capacity,proto3" json:"capacity,omitempty"`
	// string fwrev = 10; // The firmware revision of the module.
	Loc                  *ScmModule_Location `protobuf:"bytes,3,opt,name=loc,proto3" json:"loc,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ScmModule) Reset()         { *m = ScmModule{} }
func (m *ScmModule) String() string { return proto.CompactTextString(m) }
func (*ScmModule) ProtoMessage()    {}
func (*ScmModule) Descriptor() ([]byte, []int) {
	return fileDescriptor_storage_scm_0bc9936a1221be65, []int{0}
}
func (m *ScmModule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScmModule.Unmarshal(m, b)
}
func (m *ScmModule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScmModule.Marshal(b, m, deterministic)
}
func (dst *ScmModule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScmModule.Merge(dst, src)
}
func (m *ScmModule) XXX_Size() int {
	return xxx_messageInfo_ScmModule.Size(m)
}
func (m *ScmModule) XXX_DiscardUnknown() {
	xxx_messageInfo_ScmModule.DiscardUnknown(m)
}

var xxx_messageInfo_ScmModule proto.InternalMessageInfo

func (m *ScmModule) GetPhysicalid() uint32 {
	if m != nil {
		return m.Physicalid
	}
	return 0
}

func (m *ScmModule) GetCapacity() uint64 {
	if m != nil {
		return m.Capacity
	}
	return 0
}

func (m *ScmModule) GetLoc() *ScmModule_Location {
	if m != nil {
		return m.Loc
	}
	return nil
}

type ScmModule_Location struct {
	Channel              uint32   `protobuf:"varint,1,opt,name=channel,proto3" json:"channel,omitempty"`
	Channelpos           uint32   `protobuf:"varint,2,opt,name=channelpos,proto3" json:"channelpos,omitempty"`
	Memctrlr             uint32   `protobuf:"varint,3,opt,name=memctrlr,proto3" json:"memctrlr,omitempty"`
	Socket               uint32   `protobuf:"varint,4,opt,name=socket,proto3" json:"socket,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ScmModule_Location) Reset()         { *m = ScmModule_Location{} }
func (m *ScmModule_Location) String() string { return proto.CompactTextString(m) }
func (*ScmModule_Location) ProtoMessage()    {}
func (*ScmModule_Location) Descriptor() ([]byte, []int) {
	return fileDescriptor_storage_scm_0bc9936a1221be65, []int{0, 0}
}
func (m *ScmModule_Location) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScmModule_Location.Unmarshal(m, b)
}
func (m *ScmModule_Location) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScmModule_Location.Marshal(b, m, deterministic)
}
func (dst *ScmModule_Location) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScmModule_Location.Merge(dst, src)
}
func (m *ScmModule_Location) XXX_Size() int {
	return xxx_messageInfo_ScmModule_Location.Size(m)
}
func (m *ScmModule_Location) XXX_DiscardUnknown() {
	xxx_messageInfo_ScmModule_Location.DiscardUnknown(m)
}

var xxx_messageInfo_ScmModule_Location proto.InternalMessageInfo

func (m *ScmModule_Location) GetChannel() uint32 {
	if m != nil {
		return m.Channel
	}
	return 0
}

func (m *ScmModule_Location) GetChannelpos() uint32 {
	if m != nil {
		return m.Channelpos
	}
	return 0
}

func (m *ScmModule_Location) GetMemctrlr() uint32 {
	if m != nil {
		return m.Memctrlr
	}
	return 0
}

func (m *ScmModule_Location) GetSocket() uint32 {
	if m != nil {
		return m.Socket
	}
	return 0
}

// ScmMount represents mounted AppDirect region made up of SCM module set.
type ScmMount struct {
	Mntpoint             string       `protobuf:"bytes,1,opt,name=mntpoint,proto3" json:"mntpoint,omitempty"`
	Modules              []*ScmModule `protobuf:"bytes,2,rep,name=modules,proto3" json:"modules,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ScmMount) Reset()         { *m = ScmMount{} }
func (m *ScmMount) String() string { return proto.CompactTextString(m) }
func (*ScmMount) ProtoMessage()    {}
func (*ScmMount) Descriptor() ([]byte, []int) {
	return fileDescriptor_storage_scm_0bc9936a1221be65, []int{1}
}
func (m *ScmMount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScmMount.Unmarshal(m, b)
}
func (m *ScmMount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScmMount.Marshal(b, m, deterministic)
}
func (dst *ScmMount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScmMount.Merge(dst, src)
}
func (m *ScmMount) XXX_Size() int {
	return xxx_messageInfo_ScmMount.Size(m)
}
func (m *ScmMount) XXX_DiscardUnknown() {
	xxx_messageInfo_ScmMount.DiscardUnknown(m)
}

var xxx_messageInfo_ScmMount proto.InternalMessageInfo

func (m *ScmMount) GetMntpoint() string {
	if m != nil {
		return m.Mntpoint
	}
	return ""
}

func (m *ScmMount) GetModules() []*ScmModule {
	if m != nil {
		return m.Modules
	}
	return nil
}

// ScmModuleResult represents operation state for specific SCM/PM module.
//
// TODO: replace identifier with serial when returned in scan
type ScmModuleResult struct {
	Loc                  *ScmModule_Location `protobuf:"bytes,1,opt,name=loc,proto3" json:"loc,omitempty"`
	State                *ResponseState      `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ScmModuleResult) Reset()         { *m = ScmModuleResult{} }
func (m *ScmModuleResult) String() string { return proto.CompactTextString(m) }
func (*ScmModuleResult) ProtoMessage()    {}
func (*ScmModuleResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_storage_scm_0bc9936a1221be65, []int{2}
}
func (m *ScmModuleResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScmModuleResult.Unmarshal(m, b)
}
func (m *ScmModuleResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScmModuleResult.Marshal(b, m, deterministic)
}
func (dst *ScmModuleResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScmModuleResult.Merge(dst, src)
}
func (m *ScmModuleResult) XXX_Size() int {
	return xxx_messageInfo_ScmModuleResult.Size(m)
}
func (m *ScmModuleResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ScmModuleResult.DiscardUnknown(m)
}

var xxx_messageInfo_ScmModuleResult proto.InternalMessageInfo

func (m *ScmModuleResult) GetLoc() *ScmModule_Location {
	if m != nil {
		return m.Loc
	}
	return nil
}

func (m *ScmModuleResult) GetState() *ResponseState {
	if m != nil {
		return m.State
	}
	return nil
}

// ScmMountResult represents operation state for specific SCM mount point.
type ScmMountResult struct {
	Mntpoint             string         `protobuf:"bytes,1,opt,name=mntpoint,proto3" json:"mntpoint,omitempty"`
	State                *ResponseState `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ScmMountResult) Reset()         { *m = ScmMountResult{} }
func (m *ScmMountResult) String() string { return proto.CompactTextString(m) }
func (*ScmMountResult) ProtoMessage()    {}
func (*ScmMountResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_storage_scm_0bc9936a1221be65, []int{3}
}
func (m *ScmMountResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScmMountResult.Unmarshal(m, b)
}
func (m *ScmMountResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScmMountResult.Marshal(b, m, deterministic)
}
func (dst *ScmMountResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScmMountResult.Merge(dst, src)
}
func (m *ScmMountResult) XXX_Size() int {
	return xxx_messageInfo_ScmMountResult.Size(m)
}
func (m *ScmMountResult) XXX_DiscardUnknown() {
	xxx_messageInfo_ScmMountResult.DiscardUnknown(m)
}

var xxx_messageInfo_ScmMountResult proto.InternalMessageInfo

func (m *ScmMountResult) GetMntpoint() string {
	if m != nil {
		return m.Mntpoint
	}
	return ""
}

func (m *ScmMountResult) GetState() *ResponseState {
	if m != nil {
		return m.State
	}
	return nil
}

type ScanScmReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ScanScmReq) Reset()         { *m = ScanScmReq{} }
func (m *ScanScmReq) String() string { return proto.CompactTextString(m) }
func (*ScanScmReq) ProtoMessage()    {}
func (*ScanScmReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_storage_scm_0bc9936a1221be65, []int{4}
}
func (m *ScanScmReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScanScmReq.Unmarshal(m, b)
}
func (m *ScanScmReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScanScmReq.Marshal(b, m, deterministic)
}
func (dst *ScanScmReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScanScmReq.Merge(dst, src)
}
func (m *ScanScmReq) XXX_Size() int {
	return xxx_messageInfo_ScanScmReq.Size(m)
}
func (m *ScanScmReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ScanScmReq.DiscardUnknown(m)
}

var xxx_messageInfo_ScanScmReq proto.InternalMessageInfo

type FormatScmReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FormatScmReq) Reset()         { *m = FormatScmReq{} }
func (m *FormatScmReq) String() string { return proto.CompactTextString(m) }
func (*FormatScmReq) ProtoMessage()    {}
func (*FormatScmReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_storage_scm_0bc9936a1221be65, []int{5}
}
func (m *FormatScmReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FormatScmReq.Unmarshal(m, b)
}
func (m *FormatScmReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FormatScmReq.Marshal(b, m, deterministic)
}
func (dst *FormatScmReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FormatScmReq.Merge(dst, src)
}
func (m *FormatScmReq) XXX_Size() int {
	return xxx_messageInfo_FormatScmReq.Size(m)
}
func (m *FormatScmReq) XXX_DiscardUnknown() {
	xxx_messageInfo_FormatScmReq.DiscardUnknown(m)
}

var xxx_messageInfo_FormatScmReq proto.InternalMessageInfo

type UpdateScmReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateScmReq) Reset()         { *m = UpdateScmReq{} }
func (m *UpdateScmReq) String() string { return proto.CompactTextString(m) }
func (*UpdateScmReq) ProtoMessage()    {}
func (*UpdateScmReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_storage_scm_0bc9936a1221be65, []int{6}
}
func (m *UpdateScmReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateScmReq.Unmarshal(m, b)
}
func (m *UpdateScmReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateScmReq.Marshal(b, m, deterministic)
}
func (dst *UpdateScmReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateScmReq.Merge(dst, src)
}
func (m *UpdateScmReq) XXX_Size() int {
	return xxx_messageInfo_UpdateScmReq.Size(m)
}
func (m *UpdateScmReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateScmReq.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateScmReq proto.InternalMessageInfo

type BurninScmReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BurninScmReq) Reset()         { *m = BurninScmReq{} }
func (m *BurninScmReq) String() string { return proto.CompactTextString(m) }
func (*BurninScmReq) ProtoMessage()    {}
func (*BurninScmReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_storage_scm_0bc9936a1221be65, []int{7}
}
func (m *BurninScmReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BurninScmReq.Unmarshal(m, b)
}
func (m *BurninScmReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BurninScmReq.Marshal(b, m, deterministic)
}
func (dst *BurninScmReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BurninScmReq.Merge(dst, src)
}
func (m *BurninScmReq) XXX_Size() int {
	return xxx_messageInfo_BurninScmReq.Size(m)
}
func (m *BurninScmReq) XXX_DiscardUnknown() {
	xxx_messageInfo_BurninScmReq.DiscardUnknown(m)
}

var xxx_messageInfo_BurninScmReq proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ScmModule)(nil), "mgmt.ScmModule")
	proto.RegisterType((*ScmModule_Location)(nil), "mgmt.ScmModule.Location")
	proto.RegisterType((*ScmMount)(nil), "mgmt.ScmMount")
	proto.RegisterType((*ScmModuleResult)(nil), "mgmt.ScmModuleResult")
	proto.RegisterType((*ScmMountResult)(nil), "mgmt.ScmMountResult")
	proto.RegisterType((*ScanScmReq)(nil), "mgmt.ScanScmReq")
	proto.RegisterType((*FormatScmReq)(nil), "mgmt.FormatScmReq")
	proto.RegisterType((*UpdateScmReq)(nil), "mgmt.UpdateScmReq")
	proto.RegisterType((*BurninScmReq)(nil), "mgmt.BurninScmReq")
}

func init() { proto.RegisterFile("storage_scm.proto", fileDescriptor_storage_scm_0bc9936a1221be65) }

var fileDescriptor_storage_scm_0bc9936a1221be65 = []byte{
	// 334 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xbd, 0x4e, 0xc3, 0x30,
	0x14, 0x85, 0x15, 0x5a, 0xfa, 0x73, 0xfb, 0x27, 0x8c, 0x84, 0xa2, 0x0e, 0xa8, 0xca, 0x94, 0x32,
	0x64, 0x28, 0x6f, 0xc0, 0xc0, 0x04, 0x03, 0x8e, 0x10, 0x23, 0x32, 0xae, 0xd5, 0x46, 0xc4, 0xbe,
	0x26, 0xbe, 0x91, 0xe8, 0x43, 0xf3, 0x0e, 0x28, 0x6e, 0x7e, 0x2a, 0x06, 0xd4, 0x2d, 0xe7, 0xdc,
	0x13, 0x7f, 0xe7, 0x5a, 0x86, 0x2b, 0x47, 0x58, 0x88, 0x9d, 0x7a, 0x77, 0x52, 0x27, 0xb6, 0x40,
	0x42, 0xd6, 0xd7, 0x3b, 0x4d, 0xcb, 0xa9, 0x44, 0xad, 0xd1, 0x1c, 0xbd, 0xe8, 0x27, 0x80, 0x71,
	0x2a, 0xf5, 0x33, 0x6e, 0xcb, 0x5c, 0xb1, 0x5b, 0x00, 0xbb, 0x3f, 0xb8, 0x4c, 0x8a, 0x3c, 0xdb,
	0x86, 0xc1, 0x2a, 0x88, 0x67, 0xfc, 0xc4, 0x61, 0x4b, 0x18, 0x49, 0x61, 0x85, 0xcc, 0xe8, 0x10,
	0x5e, 0xac, 0x82, 0xb8, 0xcf, 0x5b, 0xcd, 0xee, 0xa0, 0x97, 0xa3, 0x0c, 0x7b, 0xab, 0x20, 0x9e,
	0x6c, 0xc2, 0xa4, 0x62, 0x25, 0xed, 0xc9, 0xc9, 0x13, 0x4a, 0x41, 0x19, 0x1a, 0x5e, 0x85, 0x96,
	0xdf, 0x30, 0x6a, 0x0c, 0x16, 0xc2, 0x50, 0xee, 0x85, 0x31, 0x2a, 0xaf, 0x81, 0x8d, 0xac, 0xda,
	0xd4, 0x9f, 0x16, 0x9d, 0xe7, 0xcd, 0xf8, 0x89, 0x53, 0xb5, 0xd1, 0x4a, 0x4b, 0x2a, 0xf2, 0xc2,
	0x63, 0x67, 0xbc, 0xd5, 0xec, 0x06, 0x06, 0x0e, 0xe5, 0xa7, 0xa2, 0xb0, 0xef, 0x27, 0xb5, 0x8a,
	0x5e, 0x60, 0xe4, 0x4b, 0x95, 0x86, 0xfc, 0xff, 0x86, 0x2c, 0x66, 0x86, 0x3c, 0x7a, 0xcc, 0x5b,
	0xcd, 0xd6, 0x30, 0xd4, 0xbe, 0x79, 0x05, 0xee, 0xc5, 0x93, 0xcd, 0xe2, 0xcf, 0x46, 0xbc, 0x99,
	0x47, 0x7b, 0x58, 0x74, 0xae, 0x72, 0x65, 0x4e, 0xcd, 0x5d, 0x04, 0x67, 0xdc, 0x05, 0x5b, 0xc3,
	0xa5, 0x23, 0x41, 0xca, 0x2f, 0x38, 0xd9, 0x5c, 0x1f, 0xd3, 0x5c, 0x39, 0x8b, 0xc6, 0xa9, 0xb4,
	0x1a, 0xf1, 0x63, 0x22, 0x7a, 0x83, 0x79, 0x53, 0xbe, 0x06, 0xfd, 0xbf, 0xc2, 0xd9, 0x07, 0x4f,
	0x01, 0x52, 0x29, 0x4c, 0x2a, 0x35, 0x57, 0x5f, 0xd1, 0x1c, 0xa6, 0x8f, 0x58, 0x68, 0x41, 0x9d,
	0x7e, 0xb5, 0x5b, 0x41, 0xaa, 0xd3, 0x0f, 0x65, 0x61, 0xb2, 0x3a, 0xff, 0x31, 0xf0, 0x4f, 0xe9,
	0xfe, 0x37, 0x00, 0x00, 0xff, 0xff, 0x10, 0x03, 0x88, 0x7d, 0x73, 0x02, 0x00, 0x00,
}
