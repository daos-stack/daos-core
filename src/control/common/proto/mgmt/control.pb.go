// Code generated by protoc-gen-go. DO NOT EDIT.
// source: control.proto

package mgmt

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MgmtCtlClient is the client API for MgmtCtl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MgmtCtlClient interface {
	// Prepare nonvolatile storage devices for use with DAOS
	StoragePrepare(ctx context.Context, in *StoragePrepareReq, opts ...grpc.CallOption) (*StoragePrepareResp, error)
	// Retrieve details of nonvolatile storage on server
	StorageScan(ctx context.Context, in *StorageScanReq, opts ...grpc.CallOption) (*StorageScanResp, error)
	// Format nonvolatile storage devices for use with DAOS
	StorageFormat(ctx context.Context, in *StorageFormatReq, opts ...grpc.CallOption) (MgmtCtl_StorageFormatClient, error)
	// Update nonvolatile storage device firmware
	StorageUpdate(ctx context.Context, in *StorageUpdateReq, opts ...grpc.CallOption) (MgmtCtl_StorageUpdateClient, error)
	// Perform burn-in testing to verify nonvolatile storage devices
	StorageBurnIn(ctx context.Context, in *StorageBurnInReq, opts ...grpc.CallOption) (MgmtCtl_StorageBurnInClient, error)
	// Fetch FIO configuration file specifying burn-in jobs/workloads
	FetchFioConfigPaths(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (MgmtCtl_FetchFioConfigPathsClient, error)
	// List features supported on remote storage server/DAOS system
	ListFeatures(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (MgmtCtl_ListFeaturesClient, error)
}

type mgmtCtlClient struct {
	cc *grpc.ClientConn
}

func NewMgmtCtlClient(cc *grpc.ClientConn) MgmtCtlClient {
	return &mgmtCtlClient{cc}
}

func (c *mgmtCtlClient) StoragePrepare(ctx context.Context, in *StoragePrepareReq, opts ...grpc.CallOption) (*StoragePrepareResp, error) {
	out := new(StoragePrepareResp)
	err := c.cc.Invoke(ctx, "/mgmt.MgmtCtl/StoragePrepare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mgmtCtlClient) StorageScan(ctx context.Context, in *StorageScanReq, opts ...grpc.CallOption) (*StorageScanResp, error) {
	out := new(StorageScanResp)
	err := c.cc.Invoke(ctx, "/mgmt.MgmtCtl/StorageScan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mgmtCtlClient) StorageFormat(ctx context.Context, in *StorageFormatReq, opts ...grpc.CallOption) (MgmtCtl_StorageFormatClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MgmtCtl_serviceDesc.Streams[0], "/mgmt.MgmtCtl/StorageFormat", opts...)
	if err != nil {
		return nil, err
	}
	x := &mgmtCtlStorageFormatClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MgmtCtl_StorageFormatClient interface {
	Recv() (*StorageFormatResp, error)
	grpc.ClientStream
}

type mgmtCtlStorageFormatClient struct {
	grpc.ClientStream
}

func (x *mgmtCtlStorageFormatClient) Recv() (*StorageFormatResp, error) {
	m := new(StorageFormatResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mgmtCtlClient) StorageUpdate(ctx context.Context, in *StorageUpdateReq, opts ...grpc.CallOption) (MgmtCtl_StorageUpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MgmtCtl_serviceDesc.Streams[1], "/mgmt.MgmtCtl/StorageUpdate", opts...)
	if err != nil {
		return nil, err
	}
	x := &mgmtCtlStorageUpdateClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MgmtCtl_StorageUpdateClient interface {
	Recv() (*StorageUpdateResp, error)
	grpc.ClientStream
}

type mgmtCtlStorageUpdateClient struct {
	grpc.ClientStream
}

func (x *mgmtCtlStorageUpdateClient) Recv() (*StorageUpdateResp, error) {
	m := new(StorageUpdateResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mgmtCtlClient) StorageBurnIn(ctx context.Context, in *StorageBurnInReq, opts ...grpc.CallOption) (MgmtCtl_StorageBurnInClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MgmtCtl_serviceDesc.Streams[2], "/mgmt.MgmtCtl/StorageBurnIn", opts...)
	if err != nil {
		return nil, err
	}
	x := &mgmtCtlStorageBurnInClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MgmtCtl_StorageBurnInClient interface {
	Recv() (*StorageBurnInResp, error)
	grpc.ClientStream
}

type mgmtCtlStorageBurnInClient struct {
	grpc.ClientStream
}

func (x *mgmtCtlStorageBurnInClient) Recv() (*StorageBurnInResp, error) {
	m := new(StorageBurnInResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mgmtCtlClient) FetchFioConfigPaths(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (MgmtCtl_FetchFioConfigPathsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MgmtCtl_serviceDesc.Streams[3], "/mgmt.MgmtCtl/FetchFioConfigPaths", opts...)
	if err != nil {
		return nil, err
	}
	x := &mgmtCtlFetchFioConfigPathsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MgmtCtl_FetchFioConfigPathsClient interface {
	Recv() (*FilePath, error)
	grpc.ClientStream
}

type mgmtCtlFetchFioConfigPathsClient struct {
	grpc.ClientStream
}

func (x *mgmtCtlFetchFioConfigPathsClient) Recv() (*FilePath, error) {
	m := new(FilePath)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mgmtCtlClient) ListFeatures(ctx context.Context, in *EmptyReq, opts ...grpc.CallOption) (MgmtCtl_ListFeaturesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MgmtCtl_serviceDesc.Streams[4], "/mgmt.MgmtCtl/ListFeatures", opts...)
	if err != nil {
		return nil, err
	}
	x := &mgmtCtlListFeaturesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MgmtCtl_ListFeaturesClient interface {
	Recv() (*Feature, error)
	grpc.ClientStream
}

type mgmtCtlListFeaturesClient struct {
	grpc.ClientStream
}

func (x *mgmtCtlListFeaturesClient) Recv() (*Feature, error) {
	m := new(Feature)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MgmtCtlServer is the server API for MgmtCtl service.
type MgmtCtlServer interface {
	// Prepare nonvolatile storage devices for use with DAOS
	StoragePrepare(context.Context, *StoragePrepareReq) (*StoragePrepareResp, error)
	// Retrieve details of nonvolatile storage on server
	StorageScan(context.Context, *StorageScanReq) (*StorageScanResp, error)
	// Format nonvolatile storage devices for use with DAOS
	StorageFormat(*StorageFormatReq, MgmtCtl_StorageFormatServer) error
	// Update nonvolatile storage device firmware
	StorageUpdate(*StorageUpdateReq, MgmtCtl_StorageUpdateServer) error
	// Perform burn-in testing to verify nonvolatile storage devices
	StorageBurnIn(*StorageBurnInReq, MgmtCtl_StorageBurnInServer) error
	// Fetch FIO configuration file specifying burn-in jobs/workloads
	FetchFioConfigPaths(*EmptyReq, MgmtCtl_FetchFioConfigPathsServer) error
	// List features supported on remote storage server/DAOS system
	ListFeatures(*EmptyReq, MgmtCtl_ListFeaturesServer) error
}

func RegisterMgmtCtlServer(s *grpc.Server, srv MgmtCtlServer) {
	s.RegisterService(&_MgmtCtl_serviceDesc, srv)
}

func _MgmtCtl_StoragePrepare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoragePrepareReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MgmtCtlServer).StoragePrepare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mgmt.MgmtCtl/StoragePrepare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MgmtCtlServer).StoragePrepare(ctx, req.(*StoragePrepareReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MgmtCtl_StorageScan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StorageScanReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MgmtCtlServer).StorageScan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mgmt.MgmtCtl/StorageScan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MgmtCtlServer).StorageScan(ctx, req.(*StorageScanReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MgmtCtl_StorageFormat_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StorageFormatReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MgmtCtlServer).StorageFormat(m, &mgmtCtlStorageFormatServer{stream})
}

type MgmtCtl_StorageFormatServer interface {
	Send(*StorageFormatResp) error
	grpc.ServerStream
}

type mgmtCtlStorageFormatServer struct {
	grpc.ServerStream
}

func (x *mgmtCtlStorageFormatServer) Send(m *StorageFormatResp) error {
	return x.ServerStream.SendMsg(m)
}

func _MgmtCtl_StorageUpdate_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StorageUpdateReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MgmtCtlServer).StorageUpdate(m, &mgmtCtlStorageUpdateServer{stream})
}

type MgmtCtl_StorageUpdateServer interface {
	Send(*StorageUpdateResp) error
	grpc.ServerStream
}

type mgmtCtlStorageUpdateServer struct {
	grpc.ServerStream
}

func (x *mgmtCtlStorageUpdateServer) Send(m *StorageUpdateResp) error {
	return x.ServerStream.SendMsg(m)
}

func _MgmtCtl_StorageBurnIn_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StorageBurnInReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MgmtCtlServer).StorageBurnIn(m, &mgmtCtlStorageBurnInServer{stream})
}

type MgmtCtl_StorageBurnInServer interface {
	Send(*StorageBurnInResp) error
	grpc.ServerStream
}

type mgmtCtlStorageBurnInServer struct {
	grpc.ServerStream
}

func (x *mgmtCtlStorageBurnInServer) Send(m *StorageBurnInResp) error {
	return x.ServerStream.SendMsg(m)
}

func _MgmtCtl_FetchFioConfigPaths_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EmptyReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MgmtCtlServer).FetchFioConfigPaths(m, &mgmtCtlFetchFioConfigPathsServer{stream})
}

type MgmtCtl_FetchFioConfigPathsServer interface {
	Send(*FilePath) error
	grpc.ServerStream
}

type mgmtCtlFetchFioConfigPathsServer struct {
	grpc.ServerStream
}

func (x *mgmtCtlFetchFioConfigPathsServer) Send(m *FilePath) error {
	return x.ServerStream.SendMsg(m)
}

func _MgmtCtl_ListFeatures_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EmptyReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MgmtCtlServer).ListFeatures(m, &mgmtCtlListFeaturesServer{stream})
}

