package dal

import (
	"api/global"
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

func (c *commentDal) PostComment(ctx context.Context, comment *Comment) error {
	return global.MysqlDB.WithContext(ctx).Create(&comment).Error
}

func (c *commentDal) DelComment(ctx context.Context, comment *Comment) error {
	return global.MysqlDB.WithContext(ctx).Delete(&comment).Error
}

func (c *commentDal) GetCommentList(ctx context.Context, videoId int64) ([]*Comment, error) {
	res := make([]*Comment, 0)
	if err := global.MysqlDB.WithContext(ctx).Where("video_id = ?", videoId).Order("created_at desc").Limit(30).Find(&res).Error; err != nil {
		return []*Comment{}, err
	}
	return res, nil
}
