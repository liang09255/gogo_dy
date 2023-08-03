package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/service"
	"net/http"
)

type UserRegisterRequest struct {
	Password string // 密码，最长32个字符
	Username string // 注册用户名，最长32个字符
}

type UserRegisterResponse struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
	Token      string `json:"token"`       // 用户鉴权token
	UserID     int64  `json:"user_id"`     // 用户id
}

type UserLoginRequest struct {
	Password string `json:"password"` // 密码，最长32个字符
	Username string `json:"username"` // 注册用户名，最长32个字符
}

type UserLoginResponse struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	Token      *string `json:"token"`       // 用户鉴权token
	UserID     *int64  `json:"user_id"`     // 用户id
}

type UserInfoRequest struct {
	Token  string `json:"token"`   // 用户鉴权token
	UserID string `json:"user_id"` // 用户id
}

type UserInfoResponse struct {
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	User       *User   `json:"user"`        // 用户信息
}

// User
type User struct {
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

type user struct{}

var UserOperate = &user{}

func RegUser(h *server.Hertz) {
	h.POST("/douyin/user/register/", UserOperate.Register)
	h.POST("/douyin/user/login/", UserOperate.login)
}

func (e *user) Register(c context.Context, ctx *app.RequestContext) {
	var req UserRegisterRequest
	//获取请求参数
	err := ctx.BindAndValidate(&req)

	var resp UserRegisterResponse

	if err != nil {
		resp.StatusCode = FailCode
		hlog.CtxErrorf(c, "User register error: %v", err)
		return
	}

	err = service.Registerinfo.Register(req.Username, req.Password)

	//用户名已经存在
	if err != nil {
		hlog.CtxErrorf(c, "User register error: %v", err)
		ctx.JSON(http.StatusOK, UserRegisterResponse{
			StatusCode: FailCode,
			StatusMsg:  "Register fail",
			UserID:     0,
			Token:      string(0),
		})
		return
	}

	//注册用户成功，返回对应信息
	ctx.JSON(http.StatusOK, UserRegisterResponse{
		StatusCode: SuccessCode,
		StatusMsg:  "Register successfully",
		UserID:     1,
		Token:      req.Username + req.Password,
	})
}

func (e *user) login(c context.Context, ctx *app.RequestContext) {
	var req UserLoginRequest
	var resp UserLoginResponse
	err := ctx.BindAndValidate(&req)
	if err != nil {
		resp.StatusCode = FailCode
		hlog.CtxErrorf(c, "User Login error: %v", err)
		return
	}

}
