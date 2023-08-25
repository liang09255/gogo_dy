package repo

import (
	"context"
	"video/internal/model"
)

type VideoRepo interface {
	// 增加点赞数
	AddFavoriteCount(ctx context.Context, id int64, count int64) error
	// 减少点赞数
	CancelFavoriteCount(ctx context.Context, id int64, count int64) error
	// 增加评论数
	AddCommentCount(ctx context.Context, id int64, count int64) error
	// 删除评论数
	ReduceCommentCount(ctx context.Context, id int64, count int64) error
	// 批量获得视频信息
	MGetVideoInfo(ctx context.Context, ids []int64) ([]model.Video, error)
	GetVideoInfo(ctx context.Context, videoId int64) (model.Video, error)
}
