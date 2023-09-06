package main

import (
	"common/ggConfig"
	"common/ggShutDown"
	"user/internal/dal"
	"user/router"
)

func main() {

	// 开始每日关注统计同步任务
	fsd := dal.NewFollowStatusDal()
	cronJob := fsd.StartDailyFollowStatsSync()
	defer cronJob.Stop() // 在 main 结束时确保 cron 停止

	dal.Init()

	userServerConfig := ggConfig.Config.UserServer
	// 注册grpc服务
	gc := router.StartGrpc()

	// 将grpc服务注册到etcd
	router.RegisterEtcdServer()
	//
	//// 初始化kafka
	//kafkaCloseFunc := utils.InitKafkaWriter()
	//
	//// 初始化kafka消费者
	//reader := utils.NewKafkaCacheReader()
	//go reader.DeleteCache()
	//
	// 优雅启停的时候，将grpc服务一起停掉
	stop := func() {
		gc.Stop()
		//kafkaCloseFunc()
		//reader.KR.Close()
	}

	var exit = make(chan struct{})
	<-exit
	// todo 待完善功能
	ggShutDown.ShutDown(userServerConfig.Name, userServerConfig.Addr, stop)
}
