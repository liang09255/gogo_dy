package model

// User 用户表
type User struct {
	Id            int64  `gorm:"primaryKey"`
	Name          string `gorm:"not null;unique;index"`
	Username      string `gorm:"not null;unique;index"`
	Password      string `gorm:"not null"`
	FollowCount   int64
	FollowerCount int64
}

// Video 视频
type Video struct {
	Id            int64 `gorm:"primaryKey"`
	User          User  `gorm:"foreignKey:PublishId"`
	PublishId     int64
	PlayUrl       string `gorm:"not null"`
	CoverUrl      string `gorm:"not null"`
	FavoriteCount int64
	CommentCount  int64
	Title         string `gorm:"not null"`
	CreateTime    int64  `gorm:"autoCreateTime"`
}

// Comment 评论
type Comment struct {
	Id         int64 `gorm:"primaryKey"`
	User       User  `gorm:"foreignKey:UserId"`
	UserId     int64
	Video      Video `gorm:"foreignKey:VideoId"`
	VideoId    int64
	Content    string `gorm:"not null"`
	CreateTime int64  `gorm:"autoCreateTime"`
}

// Follow 关注
type Follow struct {
	Id           int64 `gorm:"primaryKey"`
	User         User  `gorm:"foreignKey:FollowId"` // 关注人
	FollowId     int64
	FollowedUser User `gorm:"foreignKey:UserId"` // 被关注人
	UserId       int64
	CreateTime   int64 `gorm:"autoCreateTime"`
}

// VideoFavorite 视频点赞
type VideoFavorite struct {
	Id         int64 `gorm:"primaryKey"`
	User       User  `gorm:"foreignKey:UserId"`
	UserId     int64
	Video      Video `gorm:"foreignKey:VideoId"`
	VideoId    int64
	CreateTime int64 `gorm:"autoCreateTime"`
}
