package service

import (
	"common/ggIDL/user"
	"common/ggIDL/video"
	"common/ggLog"
	"common/ggRPC"
	"context"
	"video/internal/domain"
)

type VideoService struct {
	video.UnsafeVideoServiceServer
	videoDomain    *domain.VideoDomain
	commentDomain  *domain.CommentDomain
	favoriteDomain *domain.FavoriteDomain
	userClient     user.UserClient
}

// 点赞/取消点赞
func (v VideoService) FavoriteAction(ctx context.Context, request *video.FavoriteActionRequest) (*video.FavoriteActionResponse, error) {
	// 点赞操作
	userid := request.UserId
	videoid := request.VideoId
	err := v.favoriteDomain.FavoriteAction(ctx, userid, videoid, request.ActionType)
	if err != nil && err.Error() != "重复记录" {
		return &video.FavoriteActionResponse{}, err
	}

	// 获得视频作者
	videoInfo, err := v.videoDomain.GetVideo(ctx, videoid)
	if err != nil {
		return &video.FavoriteActionResponse{}, err
	}

	// TODO 分布式事务，保证跨服务的数据一致
	// TODO 或者使用消息队列，解耦服务之间的依赖，写入消息然后由user服务来消费，解决顺序消费以及消息不丢失的问题即可确保数据的最终一致性
	// TODO 否则如果直接调用的话存在的问题就是，在视频表和点赞表通过事务保证了修改完毕后，如何保证调用user微服务也能成功
	// TODO 使用消息队列的好处是，前两者成功了，如果用户服务崩溃，在重启后也能达到数据一致(一段时间内数据不一致)，但是使用分布式事务的话，就是整体的修改都会丢失(保证数据的强一致)
	// 调用user客户端
	// 构造请求
	r := &user.UserFavoriteActionRequest{
		UserId:   userid,
		Action:   int64(request.ActionType),
		ToUserId: videoInfo.AuthorId,
		Count:    1,
	}
	// 增加获赞数
	_, err = v.userClient.UserFavoriteAction(ctx, r)
	if err != nil {
		// TODO 回滚前一阶段的提交
		ggLog.Error("调用User服务错误:", err)
	}
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

	// 获得视频的用户id列表,用一个map存储，重复id可以减少查询
	userIds := make([]int64, 0)
	for _, value := range favoriteListResp.VideoList {
		userIds = append(userIds, value.Author.Id)
	}

	// 获得视频作者信息
	req := &user.UserInfoRequest{
		UserId: userIds,
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

	// TODO 还需要根据user来判断是否已关注,等关注模块出来了再添加

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
	userInfo := userResp.UserInfo[0]
	comment.User.WorkCount = userInfo.WorkCount
	comment.User.FavoriteCount = userInfo.FavoriteCount
	comment.User.TotalFavorited = userInfo.TotalFavorited
	comment.User.FollowCount = userInfo.FollowerCount
	comment.User.FollowCount = userInfo.FollowCount
	comment.User.Avatar = userInfo.Avatar
	comment.User.Name = userInfo.Name
	comment.User.BackgroundImage = userInfo.BackgroundImage
	comment.User.Signature = userInfo.Signature

	// 转为返回结构
	resp := &video.CommentActionResponse{
		Comment: comment,
	}
	return resp, err
}

func (v VideoService) CommentList(ctx context.Context, request *video.CommentListRequest) (*video.CommentListResponse, error) {
	// 获得评论列表
	commentList, err := v.commentDomain.CommentList(ctx, request.VideoId)

	// 如果评论数为0，则不需要查找了
	if len(commentList) == 0 {
		return &video.CommentListResponse{CommentList: commentList}, nil
	}
	// 根据评论列表获得用户信息
	userIds := make([]int64, 0)
	for _, value := range commentList {
		userIds = append(userIds, value.User.Id)
	}
	req := &user.UserInfoRequest{
		UserId: userIds,
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

	resp := &video.CommentListResponse{
		CommentList: commentList,
	}
	return resp, err
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
