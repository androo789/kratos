package log

import "testing"

func TestStdLogger(t *testing.T) {
	logger := DefaultLogger
	//结合源码这一部分就看懂了
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)

	//直接使用log还是麻烦，使用helper定义好的info等等。这里居然还有err返回值
	_ = logger.Log(LevelInfo, "msg", "test debug")
	_ = logger.Log(LevelInfo, "msg", "test info")
	_ = logger.Log(LevelInfo, "msg", "test warn")
	_ = logger.Log(LevelInfo, "msg", "test error")
	_ = logger.Log(LevelDebug, "singular") //这行就是奇数个元素，所以补充一个元素

	logger2 := DefaultLogger
	_ = logger2.Log(LevelDebug) //~这行效果就是什么都不输出
}
