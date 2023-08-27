package model

import "gorm.io/gorm"

type Relation struct {
	gorm.Model
	Id       int64 `gorm:"primaryKey"`
	FollowId int64 // 关注者
	UserId   int64 // 被关注者
}
