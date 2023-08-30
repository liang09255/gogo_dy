package ggShutDown

import (
	"common/ggLog"
	"context"
	"time"
)

func ShutDown(srvName string, addr string, stop func()) {
	//关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if stop != nil {
		stop()
	}
	select {
	case <-ctx.Done():
		ggLog.Info("关闭服务")
	}
	ggLog.Infof("%s stop success", srvName)
}
