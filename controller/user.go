package controller

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/model"
	"main/service"
	"net/http"
)

type user struct{}

var UserOperate = &user{}

func RegUser(h *server.Hertz) {
	h.POST("/douyin/user/register/", UserOperate.Register)
	h.POST("/douyin/user/login/", UserOperate.login)
	h.GET("/douyin/user/", UserOperate.GetInfo)
}

func (e *user) Register(c context.Context, ctx *app.RequestContext) {
	var req model.UserRegisterRequest
	//获取请求参数
	//err := ctx.BindAndValidate(&req)
	req.Username = ctx.Query("username")
	req.Password = ctx.Query("password")
	//fmt.Println("注册用户信息为\n", req)

	err := service.Userinfo.Register(req.Username, req.Password)

	token := req.Username + "&" + req.Password

	//用户名已经存在
	if err != nil {
		hlog.CtxErrorf(c, "User register error: %v", err)
		ctx.JSON(http.StatusOK, model.UserRegisterResponse{
			StatusCode: FailCode,
			StatusMsg:  "Register fail",
			UserID:     0,
			Token:      token,
		})
		return
	}

	//注册用户成功，返回对应信息
	ctx.JSON(http.StatusOK, model.UserRegisterResponse{
		StatusCode: SuccessCode,
		StatusMsg:  "Register successfully",
		UserID:     1,
		Token:      token,
	})
}

func String(value string) *string {
	return &value
}

func Int64(value int64) *int64 {
	return &value
}

func (e *user) login(c context.Context, ctx *app.RequestContext) {
	var req model.UserLoginRequest

	req.Username = ctx.Query("username")
	req.Password = ctx.Query("password")

	token := req.Username + "&" + req.Password

	fmt.Println("Login:\n")
	fmt.Println("username =   token = ", req.Username, token)

	//判断登录用户的身份信息
	Id, err := service.Userinfo.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, model.UserLoginResponse{
			StatusCode: FailCode,
			StatusMsg:  String("Login Fail"),
			Token:      String(token),
			UserID:     Int64(0),
		})
		return
	}
	ctx.JSON(http.StatusOK, model.UserLoginResponse{
		StatusCode: SuccessCode,
		StatusMsg:  String("Login Successfully"),
		Token:      String(token),
		UserID:     Int64(Id),
	})
}

func (e *user) GetInfo(c context.Context, ctx *app.RequestContext) {
	userId := ctx.Query("user_id")
	token := ctx.Query("token")

	fmt.Println("GetInfo:\n")
	fmt.Println("userid =   token = ", userId, token)

	var resp model.Userinfo

	//根据id获取用户详细信息
	err := service.Userinfo.GetUserInfo(token, userId, &resp)

	if err != nil {
		ctx.JSON(http.StatusOK, model.UserInfoResponse{
			StatusCode: FailCode,
			StatusMsg:  String("Get userinfo fail"),
			User:       &resp,
		})
		hlog.CtxErrorf(c, "Get User Info error: %v", err)
		return
	}
	ctx.JSON(http.StatusOK, model.UserInfoResponse{
		StatusCode: SuccessCode,
		StatusMsg:  String("Get userinfo successfully"),
		User:       &resp,
	})
}
