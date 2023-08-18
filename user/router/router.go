package router

import (
	"common/discovery"
	"common/idl/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"user/config"
	user_service_v1 "user/pkg/service/user.service.v1"
)

type gRPCConfig struct {
	// grpc启动服务的地址
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGrpc() *grpc.Server {
	c := gRPCConfig{
		Addr: config.C.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			user.RegisterUserServiceServer(g, user_service_v1.New())
		}}

	// 接口缓存 拦截器
	s := grpc.NewServer(
	//grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	//	otelgrpc.UnaryServerInterceptor(),
	//	interceptor.New().CacheInterceptor(),
	//)),
	)

	// 注册服务
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", config.C.GC.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	// 用协程，否则main中会卡住
	go func() {
		log.Printf("grpc server started as: %s \n", config.C.GC.Addr)
		err = s.Serve(lis)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()
	return s
}

func RegisterEtcdServer() {
	// 实现grpc接口，拓展一下，使得可以识别etcd的链接
	// 创建了一个Resolver
	// Resolver实现了Build()和Scheme()，所以Resolver实现了Builder接口
	etcdRegister := discovery.NewResolver(config.C.EC.Addrs)
	resolver.Register(etcdRegister)

	// 构建grpc服务
	info := discovery.Server{
		Name:    config.C.GC.Name,
		Addr:    config.C.GC.Addr,
		Version: config.C.GC.Version,
		Weight:  config.C.GC.Weight,
	}

	// grpc服务注册
	r := discovery.NewRegister(config.C.EC.Addrs)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
