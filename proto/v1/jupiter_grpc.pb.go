// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// JupiterClient is the client API for Jupiter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JupiterClient interface {
	// Gets information about the cluster.
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
}

type jupiterClient struct {
	cc grpc.ClientConnInterface
}

func NewJupiterClient(cc grpc.ClientConnInterface) JupiterClient {
	return &jupiterClient{cc}
}

func (c *jupiterClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/jupiter.proto.v1.Jupiter/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JupiterServer is the server API for Jupiter service.
// All implementations must embed UnimplementedJupiterServer
// for forward compatibility
type JupiterServer interface {
	// Gets information about the cluster.
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	mustEmbedUnimplementedJupiterServer()
}

// UnimplementedJupiterServer must be embedded to have forward compatible implementations.
type UnimplementedJupiterServer struct {
}

func (UnimplementedJupiterServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedJupiterServer) mustEmbedUnimplementedJupiterServer() {}

// UnsafeJupiterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JupiterServer will
// result in compilation errors.
type UnsafeJupiterServer interface {
	mustEmbedUnimplementedJupiterServer()
}

func RegisterJupiterServer(s grpc.ServiceRegistrar, srv JupiterServer) {
	s.RegisterService(&Jupiter_ServiceDesc, srv)
}

func _Jupiter_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JupiterServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jupiter.proto.v1.Jupiter/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JupiterServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Jupiter_ServiceDesc is the grpc.ServiceDesc for Jupiter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Jupiter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "jupiter.proto.v1.Jupiter",
	HandlerType: (*JupiterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _Jupiter_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/jupiter.proto",
}
