package service

import (
	"context"
	"main/controller/ctlModel/userCtlModel"
	"main/controller/ctlModel/videoCtlModel"
	"main/dal"
)

type favoriteService struct{}

var FavoriteService = &favoriteService{}

// FixMe 缺少接口类型
func (h *favoriteService) GetFavoriteList(ctx context.Context, userId int64) (videos []videoCtlModel.Video, err error) {
	var videoIds []int64
	videoIds, err = dal.FavoriteDal.GetFavoriteList(ctx, userId)
	if err != nil {
		return nil, err
	}
	// FixMe 获取Video和User信息
	for _, videoId := range videoIds {
		var video = videoCtlModel.Video{
			ID:            videoId,
			Author:        userCtlModel.User{},
			PlayUrl:       "",
			CoverUrl:      "",
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    true,
			Title:         "undefined",
		}
		videos = append(videos, video)
	}
	return
}
func (h *favoriteService) PostFavoriteAction(ctx context.Context, userId int64, videoId int64, actionType int32) (string, error) {
	var err error
	msg := "没有指定的类型"
	if actionType == 1 {
		err = dal.FavoriteDal.PostFavoriteAction(ctx, userId, videoId)
		msg = "点赞成功"
	} else if actionType == 2 {
		err = dal.FavoriteDal.CancelFavoriteAction(ctx, userId, videoId)
		msg = "取消点赞成功"
	}
	if err != nil {
		return "", err
	}
	return msg, nil
}
