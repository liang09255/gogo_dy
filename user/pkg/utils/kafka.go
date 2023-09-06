package utils

import (
	"common/ggLog"
	"common/ggLogKafka"
	"context"
	"time"
	"user/internal/dal"
	"user/internal/repo"
)

var kw *ggLogKafka.KafkaWriter

func InitKafkaWriter() func() {
	kw = ggLogKafka.GetWriter("localhost:9092")
	return kw.Close
}

// SendMessageToKafka 创建一个更通用的函数，将主题作为参数传入
func SendMessageToKafka(topic string, data []byte) {
	kw.Send(ggLogKafka.LogData{
		Topic: topic,
		Data:  data,
	})
}

// 日志 使用 样例
//
// utils.SendLog(ggKafka.Info("msg", "func_name", kk.FieldMap{
//     "param": "param",
// }))

type KafkaCache struct {
	KR    *ggLogKafka.KafkaReader
	cache repo.Cache
}

func SendCache(data []byte) {
	kw.Send(ggLogKafka.LogData{
		Topic: "go-gin-grpc_cache",
		Data:  data,
	})
}

func NewKafkaCacheReader() *KafkaCache {
	kr := ggLogKafka.GetReader([]string{"localhost:9092"}, "cache_group", "go-gin-grpc_cache")
	return &KafkaCache{
		KR:    kr,
		cache: dal.Rc,
	}
}

func (c *KafkaCache) DeleteCache() {

	// 这里是kafka消费者的消费逻辑
	for {
		message, err := c.KR.R.ReadMessage(context.Background())
		if err != nil {
			ggLog.Error("读取缓存失败:", err)
			continue
		}
		// sign是相关缓存的标识
		// sign自己设定，可以设定多个
		if "sign" == string(message.Value) {
			//删除任务相关的缓存key
			fields, _ := c.cache.HKeys(context.Background(), "sign")
			// 延迟删除缓存
			time.Sleep(1 * time.Second)
			c.cache.Delete(context.Background(), fields)
		}
	}
}

//// go reader.ProcessRelationActionMessages()
//func (c *KafkaCache) ProcessRelationActionMessages(relationDomain *domain.RelationDomain) {
//	for {
//		message, err := c.KR.R.ReadMessage(context.Background())
//		if err != nil {
//			ggLog.Error("读取关系操作消息失败:", err)
//			continue
//		}
//		var msgData model.RelationActionMessage
//		err = json.Unmarshal(message.Value, &msgData)
//		if err != nil {
//			ggLog.Error("解析消息失败:", err)
//			continue
//		}
//
//		// 根据消息内容异步更新关注统计表
//		go relationDomain.UpdateFollowStats(msgData)
//	}
//}
