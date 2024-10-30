// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: api/cloud/v1alpha1/cloud.proto

package v1alpha1

import (
	context "context"
	common "github.com/f-rambo/ship/api/common"
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
	CloudInterface_Ping_FullMethodName                      = "/cloud.v1alpha1.CloudInterface/Ping"
	CloudInterface_NodeInit_FullMethodName                  = "/cloud.v1alpha1.CloudInterface/NodeInit"
	CloudInterface_InstallKubeadmKubeletCriO_FullMethodName = "/cloud.v1alpha1.CloudInterface/InstallKubeadmKubeletCriO"
	CloudInterface_KubeadmJoin_FullMethodName               = "/cloud.v1alpha1.CloudInterface/KubeadmJoin"
	CloudInterface_KubeadmInit_FullMethodName               = "/cloud.v1alpha1.CloudInterface/KubeadmInit"
	CloudInterface_KubeadmReset_FullMethodName              = "/cloud.v1alpha1.CloudInterface/KubeadmReset"
	CloudInterface_KubeadmUpgrade_FullMethodName            = "/cloud.v1alpha1.CloudInterface/KubeadmUpgrade"
)

// CloudInterfaceClient is the client API for CloudInterface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CloudInterfaceClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*common.Msg, error)
	// node init
	NodeInit(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error)
	// InstallKubeadmKubeletCriO
	InstallKubeadmKubeletCriO(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error)
	// KubeadmJoin
	KubeadmJoin(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error)
	// KubeadmInit
	KubeadmInit(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error)
	// KubeadmReset
	KubeadmReset(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error)
	// KubeadmUpgrade
	KubeadmUpgrade(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error)
}

type cloudInterfaceClient struct {
	cc grpc.ClientConnInterface
}

func NewCloudInterfaceClient(cc grpc.ClientConnInterface) CloudInterfaceClient {
	return &cloudInterfaceClient{cc}
}

