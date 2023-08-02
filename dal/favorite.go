package dal

import (
	"context"
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	ID      int64 `gorm:"primarykey" json:"id"`
	UserId  int64 `gorm:"not null" json:"user_id"`
	VideoId int64 `gorm:"not null" json:"video_id"`
}

func GetFavoriteList(ctx context.Context, userId int64, videoId int64) error {
	//  FixMe 依赖于 Video 表 暂未实现获取点赞列表功能
	return nil
}

func PostFavoriteAction(ctx context.Context, userId int64, videoId int64) error {

	var cnt int64 = 0
	if err := DB.WithContext(ctx).Model(&Favorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Count(&cnt).Error; err != nil {
		return err
	}

	if cnt != 0 {
		return nil
	}

	// FixMe   依赖于Video 表 暂未实现点赞功能
	return nil
}
