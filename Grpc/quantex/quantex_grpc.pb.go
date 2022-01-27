// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/quantex.proto

package GrpcQuantex

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// QuantexServiceClient is the client API for QuantexService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QuantexServiceClient interface {
	UpdateOrder(ctx context.Context, in *UpdateOrderRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type quantexServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQuantexServiceClient(cc grpc.ClientConnInterface) QuantexServiceClient {
	return &quantexServiceClient{cc}
}

func (c *quantexServiceClient) UpdateOrder(ctx context.Context, in *UpdateOrderRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/GrpcQuantex.QuantexService/UpdateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QuantexServiceServer is the server API for QuantexService service.
// All implementations should embed UnimplementedQuantexServiceServer
// for forward compatibility
type QuantexServiceServer interface {
	UpdateOrder(context.Context, *UpdateOrderRequest) (*empty.Empty, error)
}

// UnimplementedQuantexServiceServer should be embedded to have forward compatible implementations.
type UnimplementedQuantexServiceServer struct {
}

func (UnimplementedQuantexServiceServer) UpdateOrder(context.Context, *UpdateOrderRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrder not implemented")
}

// UnsafeQuantexServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QuantexServiceServer will
// result in compilation errors.
type UnsafeQuantexServiceServer interface {
	mustEmbedUnimplementedQuantexServiceServer()
}

func RegisterQuantexServiceServer(s grpc.ServiceRegistrar, srv QuantexServiceServer) {
	s.RegisterService(&QuantexService_ServiceDesc, srv)
}

func _QuantexService_UpdateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuantexServiceServer).UpdateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GrpcQuantex.QuantexService/UpdateOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuantexServiceServer).UpdateOrder(ctx, req.(*UpdateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// QuantexService_ServiceDesc is the grpc.ServiceDesc for QuantexService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QuantexService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "GrpcQuantex.QuantexService",
	HandlerType: (*QuantexServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateOrder",
			Handler:    _QuantexService_UpdateOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/quantex.proto",
}
