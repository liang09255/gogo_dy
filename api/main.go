package main

import (
	"api/controller"
	"api/controller/middleware"
	"api/dal"
	"api/global"
	"github.com/cloudwego/hertz/pkg/app/server"
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
