// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: api/cluster/v1alpha1/cluster.proto

package v1alpha1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ClusterInterface_Ping_FullMethodName = "/cluster.v1alpha1.ClusterInterface/Ping"
)

// ClusterInterfaceClient is the client API for ClusterInterface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClusterInterfaceClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Msg, error)
}

type clusterInterfaceClient struct {
	cc grpc.ClientConnInterface
}

func NewClusterInterfaceClient(cc grpc.ClientConnInterface) ClusterInterfaceClient {
	return &clusterInterfaceClient{cc}
}

func (c *clusterInterfaceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Msg, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Msg)
	err := c.cc.Invoke(ctx, ClusterInterface_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClusterInterfaceServer is the server API for ClusterInterface service.
// All implementations must embed UnimplementedClusterInterfaceServer
// for forward compatibility.
type ClusterInterfaceServer interface {
	Ping(context.Context, *emptypb.Empty) (*Msg, error)
	mustEmbedUnimplementedClusterInterfaceServer()
}

// UnimplementedClusterInterfaceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedClusterInterfaceServer struct{}

func (UnimplementedClusterInterfaceServer) Ping(context.Context, *emptypb.Empty) (*Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedClusterInterfaceServer) mustEmbedUnimplementedClusterInterfaceServer() {}
func (UnimplementedClusterInterfaceServer) testEmbeddedByValue()                          {}

// UnsafeClusterInterfaceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClusterInterfaceServer will
// result in compilation errors.
type UnsafeClusterInterfaceServer interface {
	mustEmbedUnimplementedClusterInterfaceServer()
}

func RegisterClusterInterfaceServer(s grpc.ServiceRegistrar, srv ClusterInterfaceServer) {
	// If the following call pancis, it indicates UnimplementedClusterInterfaceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ClusterInterface_ServiceDesc, srv)
}

func _ClusterInterface_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClusterInterfaceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClusterInterface_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClusterInterfaceServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ClusterInterface_ServiceDesc is the grpc.ServiceDesc for ClusterInterface service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClusterInterface_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cluster.v1alpha1.ClusterInterface",
	HandlerType: (*ClusterInterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _ClusterInterface_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/cluster/v1alpha1/cluster.proto",
}
