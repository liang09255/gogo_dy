package model

import "gorm.io/gorm"

type VideoDetail struct {
	gorm.Model
	Id            int64 `gorm:"primaryKey"` // 映射vid
	FavoriteCount int64 `gorm:"favorite_count"`
	CommentCount  int64 `gorm:"comment_count"`
}
