// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/services.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EthereumServiceClient is the client API for EthereumService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EthereumServiceClient interface {
	IsWaitingPermit(ctx context.Context, in *WalletAddress, opts ...grpc.CallOption) (*IsWaitingPermitOutput, error)
	GetBalance(ctx context.Context, in *GetBalanceInput, opts ...grpc.CallOption) (*GetBalanceOutput, error)
	GetAllowance(ctx context.Context, in *GetAllowanceInput, opts ...grpc.CallOption) (*GetAllowanceOutput, error)
}

type ethereumServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEthereumServiceClient(cc grpc.ClientConnInterface) EthereumServiceClient {
	return &ethereumServiceClient{cc}
}

func (c *ethereumServiceClient) IsWaitingPermit(ctx context.Context, in *WalletAddress, opts ...grpc.CallOption) (*IsWaitingPermitOutput, error) {
	out := new(IsWaitingPermitOutput)
	err := c.cc.Invoke(ctx, "/pb.EthereumService/IsWaitingPermit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ethereumServiceClient) GetBalance(ctx context.Context, in *GetBalanceInput, opts ...grpc.CallOption) (*GetBalanceOutput, error) {
	out := new(GetBalanceOutput)
	err := c.cc.Invoke(ctx, "/pb.EthereumService/GetBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ethereumServiceClient) GetAllowance(ctx context.Context, in *GetAllowanceInput, opts ...grpc.CallOption) (*GetAllowanceOutput, error) {
	out := new(GetAllowanceOutput)
	err := c.cc.Invoke(ctx, "/pb.EthereumService/GetAllowance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EthereumServiceServer is the server API for EthereumService service.
// All implementations must embed UnimplementedEthereumServiceServer
// for forward compatibility
type EthereumServiceServer interface {
	IsWaitingPermit(context.Context, *WalletAddress) (*IsWaitingPermitOutput, error)
	GetBalance(context.Context, *GetBalanceInput) (*GetBalanceOutput, error)
	GetAllowance(context.Context, *GetAllowanceInput) (*GetAllowanceOutput, error)
	mustEmbedUnimplementedEthereumServiceServer()
}

// UnimplementedEthereumServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEthereumServiceServer struct {
}

func (UnimplementedEthereumServiceServer) IsWaitingPermit(context.Context, *WalletAddress) (*IsWaitingPermitOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsWaitingPermit not implemented")
}
func (UnimplementedEthereumServiceServer) GetBalance(context.Context, *GetBalanceInput) (*GetBalanceOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
func (UnimplementedEthereumServiceServer) GetAllowance(context.Context, *GetAllowanceInput) (*GetAllowanceOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllowance not implemented")
}
func (UnimplementedEthereumServiceServer) mustEmbedUnimplementedEthereumServiceServer() {}

// UnsafeEthereumServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EthereumServiceServer will
// result in compilation errors.
type UnsafeEthereumServiceServer interface {
	mustEmbedUnimplementedEthereumServiceServer()
}

func RegisterEthereumServiceServer(s grpc.ServiceRegistrar, srv EthereumServiceServer) {
	s.RegisterService(&EthereumService_ServiceDesc, srv)
}

func _EthereumService_IsWaitingPermit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WalletAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EthereumServiceServer).IsWaitingPermit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.EthereumService/IsWaitingPermit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EthereumServiceServer).IsWaitingPermit(ctx, req.(*WalletAddress))
	}
	return interceptor(ctx, in, info, handler)
}

func _EthereumService_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EthereumServiceServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.EthereumService/GetBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EthereumServiceServer).GetBalance(ctx, req.(*GetBalanceInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _EthereumService_GetAllowance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllowanceInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EthereumServiceServer).GetAllowance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.EthereumService/GetAllowance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EthereumServiceServer).GetAllowance(ctx, req.(*GetAllowanceInput))
	}
	return interceptor(ctx, in, info, handler)
}

// EthereumService_ServiceDesc is the grpc.ServiceDesc for EthereumService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EthereumService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.EthereumService",
	HandlerType: (*EthereumServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsWaitingPermit",
			Handler:    _EthereumService_IsWaitingPermit_Handler,
		},
		{
			MethodName: "GetBalance",
			Handler:    _EthereumService_GetBalance_Handler,
		},
		{
			MethodName: "GetAllowance",
			Handler:    _EthereumService_GetAllowance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/services.proto",
}

