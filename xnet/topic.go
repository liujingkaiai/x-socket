package xnet

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrTagRecordNotFound = errors.New("tag not exists")
)

type Topic struct {
	channels    sync.Map
	connManager *ConnManager
	m           sync.RWMutex
}

func NewTopic(manager *ConnManager) *Topic {
	return &Topic{
		channels:    sync.Map{},
		connManager: manager,
	}
}

func (t *Topic) AddTag(tag string) chan Message {
	var sub *subscriber
	t.m.Lock()
	old, ok := t.channels.Load(tag)
	if !ok {
		sub = &subscriber{
			msgChan:   make(chan Message, 0),
			closeChan: make(chan struct{}, 0),
		}
	} else {
		sub = old.(*subscriber)
	}
	t.m.Unlock()
	go func(sub *subscriber, connManager *ConnManager) {
		for {
			select {
			case msg := <-sub.msgChan:
				sub.conns.Range(func(k, v any) bool {
					id, ok := k.(string)
					if !ok {
						return true
					}
					conn, err := connManager.Get(id)
					if err == nil {
						conn.SendBufferMesg(msg.GetMsgId(), msg.GetData())
					}
					return true
				})
			case <-sub.closeChan:
				sub.Close()
			}
		}
	}(sub, t.connManager)
	return sub.msgChan
}

func (t *Topic) AddSubscribe(tag string, conn Connection) error {
	sub, ok := t.getsubscribe(tag)
	if !ok {
		return fmt.Errorf("error:%w tag:%s ", ErrTagRecordNotFound, tag)
	}
	sub.Add(conn)
	return nil
}

func (t *Topic) getsubscribe(tag string) (*subscriber, bool) {
	sub, ok := t.channels.Load(tag)
	if !ok {
		return nil, ok
	}
	return sub.(*subscriber), ok
}

type subscriber struct {
	msgChan   chan Message
	closeChan chan struct{}
	conns     sync.Map
	m         sync.RWMutex
}

func (s *subscriber) SetChan(ch chan Message) {
	s.msgChan = ch
}

func (s *subscriber) Add(conn Connection) {
	s.conns.Store(conn.GetConnId(), struct{}{})
}

func (s *subscriber) Close() {
	close(s.msgChan)
	s.conns.Range(func(k, v any) bool {
		s.conns.Delete(k)
		return true
	})
}
