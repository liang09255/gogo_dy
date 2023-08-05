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
func (h *favoriteService) PostFavoriteAction(ctx context.Context, userId int64, videoId int64) error {
	err := dal.FavoriteDal.PostFavoriteAction(ctx, userId, videoId)
	if err != nil {
		return err
	}
	return nil
}

func (h *favoriteService) PostCancelAction(ctx context.Context, userId int64, videoId int64) error {
	err := dal.FavoriteDal.CancelFavoriteAction(ctx, userId, videoId)
	if err != nil {
		return err
	}
	return nil
}
