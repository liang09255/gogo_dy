package ctlFunc

import (
	"main/controller/ctlModel/baseCtlModel"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

func Response(ctx *app.RequestContext, data interface{}) {
	ctx.JSON(http.StatusOK, data)
}

func BaseSuccessResp(ctx *app.RequestContext, msg ...string) {
	Response(ctx, baseCtlModel.NewBaseSuccessResp(msg...))
}

func BaseFailedResp(ctx *app.RequestContext, msg ...string) {
	Response(ctx, baseCtlModel.NewBaseFailedResp(msg...))
}

// TODo 以下内容计划废弃

func ResponseWithData(ctx *app.RequestContext, msg string, data interface{}) {
	Response(ctx, BaseResponseWithData{
		BaseResp: baseCtlModel.BaseResp{
			StatusCode: 0,
			StatusMsg:  "",
		},
		Data: data,
	})
}

type BaseResponseWithData struct {
	baseCtlModel.BaseResp
	Data interface{} `json:"data"`
}
