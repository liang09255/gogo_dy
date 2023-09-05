package kafkax

import "github.com/IBM/sarama"

// 同步生产者
func NewSyncProducer(addr []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 设置分区随机 - 如果需要设置根据key分配分区的话可以修改
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 设置成功会有响应
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		return producer, err
	}

	return producer, nil
}

// 异步生产者，先写缓冲区然后就返回
func NewAsyncProducer(addr []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 同步生产者需要开启，因为要同步返回发送成功或失败
	// 异步生产者可以考虑只开启Error，但是需要把返回值从error中读取出来，否则会阻塞
	//config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(addr, config)
	if err != nil {
		return producer, err
	}
	return producer, nil
}
