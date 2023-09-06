package service

import (
	"common/ggIDL/user"
	"common/ggIDL/video"
	"common/ggLog"
	"common/ggRPC"
	"context"
	"github.com/sourcegraph/conc"
	"video/internal/domain"
)

type VideoService struct {
	video.UnsafeVideoServiceServer
	videoDomain    *domain.VideoDomain
	commentDomain  *domain.CommentDomain
	favoriteDomain *domain.FavoriteDomain
	hotDomain      *domain.HotDomain
	userClient     user.UserClient
}

var _ video.VideoServiceServer = (*VideoService)(nil)

func New() *VideoService {
	return &VideoService{
		videoDomain:    domain.NewVideoDomain(),
		favoriteDomain: domain.NewFavoriteDomain(),
		commentDomain:  domain.NewCommentDomain(),
		userClient:     ggRPC.GetUserClient(),
	}
}

func (v VideoService) UploadVideo(ctx context.Context, request *video.UploadVideoRequest) (*video.UploadVideoResponse, error) {
	return &video.UploadVideoResponse{}, v.videoDomain.NewVideo(ctx, request.UserId, request.PlayUrl, request.CoverUrl, request.Title)
}

func (v VideoService) VideoList(ctx context.Context, request *video.VideoListRequest) (*video.VideoListResponse, error) {
	videos, err := v.videoDomain.PublishList(ctx, request.UserId)
	if err != nil {
		return &video.VideoListResponse{}, err
	}
	vids := make([]int64, len(videos))
	for _, v := range videos {
		vids = append(vids, v.Id)
	}
	// 获得视频的点赞数
	favoriteCountMap := v.favoriteDomain.GetFavoriteCountByVideoID(ctx, vids)
	for _, v := range videos {
		v.FavoriteCount = favoriteCountMap[v.Id]
	}
	resp := &video.VideoListResponse{
		VideoList: videos,
	}
	return resp, err
}

func (v VideoService) Feed(ctx context.Context, request *video.FeedRequest) (*video.FeedResponse, error) {
	videos, nextTime, err := v.videoDomain.Feed(ctx, request.Latest)
	if err != nil {
		return &video.FeedResponse{}, err
	}
	if len(videos) == 0 {
		ggLog.Info("没有更多视频了")
		return &video.FeedResponse{NextTime: nextTime}, nil
	}
	// 获取视频的点赞数
	vids := make([]int64, len(videos))
	for _, v := range videos {
		vids = append(vids, v.Id)
	}

	wg := conc.NewWaitGroup()
	errChan := make(chan error)
	downChan := make(chan struct{})
	favoriteCountMap := make(map[int64]int64)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg.Go(func() {
		// 查询视频的点赞数
		favoriteCountMap = v.favoriteDomain.GetFavoriteCountByVideoID(ctx, vids)
	})
	favoriteMap := make(map[int64]bool)
	wg.Go(func() {
		// 查询是否已经点赞
		favoriteMap = v.favoriteDomain.CheckFavorite(ctx, request.UserId, vids)
	})
	commentCountMap := make(map[int64]int64)
	wg.Go(func() {
		// 获取视频的评论数
		commentCountMap = v.commentDomain.GetCommentCountByVideoID(ctx, vids)
	})
	userInfoMap := make(map[int64]*user.UserInfoModel)
	wg.Go(func() {
		// 获得视频作者信息
		msg := &user.UserInfoRequest{
			MyId:   request.UserId,
			UserId: getUserIDFromVideoList(videos),
		}
		infos, err := v.userClient.MGetUserInfo(ctx, msg)
		if err != nil {
			errChan <- err
		}

		for _, v := range infos.UserInfo {
			userInfoMap[v.Id] = v
		}
	})
	go func() {
		wg.Wait()
		downChan <- struct{}{}
	}()
	select {
	case err = <-errChan:
		ggLog.Error("获取视频信息失败:", err)
		return nil, err
	case <-downChan:
	}

	// 加进返回值
	for _, v := range videos {
		v.Author = userInfoMap[v.Author.Id]
		v.CommentCount = commentCountMap[v.Id]
		v.IsFavorite = favoriteMap[v.Id]
		v.FavoriteCount = favoriteCountMap[v.Id]
	}
	resp := &video.FeedResponse{
		VideoList: videos,
		NextTime:  nextTime,
	}
	// 统计热点
	go func() {
		v.hotDomain.HotStatistics(request.UserId, vids)
	}()
	return resp, err
}

// FavoriteAction 点赞/取消点赞
func (v VideoService) FavoriteAction(ctx context.Context, request *video.FavoriteActionRequest) (*video.FavoriteActionResponse, error) {
	// 点赞操作
	userid := request.UserId
	videoid := request.VideoId
	// 获得被点赞的视频作者的id

	err := v.favoriteDomain.FavoriteAction(ctx, userid, videoid, request.ActionType)
	return &video.FavoriteActionResponse{}, err
}

