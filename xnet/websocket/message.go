package websocket

type Message struct {
	Id      uint32
	Len     uint32
	DataLen uint32

	Data []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:   id,
		Data: data,
		Len:  uint32(len(data)),
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.Len
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(len uint32) {
	m.Len = len
}
