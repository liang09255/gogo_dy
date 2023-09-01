package dal

import (
	"common/ggLog"
	"context"
	"errors"
	"gorm.io/gorm"
	"video/internal/database"
	"video/internal/model"
	"video/internal/repo"
)

type CommentDal struct {
	conn *database.GormConn
}

var _ repo.CommentRepo = (*CommentDal)(nil)

func NewCommentDao() *CommentDal {
	return &CommentDal{
		conn: database.New(),
	}
}

// AddComment 增加评论
func (cd *CommentDal) AddComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	t := cd.conn.WithContext(ctx).Model(&model.Comment{}).Create(&comment)
	return comment, t.Error
}

// DeleteComment 删除评论
func (cd *CommentDal) DeleteComment(ctx context.Context, comment *model.Comment) error {
	return cd.conn.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", comment.ID).Delete(&comment).Error
}

// GetCommentList 评论列表
func (cd *CommentDal) GetCommentList(ctx context.Context, videoId int64) ([]model.Comment, error) {
	res := make([]model.Comment, 0)
	err := cd.conn.WithContext(ctx).Where("video_id = ?", videoId).Order("created_at desc").
		Limit(30).
		Find(&res).Error

	// 如果查无记录，不会出错
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return res, nil
	}

	return res, err
}

// MGetCommentCount 批量获取评论数
func (cd *CommentDal) MGetCommentCount(ctx context.Context, videoId []int64) map[int64]int64 {
	res := make(map[int64]int64)
	queryRes := make([]struct {
		VideoId int64
		Count   int64
	}, 0)
	err := cd.conn.WithContext(ctx).Model(&model.Comment{}).Select("video_id, count(*) as count").
		Where("video_id in ?", videoId).
		Group("video_id").
		Scan(&queryRes).Error
	if err != nil {
		ggLog.Errorf("获取评论数错误:%v", err)
		return nil
	}
	for _, value := range queryRes {
		res[value.VideoId] = value.Count
	}
	return res
}
