// Code generated by protoc-gen-go. DO NOT EDIT.
// source: control.proto

package ctl

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

func init() { proto.RegisterFile("control.proto", fileDescriptor_0c5120591600887d) }

var fileDescriptor_0c5120591600887d = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xcd, 0x4a, 0xc3, 0x40,
	0x10, 0xc7, 0x5b, 0x14, 0x85, 0xd5, 0x78, 0x98, 0xc6, 0x2a, 0x39, 0xf6, 0x01, 0x42, 0xd1, 0x83,
	0xe0, 0xd1, 0x4a, 0x0f, 0xa2, 0x52, 0x9b, 0x07, 0x90, 0x75, 0x1d, 0x4a, 0x30, 0xc9, 0xa6, 0xb3,
	0xd3, 0x4a, 0x9e, 0xc3, 0x17, 0x96, 0xfd, 0xa8, 0x26, 0x4d, 0x8f, 0xf3, 0xdb, 0xff, 0x07, 0x93,
	0x89, 0x88, 0x94, 0xae, 0x98, 0x74, 0x91, 0xd6, 0xa4, 0x59, 0xc3, 0x91, 0xe2, 0x22, 0x89, 0x0c,
	0x6b, 0x92, 0x2b, 0xf4, 0x2c, 0x81, 0x0a, 0xf9, 0x5b, 0xd3, 0xd7, 0xbb, 0x51, 0xb2, 0x0a, 0xec,
	0xdc, 0x34, 0x86, 0xb1, 0xf4, 0xd3, 0xcd, 0xcf, 0xb1, 0x38, 0x7d, 0x59, 0x95, 0x3c, 0xe3, 0x02,
	0x66, 0xe2, 0x22, 0xf3, 0xf6, 0x05, 0x61, 0x2d, 0x09, 0x61, 0x9c, 0x2a, 0x2e, 0xd2, 0x2e, 0x5c,
	0xe2, 0x3a, 0xb9, 0x3a, 0xc8, 0x4d, 0x3d, 0x19, 0xc0, 0xbd, 0x38, 0x0b, 0x3c, 0x53, 0xb2, 0x82,
	0x51, 0x5b, 0x69, 0x89, 0xb5, 0xc7, 0x7d, 0xe8, 0xbc, 0x0f, 0x22, 0x0a, 0x70, 0xae, 0xa9, 0x94,
	0x0c, 0x97, 0x6d, 0xa1, 0x67, 0xd6, 0x3f, 0x3e, 0x84, 0x6d, 0xc2, 0x74, 0xe8, 0xfa, 0xdd, 0x82,
	0x6f, 0x1b, 0xa4, 0x66, 0xd7, 0xff, 0x4f, 0x5a, 0xfd, 0x6d, 0xe8, 0xfa, 0xef, 0x84, 0xf0, 0x30,
	0x63, 0x5d, 0x03, 0xb4, 0x54, 0x16, 0x58, 0xe7, 0xa8, 0xc7, 0xfe, 0x96, 0x0e, 0x4c, 0x12, 0x43,
	0x57, 0x25, 0x89, 0xf7, 0x4b, 0x03, 0x74, 0xde, 0x27, 0x11, 0xbf, 0xfa, 0x2b, 0x3d, 0xe7, 0x86,
	0x17, 0xa4, 0xb7, 0xf9, 0x27, 0x92, 0x81, 0x6b, 0xa7, 0xdf, 0xcd, 0xf6, 0x6d, 0x89, 0xeb, 0x0d,
	0x1a, 0x0e, 0xeb, 0x77, 0x5f, 0xea, 0xa2, 0x99, 0x0c, 0x60, 0x2e, 0x20, 0x64, 0xd9, 0xaf, 0xfa,
	0x88, 0xdb, 0x5c, 0xa1, 0x09, 0x57, 0xf4, 0x53, 0x38, 0x81, 0xcb, 0x89, 0x7b, 0xdc, 0xa5, 0x4c,
	0x87, 0x1f, 0x27, 0xee, 0xe7, 0xb8, 0xfd, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xd1, 0xd3, 0xc5, 0x2b,
	0x63, 0x02, 0x00, 0x00,
}

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
	// Retrieve details of nonvolatile storage on server, including health info
	StorageScan(ctx context.Context, in *StorageScanReq, opts ...grpc.CallOption) (*StorageScanResp, error)
	// Format nonvolatile storage devices for use with DAOS
	StorageFormat(ctx context.Context, in *StorageFormatReq, opts ...grpc.CallOption) (MgmtCtl_StorageFormatClient, error)
	// Query DAOS system status
	SystemQuery(ctx context.Context, in *SystemQueryReq, opts ...grpc.CallOption) (*SystemQueryResp, error)
	// Stop DAOS system (shutdown data-plane instances)
	SystemStop(ctx context.Context, in *SystemStopReq, opts ...grpc.CallOption) (*SystemStopResp, error)
	// Start DAOS system (restart data-plane instances)
	SystemStart(ctx context.Context, in *SystemStartReq, opts ...grpc.CallOption) (*SystemStartResp, error)
	// Retrieve a list of supported fabric providers
	NetworkListProviders(ctx context.Context, in *ProviderListRequest, opts ...grpc.CallOption) (*ProviderListReply, error)
	// Perform a fabric scan to determine the available provider, device, NUMA node combinations
	NetworkScanDevices(ctx context.Context, in *DeviceScanRequest, opts ...grpc.CallOption) (MgmtCtl_NetworkScanDevicesClient, error)
}

