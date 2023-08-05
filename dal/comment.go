package dal

import (
	"context"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID      int64  `gorm:"primarykey" json:"id"`
	UserId  int64  `gorm:"not null" json:"user_id"`
	VideoId int64  `gorm:"not null" json:"video_id"`
	Content string `gorm:"not null" json:"content"`
}

type commentDal struct{}

var CommentDal = &commentDal{}

func (c *commentDal) PostCommentAction(ctx context.Context, comment *Comment) error {
	if err := DB.WithContext(ctx).Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

func (c *commentDal) GetCommentList(ctx context.Context, userId int64, videoId int64) ([]*Comment, error) {
	res := make([]*Comment, 0)
	if err := DB.WithContext(ctx).Where("video_id = ? and user_id = ?", videoId, userId).Order("created_at desc").Limit(30).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
