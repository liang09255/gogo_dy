package router

import (
	"chat/pkg/service"
	"common/ggConfig"
	"common/ggDiscovery"
	"common/ggIDL/chat"
	"common/ggIP"
	"common/ggLog"
	"common/grpcInterceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net"
)

func StartGrpc() *grpc.Server {
	chatServerConfig := ggConfig.Config.ChatServer

	var recoveryFunc recovery.RecoveryHandlerFunc
	recoveryFunc = func(p any) (err error) {
		panicErr := recovery.DefaultRecovery(p)
		ggLog.Error(panicErr.Error())
		return panicErr
	}
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(recoveryFunc),
	}
	interceptor := grpc.UnaryInterceptor(recovery.UnaryServerInterceptor(opts...))

	// 接口缓存 拦截器
	g := grpc.NewServer(interceptor)
	chat.RegisterChatServer(g, service.New())

	lis, err := net.Listen("tcp", chatServerConfig.Addr)
	if err != nil {
		ggLog.Fatalf("端口监听失败: %s", err.Error())
	}

	go func() {
		ggLog.Infof("grpc server started at %s", chatServerConfig.Addr)
		if err = g.Serve(lis); err != nil {
			ggLog.Fatalf("grpc server started error: %s", err.Error())
		}
	}()

	return g
}

func RegisterEtcdServer() {
	etcdConfig := ggConfig.Config.Etcd
	chatServiceConfig := ggConfig.Config.ChatServer
	// 实现grpc接口，拓展一下，使得可以识别etcd的链接
	// 创建了一个Resolver
	// Resolver实现了Build()和Scheme()，所以Resolver实现了Builder接口
	etcdRegister := ggDiscovery.NewResolver(etcdConfig.Addrs)
	resolver.Register(etcdRegister)

	chatServiceConfig.Addr = ggIP.GetIP() + ":" + chatServiceConfig.Port

	// 构建grpc服务
	info := ggDiscovery.Server{
		Name:    chatServiceConfig.Name,
		Addr:    chatServiceConfig.Addr,
		Version: "v1.0.0",
		Weight:  1,
	}

	// grpc服务注册
	r := ggDiscovery.NewRegister(etcdConfig.Addrs)
	_, err := r.Register(info, 2)
	if err != nil {
		ggLog.Fatalf("grpc server register error: %s", err.Error())
	}
}
