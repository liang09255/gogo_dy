package controller

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"main/controller/middleware"
)

func Init(h *server.Hertz) {
	dy := h.Group("/douyin/")
	// 用户接口
	userGroup := dy.Group("user/")
	userGroup.POST("register/", User.Register)
	userGroup.POST("login/", User.Login)
	userGroup.GET("", middleware.Jwt.MiddlewareFunc(), User.UserInfo)
	// 视频接口
	dy.GET("feed/", Video.Feed)
	dy.POST("publish/action/", Video.PublishAction)
	dy.GET("publish/list/", Video.Publishlist)
	dy.GET("find/", Video.PlayVideo)
	// 点赞接口
	favoriteGroup := dy.Group("favorite/")
	favoriteGroup.POST("action/", Favorite.Action)
	favoriteGroup.GET("list/", Favorite.List)
	// 评论接口
	commentGroup := dy.Group("comment/")
	commentGroup.POST("action/", Comment.Action)
	commentGroup.GET("list/", Comment.List)
	// 关注接口
	relationGroup := dy.Group("relation/")
	relationGroup.POST("action/", Relation.Action)
	relationGroup.GET("follow/list/", Relation.FollowList)
	relationGroup.GET("follower/list/", Relation.FollowerList)
	relationGroup.GET("friend/list/", Relation.FriendList)
	// 聊天接口
	messageGroup := dy.Group("message/", middleware.Jwt.MiddlewareFunc())
	messageGroup.POST("action/", Message.Action)
	messageGroup.GET("chat/", Message.Chat)
}
