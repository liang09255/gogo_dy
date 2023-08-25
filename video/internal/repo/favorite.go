package repo

import "context"

type FavoriteRepo interface {
	GetFavoriteList(ctx context.Context, userId int64) ([]int64, error)
	PostFavoriteAction(ctx context.Context, userId int64, videoId int64) error
	CancelFavoriteAction(ctx context.Context, userId int64, videoId int64) error
}
