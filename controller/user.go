package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/controller/ctlFunc"
	"main/controller/ctlModel"
	"main/controller/middleware"
	"main/service"
)

type user struct{}

var User = &user{}

func (u *user) Register(c context.Context, ctx *app.RequestContext) {
	//获取参数username，password
	var reqObj ctlModel.UserRegisterReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}
	username := reqObj.Username
	password := reqObj.Password

	//调service注册，拿到token和user_id
	LoginResponse, err := service.UserService.Register(username, password)
	if err != nil {
		hlog.CtxErrorf(c, "user register error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "user register error")
		return
	}
	LoginSuccessResponse(ctx, "success", *LoginResponse)
}

func (u *user) Login(c context.Context, ctx *app.RequestContext) {
	middleware.Jwt.LoginHandler(c, ctx)
}

func (u *user) UserInfo(c context.Context, ctx *app.RequestContext) {
	var reqObj ctlModel.UserInfoReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}
	//获取参数username，password
	userId := reqObj.UserId

	//调service查用户信息
	UserInfoResponse, err := service.UserService.GetUserInfo(userId)
	if err != nil {
		hlog.CtxErrorf(c, "get userinfo error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "get userinfo error")
		return
	}
	UserInfoSuccessResponse(ctx, "success", *UserInfoResponse)
}

// 封装登录接口返回信息
type MessageUserResponse struct {
	ctlModel.BaseResp
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

func LoginSuccessResponse(ctx *app.RequestContext, msg string, response service.LoginResponse) {
	ctlFunc.Response(ctx, MessageUserResponse{
		BaseResp: ctlModel.BaseResp{StatusCode: ctlFunc.SuccessCode, StatusMsg: msg},
		UserId:   response.UserId,
		Token:    response.Token,
	})
}

// 封装用户信息查询接口返回信息(不全)
type UserInfoResponse struct {
	ctlModel.BaseResp
	UserInfoList interface{} `json:"user"`
}

func UserInfoSuccessResponse(ctx *app.RequestContext, msg string, data interface{}) {
	ctlFunc.Response(ctx, UserInfoResponse{
		BaseResp:     ctlModel.BaseResp{StatusCode: ctlFunc.SuccessCode, StatusMsg: msg},
		UserInfoList: data,
	})
}
