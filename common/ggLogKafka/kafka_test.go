package ggLogKafka

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	w := GetWriter("localhost:9092")
	msg := make(map[string]string)
	msg["projectCode"] = "7777"
	bytes, _ := json.Marshal(msg)
	w.Send(LogData{
		Topic: "test_log",
		Data:  bytes,
	})
	time.Sleep(2 * time.Second)
}

func TestConsumer(t *testing.T) {
	reader := GetReader([]string{"localhost:9092"}, "group1", "test_log")
	reader.readMsg()
	//select {}
}
