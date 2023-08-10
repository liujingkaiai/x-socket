package xnet

import "sync"

type States uint8

const (
	Start States = iota
	Close
)

type Server interface {
	Service
	//获取配置
	GetServerOpt() *ServerOption
	//获取管理连接器
	GetConnManager() *ConnManager
	//鉴权认证
	SetAcceptFunc(func([]byte) (bool, IdCreater))
	//获取鉴权函数
	GetAcceptFunc() func([]byte) (bool, IdCreater)
	GetPool() *sync.Pool
	//设置路由管理
	SetDistpatcher(Dispatcher)
	//获取路由管理
	GetDispatcher() Dispatcher
	//设置pool size
	SetPoolSize(uint32)
	//获取服务运行状态
	GetStates() States
	//
}

type IdCreater interface {
	Id() string
}

type Service interface {
	//获取ID
	GetID() int64
	//启动
	Start()
	//运行
	Serve()
	//停止
	Stop()
}