func (c *cloudInterfaceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*common.Msg, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(common.Msg)
	err := c.cc.Invoke(ctx, CloudInterface_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudInterfaceClient) NodeInit(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(common.Msg)
	err := c.cc.Invoke(ctx, CloudInterface_NodeInit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudInterfaceClient) InstallKubeadmKubeletCriO(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(common.Msg)
	err := c.cc.Invoke(ctx, CloudInterface_InstallKubeadmKubeletCriO_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudInterfaceClient) KubeadmJoin(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(common.Msg)
	err := c.cc.Invoke(ctx, CloudInterface_KubeadmJoin_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudInterfaceClient) KubeadmInit(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(common.Msg)
	err := c.cc.Invoke(ctx, CloudInterface_KubeadmInit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudInterfaceClient) KubeadmReset(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(common.Msg)
	err := c.cc.Invoke(ctx, CloudInterface_KubeadmReset_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cloudInterfaceClient) KubeadmUpgrade(ctx context.Context, in *Cloud, opts ...grpc.CallOption) (*common.Msg, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(common.Msg)
	err := c.cc.Invoke(ctx, CloudInterface_KubeadmUpgrade_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CloudInterfaceServer is the server API for CloudInterface service.
// All implementations must embed UnimplementedCloudInterfaceServer
// for forward compatibility.
type CloudInterfaceServer interface {
	Ping(context.Context, *emptypb.Empty) (*common.Msg, error)
	// node init
	NodeInit(context.Context, *Cloud) (*common.Msg, error)
	// InstallKubeadmKubeletCriO
	InstallKubeadmKubeletCriO(context.Context, *Cloud) (*common.Msg, error)
	// KubeadmJoin
	KubeadmJoin(context.Context, *Cloud) (*common.Msg, error)
	// KubeadmInit
	KubeadmInit(context.Context, *Cloud) (*common.Msg, error)
	// KubeadmReset
	KubeadmReset(context.Context, *Cloud) (*common.Msg, error)
	// KubeadmUpgrade
	KubeadmUpgrade(context.Context, *Cloud) (*common.Msg, error)
	mustEmbedUnimplementedCloudInterfaceServer()
}

// UnimplementedCloudInterfaceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCloudInterfaceServer struct{}

func (UnimplementedCloudInterfaceServer) Ping(context.Context, *emptypb.Empty) (*common.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedCloudInterfaceServer) NodeInit(context.Context, *Cloud) (*common.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodeInit not implemented")
}
func (UnimplementedCloudInterfaceServer) InstallKubeadmKubeletCriO(context.Context, *Cloud) (*common.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InstallKubeadmKubeletCriO not implemented")
}
func (UnimplementedCloudInterfaceServer) KubeadmJoin(context.Context, *Cloud) (*common.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KubeadmJoin not implemented")
}
func (UnimplementedCloudInterfaceServer) KubeadmInit(context.Context, *Cloud) (*common.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KubeadmInit not implemented")
}
func (UnimplementedCloudInterfaceServer) KubeadmReset(context.Context, *Cloud) (*common.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KubeadmReset not implemented")
}
func (UnimplementedCloudInterfaceServer) KubeadmUpgrade(context.Context, *Cloud) (*common.Msg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KubeadmUpgrade not implemented")
}
func (UnimplementedCloudInterfaceServer) mustEmbedUnimplementedCloudInterfaceServer() {}
func (UnimplementedCloudInterfaceServer) testEmbeddedByValue()                        {}

// UnsafeCloudInterfaceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CloudInterfaceServer will
// result in compilation errors.
type UnsafeCloudInterfaceServer interface {
	mustEmbedUnimplementedCloudInterfaceServer()
}

func RegisterCloudInterfaceServer(s grpc.ServiceRegistrar, srv CloudInterfaceServer) {
	// If the following call pancis, it indicates UnimplementedCloudInterfaceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CloudInterface_ServiceDesc, srv)
}

func _CloudInterface_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudInterfaceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudInterface_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudInterfaceServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudInterface_NodeInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cloud)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudInterfaceServer).NodeInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudInterface_NodeInit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudInterfaceServer).NodeInit(ctx, req.(*Cloud))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudInterface_InstallKubeadmKubeletCriO_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cloud)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudInterfaceServer).InstallKubeadmKubeletCriO(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudInterface_InstallKubeadmKubeletCriO_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudInterfaceServer).InstallKubeadmKubeletCriO(ctx, req.(*Cloud))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudInterface_KubeadmJoin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cloud)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudInterfaceServer).KubeadmJoin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudInterface_KubeadmJoin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudInterfaceServer).KubeadmJoin(ctx, req.(*Cloud))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudInterface_KubeadmInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cloud)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudInterfaceServer).KubeadmInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudInterface_KubeadmInit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudInterfaceServer).KubeadmInit(ctx, req.(*Cloud))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudInterface_KubeadmReset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cloud)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudInterfaceServer).KubeadmReset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudInterface_KubeadmReset_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudInterfaceServer).KubeadmReset(ctx, req.(*Cloud))
	}
	return interceptor(ctx, in, info, handler)
}

func _CloudInterface_KubeadmUpgrade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cloud)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CloudInterfaceServer).KubeadmUpgrade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CloudInterface_KubeadmUpgrade_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CloudInterfaceServer).KubeadmUpgrade(ctx, req.(*Cloud))
	}
	return interceptor(ctx, in, info, handler)
}

// CloudInterface_ServiceDesc is the grpc.ServiceDesc for CloudInterface service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CloudInterface_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cloud.v1alpha1.CloudInterface",
	HandlerType: (*CloudInterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _CloudInterface_Ping_Handler,
		},
		{
			MethodName: "NodeInit",
			Handler:    _CloudInterface_NodeInit_Handler,
		},
		{
			MethodName: "InstallKubeadmKubeletCriO",
			Handler:    _CloudInterface_InstallKubeadmKubeletCriO_Handler,
		},
		{
			MethodName: "KubeadmJoin",
			Handler:    _CloudInterface_KubeadmJoin_Handler,
		},
		{
			MethodName: "KubeadmInit",
			Handler:    _CloudInterface_KubeadmInit_Handler,
		},
		{
			MethodName: "KubeadmReset",
			Handler:    _CloudInterface_KubeadmReset_Handler,
		},
		{
			MethodName: "KubeadmUpgrade",
			Handler:    _CloudInterface_KubeadmUpgrade_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/cloud/v1alpha1/cloud.proto",
}