type MgmtCtl_ListFeaturesServer interface {
	Send(*Feature) error
	grpc.ServerStream
}

type mgmtCtlListFeaturesServer struct {
	grpc.ServerStream
}

func (x *mgmtCtlListFeaturesServer) Send(m *Feature) error {
	return x.ServerStream.SendMsg(m)
}

var _MgmtCtl_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mgmt.MgmtCtl",
	HandlerType: (*MgmtCtlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoragePrepare",
			Handler:    _MgmtCtl_StoragePrepare_Handler,
		},
		{
			MethodName: "StorageScan",
			Handler:    _MgmtCtl_StorageScan_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StorageFormat",
			Handler:       _MgmtCtl_StorageFormat_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StorageUpdate",
			Handler:       _MgmtCtl_StorageUpdate_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StorageBurnIn",
			Handler:       _MgmtCtl_StorageBurnIn_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "FetchFioConfigPaths",
			Handler:       _MgmtCtl_FetchFioConfigPaths_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListFeatures",
			Handler:       _MgmtCtl_ListFeatures_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "control.proto",
}

func init() { proto.RegisterFile("control.proto", fileDescriptor_control_f325a27f61b02b5d) }

var fileDescriptor_control_f325a27f61b02b5d = []byte{
	// 261 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xcd, 0x4a, 0xc3, 0x40,
	0x14, 0x85, 0x15, 0x45, 0x61, 0x6c, 0xb2, 0x18, 0x7f, 0x2a, 0x59, 0xfa, 0x00, 0xc1, 0x9f, 0x95,
	0xe0, 0xca, 0xda, 0x01, 0x41, 0xa1, 0x58, 0x7c, 0x80, 0x31, 0xde, 0xa6, 0x81, 0xcc, 0xdc, 0x71,
	0xe6, 0x76, 0xe1, 0x4b, 0xfa, 0x4c, 0x32, 0xb9, 0x53, 0x68, 0x9a, 0x2c, 0xcf, 0x77, 0xcf, 0xf9,
	0x36, 0x57, 0x64, 0x15, 0x5a, 0xf2, 0xd8, 0x96, 0xce, 0x23, 0xa1, 0x3c, 0x36, 0xb5, 0xa1, 0x62,
	0x52, 0xa1, 0x31, 0x68, 0x99, 0x15, 0x59, 0x20, 0xf4, 0xba, 0x86, 0x14, 0xf3, 0x15, 0x68, 0xda,
	0x78, 0x08, 0x9c, 0xef, 0xff, 0x8e, 0xc4, 0xe9, 0x7b, 0x6d, 0x68, 0x46, 0xad, 0x9c, 0x8b, 0x7c,
	0xc9, 0xe5, 0x85, 0x07, 0xa7, 0x3d, 0xc8, 0x69, 0x19, 0x8d, 0x65, 0x9f, 0x7e, 0xc0, 0x4f, 0x71,
	0x3d, 0x7e, 0x08, 0xee, 0xe6, 0x40, 0x3e, 0x89, 0xb3, 0xc4, 0x97, 0x95, 0xb6, 0xf2, 0xa2, 0x57,
	0x8d, 0x28, 0x0a, 0x2e, 0x47, 0x68, 0xb7, 0x7e, 0x11, 0x59, 0x82, 0x0a, 0xbd, 0xd1, 0x24, 0xaf,
	0x7a, 0x4d, 0x86, 0xd1, 0x30, 0x1d, 0xe5, 0xd1, 0x71, 0x7b, 0xb8, 0x63, 0xf9, 0x74, 0xdf, 0x9a,
	0x60, 0xcf, 0xc2, 0x70, 0x68, 0xd9, 0xf2, 0x81, 0xe5, 0x79, 0xe3, 0xed, 0xab, 0xdd, 0xb3, 0x30,
	0x1c, 0x5a, 0xb6, 0x3c, 0x59, 0x1e, 0xc5, 0xb9, 0x02, 0xaa, 0xd6, 0xaa, 0xc1, 0x19, 0xda, 0x55,
	0x53, 0x2f, 0x34, 0xad, 0x83, 0xcc, 0x79, 0x33, 0x37, 0x8e, 0x7e, 0xa3, 0x23, 0x65, 0xd5, 0xb4,
	0x10, 0x0b, 0xdd, 0xf4, 0x4e, 0x4c, 0xde, 0x9a, 0x40, 0x2a, 0xfd, 0x6c, 0xb0, 0xc9, 0xd2, 0x86,
	0xef, 0x71, 0xf2, 0x75, 0xd2, 0xfd, 0xf5, 0xe1, 0x3f, 0x00, 0x00, 0xff, 0xff, 0x7f, 0xef, 0xdd,
	0xf6, 0x1b, 0x02, 0x00, 0x00,
}
