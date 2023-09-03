package mq

import (
	"common/ggConfig"
	"common/kafkax"
	"video/internal/dal"
	"video/internal/repo"
)

//var (
//	producer sarama.SyncProducer
//	consumer sarama.Consumer
//	data     chan *ProducerData
//	// 记录已初始化的分区消费者,防止重复初始化
//	partitionComsemer map[string]sarama.PartitionConsumer
//)

var Kafka *kafkax.Kafka

var kd *KafkaDomain

// 操控数据库方法
type KafkaDomain struct {
	// 获得数据库连接
	favoriteRepo    repo.FavoriteRepo
	videoDetailRepo repo.VideoDetailRepo
}

func NewKafkaDomain() {
	kd = &KafkaDomain{
		favoriteRepo:    dal.NewFavoriteDao(),
		videoDetailRepo: dal.NewVideoDetailDao(),
	}
}

const retries = 5

func InitKafka() error {
	// 获得生产者
	var err error
	Kafka, err = kafkax.NewKafka([]string{ggConfig.Config.Kafka.Addr}, nil)
	NewKafkaDomain()
	// 初始化管道
	batchChan = make(chan FavoriteMessage, 1e6)
	// 开启统计协程
	go batchVideoFavoriteTask()
	// 开启定时任务
	go batchTickerTask()

	if err != nil {
		return err
	}

	return nil

}
