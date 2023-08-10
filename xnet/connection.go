package xnet

type Connection interface {
	Start()
	StartReader()
	StartWriter()
	Stop()
	GetConnId() uint32
}
