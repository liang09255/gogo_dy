package constant

import "time"

const (
	// 用户获赞数
	UserGetFavoriteCountKey = "user.getFavorite.%d"
	// 用户喜欢数
	UserFavoriteCountKey = "user.favorite.%d"
	// 视频点赞数
	VideoFavoriteCountKey = "video.favorite.%d"
	// 评论列表缓存
	CommentListKey = "comment.list.%d"
)

const (
	CommentExpireTime = 60 * time.Second
)
