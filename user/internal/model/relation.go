package model

import (
	"gorm.io/gorm"
	"time"
)

type Relation struct {
	gorm.Model
	Id       int64 `gorm:"primaryKey"`
	FollowId int64 `gorm:"uniqueIndex:delete_at"` // 关注者
	UserId   int64 `gorm:"uniqueIndex:delete_at"` // 被关注者
}

type RelationActionMessage struct {
	UserID   int64 `json:"userID"`
	TargetID int64 `json:"targetID"`
	Action   int32 `json:"action"`
}

type FollowStats struct {
	UserID         int64     `gorm:"primaryKey"`
	FollowerCount  int64     `gorm:"column:follower_count"`
	FollowingCount int64     `gorm:"column:following_count"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}