type mgmtCtlClient struct {
	cc *grpc.ClientConn
}

func NewMgmtCtlClient(cc *grpc.ClientConn) MgmtCtlClient {
	return &mgmtCtlClient{cc}
}

func (c *mgmtCtlClient) StoragePrepare(ctx context.Context, in *StoragePrepareReq, opts ...grpc.CallOption) (*StoragePrepareResp, error) {
	out := new(StoragePrepareResp)
	err := c.cc.Invoke(ctx, "/ctl.MgmtCtl/StoragePrepare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mgmtCtlClient) StorageScan(ctx context.Context, in *StorageScanReq, opts ...grpc.CallOption) (*StorageScanResp, error) {
	out := new(StorageScanResp)
	err := c.cc.Invoke(ctx, "/ctl.MgmtCtl/StorageScan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mgmtCtlClient) StorageFormat(ctx context.Context, in *StorageFormatReq, opts ...grpc.CallOption) (MgmtCtl_StorageFormatClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MgmtCtl_serviceDesc.Streams[0], "/ctl.MgmtCtl/StorageFormat", opts...)
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

func (c *mgmtCtlClient) SystemQuery(ctx context.Context, in *SystemQueryReq, opts ...grpc.CallOption) (*SystemQueryResp, error) {
	out := new(SystemQueryResp)
	err := c.cc.Invoke(ctx, "/ctl.MgmtCtl/SystemQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mgmtCtlClient) SystemStop(ctx context.Context, in *SystemStopReq, opts ...grpc.CallOption) (*SystemStopResp, error) {
	out := new(SystemStopResp)
	err := c.cc.Invoke(ctx, "/ctl.MgmtCtl/SystemStop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mgmtCtlClient) SystemStart(ctx context.Context, in *SystemStartReq, opts ...grpc.CallOption) (*SystemStartResp, error) {
	out := new(SystemStartResp)
	err := c.cc.Invoke(ctx, "/ctl.MgmtCtl/SystemStart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mgmtCtlClient) NetworkListProviders(ctx context.Context, in *ProviderListRequest, opts ...grpc.CallOption) (*ProviderListReply, error) {
	out := new(ProviderListReply)
	err := c.cc.Invoke(ctx, "/ctl.MgmtCtl/NetworkListProviders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mgmtCtlClient) NetworkScanDevices(ctx context.Context, in *DeviceScanRequest, opts ...grpc.CallOption) (MgmtCtl_NetworkScanDevicesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MgmtCtl_serviceDesc.Streams[1], "/ctl.MgmtCtl/NetworkScanDevices", opts...)
	if err != nil {
		return nil, err
	}
	x := &mgmtCtlNetworkScanDevicesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MgmtCtl_NetworkScanDevicesClient interface {
	Recv() (*DeviceScanReply, error)
	grpc.ClientStream
}

type mgmtCtlNetworkScanDevicesClient struct {
	grpc.ClientStream
}

func (x *mgmtCtlNetworkScanDevicesClient) Recv() (*DeviceScanReply, error) {
	m := new(DeviceScanReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MgmtCtlServer is the server API for MgmtCtl service.
type MgmtCtlServer interface {
	// Prepare nonvolatile storage devices for use with DAOS
	StoragePrepare(context.Context, *StoragePrepareReq) (*StoragePrepareResp, error)
	// Retrieve details of nonvolatile storage on server, including health info
	StorageScan(context.Context, *StorageScanReq) (*StorageScanResp, error)
	// Format nonvolatile storage devices for use with DAOS
	StorageFormat(*StorageFormatReq, MgmtCtl_StorageFormatServer) error
	// Query DAOS system status
	SystemQuery(context.Context, *SystemQueryReq) (*SystemQueryResp, error)
	// Stop DAOS system (shutdown data-plane instances)
	SystemStop(context.Context, *SystemStopReq) (*SystemStopResp, error)
	// Start DAOS system (restart data-plane instances)
	SystemStart(context.Context, *SystemStartReq) (*SystemStartResp, error)
	// Retrieve a list of supported fabric providers
	NetworkListProviders(context.Context, *ProviderListRequest) (*ProviderListReply, error)
	// Perform a fabric scan to determine the available provider, device, NUMA node combinations
	NetworkScanDevices(*DeviceScanRequest, MgmtCtl_NetworkScanDevicesServer) error
}

// UnimplementedMgmtCtlServer can be embedded to have forward compatible implementations.
type UnimplementedMgmtCtlServer struct {
}

func (*UnimplementedMgmtCtlServer) StoragePrepare(ctx context.Context, req *StoragePrepareReq) (*StoragePrepareResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoragePrepare not implemented")
}
func (*UnimplementedMgmtCtlServer) StorageScan(ctx context.Context, req *StorageScanReq) (*StorageScanResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StorageScan not implemented")
}
func (*UnimplementedMgmtCtlServer) StorageFormat(req *StorageFormatReq, srv MgmtCtl_StorageFormatServer) error {
	return status.Errorf(codes.Unimplemented, "method StorageFormat not implemented")
}
func (*UnimplementedMgmtCtlServer) SystemQuery(ctx context.Context, req *SystemQueryReq) (*SystemQueryResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SystemQuery not implemented")
}
func (*UnimplementedMgmtCtlServer) SystemStop(ctx context.Context, req *SystemStopReq) (*SystemStopResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SystemStop not implemented")
}
func (*UnimplementedMgmtCtlServer) SystemStart(ctx context.Context, req *SystemStartReq) (*SystemStartResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SystemStart not implemented")
}
func (*UnimplementedMgmtCtlServer) NetworkListProviders(ctx context.Context, req *ProviderListRequest) (*ProviderListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NetworkListProviders not implemented")
}
func (*UnimplementedMgmtCtlServer) NetworkScanDevices(req *DeviceScanRequest, srv MgmtCtl_NetworkScanDevicesServer) error {
	return status.Errorf(codes.Unimplemented, "method NetworkScanDevices not implemented")
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
		FullMethod: "/ctl.MgmtCtl/StoragePrepare",
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
		FullMethod: "/ctl.MgmtCtl/StorageScan",
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

func _MgmtCtl_SystemQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SystemQueryReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MgmtCtlServer).SystemQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ctl.MgmtCtl/SystemQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MgmtCtlServer).SystemQuery(ctx, req.(*SystemQueryReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MgmtCtl_SystemStop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SystemStopReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MgmtCtlServer).SystemStop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ctl.MgmtCtl/SystemStop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MgmtCtlServer).SystemStop(ctx, req.(*SystemStopReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MgmtCtl_SystemStart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SystemStartReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MgmtCtlServer).SystemStart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ctl.MgmtCtl/SystemStart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MgmtCtlServer).SystemStart(ctx, req.(*SystemStartReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _MgmtCtl_NetworkListProviders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProviderListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MgmtCtlServer).NetworkListProviders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ctl.MgmtCtl/NetworkListProviders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MgmtCtlServer).NetworkListProviders(ctx, req.(*ProviderListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MgmtCtl_NetworkScanDevices_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DeviceScanRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MgmtCtlServer).NetworkScanDevices(m, &mgmtCtlNetworkScanDevicesServer{stream})
}

type MgmtCtl_NetworkScanDevicesServer interface {
	Send(*DeviceScanReply) error
	grpc.ServerStream
}

type mgmtCtlNetworkScanDevicesServer struct {
	grpc.ServerStream
}

func (x *mgmtCtlNetworkScanDevicesServer) Send(m *DeviceScanReply) error {
	return x.ServerStream.SendMsg(m)
}

var _MgmtCtl_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ctl.MgmtCtl",
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
		{
			MethodName: "SystemQuery",
			Handler:    _MgmtCtl_SystemQuery_Handler,
		},
		{
			MethodName: "SystemStop",
			Handler:    _MgmtCtl_SystemStop_Handler,
		},
		{
			MethodName: "SystemStart",
			Handler:    _MgmtCtl_SystemStart_Handler,
		},
		{
			MethodName: "NetworkListProviders",
			Handler:    _MgmtCtl_NetworkListProviders_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StorageFormat",
			Handler:       _MgmtCtl_StorageFormat_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "NetworkScanDevices",
			Handler:       _MgmtCtl_NetworkScanDevices_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "control.proto",
}