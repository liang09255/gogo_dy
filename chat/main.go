package main

import (
	"chat/internal/dal"
	"chat/router"
	"common/ggConfig"
	"common/ggShutDown"
)

func main() {
	dal.Init()

	chatServerConfig := ggConfig.Config.ChatServer
	// 注册grpc服务
	gc := router.StartGrpc()

	// 将grpc服务注册到etcd
	router.RegisterEtcdServer()

	var exit = make(chan struct{})
	<-exit
	// todo 待完善功能
	ggShutDown.ShutDown(chatServerConfig.Name, chatServerConfig.Addr, gc.Stop)
}
