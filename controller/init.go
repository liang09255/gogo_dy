package controller

import "github.com/cloudwego/hertz/pkg/app/server"

func Init(h *server.Hertz) {
	RegExample(h)
	RegMessageAction(h)
	RegMessageChat(h)
}
