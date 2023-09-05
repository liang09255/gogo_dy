package kafkax

import (
	"common/ggLog"
	"fmt"
	"github.com/IBM/sarama"
	"strings"
	"time"
)

type Kafka struct {
	// 生产者数据
	Data chan *ProducerMsg

	// 消费者消费实例
	Handler chan *CustomerHandler

	// 连接地址
	Broker []string

	Config *sarama.Config

	Producer sarama.SyncProducer
	Consumer sarama.Consumer

	// 记录已经初始化的分区消费者
	//Pc map[string]sarama.PartitionConsumer

	Retire int
}

type ProducerMsg struct {
	Topic     string
	Partition int32
	Data      []byte
}

// 消费者读取目标以及处理
type CustomerHandler struct {
	Topic     string
	Partition int32
	Offset    int64
	// 自定义的数据消费操作函数
	CustomFn func(message *sarama.ConsumerMessage) error
	// 自定义错误处理函数
	ErrorFn func(*CustomerHandler, error)
}

func NewKafka(addr []string, config *sarama.Config) (k *Kafka, err error) {
	if config == nil {
		config = sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		// 是否等待成功和失败后的响应，同步生产者需要等待
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true

		// 提交offset的间隔时间，每秒提交一次
		// 事实上手动提交可能会更好一点，后续有空再改
		config.Consumer.Offsets.AutoCommit.Enable = true
		config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	}

	if len(addr) == 0 {
		return nil, fmt.Errorf("broker is empty")
	}

	k = &Kafka{
		Config:  config,
		Broker:  addr,
		Handler: make(chan *CustomerHandler),
		Data:    make(chan *ProducerMsg, 100000),
		//Pc:      make(map[string]sarama.PartitionConsumer),
		// 先设置重试次数5次，可以改为参数配置
		Retire: 5,
	}

	k.Producer, err = sarama.NewSyncProducer(k.Broker, k.Config)
	if err != nil {
		return k, err
	}

	k.Consumer, err = sarama.NewConsumer(k.Broker, k.Config)
	if err != nil {
		return k, err
	}

	//defer k.producer.Close()
	//defer k.consumer.Close()

	go k.produce()
	go k.consume()

	return

}

// 使用管道发送主要是为了快，做了异步性能优化，而且进行解耦，业务程序只需要写入管道，然后同步生产者去读取再传入kafka
// 而且设置管道的容量，可以避免并发过高的流量
// 但这样处理可能有个问题，就是如果消息发送失败了，业务逻辑也没有去管他，所以要设置好重试机制，如果不允许发送失败的话
// 需要自己控制生产者去发送，并且接收返回结果，设置重试机制
func (k *Kafka) SendMessage(msg *ProducerMsg) error {
	if msg == nil || k.Producer == nil {
		return fmt.Errorf("msg is nil or producer is nil")
	}
	k.Data <- msg
	return nil
}

func (k *Kafka) AddConsumerHandler(h *CustomerHandler) (err error) {
	if h == nil || k.Consumer == nil {
		return fmt.Errorf("h is nil or consumer is nil")
	}
	if h.Offset == 0 {
		h.Offset = sarama.OffsetNewest
	}
	k.Handler <- h
	return err
}

func (k *Kafka) produce() {
	for {
		select {
		case p := <-k.Data:
			msg := &sarama.ProducerMessage{
				Topic:     p.Topic,
				Partition: p.Partition,
				Value:     sarama.ByteEncoder(p.Data),
			}
			for i := 0; i < k.Retire; i++ {
				_, _, err := k.Producer.SendMessage(msg)
				if err == nil {
					break
				}
				ggLog.Error("发送消息错误:%v", p)
				// 可以考虑做持久化发送错误消息，然后由后台协程进行处理
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (k *Kafka) consume() {
	for {
		select {
		// 每个分区消费者添加一个分区消费实例
		case h := <-k.Handler:
			// 检查是否已经创建过 ConsumerPartition
			//key := fmt.Sprintf("%s:%d", h.Topic, h.Partition)
			//_, ok := k.Pc[key]
			//if ok {
			//	// 如果已经创建，则跳过
			//	continue
			//}

			var pc sarama.PartitionConsumer
			pc, err := k.Consumer.ConsumePartition(h.Topic, h.Partition, sarama.OffsetNewest)
			if err != nil {
				// outside 一般就是偏移量超出了范围
				if strings.Contains(err.Error(), "outside") {
					pc, err = k.Consumer.ConsumePartition(h.Topic, h.Partition, sarama.OffsetNewest)
					if err != nil {
						h.ErrorFn(h, err)
						break
					}
				} else {
					h.ErrorFn(h, err)
					break
				}
			}
			//k.Pc[key] = pc
			// 等待 生产者产生对应数据 然后消费
			go func() {
				for {
					select {
					case msg := <-pc.Messages():
						err = h.CustomFn(msg)
						if err != nil {
							h.ErrorFn(h, err)
							break
						}
					case err = <-pc.Errors():
						if err != nil {
							h.ErrorFn(h, err)
							break
						}
					}
				}
			}()
		}
	}
}
