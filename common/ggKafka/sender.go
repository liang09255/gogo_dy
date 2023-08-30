package ggKafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

type LogData struct {
	// json数据
	Data  []byte
	Topic string
}

type KafkaWriter struct {
	w    *kafka.Writer
	Data chan LogData
}

// GetWriter 初始化客户端
func GetWriter(addr string) *KafkaWriter {
	w := &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Balancer: &kafka.LeastBytes{},
	}
	k := &KafkaWriter{w: w, Data: make(chan LogData, 100)}
	go k.sendMsg()
	return k
}

func (kw *KafkaWriter) Send(msg LogData) {
	kw.Data <- msg
}

func (kw *KafkaWriter) sendMsg() {
	for {
		select {
		case data := <-kw.Data:
			msg := kafka.Message{
				Topic: data.Topic,
				Key:   []byte("logMsg"),
				Value: data.Data,
			}
			var err error
			const retries = 3
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			for i := 0; i < retries; i++ {
				err = kw.w.WriteMessages(ctx, msg)
				if err == nil {
					break
				}
				if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					continue
				}
				if err != nil {
					zap.L().Error("kafka send log writer msg err", zap.Error(err))
				}
			}
		}
	}
}

func (kw *KafkaWriter) Close() {
	if kw.w != nil {
		kw.w.Close()
	}
}
