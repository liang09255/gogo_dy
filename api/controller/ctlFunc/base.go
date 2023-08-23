package ctlFunc

import (
	"api/controller/ctlModel/baseCtlModel"
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
