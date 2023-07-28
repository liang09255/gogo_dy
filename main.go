package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"main/controller"
	"main/dal"
	"main/global"
)

func main() {
	global.Init()
	dal.Init()

	h := server.Default()
	controller.Init(h)
	h.Spin()
}
