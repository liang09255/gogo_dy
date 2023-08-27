package utils

import (
	"common/ggKafka"
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

// user example
//
// utils.SendLog(ggKafka.Info("msg", "func_name", kk.FieldMap{
//     "param": "param",
// }))
