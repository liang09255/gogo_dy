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

type Userinfo struct {
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	ID              int64  `json:"id"`               // 用户id
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Name            string `json:"name"`             // 用户名称
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  string `json:"total_favorited"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`       // 作品数
}

// Video 视频
type Video struct {
	Id            int64  `gorm:"primaryKey"`
	AuthorId      int64  `gorm:"not null"`
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

// Relation 数据库用户关系表
type Relation struct {
	Id           int64 `gorm:"primaryKey"`
	User         User  `gorm:"foreignKey:FollowId"` // 关注人
	FollowId     int64
	FollowedUser User `gorm:"foreignKey:UserId"` // 被关注人
	UserId       int64
	CreateTime   int64 `gorm:"autoCreateTime"`
}
