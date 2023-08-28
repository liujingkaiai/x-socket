package websocket

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liujingkaiai/x-socket/xnet"
)

type Connection struct {
	id         string
	conn       *websocket.Conn
	isClose    bool
	closeC     chan struct{}
	msg        chan []byte
	m          *sync.RWMutex
	attributes map[string]any
	wsServer   xnet.Server
}

func NewConnection(id uint32, conn *websocket.Conn, server xnet.Server) xnet.Connection {
	return &Connection{
		id:         strconv.Itoa(int(id)),
		conn:       conn,
		closeC:     make(chan struct{}),
		wsServer:   server,
		msg:        make(chan []byte),
		attributes: make(map[string]any, 0),
	}
}

func (c *Connection) StartReader() {
	defer c.Stop()
	//阻塞读取业务

	//设置消息最大size
	c.conn.SetReadLimit(int64(c.wsServer.GetServerOpt().MaxMessageSize))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		c.conn.SetReadDeadline(time.Now().Add(c.wsServer.GetServerOpt().PongWait))
		return nil
	})

	//如果存在鉴权
	if c.wsServer.GetAcceptFunc() != nil {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("accept read msg err:", err)
			return
		}
		acceptFunc := c.wsServer.GetAcceptFunc()
		suc, IdCreater := acceptFunc(data)
		//鉴权失败
		if !suc {
			c.conn.Close()
			return
		}
		c.id = IdCreater.Id()
	}

	c.wsServer.GetConnManager().Add(c)

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
			c.conn.WriteMessage(websocket.PongMessage, []byte{})
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

	for {
		select {
		case msg := <-c.msg:
			if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				fmt.Println("Send data error,", err, "Conn Writer exit")
				return
			}
		case <-c.closeC:
			return
		}
	}
}

func (c *Connection) Start() {
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
	close(c.msg)
}

func (c *Connection) GetConnId() string {
	return c.id
}

func (c *Connection) SendMsg(id uint32, data []byte) error {
	wire := Wire{}
	m, err := wire.Pack(NewMessage(id, data))
	if err != nil {
		// 打包失败
		return err
	}
	c.msg <- m
	return nil
}

func (c *Connection) SetAttribute(key string, val any) {
	c.m.Lock()
	c.attributes[key] = val
	c.m.Unlock()
}

func (c *Connection) GetAttribute(key string) (any, bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	val, ok := c.attributes[key]
	return val, ok
}
