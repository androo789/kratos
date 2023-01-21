module github.com/go-kratos/kratos/contrib/log/zap/v2

go 1.16

require (
	//这里是对应了git的tag
	github.com/go-kratos/kratos/v2 v2.5.3
	go.uber.org/zap v1.23.0
)
//但是我fork再clone之后，我这个git里面已经没有tag了额，怎么还能生效
//其实用了replace中的版本号，上面的v2.5.3根本没有意义
//replace中又没有版本号，所以就是当前的代码就完事了
replace github.com/go-kratos/kratos/v2 => ../../../
