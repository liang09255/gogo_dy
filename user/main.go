package main

import (
	"common/ggConfig"
	"common/ggShutDown"
	"user/pkg/utils"
	"user/router"
)

func main() {
	userServerConfig := ggConfig.Config.UserServer
	// 注册grpc服务
	gc := router.StartGrpc()

	// 将grpc服务注册到etcd
	router.RegisterEtcdServer()

	// 初始化kafka
	kafkaCloseFunc := utils.InitKafkaWriter()

	// 优雅启停的时候，将grpc服务一起停掉
	stop := func() {
		gc.Stop()
		kafkaCloseFunc()
	}

	var exit = make(chan struct{})
	<-exit
	// todo 待完善功能
	ggShutDown.ShutDown(userServerConfig.Name, userServerConfig.Addr, stop)
}
