package xnet

// 路由
type Dispatcher interface {
	//添加路由
	HandleFunc(uint8, HandleFunc)
	// 执行函数
	Handle(Request)
	//开启协程池
	StartWookerPool(poosize uint32, queueLent uint32)
	SentToTaskQueue(Request)
}
