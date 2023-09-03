package kafkax

import (
	"fmt"
	"testing"
)

var topic = "TOPIC"

func TestProducer(t *testing.T) {
	k, err := NewKafka([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println(err)
	}
	k.SendMessage(&ProducerMsg{
		Topic: topic,
		Data:  []byte("124124"),
	})

}

func SendMessage() {

}
