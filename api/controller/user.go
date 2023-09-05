package controller

import (
	"context"
	"time"

	"api/controller/ctlFunc"
	"api/controller/ctlModel/baseCtlModel"
	"api/controller/ctlModel/userCtlModel"
	"api/controller/middleware"
	"api/global"
	userRPC "common/ggIDL/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/jinzhu/copier"
)

type user struct{}

var User = &user{}

func (u *user) Register(c context.Context, ctx *app.RequestContext) {
	//获取参数username，password
	var reqObj userCtlModel.RegisterReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err)
		return
	}

	msg := &userRPC.RegisterRequest{
		Username: reqObj.Username,
		Password: reqObj.Password,
	}
	// 超时控制
	c, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	// 调grpc service注册，拿到token和user_id
	RegisterResponse, err := global.UserClient.Register(c, msg)
	if err != nil {
		hlog.CtxErrorf(c, "user register error: %v", err)
		ctlFunc.BaseFailedResp(ctx, baseCtlModel.ServerInternal.WithDetails(err.Error()))
		return
	}

	// 封装返回信息
	registerResponse := &userCtlModel.RegisterResponse{}
	if err := copier.Copy(registerResponse, RegisterResponse); err != nil {
		hlog.CtxErrorf(c, "user register error: %v", err)
		ctlFunc.BaseFailedResp(ctx, err)
		return
	}

	// 添加token
	registerResponse.Token, err = middleware.ReleaseToken(registerResponse.UserId)
	if err != nil {
		hlog.CtxErrorf(c, "user register error: %v", err)
		ctlFunc.BaseFailedResp(ctx, baseCtlModel.InvalidToken.WithDetails(err.Error()))
		return
	}
	var response = userCtlModel.RegisterResp{
		APIBaseResp:      baseCtlModel.NewBaseSuccessResp(),
		RegisterResponse: *registerResponse,
	}
	ctlFunc.Response(ctx, &response)
}

func (u *user) Login(c context.Context, ctx *app.RequestContext) {
	middleware.Jwt.LoginHandler(c, ctx)
}

func (u *user) UserInfo(c context.Context, ctx *app.RequestContext) {
	var myID = middleware.GetUserID(ctx)
	var reqObj userCtlModel.InfoReq
	if err := ctx.BindAndValidate(&reqObj); err != nil {
		ctlFunc.BaseFailedResp(ctx, err)
		return
	}

	var msg = &userRPC.UserInfoRequest{
		UserId: []int64{reqObj.UserId},
		MyId:   myID,
	}

	UserInfoResponse, err := global.UserClient.MGetUserInfo(c, msg)
	if err != nil {
		hlog.CtxErrorf(c, "user register error: %v", err)
		ctlFunc.BaseFailedResp(ctx, baseCtlModel.ServerInternal.WithDetails(err.Error()))
		return
	}

	hlog.Infof("UserInfoResponse: %v", UserInfoResponse)

	userInfo := &userCtlModel.User{}
	if err := copier.Copy(userInfo, UserInfoResponse.UserInfo[0]); err != nil {
		hlog.CtxErrorf(c, "user register error: %v", err)
		ctlFunc.BaseFailedResp(ctx, baseCtlModel.UserNotFound.WithDetails(err.Error()))
		return
	}

	// 封装返回信息
	var resp = &userCtlModel.InfoResp{
		APIBaseResp: baseCtlModel.NewBaseSuccessResp(),
		User:        *userInfo,
	}
	ctlFunc.Response(ctx, resp)
}
