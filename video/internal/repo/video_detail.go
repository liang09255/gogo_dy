package repo

import (
	"context"
	"video/internal/model"
)

type VideoDetailRepo interface {
	AddFavoriteCount(ctx context.Context, vid int64, count int64) error
	SubFavoriteCount(ctx context.Context, vid int64, count int64) error
	AddCommentCount(ctx context.Context, vid int64, count int64) error
	SubCommentCount(ctx context.Context, vid int64, count int64) error
	BatchInsertFavorite(ctx context.Context, m map[int64]int) (faultCount int, err error)
	BatchInsertComment(ctx context.Context, m map[int64]int) (faultCount int, err error)
	GetVideoDetail(ctx context.Context, vid int64) (*model.VideoDetail, error)
}
