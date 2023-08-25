package dal

import (
	"context"
	"gorm.io/gorm"
	"video/internal/database"
	"video/internal/model"
	"video/internal/repo"
)

type VideoDal struct {
	conn *database.GormConn
}

var _ repo.VideoRepo = (*VideoDal)(nil)

func NewVideoDao() *VideoDal {
	return &VideoDal{
		conn: database.New(),
	}
}

// PostFavoriteCount 增加点赞数
func (v *VideoDal) AddFavoriteCount(ctx context.Context, id int64, count int64) error {
	return v.conn.WithContext(ctx).Model(&model.Video{}).
		Where("id = ?", id).
		Update("favorite_count", gorm.Expr("favorite_count + ?", count)).Error
}

// CancelFavoriteCount 取消点赞数
func (v *VideoDal) CancelFavoriteCount(ctx context.Context, id int64, count int64) error {
	return v.conn.WithContext(ctx).Model(&model.Video{}).
		Where("id = ? and favorite_count > ?", id, count).
		Update("favorite_count", gorm.Expr("favorite_count - ?", count)).Error
}

// AddCommentCount 增加评论数
func (v *VideoDal) AddCommentCount(ctx context.Context, id int64, count int64) error {
	return v.conn.WithContext(ctx).Model(&model.Video{}).
		Where("id = ?", id).
		Update("comment_count", gorm.Expr("comment_count + ?", count)).Error
}

// ReduceCommentCount 减少评论数
func (v *VideoDal) ReduceCommentCount(ctx context.Context, id int64, count int64) error {
	return v.conn.WithContext(ctx).Model(&model.Video{}).
		Where("id = ? and comment_count > ?", id, count).
		Update("comment_count", gorm.Expr("comment_count - ?", count)).Error
}

// MGetVideoInfo 批量获得视频信息
func (v *VideoDal) MGetVideoInfo(ctx context.Context, ids []int64) ([]model.Video, error) {
	var videoList []model.Video
	t := v.conn.WithContext(ctx).Model(&model.Video{}).
		Where("id in ?", ids).
		Find(&videoList)

	return videoList, t.Error
}
