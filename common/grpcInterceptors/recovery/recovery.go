package recovery

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"runtime"
)

// RecoveryHandlerFunc 从panic `p` 中 recovery 并且返回一个error的方法
type RecoveryHandlerFunc func(p any) (err error)

// RecoveryHandlerFuncContext 在RecoveryHandlerFunc的基础上，加上了ctx，用于提取请求作用域的元数据和上下文值
type RecoveryHandlerFuncContext func(ctx context.Context, p any) (err error)

type PanicError struct {
	Panic any
	Stack []string
}

// UnaryServerInterceptor 一元拦截器
// 传入一个配置列表，配置包括一个recover方法，加载后返回一个grpc的拦截器，加入panic捕获，当捕获到panic则调用自定义的recover方法
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	// 加载配置
	o := evaluateOptions(opts)

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(ctx, r, o.recoveryHandlerFunc)
			}
		}()
		return handler(ctx, req)
	}
}

// TODO 流拦截器 有需要用到流RPC再来实现

// recoverFrom 恢复方法，若有自定义方法则直接使用，否则用默认的方法
func recoverFrom(ctx context.Context, p any, r RecoveryHandlerFuncContext) error {
	// 如果自定义了方法，则调用恢复方法
	if r != nil {
		return r(ctx, p)
	}
	// 否则使用默认的
	return DefaultRecovery(p)
}

// 默认的一个处理方法，捕捉了panic构造成了一个结构体
func DefaultRecovery(p any) *PanicError {
	stack := make([]string, 64<<10)
	//stack = stack[:runtime.Stack(stack, false)]
	//stack = stack[:runtime.Stack(stack, false)]
	// 获得调用栈
	for i := 3; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = append(stack, fmt.Sprint(pc, file, line))
	}
	return &PanicError{Panic: p, Stack: stack}
}

// Error 实现error接口
func (e *PanicError) Error() string {
	return fmt.Sprintf("panic caught :%v\n", e.Panic)
}
