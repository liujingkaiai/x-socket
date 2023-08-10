package websocket

import (
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liujingkaiai/x-socket/xnet"
)

type Connection struct {
	id      uint32
	conn    *websocket.Conn
	isClose bool
	closeC  chan struct{}
	msg     chan []byte

	wsServer xnet.Server
}

func NewConnection(id uint32, conn *websocket.Conn, server xnet.Server) xnet.Connection {
	return &Connection{
		id:       id,
		conn:     conn,
		closeC:   make(chan struct{}),
		wsServer: server,
		msg:      make(chan []byte),
	}
}

func (c *Connection) StartReader() {
	defer c.Stop()
	//阻塞读取业务

	//设置消息最大size
	c.conn.SetReadLimit(int64(c.wsServer.GetServerOpt().MaxMessageSize))
	c.conn.SetPongHandler(func(appData string) error {
		fmt.Println("sdsd")
		c.conn.SetReadDeadline(time.Now().Add(c.wsServer.GetServerOpt().PongWait))
		return nil
	})

	for {
		t, data, err := c.conn.ReadMessage()
		if err != nil {
			if err == io.EOF {
				c.Stop()
			}
			return
		}

		switch t {
		case websocket.PingMessage:
			fmt.Println(2323)
			c.conn.WriteMessage(websocket.PongMessage, nil)
		case websocket.BinaryMessage:
			wire := Wire{}
			msg, err := wire.Unpack(data)
			//如果消息解包失败，不处理
			if err != nil {
				fmt.Println("pack msg err:", err)
				return
			}
			req := NewRequest(msg, c)
			// 默认server 关闭不使用连接池
			if c.wsServer.GetStates() == xnet.Close {
				go c.wsServer.GetDispatcher().Handle(req)
			} else {
				go c.wsServer.GetDispatcher().SentToTaskQueue(req)
			}

		case websocket.TextMessage:
			c.conn.WriteMessage(t, data)
		}

	}
}

func (c *Connection) StartWriter() {

}

func (c *Connection) Start() {
	c.wsServer.GetConnManager().Add(c)
	go c.StartReader()
	go c.StartWriter()
	select {}
}

func (c *Connection) Stop() {
	fmt.Println("remote addr:", c.conn.RemoteAddr(), " is closed")
	if c.isClose {
		return
	}
	c.isClose = true
	c.conn.Close()
	c.wsServer.GetConnManager().Remove(c)
	close(c.closeC)
}

func (c *Connection) GetConnId() uint32 {
	return c.id
}
