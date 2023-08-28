package xnet

type Connection interface {
	Start()
	StartReader()
	StartWriter()
	Stop()
	GetConnId() string
	//将消息打包发送给写writer
	SendMsg(uint32, []byte) error
	SetAttribute(key string, val any)
	GetAttribute(key string) (any, bool)
}
