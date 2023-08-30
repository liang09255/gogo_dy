package utils

import (
	"common/ggKafka"
	"common/ggLog"
	"context"
	"time"
	"user/internal/dal"
	"user/internal/repo"
)

var kw *ggKafka.KafkaWriter

func InitKafkaWriter() func() {
	kw = ggKafka.GetWriter("localhost:9092")
	return kw.Close
}

func SendLog(data []byte) {
	kw.Send(ggKafka.LogData{
		Topic: "test_log",
		Data:  data,
	})
}

// 日志 使用 样例
//
// utils.SendLog(ggKafka.Info("msg", "func_name", kk.FieldMap{
//     "param": "param",
// }))

type KafkaCache struct {
	KR    *ggKafka.KafkaReader
	cache repo.Cache
}

func SendCache(data []byte) {
	kw.Send(ggKafka.LogData{
		Topic: "go-gin-grpc_cache",
		Data:  data,
	})
}

func NewKafkaCacheReader() *KafkaCache {
	kr := ggKafka.GetReader([]string{"localhost:9092"}, "cache_group", "go-gin-grpc_cache")
	return &KafkaCache{
		KR:    kr,
		cache: dal.Rc,
	}
}

func (c *KafkaCache) DeleteCache() {
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
