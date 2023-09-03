package mq

import (
	"common/kafkax"
	"encoding/json"
	"video/pkg/constant"
)

type FavoriteMessage struct {
	Uid int64
	Vid int64
	// 1为新增，2为删除
	Method int64
}

// 生产者，插入数据,每次调用都是直接添加方法
func AddFavoriteMessage(m *FavoriteMessage) error {
	// 添加一条插入点赞记录消息
	// 消息需要记录，xx人对xx视频进行了点赞
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	msg := &kafkax.ProducerMsg{
		Topic: constant.FavoriteTopic,
		Data:  data,
	}
	err = Kafka.SendMessage(msg)
	return err
}
