package main

import (
	"fmt"
	"log"

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

	// 启动一个协程用于接收和处理服务器发送的消息
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("读取消息失败：", err)
				return
			}
			fmt.Printf("收到消息：%s\n", msg)
		}
	}()

	// 循环发送消息到服务器
	for {

		conn.WriteMessage(websocket.PingMessage, []byte("ping"))

		var msg string
		fmt.Print("输入消息：")
		fmt.Scanln(&msg)
		wire := ws.Wire{}

		bt, err := wire.Pack(ws.NewMessage(0, []byte(msg)))
		if err != nil {
			fmt.Println(err)
		}
		err = conn.WriteMessage(websocket.BinaryMessage, bt)
		if err != nil {
			log.Println("发送消息失败：", err)
			return
		}
	}
}
