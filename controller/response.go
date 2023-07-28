package controller

import (
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

const (
	SuccessCode = 0
	FailCode    = 1
)

type BaseResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func Response(ctx *app.RequestContext, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func BaseSuccessResponse(ctx *app.RequestContext, msg string) {
	Response(ctx, BaseResponse{StatusCode: SuccessCode, StatusMsg: msg})
}

func BaseFailResponse(ctx *app.RequestContext, msg string) {
	Response(ctx, BaseResponse{StatusCode: FailCode, StatusMsg: msg})
}
