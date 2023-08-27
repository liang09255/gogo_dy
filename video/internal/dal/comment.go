package dal

import (
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

// 增加评论
func (cd *CommentDal) AddComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	t := cd.conn.WithContext(ctx).Model(&model.Comment{}).Create(&comment)
	return comment, t.Error
}

// 删除评论
func (cd *CommentDal) DeleteComment(ctx context.Context, comment *model.Comment) error {
	return cd.conn.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", comment.ID).Delete(&comment).Error
}

// 评论列表
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
