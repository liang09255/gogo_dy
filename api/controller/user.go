package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/controller/ctlFunc"
	"main/controller/ctlModel/baseCtlModel"
	"main/controller/ctlModel/userCtlModel"
	"main/controller/middleware"
	"main/service"
)

type user struct{}

var User = &user{}

func (u *user) Register(c context.Context, ctx *app.RequestContext) {
	//获取参数username，password
	var reqObj userCtlModel.RegisterReq
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

	// 封装返回信息
	var resp = userCtlModel.RegisterResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		RegisterResponse: userCtlModel.RegisterResponse{
			UserId: LoginResponse.UserId,
			Token:  LoginResponse.Token,
		},
	}
	ctlFunc.Response(ctx, &resp)
}

func (u *user) Login(c context.Context, ctx *app.RequestContext) {
	middleware.Jwt.LoginHandler(c, ctx)
}

func (u *user) UserInfo(c context.Context, ctx *app.RequestContext) {
	var myID = middleware.GetUserID(ctx)
	var reqObj userCtlModel.InfoReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err.Error())
		return
	}
	//获取参数username，password
	userId := reqObj.UserId

	//调service查用户信息
	UserInfoResponse, err := service.UserService.GetUserInfo(userId, myID)
	if err != nil {
		hlog.CtxErrorf(c, "get userinfo error: %v", err)
		ctlFunc.BaseFailedResp(ctx, "get userinfo error")
		return
	}

	// 封装返回信息
	var resp = &userCtlModel.InfoResp{
		BaseResp: baseCtlModel.NewBaseSuccessResp(),
		User: userCtlModel.User{
			ID:             UserInfoResponse.ID,
			Username:       UserInfoResponse.Username,
			FollowCount:    UserInfoResponse.FollowCount,
			FollowerCount:  UserInfoResponse.FollowerCount,
			IsFollow:       UserInfoResponse.IsFollow,
			TotalFavorited: UserInfoResponse.TotalFavorited,
			WorkCount:      UserInfoResponse.WorkCount,
			FavoriteCount:  UserInfoResponse.FavoriteCount,
		},
	}
	ctlFunc.Response(ctx, resp)
}
