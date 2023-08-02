package service

import (
	"context"
	"main/dal"
)

type favoriteService struct{}

var FavoriteService = &favoriteService{}

func (h *favoriteService) GetFavoriteList(ctx context.Context, userId int64, videoId int64) error {
	dal.GetFavoriteList(ctx, userId, videoId)
	//if err != nil {
	//	return nil, err
	//}
	return nil
}
func (h *favoriteService) PostFavoriteAction(ctx context.Context, userId int64, videoId int64) error {
	dal.PostFavoriteAction(ctx, userId, videoId)
	//if err != nil {
	//	return nil, err
	//}
	return nil
}
