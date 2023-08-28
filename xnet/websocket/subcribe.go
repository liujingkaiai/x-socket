package websocket

import "sync"

type Subscribe struct {
	name       string
	mux        *sync.RWMutex
	connetions map[string]struct{}
}

func NewSubscribe(name string) *Subscribe {
	return &Subscribe{
		name: name,
	}
}

func (s *Subscribe) Tag() string {
	return s.name
}
