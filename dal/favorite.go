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

type favoriteDal struct{}

var FavoriteDal = &favoriteDal{}

func (f *favoriteDal) GetFavoriteList(ctx context.Context, userId int64) ([]int64, error) {
	ids := make([]int64, 0)
	if err := DB.WithContext(ctx).Where("id in (?)",
		DB.WithContext(ctx).Table("favorites").Select("video_id").Where("user_id = ? AND deleted_at is null", userId).Find(&ids)).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (f *favoriteDal) PostFavoriteAction(ctx context.Context, userId int64, videoId int64) error {
	favorite := Favorite{
		UserId:  userId,
		VideoId: videoId,
	}
	var cnt int64 = 0
	if err := DB.WithContext(ctx).Model(&Favorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt != 0 {
		return nil
	}
	if err := DB.WithContext(ctx).Create(&favorite).Error; err != nil {
		return err
	}
	return nil
}
func (c *favoriteDal) CancelFavoriteAction(ctx context.Context, userId int64, videoId int64) error {
	var cnt int64 = 0
	if err := DB.WithContext(ctx).Model(&Favorite{}).Where("user_id = ? and video_id = ?", userId, videoId).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 {
		return nil
	}

	if err := DB.WithContext(ctx).Where("user_id = ? and video_id = ?", userId, videoId).Delete(&Favorite{}).Error; err != nil {
		return err
	}
	return nil
}
