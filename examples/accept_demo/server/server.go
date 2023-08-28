package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/liujingkaiai/x-socket/xnet"
	"github.com/liujingkaiai/x-socket/xnet/websocket"
)

const (
	Send uint8 = iota
	Recive
)

// id 生成器
type IdCreater struct {
	data []byte
}

func newIdCreate(data []byte) *IdCreater {
	return &IdCreater{
		data: data,
	}
}

func (id *IdCreater) Id() string {
	return string(id.data)
}

type ChatHandler struct{}

func (c *ChatHandler) ChatTo(req xnet.Request) {

}

func main() {
	//实例化wsserver
	wsserver := websocket.Default()
	//设置任务处理池
	wsserver.SetPoolSize(10)
	//设置鉴权握手函数
	wsserver.SetAcceptFunc(func(b []byte) (bool, xnet.IdCreater) {
		return true, newIdCreate(b)
	})

	//开启任务处理
	go func() {
		wsserver.Start()
		handler := &ChatHandler{}
		wsserver.GetDispatcher().HandleFunc(Send, handler.ChatTo)
		http.HandleFunc("/ws", wsserver.ServeWs)
		//fmt.Println("server starting at 7772")
		log.Fatal(http.ListenAndServe(":7772", nil))
	}()

	for {
		var id string
		var msg string
		fmt.Println("请输入你要发送的id")
		fmt.Scanf("%s\n", &id)
		fmt.Println("请输入你要发送的消息")
		fmt.Scanf("%s\n", &msg)

		fmt.Printf("正在发送:%s => 用户ID:%s\n", msg, id)
		if err := wsserver.ChatWith(id, 0, []byte(msg)); err != nil {
			fmt.Println("发送失败 err:", err)
		}
	}
}
