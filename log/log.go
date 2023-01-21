package log

import (
	"context"
	"log"
)

// DefaultLogger is default logger.
var DefaultLogger = NewStdLogger(log.Writer())

// Logger is a logger interface.
type Logger interface {
	Log(level Level, keyvals ...interface{}) error
}

type logger struct {
	logger    Logger
	prefix    []interface{}
	hasValuer bool
	ctx       context.Context
}

func (c *logger) Log(level Level, keyvals ...interface{}) error {
	kvs := make([]interface{}, 0, len(c.prefix)+len(keyvals))
	kvs = append(kvs, c.prefix...)
	if c.hasValuer {
		bindValues(c.ctx, kvs)
	}
	kvs = append(kvs, keyvals...)
	//~真实打印的时候又用到了之前传入的l
	if err := c.logger.Log(level, kvs...); err != nil {
		return err
	}
	return nil
}

// With with logger fields.
// ~返回的是这里的小logger，等于是在传入的l基础上再封装更多东西
func With(l Logger, kv ...interface{}) Logger {
	//~第一次进来肯定不是一个logger，就初始化。with可以多次调用
	c, ok := l.(*logger)
	if !ok {
		//默认的ctx是background，暂时还看不出来有什么用
		return &logger{logger: l, prefix: kv, hasValuer: containsValuer(kv), ctx: context.Background()}
	}
	kvs := make([]interface{}, 0, len(c.prefix)+len(kv))
	kvs = append(kvs, c.prefix...)
	//~偶数位置现在是一个valuer，后面会被替换为具体的值，比如代码位置，时间等等
	//~这个设计秒啊
	kvs = append(kvs, kv...)
	return &logger{
		logger:    c.logger,
		prefix:    kvs,
		hasValuer: containsValuer(kvs),
		ctx:       c.ctx,
	}
}

// WithContext returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
// ~上面没有context，这函数作用仅仅是传入context。别的不变。因为context在不同请求中是不同的
func WithContext(ctx context.Context, l Logger) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{logger: l, ctx: ctx}
	}
	return &logger{
		logger:    c.logger,
		prefix:    c.prefix,
		hasValuer: c.hasValuer,
		ctx:       ctx,
	}
}
