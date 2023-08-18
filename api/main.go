package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"main/controller"
	"main/controller/middleware"
	"main/dal"
	"main/global"
)

func main() {
	global.Init()
	dal.Init()
	middleware.Init()

	// 1024code 只能使用8080
	h := server.Default(
		server.WithHostPorts("0.0.0.0:8080"),
		server.WithMaxRequestBodySize(30<<20), // 最大30MB
	)
	controller.Init(h)
	h.Spin()
}
