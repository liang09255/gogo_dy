package mq

import (
	"context"
	"sync"
)

var (
	// vid - commentCount
	batchCommentMap  map[int64]int
	batchCommentChan chan CommentMessage
	muComment        sync.Mutex
)

func SendCommentMsg(msg CommentMessage) {
	batchCommentChan <- msg
}

func batchCommentTask() {
	for {
		select {
		case msg := <-batchCommentChan:
			muComment.Lock()
			if _, ok := batchCommentMap[msg.Vid]; ok {
				if msg.Method == 1 {
					batchCommentMap[msg.Vid]++
				} else {
					batchCommentMap[msg.Vid]--
				}
			} else {
				if msg.Method == 1 {
					batchCommentMap[msg.Vid] = 1
				} else {
					batchCommentMap[msg.Vid] = -1
				}
			}
			muComment.Unlock()
		}
	}
}

func TimeCommentTask() {
	muComment.Lock()
	if len(batchCommentMap) == 0 {
		muComment.Unlock()
		return
	}

	// 开启协程写入mysql,避免阻塞，交给mysql去处理
	go kd.videoDetailRepo.BatchInsertComment(context.Background(), batchCommentMap)
	batchCommentMap = make(map[int64]int)
	muComment.Unlock()
}

//func AddCommentConsumer() error {
//	partitionList, err := Kafka.Consumer.Partitions(constant.CommentTopic)
//	if err != nil {
//		return err
//	}
//
//	for patition := range partitionList {
//		h := &kafkax.CustomerHandler{
//			Topic:     constant.CommentTopic,
//			Partition: int32(patition),
//			CustomFn:  CommentHandler,
//			ErrorFn:   ErrorFn,
//		}
//		err = Kafka.AddConsumerHandler(h)
//	}
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func CommentHandler(msg *sarama.ConsumerMessage) error {
//	commentMsg := &CommentMessage{}
//	err := json.Unmarshal(msg.Value,commentMsg)
//
//	if err != nil{
//		return err
//	}
//
//	// 写入管道
//}
