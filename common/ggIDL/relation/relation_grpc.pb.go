// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: relation/relation.proto

package relation

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

// RelationClient is the client API for Relation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelationClient interface {
	Action(ctx context.Context, in *ActionRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	FollowList(ctx context.Context, in *FollowListRequest, opts ...grpc.CallOption) (*FollowListResponse, error)
	FollowerList(ctx context.Context, in *FollowerListRequest, opts ...grpc.CallOption) (*FollowerListResponse, error)
	FriendList(ctx context.Context, in *FriendListRequest, opts ...grpc.CallOption) (*FriendListResponse, error)
}

type relationClient struct {
	cc grpc.ClientConnInterface
}

func NewRelationClient(cc grpc.ClientConnInterface) RelationClient {
	return &relationClient{cc}
}

func (c *relationClient) Action(ctx context.Context, in *ActionRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/Action", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationClient) FollowList(ctx context.Context, in *FollowListRequest, opts ...grpc.CallOption) (*FollowListResponse, error) {
	out := new(FollowListResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/FollowList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationClient) FollowerList(ctx context.Context, in *FollowerListRequest, opts ...grpc.CallOption) (*FollowerListResponse, error) {
	out := new(FollowerListResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/FollowerList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationClient) FriendList(ctx context.Context, in *FriendListRequest, opts ...grpc.CallOption) (*FriendListResponse, error) {
	out := new(FriendListResponse)
	err := c.cc.Invoke(ctx, "/relation.Relation/FriendList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelationServer is the server API for Relation service.
// All implementations must embed UnimplementedRelationServer
// for forward compatibility
type RelationServer interface {
	Action(context.Context, *ActionRequest) (*ActionResponse, error)
	FollowList(context.Context, *FollowListRequest) (*FollowListResponse, error)
	FollowerList(context.Context, *FollowerListRequest) (*FollowerListResponse, error)
	FriendList(context.Context, *FriendListRequest) (*FriendListResponse, error)
	mustEmbedUnimplementedRelationServer()
}

// UnimplementedRelationServer must be embedded to have forward compatible implementations.
type UnimplementedRelationServer struct {
}

func (UnimplementedRelationServer) Action(context.Context, *ActionRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Action not implemented")
}
func (UnimplementedRelationServer) FollowList(context.Context, *FollowListRequest) (*FollowListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FollowList not implemented")
}
func (UnimplementedRelationServer) FollowerList(context.Context, *FollowerListRequest) (*FollowerListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FollowerList not implemented")
}
func (UnimplementedRelationServer) FriendList(context.Context, *FriendListRequest) (*FriendListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FriendList not implemented")
}
func (UnimplementedRelationServer) mustEmbedUnimplementedRelationServer() {}

// UnsafeRelationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelationServer will
// result in compilation errors.
type UnsafeRelationServer interface {
	mustEmbedUnimplementedRelationServer()
}

func RegisterRelationServer(s grpc.ServiceRegistrar, srv RelationServer) {
	s.RegisterService(&Relation_ServiceDesc, srv)
}

func _Relation_Action_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).Action(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/Action",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).Action(ctx, req.(*ActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relation_FollowList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).FollowList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/FollowList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).FollowList(ctx, req.(*FollowListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relation_FollowerList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FollowerListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).FollowerList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/FollowerList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).FollowerList(ctx, req.(*FollowerListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relation_FriendList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FriendListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelationServer).FriendList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relation.Relation/FriendList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelationServer).FriendList(ctx, req.(*FriendListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Relation_ServiceDesc is the grpc.ServiceDesc for Relation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Relation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "relation.Relation",
	HandlerType: (*RelationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Action",
			Handler:    _Relation_Action_Handler,
		},
		{
			MethodName: "FollowList",
			Handler:    _Relation_FollowList_Handler,
		},
		{
			MethodName: "FollowerList",
			Handler:    _Relation_FollowerList_Handler,
		},
		{
			MethodName: "FriendList",
			Handler:    _Relation_FriendList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "relation/relation.proto",
}
