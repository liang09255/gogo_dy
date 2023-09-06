package ggLogKafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaReader struct {
	R *kafka.Reader
}

func GetReader(brokers []string, groupId, topic string) *KafkaReader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupId, // 同一组下的consumer 协同工作 共同消费topic队列中的内容
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	k := &KafkaReader{
		R: r,
	}
	// go k.readMsg()
	return k
}

func (kr *KafkaReader) readMsg() {
	for {
		m, err := kr.R.ReadMessage(context.Background())
		if err != nil {
			zap.L().Error("kafka receiver read msg err", zap.Error(err))
			continue
		}
		fmt.Printf("message at top %s / offset %d: %s = %s\n", m.Topic, m.Offset, string(m.Key), string(m.Value))
	}
}

func (kr *KafkaReader) Close() {
	kr.R.Close()
}
