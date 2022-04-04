// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: api/system/v1/resource.proto

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

// ResourceClient is the client API for Resource service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResourceClient interface {
	CreateMenu(ctx context.Context, in *MenuRequest, opts ...grpc.CallOption) (*IDReply, error)
	UpdateMenu(ctx context.Context, in *MenuRequest, opts ...grpc.CallOption) (*IDReply, error)
	DeleteMenu(ctx context.Context, in *IDsRequest, opts ...grpc.CallOption) (*EmptyReply, error)
	GetMenuTree(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*MenuReply, error)
	GetMenuTreeByRole(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*MenuReply, error)
	GetRouteTree(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*RouterReply, error)
	GetRouteTreeByRole(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*RouterReply, error)
	EditRoutePolicy(ctx context.Context, in *RouterRequest, opts ...grpc.CallOption) (*EmptyReply, error)
}

type resourceClient struct {
	cc grpc.ClientConnInterface
}

func NewResourceClient(cc grpc.ClientConnInterface) ResourceClient {
	return &resourceClient{cc}
}

func (c *resourceClient) CreateMenu(ctx context.Context, in *MenuRequest, opts ...grpc.CallOption) (*IDReply, error) {
	out := new(IDReply)
	err := c.cc.Invoke(ctx, "/api.system.v1.Resource/CreateMenu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) UpdateMenu(ctx context.Context, in *MenuRequest, opts ...grpc.CallOption) (*IDReply, error) {
	out := new(IDReply)
	err := c.cc.Invoke(ctx, "/api.system.v1.Resource/UpdateMenu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) DeleteMenu(ctx context.Context, in *IDsRequest, opts ...grpc.CallOption) (*EmptyReply, error) {
	out := new(EmptyReply)
	err := c.cc.Invoke(ctx, "/api.system.v1.Resource/DeleteMenu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) GetMenuTree(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*MenuReply, error) {
	out := new(MenuReply)
	err := c.cc.Invoke(ctx, "/api.system.v1.Resource/GetMenuTree", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) GetMenuTreeByRole(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*MenuReply, error) {
	out := new(MenuReply)
	err := c.cc.Invoke(ctx, "/api.system.v1.Resource/GetMenuTreeByRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) GetRouteTree(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*RouterReply, error) {
	out := new(RouterReply)
	err := c.cc.Invoke(ctx, "/api.system.v1.Resource/GetRouteTree", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) GetRouteTreeByRole(ctx context.Context, in *IDRequest, opts ...grpc.CallOption) (*RouterReply, error) {
	out := new(RouterReply)
	err := c.cc.Invoke(ctx, "/api.system.v1.Resource/GetRouteTreeByRole", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceClient) EditRoutePolicy(ctx context.Context, in *RouterRequest, opts ...grpc.CallOption) (*EmptyReply, error) {
	out := new(EmptyReply)
	err := c.cc.Invoke(ctx, "/api.system.v1.Resource/EditRoutePolicy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceServer is the server API for Resource service.
// All implementations must embed UnimplementedResourceServer
// for forward compatibility
type ResourceServer interface {
	CreateMenu(context.Context, *MenuRequest) (*IDReply, error)
	UpdateMenu(context.Context, *MenuRequest) (*IDReply, error)
	DeleteMenu(context.Context, *IDsRequest) (*EmptyReply, error)
	GetMenuTree(context.Context, *EmptyRequest) (*MenuReply, error)
	GetMenuTreeByRole(context.Context, *IDRequest) (*MenuReply, error)
	GetRouteTree(context.Context, *EmptyRequest) (*RouterReply, error)
	GetRouteTreeByRole(context.Context, *IDRequest) (*RouterReply, error)
	EditRoutePolicy(context.Context, *RouterRequest) (*EmptyReply, error)
	mustEmbedUnimplementedResourceServer()
}

// UnimplementedResourceServer must be embedded to have forward compatible implementations.
type UnimplementedResourceServer struct {
}

func (UnimplementedResourceServer) CreateMenu(context.Context, *MenuRequest) (*IDReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMenu not implemented")
}
func (UnimplementedResourceServer) UpdateMenu(context.Context, *MenuRequest) (*IDReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMenu not implemented")
}
func (UnimplementedResourceServer) DeleteMenu(context.Context, *IDsRequest) (*EmptyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMenu not implemented")
}
func (UnimplementedResourceServer) GetMenuTree(context.Context, *EmptyRequest) (*MenuReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMenuTree not implemented")
}
func (UnimplementedResourceServer) GetMenuTreeByRole(context.Context, *IDRequest) (*MenuReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMenuTreeByRole not implemented")
}
func (UnimplementedResourceServer) GetRouteTree(context.Context, *EmptyRequest) (*RouterReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRouteTree not implemented")
}
func (UnimplementedResourceServer) GetRouteTreeByRole(context.Context, *IDRequest) (*RouterReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRouteTreeByRole not implemented")
}
func (UnimplementedResourceServer) EditRoutePolicy(context.Context, *RouterRequest) (*EmptyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditRoutePolicy not implemented")
}
func (UnimplementedResourceServer) mustEmbedUnimplementedResourceServer() {}

// UnsafeResourceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResourceServer will
// result in compilation errors.
type UnsafeResourceServer interface {
	mustEmbedUnimplementedResourceServer()
}

func RegisterResourceServer(s grpc.ServiceRegistrar, srv ResourceServer) {
	s.RegisterService(&Resource_ServiceDesc, srv)
}

func _Resource_CreateMenu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MenuRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).CreateMenu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.system.v1.Resource/CreateMenu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).CreateMenu(ctx, req.(*MenuRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_UpdateMenu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MenuRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).UpdateMenu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.system.v1.Resource/UpdateMenu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).UpdateMenu(ctx, req.(*MenuRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_DeleteMenu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).DeleteMenu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.system.v1.Resource/DeleteMenu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).DeleteMenu(ctx, req.(*IDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_GetMenuTree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).GetMenuTree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.system.v1.Resource/GetMenuTree",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).GetMenuTree(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_GetMenuTreeByRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).GetMenuTreeByRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.system.v1.Resource/GetMenuTreeByRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).GetMenuTreeByRole(ctx, req.(*IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_GetRouteTree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).GetRouteTree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.system.v1.Resource/GetRouteTree",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).GetRouteTree(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_GetRouteTreeByRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).GetRouteTreeByRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.system.v1.Resource/GetRouteTreeByRole",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).GetRouteTreeByRole(ctx, req.(*IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Resource_EditRoutePolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RouterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceServer).EditRoutePolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.system.v1.Resource/EditRoutePolicy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceServer).EditRoutePolicy(ctx, req.(*RouterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Resource_ServiceDesc is the grpc.ServiceDesc for Resource service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Resource_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.system.v1.Resource",
	HandlerType: (*ResourceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMenu",
			Handler:    _Resource_CreateMenu_Handler,
		},
		{
			MethodName: "UpdateMenu",
			Handler:    _Resource_UpdateMenu_Handler,
		},
		{
			MethodName: "DeleteMenu",
			Handler:    _Resource_DeleteMenu_Handler,
		},
		{
			MethodName: "GetMenuTree",
			Handler:    _Resource_GetMenuTree_Handler,
		},
		{
			MethodName: "GetMenuTreeByRole",
			Handler:    _Resource_GetMenuTreeByRole_Handler,
		},
		{
			MethodName: "GetRouteTree",
			Handler:    _Resource_GetRouteTree_Handler,
		},
		{
			MethodName: "GetRouteTreeByRole",
			Handler:    _Resource_GetRouteTreeByRole_Handler,
		},
		{
			MethodName: "EditRoutePolicy",
			Handler:    _Resource_EditRoutePolicy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/system/v1/resource.proto",
}
