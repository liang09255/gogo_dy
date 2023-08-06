package ctlFunc

import (
	"main/controller/ctlModel"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	SuccessCode = 0
	FailCode    = 1
)

func Response(ctx *app.RequestContext, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func BaseSuccessResp(ctx *app.RequestContext, msg string) {
	Response(ctx, ctlModel.BaseResp{
		StatusCode: SuccessCode,
		StatusMsg:  msg,
	})
}

func BaseFailedResp(ctx *app.RequestContext, msg string) {
	Response(ctx, ctlModel.BaseResp{
		StatusCode: FailCode,
		StatusMsg:  msg,
	})
}
