package xnet

type Request interface {
	//获取消息id
	GetMsgID() uint32
	//获取消息内容
	GetData() []byte
	//获取对应的链接
	GetConnection() Connection
}
