package xnet

type Wrie interface {
	GetHeadLen() uint32
	Pack(Message) ([]byte, error)
	Unpack([]byte) (Message, error)
}

type CodeC interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}
