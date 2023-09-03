package model

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID         int64  `gorm:"primarykey" json:"message_id"`
	ToUserID   int64  `gorm:"type:bigint;column:to_user_id;index:delete_at" json:"to_user_id"`
	FromUserID int64  `gorm:"type:bigint;column:from_user_id;index:delete_at" json:"from_user_id"` // 存储发送者的ID,可以从鉴权token中得到
	ActionType int32  `gorm:"type:int;column:action_type" json:"action_type"`
	Content    string `gorm:"type:text;column:content" json:"content"`
}
