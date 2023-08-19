package shutDown

import (
	"context"
	"log"
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
		log.Println("关闭服务")
	}
	log.Printf("%s stop success...", srvName)
}
