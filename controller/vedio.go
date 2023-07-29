package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"time"
)

type vedio struct{}

var Vedio = &vedio{}

type Video struct {
	ID       int64     `gorm:"column:id;type:int;primaryKey;autoIncrement:true" json:"id,string"`
	AuthorID int64     `gorm:"column:author_id;type:int" json:"author_id"`
	PlayURL  string    `gorm:"column:play_url;type:varchar(255)" json:"play_url"`
	CoverURL string    `gorm:"column:cover_url;type:varchar(255)" json:"cover_url"`
	Time     time.Time `gorm:"column:time;type:datetime" json:"time"`
	Title    string    `gorm:"column:title;type:varchar(255)" json:"title"`
}

func RegVedio(h *server.Hertz) {
	h.POST("/douyin/feed/", Vedio.Feed)

}
func (e *vedio) Feed(c context.Context, ctx *app.RequestContext) {

}

func (e *vedio) Publishaction(c context.Context, ctx *app.RequestContext) {

}

func (e *vedio) Publishlist(c context.Context, ctx *app.RequestContext) {

}

//func (e *example) Hello(c context.Context, ctx *app.RequestContext) {
//	name, ok := ctx.GetQuery("name")
//	if !ok {
//		BaseFailResponse(ctx, "name is required")
//		return
//	}
//	msg, err := service.HelloService.SayHello(name)
//	if err != nil {
//		hlog.CtxErrorf(c, "say hello error: %v", err)
//		BaseFailResponse(ctx, "say hello error")
//		return
//	}
//	BaseSuccessResponse(ctx, msg)
//}
//
//func (e *example) Ping(c context.Context, ctx *app.RequestContext) {
//	BaseSuccessResponse(ctx, "pong")
//}
