package service

import (
	"context"
	"main/dal"
)

type favoriteService struct{}

var FavoriteService = &favoriteService{}

// FixMe 缺少接口类型
func (h *favoriteService) GetFavoriteList(ctx context.Context, userId int64, videoId int64) (interface{}, error) {
	videoIds, err := dal.FavoriteDal.GetFavoriteList(ctx, userId)
	// FixMe 缺少对Video表的查询
	if err != nil {
		return nil, err
	}
	return videoIds, nil
}
func (h *favoriteService) PostFavoriteAction(ctx context.Context, userId int64, videoId int64, actionType int64) (string, error) {
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
