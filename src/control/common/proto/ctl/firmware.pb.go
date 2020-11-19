// Code generated by protoc-gen-go. DO NOT EDIT.
// source: firmware.proto

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

type FirmwareUpdateReq_DeviceType int32

const (
	FirmwareUpdateReq_SCM  FirmwareUpdateReq_DeviceType = 0
	FirmwareUpdateReq_NVMe FirmwareUpdateReq_DeviceType = 1
)

var FirmwareUpdateReq_DeviceType_name = map[int32]string{
	0: "SCM",
	1: "NVMe",
}

var FirmwareUpdateReq_DeviceType_value = map[string]int32{
	"SCM":  0,
	"NVMe": 1,
}

func (x FirmwareUpdateReq_DeviceType) String() string {
	return proto.EnumName(FirmwareUpdateReq_DeviceType_name, int32(x))
}

func (FirmwareUpdateReq_DeviceType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_138455e383c002dd, []int{4, 0}
}

type FirmwareQueryReq struct {
	QueryScm             bool     `protobuf:"varint,1,opt,name=queryScm,proto3" json:"queryScm,omitempty"`
	QueryNvme            bool     `protobuf:"varint,2,opt,name=queryNvme,proto3" json:"queryNvme,omitempty"`
	DeviceIDs            []string `protobuf:"bytes,3,rep,name=deviceIDs,proto3" json:"deviceIDs,omitempty"`
	ModelID              string   `protobuf:"bytes,4,opt,name=modelID,proto3" json:"modelID,omitempty"`
	FirmwareRev          string   `protobuf:"bytes,5,opt,name=firmwareRev,proto3" json:"firmwareRev,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FirmwareQueryReq) Reset()         { *m = FirmwareQueryReq{} }
func (m *FirmwareQueryReq) String() string { return proto.CompactTextString(m) }
func (*FirmwareQueryReq) ProtoMessage()    {}
func (*FirmwareQueryReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_138455e383c002dd, []int{0}
}

func (m *FirmwareQueryReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirmwareQueryReq.Unmarshal(m, b)
}
func (m *FirmwareQueryReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirmwareQueryReq.Marshal(b, m, deterministic)
}
func (m *FirmwareQueryReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirmwareQueryReq.Merge(m, src)
}
func (m *FirmwareQueryReq) XXX_Size() int {
	return xxx_messageInfo_FirmwareQueryReq.Size(m)
}
func (m *FirmwareQueryReq) XXX_DiscardUnknown() {
	xxx_messageInfo_FirmwareQueryReq.DiscardUnknown(m)
}

var xxx_messageInfo_FirmwareQueryReq proto.InternalMessageInfo

func (m *FirmwareQueryReq) GetQueryScm() bool {
	if m != nil {
		return m.QueryScm
	}
	return false
}

func (m *FirmwareQueryReq) GetQueryNvme() bool {
	if m != nil {
		return m.QueryNvme
	}
	return false
}

func (m *FirmwareQueryReq) GetDeviceIDs() []string {
	if m != nil {
		return m.DeviceIDs
	}
	return nil
}

func (m *FirmwareQueryReq) GetModelID() string {
	if m != nil {
		return m.ModelID
	}
	return ""
}

func (m *FirmwareQueryReq) GetFirmwareRev() string {
	if m != nil {
		return m.FirmwareRev
	}
	return ""
}

type ScmFirmwareQueryResp struct {
	Module               *ScmModule `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	ActiveVersion        string     `protobuf:"bytes,2,opt,name=activeVersion,proto3" json:"activeVersion,omitempty"`
	StagedVersion        string     `protobuf:"bytes,3,opt,name=stagedVersion,proto3" json:"stagedVersion,omitempty"`
	ImageMaxSizeBytes    uint32     `protobuf:"varint,4,opt,name=imageMaxSizeBytes,proto3" json:"imageMaxSizeBytes,omitempty"`
	UpdateStatus         uint32     `protobuf:"varint,5,opt,name=updateStatus,proto3" json:"updateStatus,omitempty"`
	Error                string     `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ScmFirmwareQueryResp) Reset()         { *m = ScmFirmwareQueryResp{} }
func (m *ScmFirmwareQueryResp) String() string { return proto.CompactTextString(m) }
func (*ScmFirmwareQueryResp) ProtoMessage()    {}
func (*ScmFirmwareQueryResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_138455e383c002dd, []int{1}
}

func (m *ScmFirmwareQueryResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScmFirmwareQueryResp.Unmarshal(m, b)
}
func (m *ScmFirmwareQueryResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScmFirmwareQueryResp.Marshal(b, m, deterministic)
}
func (m *ScmFirmwareQueryResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScmFirmwareQueryResp.Merge(m, src)
}
func (m *ScmFirmwareQueryResp) XXX_Size() int {
	return xxx_messageInfo_ScmFirmwareQueryResp.Size(m)
}
func (m *ScmFirmwareQueryResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ScmFirmwareQueryResp.DiscardUnknown(m)
}

var xxx_messageInfo_ScmFirmwareQueryResp proto.InternalMessageInfo

func (m *ScmFirmwareQueryResp) GetModule() *ScmModule {
	if m != nil {
		return m.Module
	}
	return nil
}

func (m *ScmFirmwareQueryResp) GetActiveVersion() string {
	if m != nil {
		return m.ActiveVersion
	}
	return ""
}

func (m *ScmFirmwareQueryResp) GetStagedVersion() string {
	if m != nil {
		return m.StagedVersion
	}
	return ""
}

func (m *ScmFirmwareQueryResp) GetImageMaxSizeBytes() uint32 {
	if m != nil {
		return m.ImageMaxSizeBytes
	}
	return 0
}

func (m *ScmFirmwareQueryResp) GetUpdateStatus() uint32 {
	if m != nil {
		return m.UpdateStatus
	}
	return 0
}

func (m *ScmFirmwareQueryResp) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type NvmeFirmwareQueryResp struct {
	Device               *NvmeController `protobuf:"bytes,1,opt,name=device,proto3" json:"device,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *NvmeFirmwareQueryResp) Reset()         { *m = NvmeFirmwareQueryResp{} }
func (m *NvmeFirmwareQueryResp) String() string { return proto.CompactTextString(m) }
func (*NvmeFirmwareQueryResp) ProtoMessage()    {}
func (*NvmeFirmwareQueryResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_138455e383c002dd, []int{2}
}

func (m *NvmeFirmwareQueryResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NvmeFirmwareQueryResp.Unmarshal(m, b)
}
func (m *NvmeFirmwareQueryResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NvmeFirmwareQueryResp.Marshal(b, m, deterministic)
}
func (m *NvmeFirmwareQueryResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NvmeFirmwareQueryResp.Merge(m, src)
}
func (m *NvmeFirmwareQueryResp) XXX_Size() int {
	return xxx_messageInfo_NvmeFirmwareQueryResp.Size(m)
}
func (m *NvmeFirmwareQueryResp) XXX_DiscardUnknown() {
	xxx_messageInfo_NvmeFirmwareQueryResp.DiscardUnknown(m)
}

var xxx_messageInfo_NvmeFirmwareQueryResp proto.InternalMessageInfo

func (m *NvmeFirmwareQueryResp) GetDevice() *NvmeController {
	if m != nil {
		return m.Device
	}
	return nil
}

type FirmwareQueryResp struct {
	ScmResults           []*ScmFirmwareQueryResp  `protobuf:"bytes,1,rep,name=scmResults,proto3" json:"scmResults,omitempty"`
	NvmeResults          []*NvmeFirmwareQueryResp `protobuf:"bytes,2,rep,name=nvmeResults,proto3" json:"nvmeResults,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *FirmwareQueryResp) Reset()         { *m = FirmwareQueryResp{} }
func (m *FirmwareQueryResp) String() string { return proto.CompactTextString(m) }
func (*FirmwareQueryResp) ProtoMessage()    {}
func (*FirmwareQueryResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_138455e383c002dd, []int{3}
}

func (m *FirmwareQueryResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirmwareQueryResp.Unmarshal(m, b)
}
func (m *FirmwareQueryResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirmwareQueryResp.Marshal(b, m, deterministic)
}
func (m *FirmwareQueryResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirmwareQueryResp.Merge(m, src)
}
func (m *FirmwareQueryResp) XXX_Size() int {
	return xxx_messageInfo_FirmwareQueryResp.Size(m)
}
func (m *FirmwareQueryResp) XXX_DiscardUnknown() {
	xxx_messageInfo_FirmwareQueryResp.DiscardUnknown(m)
}

var xxx_messageInfo_FirmwareQueryResp proto.InternalMessageInfo

func (m *FirmwareQueryResp) GetScmResults() []*ScmFirmwareQueryResp {
	if m != nil {
		return m.ScmResults
	}
	return nil
}

func (m *FirmwareQueryResp) GetNvmeResults() []*NvmeFirmwareQueryResp {
	if m != nil {
		return m.NvmeResults
	}
	return nil
}

type FirmwareUpdateReq struct {
	FirmwarePath         string                       `protobuf:"bytes,1,opt,name=firmwarePath,proto3" json:"firmwarePath,omitempty"`
	Type                 FirmwareUpdateReq_DeviceType `protobuf:"varint,2,opt,name=type,proto3,enum=ctl.FirmwareUpdateReq_DeviceType" json:"type,omitempty"`
	DeviceIDs            []string                     `protobuf:"bytes,3,rep,name=deviceIDs,proto3" json:"deviceIDs,omitempty"`
	ModelID              string                       `protobuf:"bytes,4,opt,name=modelID,proto3" json:"modelID,omitempty"`
	FirmwareRev          string                       `protobuf:"bytes,5,opt,name=firmwareRev,proto3" json:"firmwareRev,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *FirmwareUpdateReq) Reset()         { *m = FirmwareUpdateReq{} }
func (m *FirmwareUpdateReq) String() string { return proto.CompactTextString(m) }
func (*FirmwareUpdateReq) ProtoMessage()    {}
func (*FirmwareUpdateReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_138455e383c002dd, []int{4}
}

func (m *FirmwareUpdateReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirmwareUpdateReq.Unmarshal(m, b)
}
func (m *FirmwareUpdateReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirmwareUpdateReq.Marshal(b, m, deterministic)
}
func (m *FirmwareUpdateReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirmwareUpdateReq.Merge(m, src)
}
func (m *FirmwareUpdateReq) XXX_Size() int {
	return xxx_messageInfo_FirmwareUpdateReq.Size(m)
}
func (m *FirmwareUpdateReq) XXX_DiscardUnknown() {
	xxx_messageInfo_FirmwareUpdateReq.DiscardUnknown(m)
}

var xxx_messageInfo_FirmwareUpdateReq proto.InternalMessageInfo

func (m *FirmwareUpdateReq) GetFirmwarePath() string {
	if m != nil {
		return m.FirmwarePath
	}
	return ""
}

func (m *FirmwareUpdateReq) GetType() FirmwareUpdateReq_DeviceType {
	if m != nil {
		return m.Type
	}
	return FirmwareUpdateReq_SCM
}

func (m *FirmwareUpdateReq) GetDeviceIDs() []string {
	if m != nil {
		return m.DeviceIDs
	}
	return nil
}

func (m *FirmwareUpdateReq) GetModelID() string {
	if m != nil {
		return m.ModelID
	}
	return ""
}

func (m *FirmwareUpdateReq) GetFirmwareRev() string {
	if m != nil {
		return m.FirmwareRev
	}
	return ""
}

type ScmFirmwareUpdateResp struct {
	Module               *ScmModule `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty"`
	Error                string     `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ScmFirmwareUpdateResp) Reset()         { *m = ScmFirmwareUpdateResp{} }
func (m *ScmFirmwareUpdateResp) String() string { return proto.CompactTextString(m) }
func (*ScmFirmwareUpdateResp) ProtoMessage()    {}
func (*ScmFirmwareUpdateResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_138455e383c002dd, []int{5}
}

func (m *ScmFirmwareUpdateResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScmFirmwareUpdateResp.Unmarshal(m, b)
}
func (m *ScmFirmwareUpdateResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScmFirmwareUpdateResp.Marshal(b, m, deterministic)
}
func (m *ScmFirmwareUpdateResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScmFirmwareUpdateResp.Merge(m, src)
}
func (m *ScmFirmwareUpdateResp) XXX_Size() int {
	return xxx_messageInfo_ScmFirmwareUpdateResp.Size(m)
}
func (m *ScmFirmwareUpdateResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ScmFirmwareUpdateResp.DiscardUnknown(m)
}

var xxx_messageInfo_ScmFirmwareUpdateResp proto.InternalMessageInfo

func (m *ScmFirmwareUpdateResp) GetModule() *ScmModule {
	if m != nil {
		return m.Module
	}
	return nil
}

func (m *ScmFirmwareUpdateResp) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type NvmeFirmwareUpdateResp struct {
	PciAddr              string   `protobuf:"bytes,1,opt,name=pciAddr,proto3" json:"pciAddr,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NvmeFirmwareUpdateResp) Reset()         { *m = NvmeFirmwareUpdateResp{} }
func (m *NvmeFirmwareUpdateResp) String() string { return proto.CompactTextString(m) }
func (*NvmeFirmwareUpdateResp) ProtoMessage()    {}
func (*NvmeFirmwareUpdateResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_138455e383c002dd, []int{6}
}

func (m *NvmeFirmwareUpdateResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NvmeFirmwareUpdateResp.Unmarshal(m, b)
}
func (m *NvmeFirmwareUpdateResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NvmeFirmwareUpdateResp.Marshal(b, m, deterministic)
}
func (m *NvmeFirmwareUpdateResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NvmeFirmwareUpdateResp.Merge(m, src)
}
func (m *NvmeFirmwareUpdateResp) XXX_Size() int {
	return xxx_messageInfo_NvmeFirmwareUpdateResp.Size(m)
}
func (m *NvmeFirmwareUpdateResp) XXX_DiscardUnknown() {
	xxx_messageInfo_NvmeFirmwareUpdateResp.DiscardUnknown(m)
}

var xxx_messageInfo_NvmeFirmwareUpdateResp proto.InternalMessageInfo

func (m *NvmeFirmwareUpdateResp) GetPciAddr() string {
	if m != nil {
		return m.PciAddr
	}
	return ""
}

func (m *NvmeFirmwareUpdateResp) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type FirmwareUpdateResp struct {
	ScmResults           []*ScmFirmwareUpdateResp  `protobuf:"bytes,1,rep,name=scmResults,proto3" json:"scmResults,omitempty"`
	NvmeResults          []*NvmeFirmwareUpdateResp `protobuf:"bytes,2,rep,name=nvmeResults,proto3" json:"nvmeResults,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *FirmwareUpdateResp) Reset()         { *m = FirmwareUpdateResp{} }
func (m *FirmwareUpdateResp) String() string { return proto.CompactTextString(m) }
func (*FirmwareUpdateResp) ProtoMessage()    {}
func (*FirmwareUpdateResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_138455e383c002dd, []int{7}
}

func (m *FirmwareUpdateResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FirmwareUpdateResp.Unmarshal(m, b)
}
func (m *FirmwareUpdateResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FirmwareUpdateResp.Marshal(b, m, deterministic)
}
func (m *FirmwareUpdateResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FirmwareUpdateResp.Merge(m, src)
}
func (m *FirmwareUpdateResp) XXX_Size() int {
	return xxx_messageInfo_FirmwareUpdateResp.Size(m)
}
func (m *FirmwareUpdateResp) XXX_DiscardUnknown() {
	xxx_messageInfo_FirmwareUpdateResp.DiscardUnknown(m)
}

var xxx_messageInfo_FirmwareUpdateResp proto.InternalMessageInfo

func (m *FirmwareUpdateResp) GetScmResults() []*ScmFirmwareUpdateResp {
	if m != nil {
		return m.ScmResults
	}
	return nil
}

func (m *FirmwareUpdateResp) GetNvmeResults() []*NvmeFirmwareUpdateResp {
	if m != nil {
		return m.NvmeResults
	}
	return nil
}

func init() {
	proto.RegisterEnum("ctl.FirmwareUpdateReq_DeviceType", FirmwareUpdateReq_DeviceType_name, FirmwareUpdateReq_DeviceType_value)
	proto.RegisterType((*FirmwareQueryReq)(nil), "ctl.FirmwareQueryReq")
	proto.RegisterType((*ScmFirmwareQueryResp)(nil), "ctl.ScmFirmwareQueryResp")
	proto.RegisterType((*NvmeFirmwareQueryResp)(nil), "ctl.NvmeFirmwareQueryResp")
	proto.RegisterType((*FirmwareQueryResp)(nil), "ctl.FirmwareQueryResp")
	proto.RegisterType((*FirmwareUpdateReq)(nil), "ctl.FirmwareUpdateReq")
	proto.RegisterType((*ScmFirmwareUpdateResp)(nil), "ctl.ScmFirmwareUpdateResp")
	proto.RegisterType((*NvmeFirmwareUpdateResp)(nil), "ctl.NvmeFirmwareUpdateResp")
	proto.RegisterType((*FirmwareUpdateResp)(nil), "ctl.FirmwareUpdateResp")
}

func init() {
	proto.RegisterFile("firmware.proto", fileDescriptor_138455e383c002dd)
}

var fileDescriptor_138455e383c002dd = []byte{
	// 529 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xc5, 0x71, 0x9a, 0x8f, 0x49, 0x13, 0x25, 0x4b, 0x8b, 0x4c, 0x40, 0x22, 0x58, 0x08, 0x45,
	0x02, 0xe5, 0x10, 0xc4, 0x01, 0x04, 0x07, 0x68, 0x84, 0xe8, 0x21, 0x15, 0xac, 0x69, 0xaf, 0xc8,
	0xac, 0x87, 0x60, 0xc9, 0x1b, 0xbb, 0xbb, 0x9b, 0x40, 0xf8, 0x0d, 0x48, 0xfc, 0x0b, 0xfe, 0x1b,
	0x57, 0x7e, 0x01, 0xf2, 0xda, 0x4e, 0xd6, 0x75, 0x38, 0x70, 0xe8, 0xcd, 0xf3, 0xe6, 0xcd, 0xec,
	0xbc, 0x97, 0xa7, 0x40, 0xef, 0x73, 0x28, 0xf8, 0x57, 0x5f, 0xe0, 0x24, 0x11, 0xb1, 0x8a, 0x89,
	0xcd, 0x54, 0x34, 0x1c, 0x48, 0x15, 0x0b, 0x7f, 0x81, 0x1f, 0x25, 0xe3, 0x19, 0x3e, 0x24, 0x05,
	0xb4, 0x5c, 0xf3, 0x9c, 0xeb, 0xfe, 0xb2, 0xa0, 0xff, 0x26, 0x1f, 0x7f, 0xbf, 0x42, 0xb1, 0xa1,
	0x78, 0x49, 0x86, 0xd0, 0xba, 0x4c, 0xbf, 0x3d, 0xc6, 0x1d, 0x6b, 0x64, 0x8d, 0x5b, 0x74, 0x5b,
	0x93, 0xbb, 0xd0, 0xd6, 0xdf, 0x67, 0x6b, 0x8e, 0x4e, 0x4d, 0x37, 0x77, 0x40, 0xda, 0x0d, 0x70,
	0x1d, 0x32, 0x3c, 0x9d, 0x49, 0xc7, 0x1e, 0xd9, 0xe3, 0x36, 0xdd, 0x01, 0xc4, 0x81, 0x26, 0x8f,
	0x03, 0x8c, 0x4e, 0x67, 0x4e, 0x7d, 0x64, 0x8d, 0xdb, 0xb4, 0x28, 0xc9, 0x08, 0x3a, 0x85, 0x08,
	0x8a, 0x6b, 0xe7, 0x40, 0x77, 0x4d, 0xc8, 0xfd, 0x63, 0xc1, 0x91, 0xc7, 0xf8, 0x95, 0x5b, 0x65,
	0x42, 0x1e, 0x42, 0x83, 0xc7, 0xc1, 0x2a, 0x42, 0x7d, 0x6a, 0x67, 0xda, 0x9b, 0x30, 0x15, 0x4d,
	0x3c, 0xc6, 0xe7, 0x1a, 0xa5, 0x79, 0x97, 0x3c, 0x80, 0xae, 0xcf, 0x54, 0xb8, 0xc6, 0x0b, 0x14,
	0x32, 0x8c, 0x97, 0xfa, 0xf8, 0x36, 0x2d, 0x83, 0x29, 0x4b, 0x2a, 0x7f, 0x81, 0x41, 0xc1, 0xb2,
	0x33, 0x56, 0x09, 0x24, 0x8f, 0x61, 0x10, 0x72, 0x7f, 0x81, 0x73, 0xff, 0x9b, 0x17, 0x7e, 0xc7,
	0xd7, 0x1b, 0x85, 0x52, 0x4b, 0xea, 0xd2, 0x6a, 0x83, 0xb8, 0x70, 0xb8, 0x4a, 0x02, 0x5f, 0xa1,
	0xa7, 0x7c, 0xb5, 0x92, 0x5a, 0x5d, 0x97, 0x96, 0x30, 0x72, 0x04, 0x07, 0x28, 0x44, 0x2c, 0x9c,
	0x86, 0x7e, 0x2f, 0x2b, 0xdc, 0x19, 0x1c, 0xa7, 0xb6, 0x56, 0x45, 0x3f, 0x82, 0x46, 0x66, 0x6b,
	0x2e, 0xfa, 0xa6, 0x16, 0x9d, 0x72, 0x4f, 0xe2, 0xa5, 0x12, 0x71, 0x14, 0xa1, 0xa0, 0x39, 0xc5,
	0xfd, 0x61, 0xc1, 0xa0, 0xba, 0xe2, 0x19, 0x80, 0x64, 0x9c, 0xa2, 0x5c, 0x45, 0x4a, 0x3a, 0xd6,
	0xc8, 0x1e, 0x77, 0xa6, 0xb7, 0x0b, 0xef, 0x2a, 0x74, 0x6a, 0x90, 0xc9, 0x0b, 0xe8, 0xa4, 0x11,
	0x2a, 0x66, 0x6b, 0x7a, 0x76, 0xb8, 0x3d, 0xa1, 0x3a, 0x6c, 0xd2, 0xdd, 0xdf, 0xc6, 0x39, 0xe7,
	0xda, 0x83, 0x34, 0x73, 0x2e, 0x1c, 0x16, 0x3f, 0xf7, 0x3b, 0x5f, 0x7d, 0xd1, 0xba, 0xda, 0xb4,
	0x84, 0x91, 0xa7, 0x50, 0x57, 0x9b, 0x24, 0x8b, 0x5d, 0x6f, 0x7a, 0x5f, 0x3f, 0x58, 0xd9, 0x34,
	0x99, 0x69, 0xd5, 0x1f, 0x36, 0x09, 0x52, 0x4d, 0xbf, 0xc6, 0x50, 0xde, 0x03, 0xd8, 0xbd, 0x46,
	0x9a, 0x60, 0x7b, 0x27, 0xf3, 0xfe, 0x0d, 0xd2, 0x82, 0xfa, 0xd9, 0xc5, 0x1c, 0xfb, 0x96, 0x7b,
	0x0e, 0xc7, 0x86, 0x9b, 0xc5, 0x8d, 0xff, 0x91, 0xda, 0x6d, 0x2e, 0x6a, 0x66, 0x2e, 0xde, 0xc2,
	0x2d, 0xd3, 0x68, 0x63, 0xaf, 0x03, 0xcd, 0x84, 0x85, 0xaf, 0x82, 0x40, 0xe4, 0x0e, 0x16, 0xe5,
	0x3f, 0x36, 0xfd, 0xb4, 0x80, 0xec, 0x59, 0xf3, 0x7c, 0x4f, 0x38, 0x86, 0x57, 0xc3, 0xb1, 0xe3,
	0x97, 0xd2, 0xf1, 0x72, 0x5f, 0x3a, 0xee, 0x54, 0xd2, 0x61, 0x4c, 0x9b, 0xfc, 0x4f, 0x0d, 0xfd,
	0xc7, 0xf4, 0xe4, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcb, 0x8f, 0xf8, 0x73, 0xd6, 0x04, 0x00,
	0x00,
}
