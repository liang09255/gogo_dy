package main

import (
	"common/shutDown"
	"github.com/cloudwego/hertz/pkg/app/server"
	"user/config"
	"user/router"
)

func main() {
	// 注册grpc服务
	gc := router.RegisterGrpc()

	// 将grpc服务注册到etcd
	router.RegisterEtcdServer()

	h := server.Default(
		server.WithHostPorts(config.C.SC.Addr),
		server.WithMaxRequestBodySize(30<<20), // 最大30MB
	)
	h.Spin()

	shutDown.ShutDown(config.C.SC.Name, config.C.SC.Addr, gc.Stop)
}
