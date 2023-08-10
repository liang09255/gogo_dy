package userCtlModel

import (
	"main/controller/ctlModel/baseCtlModel"
)

// 注册

type RegisterReq struct {
	Username string `query:"username,required" vd:"len($)>2 && len($)<32"`
	Password string `query:"password,required" vd:"len($)>5 && len($)<32"`
}

type RegisterResp struct {
	baseCtlModel.BaseResp
	RegisterResponse
}

type RegisterResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

// 登录

type LoginReq struct {
	Username string `query:"username,required" vd:"len($)>2 && len($)<32"`
	Password string `query:"password,required" vd:"len($)>5 && len($)<32"`
}

type LoginResp struct {
	baseCtlModel.BaseResp
	LoginResponse
}

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

// 获取用户信息

type InfoReq struct {
	UserId int64 `query:"user_id,required" vd:"$>0"`
}

type InfoResp struct {
	baseCtlModel.BaseResp
	User User `json:"user"`
}

type User struct {
	ID             int64  `json:"id"`
	Username       string `json:"name"`
	FollowCount    int64  `json:"follow_count"`
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow"`
	TotalFavorited string `json:"total_favorited"` // 获点赞数量
	WorkCount      int64  `json:"work_count"`      // 作品数
	FavoriteCount  int64  `json:"favorite_count"`  //喜欢数
}
