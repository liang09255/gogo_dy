package repo

import (
	"context"
	"video/internal/model"
)

type VideoRepo interface {
	MGetVideoInfo(ctx context.Context, ids []int64) ([]*model.Video, error) // 批量获得视频信息
	GetVideoInfo(ctx context.Context, videoId int64) (*model.Video, error)
	NewVideo(ctx context.Context, userID int64, videoURL, coverURL, title string) error
	PublishList(ctx context.Context, userID int64) ([]*model.Video, error)
	Feed(ctx context.Context, latestTime int64) ([]*model.Video, error)
}
