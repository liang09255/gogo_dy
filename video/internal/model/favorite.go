package model

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	ID      int64 `gorm:"primarykey" json:"id"`
	UserId  int64 `gorm:"not null,uniqueIndex:delete_at" json:"user_id"`
	VideoId int64 `gorm:"not null,uniqueIndex:delete_at" json:"video_id"`
}
