package mq

import (
	"common/ggConfig"
	"common/ggLog"
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
	commentRepo     repo.CommentRepo
}

func NewKafkaDomain() {
	kd = &KafkaDomain{
		favoriteRepo:    dal.NewFavoriteDao(),
		videoDetailRepo: dal.NewVideoDetailDao(),
		commentRepo:     dal.NewCommentDao(),
	}
}

const retries = 5

func InitKafka() error {
	// 获得生产者
	var err error
	Kafka, err = kafkax.NewKafka([]string{ggConfig.Config.Kafka.Addr}, nil)
	if err != nil {
		return err
	}
	NewKafkaDomain()
	// 注册消费者
	err = AddVideoFavoriteConsumer()
	if err != nil {
		return err
	}
	// 初始化管道
	batchVideoFavorite = make(map[int64]int)
	batchChan = make(chan FavoriteMessage, 1e6)
	batchCommentMap = make(map[int64]int)
	batchCommentChan = make(chan CommentMessage, 1e6)
	// 开启统计协程
	go batchVideoFavoriteTask()
	go batchCommentTask()

	// 开启定时任务
	go batchTickerTask()

	ggLog.Info("Init Kafka success")
	return nil

}
