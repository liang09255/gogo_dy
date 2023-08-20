package main

import (
	"common/ggConfig"
	"common/ggShutDown"
	"user/router"
)

func main() {
	userServerConfig := ggConfig.Config.UserServer
	// 注册grpc服务
	gc := router.StartGrpc()

	// 将grpc服务注册到etcd
	router.RegisterEtcdServer()

	var exit = make(chan struct{})
	<-exit
	// todo 待完善功能
	ggShutDown.ShutDown(userServerConfig.Name, userServerConfig.Addr, gc.Stop)
}
