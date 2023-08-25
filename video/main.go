package main

import (
	"common/ggConfig"
	"common/ggShutDown"
	"video/router"
)

func main() {
	videoServerConfig := ggConfig.Config.VideoServer
	// 注册grpc服务
	gc := router.StartGrpc()

	// 服务注册
	router.RegisterEtcdServer()

	var exit = make(chan struct{})
	<-exit
	// todo 待完善功能
	ggShutDown.ShutDown(videoServerConfig.Name, videoServerConfig.Addr, gc.Stop)
}
