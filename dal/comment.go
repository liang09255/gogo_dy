package dal

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"main/global"
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

const (
	PostComment = 1
	DelComment  = 2
)

func (c *commentDal) PostCommentAction(ctx context.Context, comment *Comment, actionType int32) error {
	// FIXME 删除功能有问题
	if actionType == PostComment {
		return global.MysqlDB.WithContext(ctx).Create(&comment).Error
	} else if actionType == DelComment {
		return global.MysqlDB.WithContext(ctx).Where("id = ?", comment.ID).Delete(&Comment{}).Error
	}
	return fmt.Errorf("invalid action type for comment")
}

func (c *commentDal) GetCommentList(ctx context.Context, videoId int64) ([]*Comment, error) {
	res := make([]*Comment, 0)
	if err := global.MysqlDB.WithContext(ctx).Where("video_id = ? and deleted_at is null", videoId).Order("created_at desc").Limit(30).Find(&res).Error; err != nil {
		return []*Comment{}, err
	}
	return res, nil
}
