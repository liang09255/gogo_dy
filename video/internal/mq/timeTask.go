package mq

import "time"

// 定时处理任务
func batchTickerTask() {
	// 10s定期将map的数值消费一次
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timer.C:
			// 合并点赞任务
			go TimeFavoriteTask()
			// 合并评论任务
			go TimeCommentTask()
		}
	}
}
