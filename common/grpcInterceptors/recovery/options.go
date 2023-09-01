package recovery

import "context"

var (
	defaultOptions = &options{
		recoveryHandlerFunc: nil,
	}
)

type (
	options struct {
		recoveryHandlerFunc RecoveryHandlerFuncContext
	}
	Option func(*options)
)

func defaultRecoveryHandler(c context.Context, log any, err interface{}, stack []byte) {

}

// 加载配置
func evaluateOptions(opts []Option) *options {
	cfg := &options{
		recoveryHandlerFunc: nil,
	}
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

// WithRecoveryHandler 自定义从panic中recovery的方法
func WithRecoveryHandler(f RecoveryHandlerFunc) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = RecoveryHandlerFuncContext(func(ctx context.Context, p any) error {
			return f(p)
		})
	}
}

func WithRecoveryHandlerContext(f RecoveryHandlerFuncContext) Option {
	return func(o *options) {
		o.recoveryHandlerFunc = f
	}
}
