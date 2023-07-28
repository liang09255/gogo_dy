package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"main/service"
)

type example struct{}

var Example = &example{}

func RegExample(h *server.Hertz) {
	h.GET("/hello", Example.Hello)
	h.GET("/ping", Example.Ping)
}

func (e *example) Hello(c context.Context, ctx *app.RequestContext) {
	name, ok := ctx.GetQuery("name")
	if !ok {
		BaseFailResponse(ctx, "name is required")
		return
	}
	msg, err := service.HelloService.SayHello(name)
	if err != nil {
		hlog.CtxErrorf(c, "say hello error: %v", err)
		BaseFailResponse(ctx, "say hello error")
		return
	}
	BaseSuccessResponse(ctx, msg)
}

func (e *example) Ping(c context.Context, ctx *app.RequestContext) {
	BaseSuccessResponse(ctx, "pong")
}
