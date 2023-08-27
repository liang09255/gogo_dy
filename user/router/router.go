package router

import (
	"common/ggConfig"
	"common/ggDiscovery"
	"common/ggIDL/relation"
	"common/ggIDL/user"
	"common/ggLog"
	"common/grpcInterceptors/recovery"
	"user/pkg/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net"
)

func StartGrpc() *grpc.Server {
	userServerConfig := ggConfig.Config.UserServer

	// 接口缓存 拦截器
	// recovery中间件
	var recoveryFunc recovery.RecoveryHandlerFunc
	// 做成配置项主要是为了扩展，比如多个微服务的log不是同一个实例,或者不需要log记录，这里用闭包处理就不需要每次都传值
	// 自定义的恢复方法
	recoveryFunc = func(p any) (err error) {
		panicErr := recovery.DefaultRecovery(p)
		// 记录panic错误 - 根据需要修改显示捕获到的错误的格式
		ggLog.Panic(panicErr.Error())
		return panicErr
	}
	// 方法作为配置传入
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(recoveryFunc),
	}
	interceptor := grpc.UnaryInterceptor(
		recovery.UnaryServerInterceptor(opts...))

	// 创建grpc服务端
	g := grpc.NewServer(interceptor)


	relation.RegisterRelationServer(g, service.New2())

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
