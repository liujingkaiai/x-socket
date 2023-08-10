package xnet

type Wrie interface {
	GetHeadLen() uint32
	Pack(Message) ([]byte, error)
	Unpack([]byte) (Message, error)
}
