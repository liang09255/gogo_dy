package repo

import "context"

type FavoriteRepo interface {
	GetFavoriteList(ctx context.Context, userId int64) ([]int64, error)             // 获得点赞列表
	PostFavoriteAction(ctx context.Context, userId int64, videoId int64) error      // 点赞
	CancelFavoriteAction(ctx context.Context, userId int64, videoId int64) error    // 取消点赞
	GetFavoriteCountByVideoID(ctx context.Context, videoID []int64) map[int64]int64 // 获取视频的点赞数
	GetFavoriteCount(ctx context.Context, userID int64) int64                       // 获取用户的点赞数
}
