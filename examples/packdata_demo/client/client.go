package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	ws "github.com/liujingkaiai/x-socket/xnet/websocket"
)

func main() {
	// 连接到WebSocket服务器的URL
	url := "ws://localhost:7772/ws"

	// 建立WebSocket连接
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("连接WebSocket服务器失败：", err)
	}
	defer conn.Close()

	fmt.Println("已连接到WebSocket服务器")

	wire := ws.Wire{}

	// 启动一个协程用于接收和处理服务器发送的消息
	go func() {
		for {
			t, msg, err := conn.ReadMessage()
			if t == websocket.BinaryMessage {
				_msg, _err := wire.Unpack(msg)
				if err != nil {
					fmt.Println("unpack msg err:", _err)
					continue
				}
				fmt.Printf("接收<-------消息ID:%d 内容:%s \n", _msg.GetMsgId(), _msg.GetData())
			}
		}
	}()

	// 循环发送消息到服务器
	for {

		bt, err := wire.Pack(ws.NewMessage(0, []byte("ping")))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("发送消息===>%s\n", bt)
		err = conn.WriteMessage(websocket.BinaryMessage, bt)
		if err != nil {
			log.Println("发送消息失败：", err)
			return
		}

		time.Sleep(time.Second)
	}
}
