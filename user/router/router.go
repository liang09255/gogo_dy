package router

import (
	"common/ggConfig"
	"common/ggDiscovery"
	"common/ggIDL/user"
	"common/ggLog"
	"user/pkg/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net"
)

func StartGrpc() *grpc.Server {
	userServerConfig := ggConfig.Config.UserServer

	// 接口缓存 拦截器
	g := grpc.NewServer()
	user.RegisterUserServer(g, service.New())

	lis, err := net.Listen("tcp", userServerConfig.Addr)
	if err != nil {
		ggLog.Fatalf("端口监听失败: %s", err.Error())
	}

	go func() {
		ggLog.Infof("grpc server started at %s", userServerConfig.Addr)
		if err = g.Serve(lis); err != nil {
			ggLog.Fatalf("grpc server started error: %s", err.Error())
		}
	}()

	return g
}

func RegisterEtcdServer() {
	etcdConfig := ggConfig.Config.Etcd
	userServerConfig := ggConfig.Config.UserServer
	// 实现grpc接口，拓展一下，使得可以识别etcd的链接
	// 创建了一个Resolver
	// Resolver实现了Build()和Scheme()，所以Resolver实现了Builder接口
	etcdRegister := ggDiscovery.NewResolver(etcdConfig.Addrs)
	resolver.Register(etcdRegister)

	// 构建grpc服务
	info := ggDiscovery.Server{
		Name:    userServerConfig.Name,
		Addr:    userServerConfig.Addr,
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