func (v VideoService) FavoriteList(ctx context.Context, request *video.FavoriteListRequest) (*video.FavoriteListResponse, error) {
	// 获得favoriteList
	favoriteListResp, err := v.favoriteDomain.FavoriteList(ctx, request.UserId)
	if err != nil {
		return &video.FavoriteListResponse{}, err
	}
	// 如果喜爱列表为0，不需要查找，以免后面调用其余接口报错
	if len(favoriteListResp.VideoList) == 0 {
		return favoriteListResp, nil
	}

	// 获得视频作者信息
	req := &user.UserInfoRequest{
		UserId: getUserIDFromVideoList(favoriteListResp.VideoList),
		MyId:   request.UserId,
	}
	userResp, err := v.userClient.MGetUserInfo(ctx, req)
	if err != nil {
		ggLog.Error("调用user服务批量获得用户信息失败:", err)
		return nil, err
	}

	// 映射ID-User，因为查询返回的不一定是原来的顺序
	authorMap := make(map[int64]*user.UserInfoModel)
	for _, value := range userResp.UserInfo {
		authorMap[value.Id] = value
	}

	// 再遍历一遍赋值
	for _, value := range favoriteListResp.VideoList {
		value.Author = authorMap[value.Author.Id]
	}

	return favoriteListResp, err
}

func (v VideoService) CommentAction(ctx context.Context, request *video.CommentActionRequest) (*video.CommentActionResponse, error) {
	// 修改评论表和视频表
	comment, err := v.commentDomain.CommentAction(ctx, request)
	if err != nil {
		return &video.CommentActionResponse{}, err
	}

	// 如果为删除操作则结束了
	if request.ActionType == video.ActionType_Cancel {
		return &video.CommentActionResponse{}, err
	}

	// 从用户表获得评论用户信息
	req := &user.UserInfoRequest{
		UserId: []int64{comment.User.Id},
		MyId:   request.UserId,
	}
	userResp, err := v.userClient.MGetUserInfo(ctx, req)
	if err != nil {
		ggLog.Error("调用user服务批量获得用户信息失败:", err)
		return nil, err
	}
	comment.User = userResp.UserInfo[0]
	// 转为返回结构
	resp := &video.CommentActionResponse{
		Comment: comment,
	}
	return resp, err
}

func (v VideoService) CommentList(ctx context.Context, request *video.CommentListRequest) (*video.CommentListResponse, error) {
	// 获得评论列表
	commentList, err := v.commentDomain.CommentList(ctx, request.VideoId)
	if err != nil {
		return &video.CommentListResponse{}, err
	}
	// 如果评论数为0，则不需要查找了
	if len(commentList) == 0 {
		return &video.CommentListResponse{CommentList: commentList}, nil
	}
	// 获得评论用户信息
	req := &user.UserInfoRequest{
		UserId: getUserIDFromCommentList(commentList),
	}
	userResp, err := v.userClient.MGetUserInfo(ctx, req)
	if err != nil {
		ggLog.Error("调用user服务批量获得用户信息失败:", err)
		return nil, err
	}
	authorMap := make(map[int64]*user.UserInfoModel)
	for _, value := range userResp.UserInfo {
		authorMap[value.Id] = value
	}
	for _, value := range commentList {
		value.User = authorMap[value.User.Id]
	}
	// 转为返回结构
	resp := &video.CommentListResponse{
		CommentList: commentList,
	}
	return resp, err
}

// GetTotalFavoriteCount 获取总获得的点赞数
func (v VideoService) GetTotalFavoriteCount(ctx context.Context, request *video.GetTotalFavoriteCountRequest) (*video.GetTotalFavoriteCountResponse, error) {
	return &video.GetTotalFavoriteCountResponse{Count: v.favoriteDomain.GetFavoriteCount(ctx, request.UserId)}, nil
}

// GetTotalVideoCount 获取总视频数
func (v VideoService) GetTotalVideoCount(ctx context.Context, request *video.GetTotalVideoCountRequest) (*video.GetTotalVideoCountResponse, error) {
	return &video.GetTotalVideoCountResponse{Count: v.videoDomain.GetTotalVideoCount(ctx, request.UserId)}, nil
}

// GetTotalLikeCount 获取总喜欢的视频数
func (v VideoService) GetTotalLikeCount(ctx context.Context, request *video.GetTotalLikeCountRequest) (*video.GetTotalLikeCountResponse, error) {
	return &video.GetTotalLikeCountResponse{Count: v.favoriteDomain.GetLikeCount(ctx, request.UserId)}, nil
}
