package middleware

import (
	"context"
)

// Handler defines the handler invoked by Middleware.
// 每个请求，都是一个rep，一个resp
type Handler func(ctx context.Context, req interface{}) (interface{}, error)

// Middleware is HTTP/gRPC transport middleware.
// 中间件函数的参数就是请求，再把传递给下一个请求
// 这种设计也是妙啊，细品
type Middleware func(Handler) Handler

// Chain returns a Middleware that specifies the chained handler for endpoint.
func Chain(m ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}
