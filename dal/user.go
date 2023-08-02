package dal

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID            int64  `gorm:"primarykey" json:"user_id"`
	UserName      string `gorm:"not null" json:"user_name"`      
	Password      string `gorm:"not null" json:"password"`       
	FollowCount   int64  `gorm:"default:0" json:"follow_count"` 
	FollowerCount int64  `gorm:"default:0" json:"follower_count"` 
}
