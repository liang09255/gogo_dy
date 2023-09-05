package model

import "gorm.io/gorm"

// Video 视频
type Video struct {
	gorm.Model
	Id       int64  `gorm:"primaryKey"`
	AuthorId int64  `gorm:"not null"`
	PlayUrl  string `gorm:"not null"`
	CoverUrl string `gorm:"not null"`
	//FavoriteCount int64
	//CommentCount  int64
	Title string `gorm:"not null"`
}
