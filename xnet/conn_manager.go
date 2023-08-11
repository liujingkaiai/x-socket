package xnet

import (
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	connections map[string]Connection
	m           sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[string]Connection),
	}
}

func (cm *ConnManager) Add(conn Connection) {
	cm.m.Lock()
	defer cm.m.Unlock()
	cm.connections[conn.GetConnId()] = conn
	fmt.Println("connection add to ConnManager successfully:conn num=", cm.Len())
}

// 删除链接
func (cm *ConnManager) Remove(conn Connection) {
	cm.m.Lock()
	defer cm.m.Unlock()
	delete(cm.connections, conn.GetConnId())
	fmt.Println("connID = ", conn.GetConnId(), " remove from ConnManager sucessfully:conn num=", cm.Len())
}

// 根据Conn获取链接
func (cm *ConnManager) Get(connID string) (Connection, error) {
	cm.m.RLock()
	defer cm.m.RUnlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not FOUND!")
}

// 个数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// 清除所有链接
func (cm *ConnManager) Clear() {
	cm.m.Lock()
	defer cm.m.Unlock()

	//删除connection 并停止connection
	for connID, conn := range cm.connections {
		conn.Stop()
		//删除
		delete(cm.connections, connID)
	}
	fmt.Println("Clear All Connection succ! conn len=", cm.Len())
}
