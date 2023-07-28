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

  // 1024code 只能使用8080
	h := server.Default(server.WithHostPorts("0.0.0.0:8080"))
	controller.Init(h)
	h.Spin()
}
