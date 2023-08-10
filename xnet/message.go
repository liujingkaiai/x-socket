package xnet

type Message interface {
	GetMsgId() uint32
	GetMsgLen() uint32
	GetData() []byte
}
