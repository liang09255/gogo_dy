package router

import (
	"common/ggConfig"
	"common/ggDiscovery"
	"common/ggIDL/video"
	"common/ggIP"
	"common/ggLog"
	"common/grpcInterceptors/recovery"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net"
	"video/pkg/service"
)

func StartGrpc() *grpc.Server {
	videoServerConfig := ggConfig.Config.VideoServer

	// 接口缓存 拦截器
	// recovery
	var recoveryFunc recovery.RecoveryHandlerFunc
	// 做成配置项主要是为了扩展，比如多个微服务的log不是同一个实例,或者不需要log记录，这里用闭包处理就不需要每次都传值
	recoveryFunc = func(p any) (err error) {
		panicErr := recovery.DefaultRecovery(p)
		// 记录panic错误 - 根据需要修改显示捕获到的错误的格式
		// 调用panicErr.Stack是捕捉了调用栈 - 可以考虑输出格式，现在只输出了捕捉的panic
		ggLog.Error(panicErr.Error())
		return panicErr
	}
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(recoveryFunc),
	}
	interceptor := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			otelgrpc.UnaryServerInterceptor(),
			//interceptor.New().CacheInterceptor(), TODO 待定
			recovery.UnaryServerInterceptor(opts...),
		))

	// 创建grpc服务端
	g := grpc.NewServer(interceptor)
	video.RegisterVideoServiceServer(g, service.New())

	// 注册客户端

	lis, err := net.Listen("tcp", videoServerConfig.Addr)
	if err != nil {
		ggLog.Fatalf("端口监听失败:%s", err.Error())
	}

	go func() {
		ggLog.Infof("grpc server started at %s", videoServerConfig.Addr)
		if err = g.Serve(lis); err != nil {
			ggLog.Fatalf("grpc server started error: %s", err.Error())
		}
	}()

	return g
}

func RegisterEtcdServer() {
	etcdConfig := ggConfig.Config.Etcd
	videoServerConfig := ggConfig.Config.VideoServer
	// 实现grpc接口，拓展一下，使得可以识别etcd的链接
	// 创建了一个Resolver
	// Resolver实现了Build()和Scheme()，所以Resolver实现了Builder接口
	etcdRegister := ggDiscovery.NewResolver(etcdConfig.Addrs)
	resolver.Register(etcdRegister)

	videoServerConfig.Addr = ggIP.GetIP() + ":" + videoServerConfig.Port

	// 构建grpc服务
	info := ggDiscovery.Server{
		Name:    videoServerConfig.Name,
		Addr:    videoServerConfig.Addr,
		Version: "v1.0.0",
		Weight:  1,
	}

	// grpc服务注册
	r := ggDiscovery.NewRegister(etcdConfig.Addrs)
	_, err := r.Register(info, 2)
	if err != nil {
		ggLog.Fatalf("grpc server register error: %s", err.Error())
	}
	ggLog.Debugf("grpc server register success")
}