// PolygonServiceClient is the client API for PolygonService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PolygonServiceClient interface {
	IsWaitingPermit(ctx context.Context, in *WalletAddress, opts ...grpc.CallOption) (*IsWaitingPermitOutput, error)
	GetBalance(ctx context.Context, in *GetBalanceInput, opts ...grpc.CallOption) (*GetBalanceOutput, error)
	GetAllowance(ctx context.Context, in *GetAllowanceInput, opts ...grpc.CallOption) (*GetAllowanceOutput, error)
}

type polygonServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPolygonServiceClient(cc grpc.ClientConnInterface) PolygonServiceClient {
	return &polygonServiceClient{cc}
}

func (c *polygonServiceClient) IsWaitingPermit(ctx context.Context, in *WalletAddress, opts ...grpc.CallOption) (*IsWaitingPermitOutput, error) {
	out := new(IsWaitingPermitOutput)
	err := c.cc.Invoke(ctx, "/pb.PolygonService/IsWaitingPermit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *polygonServiceClient) GetBalance(ctx context.Context, in *GetBalanceInput, opts ...grpc.CallOption) (*GetBalanceOutput, error) {
	out := new(GetBalanceOutput)
	err := c.cc.Invoke(ctx, "/pb.PolygonService/GetBalance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *polygonServiceClient) GetAllowance(ctx context.Context, in *GetAllowanceInput, opts ...grpc.CallOption) (*GetAllowanceOutput, error) {
	out := new(GetAllowanceOutput)
	err := c.cc.Invoke(ctx, "/pb.PolygonService/GetAllowance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PolygonServiceServer is the server API for PolygonService service.
// All implementations must embed UnimplementedPolygonServiceServer
// for forward compatibility
type PolygonServiceServer interface {
	IsWaitingPermit(context.Context, *WalletAddress) (*IsWaitingPermitOutput, error)
	GetBalance(context.Context, *GetBalanceInput) (*GetBalanceOutput, error)
	GetAllowance(context.Context, *GetAllowanceInput) (*GetAllowanceOutput, error)
	mustEmbedUnimplementedPolygonServiceServer()
}

// UnimplementedPolygonServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPolygonServiceServer struct {
}

func (UnimplementedPolygonServiceServer) IsWaitingPermit(context.Context, *WalletAddress) (*IsWaitingPermitOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsWaitingPermit not implemented")
}
func (UnimplementedPolygonServiceServer) GetBalance(context.Context, *GetBalanceInput) (*GetBalanceOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalance not implemented")
}
func (UnimplementedPolygonServiceServer) GetAllowance(context.Context, *GetAllowanceInput) (*GetAllowanceOutput, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllowance not implemented")
}
func (UnimplementedPolygonServiceServer) mustEmbedUnimplementedPolygonServiceServer() {}

// UnsafePolygonServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PolygonServiceServer will
// result in compilation errors.
type UnsafePolygonServiceServer interface {
	mustEmbedUnimplementedPolygonServiceServer()
}

func RegisterPolygonServiceServer(s grpc.ServiceRegistrar, srv PolygonServiceServer) {
	s.RegisterService(&PolygonService_ServiceDesc, srv)
}

func _PolygonService_IsWaitingPermit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WalletAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolygonServiceServer).IsWaitingPermit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PolygonService/IsWaitingPermit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolygonServiceServer).IsWaitingPermit(ctx, req.(*WalletAddress))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolygonService_GetBalance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolygonServiceServer).GetBalance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PolygonService/GetBalance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolygonServiceServer).GetBalance(ctx, req.(*GetBalanceInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolygonService_GetAllowance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllowanceInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolygonServiceServer).GetAllowance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.PolygonService/GetAllowance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolygonServiceServer).GetAllowance(ctx, req.(*GetAllowanceInput))
	}
	return interceptor(ctx, in, info, handler)
}

// PolygonService_ServiceDesc is the grpc.ServiceDesc for PolygonService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PolygonService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.PolygonService",
	HandlerType: (*PolygonServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsWaitingPermit",
			Handler:    _PolygonService_IsWaitingPermit_Handler,
		},
		{
			MethodName: "GetBalance",
			Handler:    _PolygonService_GetBalance_Handler,
		},
		{
			MethodName: "GetAllowance",
			Handler:    _PolygonService_GetAllowance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/services.proto",
}
