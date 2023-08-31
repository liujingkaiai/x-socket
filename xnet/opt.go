package xnet

import "time"

type ServerOption struct {
	MachineID             int64
	ReadBufferSize        uint32
	WriteBufferSize       uint32
	PoolSize              uint32
	MaxConnectionNum      uint32
	MaxMessageSize        uint32
	PongWait              time.Duration
	MaxPackageSize        int64
	CompressionExtensions bool
}
