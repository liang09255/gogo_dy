// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: video/video.proto

package video

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

// VideoServiceClient is the client API for VideoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VideoServiceClient interface {
	// 点赞方法
	FavoriteAction(ctx context.Context, in *FavoriteActionRequest, opts ...grpc.CallOption) (*FavoriteActionResponse, error)
	// 喜欢列表方法
	FavoriteList(ctx context.Context, in *FavoriteListRequest, opts ...grpc.CallOption) (*FavoriteListResponse, error)
	// 视频评论方法
	CommentAction(ctx context.Context, in *CommentActionRequest, opts ...grpc.CallOption) (*CommentActionResponse, error)
	// 评论列表方法
	CommentList(ctx context.Context, in *CommentListRequest, opts ...grpc.CallOption) (*CommentListResponse, error)
	// 上传视频
	UploadVideo(ctx context.Context, in *UploadVideoRequest, opts ...grpc.CallOption) (*UploadVideoResponse, error)
	// 视频列表
	VideoList(ctx context.Context, in *VideoListRequest, opts ...grpc.CallOption) (*VideoListResponse, error)
	// 视频流
	Feed(ctx context.Context, in *FeedRequest, opts ...grpc.CallOption) (*FeedResponse, error)
	// 获取用户总点赞数
	GetTotalFavoriteCount(ctx context.Context, in *GetTotalFavoriteCountRequest, opts ...grpc.CallOption) (*GetTotalFavoriteCountResponse, error)
	// 获取用户总视频数
	GetTotalVideoCount(ctx context.Context, in *GetTotalVideoCountRequest, opts ...grpc.CallOption) (*GetTotalVideoCountResponse, error)
	// 获取用户总喜欢数
	GetTotalLikeCount(ctx context.Context, in *GetTotalLikeCountRequest, opts ...grpc.CallOption) (*GetTotalLikeCountResponse, error)
}

type videoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVideoServiceClient(cc grpc.ClientConnInterface) VideoServiceClient {
	return &videoServiceClient{cc}
}

