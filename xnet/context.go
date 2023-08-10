package xnet

import "context"

type Context interface {
	//获取基本的ctx
	BaseContext() context.Context
	//实现请求
	Request
	//设置request 参数
	SetRequest(Request)
	//退出中间件
	Abort()
	//设置处理函数
	SetHandelr(handlers []HandleFunc)
	//执行函数
	Next() error
	//获取退出中间件index
	GetAbortIndex() uint8
	//重置函数
	Reset()
}
