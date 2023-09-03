package mq

import (
	"common/ggLog"
	"common/kafkax"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"sync"
	"time"
	"video/pkg/constant"
)

// 增加视频的点赞记录
// vid - 点赞数
// 合并后的map
var (
	batchVideoFavorite map[int64]int
	// 消息传入管道，再通过管道写入batchVideoFavorite,然后每10s消费一次，消费的时候加锁，等解决完了再从管道中取出
	// 用管道异步处理加快了效率，但是有个问题就是如果服务宕机了，会使得channel内未消费的数据丢失
	// 但是不用的话，因为用map进行统计，频繁加锁解锁很影响性能，先用channel写着
	// 而且在这里就算数据丢失了也不是真正的丢失，因为这一步只是写入video的操作而已.如果需要正确数据可以再走一遍统计
	batchChan chan FavoriteMessage
	mu        sync.Mutex
)

func batchVideoFavoriteTask() {
	// 从chan中读取，写入batchVideoFavorite
	for {
		select {
		case msg := <-batchChan:
			mu.Lock()
			// 将msg写入map中进行统计
			if _, ok := batchVideoFavorite[msg.Vid]; ok {
				if msg.Method == 1 {
					batchVideoFavorite[msg.Vid]++
				} else {
					batchVideoFavorite[msg.Vid]--
				}
			} else {
				if msg.Method == 1 {
					batchVideoFavorite[msg.Vid] = 1
				} else {
					batchVideoFavorite[msg.Vid] = -1
				}
			}
			mu.Unlock()
		}
	}
}

// 定时处理任务
func batchTickerTask() {
	// 10s定期将map的数值消费一次
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			mu.Lock()
			if len(batchVideoFavorite) == 0 {
				mu.Unlock()
				continue
			}
			// 开启协程写入mysql，免得在这里阻塞过久,交给mysql去处理
			go kd.videoDetailRepo.BatchInsert(context.Background(), batchVideoFavorite)
			// 这里先置为空然后释放锁
			batchVideoFavorite = make(map[int64]int)
			mu.Unlock()
		}
	}
}

func AddVideoFavoriteConsumer() error {
	// 加入一个主题消费者
	// 先设置处理方法
	// 遍历分区，为每个分区创建一个分区消费者
	partitionList, err := Kafka.Consumer.Partitions(constant.FavoriteTopic)
	if err != nil {
		return err
	}
	// 加入分区消费者
	for partition := range partitionList {
		h := &kafkax.CustomerHandler{
			Topic:     constant.FavoriteTopic,
			Partition: int32(partition),
			CustomFn:  VideoFavoriteRecordHandler,
			ErrorFn:   ErrorFn,
		}
		err = Kafka.AddConsumerHandler(h)
	}

	if err != nil {
		return err
	}

	return nil

}

// 增加点赞记录处理方法
func VideoFavoriteRecordHandler(msg *sarama.ConsumerMessage) error {
	// 接收到消息后,进行处理
	// 处理流程：
	// 1. 先将消息解析出来
	favoriteMsg := &FavoriteMessage{}
	err := json.Unmarshal(msg.Value, favoriteMsg)

	if err != nil {
		return err
	}

	// 2.根据消息写入数据库
	// 获得数据库连接
	// TODO 用kafka有个问题，点赞顺序无法保证，因为是每个分区开了一个协程去读，无法保证哪个协程内的先被读取到
	// 想要顺序的话只能保证局部有序，即分区内有序
	if favoriteMsg.Method == 1 {
		// 新增
		err = kd.favoriteRepo.PostFavoriteAction(context.Background(), favoriteMsg.Uid, favoriteMsg.Vid)
	} else if favoriteMsg.Method == 2 {
		// 删除
		err = kd.favoriteRepo.CancelFavoriteAction(context.Background(), favoriteMsg.Uid, favoriteMsg.Vid)
	}

	if err != nil {
		return err
	}
	// 写入管道
	batchChan <- *favoriteMsg

	return nil
}

// 错误处理方法
func ErrorFn(h *kafkax.CustomerHandler, err error) {
	ggLog.Errorf("对主题:%s,分区:%d 进行消费错误,错误问题:%v", h.Topic, h.Partition, err)
}
