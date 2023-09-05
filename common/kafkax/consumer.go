package kafkax

import (
	"github.com/IBM/sarama"
)

func NewConsumer(addr []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()

	consumer, err := sarama.NewConsumer(addr, config)
	if err != nil {
		return consumer, err
	}

	return consumer, nil
}
