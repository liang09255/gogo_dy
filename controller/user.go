package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/middleware"
	"main/service"
	"strconv"
)

type userLoginController struct{}

var UserLoginControllerInstance = &userLoginController{}

func RegUser(h *server.Hertz) {
	h.POST("/douyin/user/register/", UserLoginControllerInstance.Register)
	h.POST("/douyin/user/login/", middleware.JwtMiddleware.LoginHandler)
	h.GET("/douyin/user/", UserLoginControllerInstance.getUserInfo)
}

func (u *userLoginController) Register(c context.Context, ctx *app.RequestContext) {
	//获取参数username，password
	username, ok := ctx.GetQuery("username")
	if !ok {
		BaseFailResponse(ctx, "username is required")
		return
	}
	password, ok := ctx.GetQuery("password")
	if !ok {
		BaseFailResponse(ctx, "password is required")
		return
	}
	//调service注册，拿到token和user_id
	LoginResponse, err := service.UserService.Register(username, password)
	if err != nil {
		hlog.CtxErrorf(c, "user register error: %v", err)
		BaseFailResponse(ctx, "user register error")
		return
	}
	LoginSuccessResponse(ctx, "success", *LoginResponse)
}

func (u *userLoginController) getUserInfo(c context.Context, ctx *app.RequestContext) {
	//获取参数username，password
	userId, ok := ctx.GetQuery("user_id")
	if !ok {
		BaseFailResponse(ctx, "user_id is required")
		return
	}
	token, ok := ctx.GetQuery("token")
	if !ok {
		BaseFailResponse(ctx, "token is required")
		return
	}
	toUserID, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		BaseFailResponse(ctx, "to_user_id must be an integer")
		return
	}
	//调service查用户信息
	UserInfoResponse, err := service.UserService.GetUserInfo(toUserID, token)
	if err != nil {
		hlog.CtxErrorf(c, "get userinfo error: %v", err)
		BaseFailResponse(ctx, "get userinfo error")
		return
	}
	UserInfoSuccessResponse(ctx, "success", *UserInfoResponse)
}

// 封装登录接口返回信息
type MessageUserResponse struct {
	BaseResponse
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

func LoginSuccessResponse(ctx *app.RequestContext, msg string, response service.LoginResponse) {
	Response(ctx, MessageUserResponse{
		BaseResponse: BaseResponse{StatusCode: SuccessCode, StatusMsg: msg},
		UserId:       response.UserId,
		Token:        response.Token,
	})
}

// 封装用户信息查询接口返回信息(不全)
type UserInfoResponse struct {
	BaseResponse
	UserInfoList interface{} `json:"user"`
}

func UserInfoSuccessResponse(ctx *app.RequestContext, msg string, data interface{}) {

	Response(ctx, UserInfoResponse{

		BaseResponse: BaseResponse{StatusCode: SuccessCode, StatusMsg: msg},

		UserInfoList: data,
	})
}
