package websocket

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/liujingkaiai/x-socket/xnet"
)

type Wire struct {
}

func (w *Wire) GetHeadLen() uint32 {
	// Len + msgId  uint32=4字节
	return 4
}

func (w *Wire) Pack(msg xnet.Message) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})
	//
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (w *Wire) Unpack(data []byte) (xnet.Message, error) {
	dataBuff := bytes.NewReader(data)
	msg := Message{}
	//获取msgid
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	// 判断是否超出长度
	if optsions.MaxPackageSize > 0 && msg.DataLen > optsions.MaxMessageSize {
		return nil, errors.New("too large msg data recv")
	}
	//获取msgid
	msg.SetData(data[4:])
	return &msg, nil
}
