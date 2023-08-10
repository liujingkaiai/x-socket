package websocket

import "github.com/liujingkaiai/x-socket/xnet"

type Request struct {
	msg  xnet.Message
	conn *Connection
}

func NewRequest(msg xnet.Message, conn *Connection) xnet.Request {
	return &Request{
		msg:  msg,
		conn: conn,
	}
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetConnection() xnet.Connection {
	return r.conn
}
