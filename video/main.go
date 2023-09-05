package main

import (
	"common/ggConfig"
	"common/ggLog"
	"common/ggShutDown"
	"video/internal/dal"
	"video/internal/mq"
	"video/router"
)

func main() {
	dal.Init()

	videoServerConfig := ggConfig.Config.VideoServer
	// 注册grpc服务
	gc := router.StartGrpc()

	// 服务注册
	router.RegisterEtcdServer()

	// 开启消息队列
	err := mq.InitKafka()
	if err != nil {
		ggLog.Error("启动kafka失败", err)
	}

	var exit = make(chan struct{})
	<-exit
	// todo 待完善功能
	ggShutDown.ShutDown(videoServerConfig.Name, videoServerConfig.Addr, gc.Stop)
}
