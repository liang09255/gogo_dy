package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID             int64  `gorm:"primarykey" json:"user_id"`
	Username       string `gorm:"not null" json:"user_name"`
	Password       string `gorm:"not null" json:"password"`
	FollowCount    int64  `gorm:"default:0" json:"follow_count"`
	FollowerCount  int64  `gorm:"default:0" json:"follower_count"`
	TotalFavorited int64  `gorm:"default:0" json:"total_favorited"` // 获点赞数量
	WorkCount      int64  `gorm:"default:0" json:"work_count"`      // 作品数
	FavoriteCount  int64  `gorm:"default:0" json:"favorite_count"`  //喜欢数
}