func (c *videoServiceClient) FavoriteAction(ctx context.Context, in *FavoriteActionRequest, opts ...grpc.CallOption) (*FavoriteActionResponse, error) {
	out := new(FavoriteActionResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/FavoriteAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) FavoriteList(ctx context.Context, in *FavoriteListRequest, opts ...grpc.CallOption) (*FavoriteListResponse, error) {
	out := new(FavoriteListResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/FavoriteList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) CommentAction(ctx context.Context, in *CommentActionRequest, opts ...grpc.CallOption) (*CommentActionResponse, error) {
	out := new(CommentActionResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/CommentAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) CommentList(ctx context.Context, in *CommentListRequest, opts ...grpc.CallOption) (*CommentListResponse, error) {
	out := new(CommentListResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/CommentList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) UploadVideo(ctx context.Context, in *UploadVideoRequest, opts ...grpc.CallOption) (*UploadVideoResponse, error) {
	out := new(UploadVideoResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/UploadVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) VideoList(ctx context.Context, in *VideoListRequest, opts ...grpc.CallOption) (*VideoListResponse, error) {
	out := new(VideoListResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/VideoList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) Feed(ctx context.Context, in *FeedRequest, opts ...grpc.CallOption) (*FeedResponse, error) {
	out := new(FeedResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/Feed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) GetTotalFavoriteCount(ctx context.Context, in *GetTotalFavoriteCountRequest, opts ...grpc.CallOption) (*GetTotalFavoriteCountResponse, error) {
	out := new(GetTotalFavoriteCountResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/GetTotalFavoriteCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) GetTotalVideoCount(ctx context.Context, in *GetTotalVideoCountRequest, opts ...grpc.CallOption) (*GetTotalVideoCountResponse, error) {
	out := new(GetTotalVideoCountResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/GetTotalVideoCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoServiceClient) GetTotalLikeCount(ctx context.Context, in *GetTotalLikeCountRequest, opts ...grpc.CallOption) (*GetTotalLikeCountResponse, error) {
	out := new(GetTotalLikeCountResponse)
	err := c.cc.Invoke(ctx, "/video.VideoService/GetTotalLikeCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VideoServiceServer is the server API for VideoService service.
// All implementations must embed UnimplementedVideoServiceServer
// for forward compatibility
type VideoServiceServer interface {
	// 点赞方法
	FavoriteAction(context.Context, *FavoriteActionRequest) (*FavoriteActionResponse, error)
	// 喜欢列表方法
	FavoriteList(context.Context, *FavoriteListRequest) (*FavoriteListResponse, error)
	// 视频评论方法
	CommentAction(context.Context, *CommentActionRequest) (*CommentActionResponse, error)
	// 评论列表方法
	CommentList(context.Context, *CommentListRequest) (*CommentListResponse, error)
	// 上传视频
	UploadVideo(context.Context, *UploadVideoRequest) (*UploadVideoResponse, error)
	// 视频列表
	VideoList(context.Context, *VideoListRequest) (*VideoListResponse, error)
	// 视频流
	Feed(context.Context, *FeedRequest) (*FeedResponse, error)
	// 获取用户总点赞数
	GetTotalFavoriteCount(context.Context, *GetTotalFavoriteCountRequest) (*GetTotalFavoriteCountResponse, error)
	// 获取用户总视频数
	GetTotalVideoCount(context.Context, *GetTotalVideoCountRequest) (*GetTotalVideoCountResponse, error)
	// 获取用户总喜欢数
	GetTotalLikeCount(context.Context, *GetTotalLikeCountRequest) (*GetTotalLikeCountResponse, error)
	mustEmbedUnimplementedVideoServiceServer()
}

// UnimplementedVideoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVideoServiceServer struct {
}

func (UnimplementedVideoServiceServer) FavoriteAction(context.Context, *FavoriteActionRequest) (*FavoriteActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FavoriteAction not implemented")
}
func (UnimplementedVideoServiceServer) FavoriteList(context.Context, *FavoriteListRequest) (*FavoriteListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FavoriteList not implemented")
}
func (UnimplementedVideoServiceServer) CommentAction(context.Context, *CommentActionRequest) (*CommentActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommentAction not implemented")
}
func (UnimplementedVideoServiceServer) CommentList(context.Context, *CommentListRequest) (*CommentListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommentList not implemented")
}
func (UnimplementedVideoServiceServer) UploadVideo(context.Context, *UploadVideoRequest) (*UploadVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadVideo not implemented")
}
func (UnimplementedVideoServiceServer) VideoList(context.Context, *VideoListRequest) (*VideoListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VideoList not implemented")
}
func (UnimplementedVideoServiceServer) Feed(context.Context, *FeedRequest) (*FeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Feed not implemented")
}
func (UnimplementedVideoServiceServer) GetTotalFavoriteCount(context.Context, *GetTotalFavoriteCountRequest) (*GetTotalFavoriteCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTotalFavoriteCount not implemented")
}
func (UnimplementedVideoServiceServer) GetTotalVideoCount(context.Context, *GetTotalVideoCountRequest) (*GetTotalVideoCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTotalVideoCount not implemented")
}
func (UnimplementedVideoServiceServer) GetTotalLikeCount(context.Context, *GetTotalLikeCountRequest) (*GetTotalLikeCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTotalLikeCount not implemented")
}
func (UnimplementedVideoServiceServer) mustEmbedUnimplementedVideoServiceServer() {}

// UnsafeVideoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VideoServiceServer will
// result in compilation errors.
type UnsafeVideoServiceServer interface {
	mustEmbedUnimplementedVideoServiceServer()
}

func RegisterVideoServiceServer(s grpc.ServiceRegistrar, srv VideoServiceServer) {
	s.RegisterService(&VideoService_ServiceDesc, srv)
}

func _VideoService_FavoriteAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FavoriteActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).FavoriteAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/FavoriteAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).FavoriteAction(ctx, req.(*FavoriteActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_FavoriteList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FavoriteListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).FavoriteList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/FavoriteList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).FavoriteList(ctx, req.(*FavoriteListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_CommentAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommentActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).CommentAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/CommentAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).CommentAction(ctx, req.(*CommentActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_CommentList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommentListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).CommentList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/CommentList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).CommentList(ctx, req.(*CommentListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_UploadVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).UploadVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/UploadVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).UploadVideo(ctx, req.(*UploadVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_VideoList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VideoListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).VideoList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/VideoList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).VideoList(ctx, req.(*VideoListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_Feed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).Feed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/Feed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).Feed(ctx, req.(*FeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_GetTotalFavoriteCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTotalFavoriteCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).GetTotalFavoriteCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/GetTotalFavoriteCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).GetTotalFavoriteCount(ctx, req.(*GetTotalFavoriteCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_GetTotalVideoCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTotalVideoCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).GetTotalVideoCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/GetTotalVideoCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).GetTotalVideoCount(ctx, req.(*GetTotalVideoCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VideoService_GetTotalLikeCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTotalLikeCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServiceServer).GetTotalLikeCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/video.VideoService/GetTotalLikeCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServiceServer).GetTotalLikeCount(ctx, req.(*GetTotalLikeCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VideoService_ServiceDesc is the grpc.ServiceDesc for VideoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VideoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "video.VideoService",
	HandlerType: (*VideoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FavoriteAction",
			Handler:    _VideoService_FavoriteAction_Handler,
		},
		{
			MethodName: "FavoriteList",
			Handler:    _VideoService_FavoriteList_Handler,
		},
		{
			MethodName: "CommentAction",
			Handler:    _VideoService_CommentAction_Handler,
		},
		{
			MethodName: "CommentList",
			Handler:    _VideoService_CommentList_Handler,
		},
		{
			MethodName: "UploadVideo",
			Handler:    _VideoService_UploadVideo_Handler,
		},
		{
			MethodName: "VideoList",
			Handler:    _VideoService_VideoList_Handler,
		},
		{
			MethodName: "Feed",
			Handler:    _VideoService_Feed_Handler,
		},
		{
			MethodName: "GetTotalFavoriteCount",
			Handler:    _VideoService_GetTotalFavoriteCount_Handler,
		},
		{
			MethodName: "GetTotalVideoCount",
			Handler:    _VideoService_GetTotalVideoCount_Handler,
		},
		{
			MethodName: "GetTotalLikeCount",
			Handler:    _VideoService_GetTotalLikeCount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "video/video.proto",
}
