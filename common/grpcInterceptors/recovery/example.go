package recovery

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var customFunc RecoveryHandlerFunc

// 主要是传入一个处理panic的方法，错误捕捉和恢复就比较容易，customFunc是这样例子的处理方法
func Example_init() {
	customFunc = func(p any) (err error) {
		return status.Errorf(codes.Unknown, "panic trigger : %v", p)
	}

	opts := []Option{
		WithRecoveryHandler(customFunc),
	}

	_ = grpc.NewServer(
		grpc.UnaryInterceptor(
			UnaryServerInterceptor(opts...),
		),
	)

}
