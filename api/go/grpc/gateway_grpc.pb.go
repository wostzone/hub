// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CapGatewayClient is the client API for CapGateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CapGatewayClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error)
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingRequest, error)
}

type capGatewayClient struct {
	cc grpc.ClientConnInterface
}

func NewCapGatewayClient(cc grpc.ClientConnInterface) CapGatewayClient {
	return &capGatewayClient{cc}
}

func (c *capGatewayClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/hiveot.grpc.CapGateway/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *capGatewayClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*PingRequest, error) {
	out := new(PingRequest)
	err := c.cc.Invoke(ctx, "/hiveot.grpc.CapGateway/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CapGatewayServer is the server API for CapGateway service.
// All implementations must embed UnimplementedCapGatewayServer
// for forward compatibility
type CapGatewayServer interface {
	Login(context.Context, *LoginRequest) (*LoginReply, error)
	Ping(context.Context, *emptypb.Empty) (*PingRequest, error)
	mustEmbedUnimplementedCapGatewayServer()
}

// UnimplementedCapGatewayServer must be embedded to have forward compatible implementations.
type UnimplementedCapGatewayServer struct {
}

func (UnimplementedCapGatewayServer) Login(context.Context, *LoginRequest) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedCapGatewayServer) Ping(context.Context, *emptypb.Empty) (*PingRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedCapGatewayServer) mustEmbedUnimplementedCapGatewayServer() {}

// UnsafeCapGatewayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CapGatewayServer will
// result in compilation errors.
type UnsafeCapGatewayServer interface {
	mustEmbedUnimplementedCapGatewayServer()
}

func RegisterCapGatewayServer(s grpc.ServiceRegistrar, srv CapGatewayServer) {
	s.RegisterService(&CapGateway_ServiceDesc, srv)
}

func _CapGateway_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CapGatewayServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hiveot.grpc.CapGateway/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CapGatewayServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CapGateway_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CapGatewayServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hiveot.grpc.CapGateway/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CapGatewayServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// CapGateway_ServiceDesc is the grpc.ServiceDesc for CapGateway service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CapGateway_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hiveot.grpc.CapGateway",
	HandlerType: (*CapGatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _CapGateway_Login_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _CapGateway_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gateway.proto",
}
