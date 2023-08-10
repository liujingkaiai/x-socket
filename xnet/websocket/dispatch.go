package websocket

import (
	"fmt"

	"github.com/liujingkaiai/x-socket/xnet"
)

type Disptach struct {
	routers  map[uint8]xnet.HandleFunc
	requests map[int]chan xnet.Request
}

func NewDispatch() xnet.Dispatcher {
	return &Disptach{
		routers:  make(map[uint8]xnet.HandleFunc),
		requests: make(map[int]chan xnet.Request),
	}
}

// 添加路由
func (d *Disptach) HandleFunc(id uint8, handler xnet.HandleFunc) {
	if _, ok := d.routers[id]; ok {
		panic(fmt.Sprintf("id:%d router already exists", id))
	}
	d.routers[id] = handler
}

func (d *Disptach) Handle(req xnet.Request) {
	//路由不合法
	if handler, ok := d.routers[uint8(req.GetMsgID())]; !ok {
		return
	} else {
		handler(req)
	}
}

func (d *Disptach) StartWookerPool(poosize uint32, queueLent uint32) {
	if poosize == 1 {
		return
	}
	for i := 0; i < int(poosize); i++ {
		d.requests[i] = make(chan xnet.Request, queueLent)
		go d.startOneWorker(i, d.requests[i])
	}

}

func (d *Disptach) startOneWorker(id int, ch chan xnet.Request) {
	fmt.Println("Worker ID = ", id, " is started....")
	for {
		select {
		case req := <-ch:
			d.Handle(req)
		}
	}
}

func (d *Disptach) SentToTaskQueue(req xnet.Request) {
	l := int(req.GetMsgID()) % len(d.requests)
	d.requests[l] <- req
}
