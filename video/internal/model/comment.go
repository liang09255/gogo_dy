package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ID      int64  `gorm:"primarykey" json:"id"`
	UserId  int64  `gorm:"not null" json:"user_id"`
	VideoId int64  `gorm:"not null" json:"video_id"`
	Content string `gorm:"not null" json:"content"`
}
