package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/liujingkaiai/x-socket/xnet"
	"github.com/liujingkaiai/x-socket/xnet/websocket"
)

const (
	PING uint8 = iota + 0
)

// 演示路由+解包用法
type PingHandler struct {
}

func (p *PingHandler) Ping(req xnet.Request) {
	fmt.Printf("rec %v %s\n", req.GetMsgID(), req.GetData())
	req.GetConnection().SendMsg(req.GetMsgID(), []byte("pong"))
}

func main() {
	//实例化wsserver
	wsserver := websocket.Default()
	wsserver.SetPoolSize(10)
	//开启任务处理
	go wsserver.Start()
	handler := &PingHandler{}
	wsserver.GetDispatcher().HandleFunc(PING, handler.Ping)
	http.HandleFunc("/ws", wsserver.ServeWs)
	fmt.Println("server starting at 7772")
	log.Fatal(http.ListenAndServe(":7772", nil))
}
