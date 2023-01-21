package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	// DefaultCaller is a Valuer that returns the file and line.
	//~我记得4可能是标识仅仅获取当前代码位置，不再向上追溯
	DefaultCaller = Caller(4)

	// DefaultTimestamp is a Valuer that returns the current wallclock time.
	// ~默认的时间value
	DefaultTimestamp = Timestamp(time.RFC3339)
)

// Valuer is returns a log value.
type Valuer func(ctx context.Context) interface{}

// Value return the function value.
func Value(ctx context.Context, v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

// Caller returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		idx := strings.LastIndexByte(file, '/')
		if idx == -1 {
			return file[idx+1:] + ":" + strconv.Itoa(line)
		}
		idx = strings.LastIndexByte(file[:idx], '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// Timestamp returns a timestamp Valuer with a custom time format.
// ~传入的context没有使用
func Timestamp(layout string) Valuer {
	return func(context.Context) interface{} {
		return time.Now().Format(layout)
	}
}

func bindValues(ctx context.Context, keyvals []interface{}) {
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			//将偶数位置替换为实际的值
			keyvals[i] = v(ctx)
		}
	}
}

//如果任意一个偶数位置是valuer，那么就需要执行valuer对应函数
//奇数偶数，从1开始数，跟一般计算机世界不一样
func containsValuer(keyvals []interface{}) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}
	return false
}
