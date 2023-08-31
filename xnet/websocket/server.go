package websocket

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/liujingkaiai/x-socket/xnet"
)

var ErrConnNotFound = errors.New("conn not found in server !")

type empty struct{}

var optsions xnet.ServerOption

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var defaultServerOpt = xnet.ServerOption{
	MachineID:        1,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	PoolSize:         0,
	MaxConnectionNum: 2048,
	MaxMessageSize:   512,
	PongWait:         time.Minute,
	MaxPackageSize:   1024,
}

type Server struct {
	id              int64
	opts            *xnet.ServerOption
	upgrader        websocket.Upgrader
	shakeFaieldFunc func(w http.ResponseWriter)
	connManager     *xnet.ConnManager
	idFacotry       *idManager
	acceptFucn      func([]byte) (bool, xnet.IdCreater)
	pool            *sync.Pool
	dispatch        xnet.Dispatcher
	states          xnet.States
}

func NewServer(opt *xnet.ServerOption) *Server {
	optsions = *opt

	return &Server{
		id:   opt.MachineID,
		opts: opt,
		upgrader: websocket.Upgrader{
			ReadBufferSize:    int(opt.ReadBufferSize),
			WriteBufferSize:   int(opt.WriteBufferSize),
			EnableCompression: opt.CompressionExtensions,
		},
		dispatch:    NewDispatch(),
		connManager: xnet.NewConnManager(),
		idFacotry:   newIdManager(int(opt.MaxConnectionNum)),
		pool: &sync.Pool{
			New: func() interface{} {
				return nil
			},
		},
	}
}

func Default() *Server {
	return NewServer(&defaultServerOpt)
}

// 停止服务
func (s *Server) Start() {
	if s.opts.PoolSize > 0 && s.dispatch != nil {
		s.states = xnet.Start
		s.dispatch.StartWookerPool(s.opts.PoolSize, 100)
	}
}

// 停止服务
func (s *Server) Stop() {
	//释放所有链接资源
	s.GetConnManager().Clear()
	fmt.Println("server is stopped")
}

// 启动一个服务
func (s *Server) Serve() {
	go s.Start()
	select {}
}

// 获取server opt
func (s *Server) GetServerOpt() *xnet.ServerOption {
	return s.opts
}

// 握手并启动链接处理
func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		if s.shakeFaieldFunc != nil {
			s.shakeFaieldFunc(w)
		} else {
			w.Write([]byte(err.Error()))
		}
		return
	}
	var id uint32
	if s.acceptFucn == nil {
		id = s.idFacotry.GetId()
	}
	wsConn := NewConnection(id, conn, s)
	wsConn.Start()
}

// 设置握手失败处理方法
func (s *Server) SetHandShakeFaildFunc(handler func(w http.ResponseWriter)) {
	s.shakeFaieldFunc = handler
}

// 获取链接id
func (s *Server) FactoryConnectionId() uint32 {
	return s.idFacotry.GetId()
}

// 移除链接id
func (s *Server) RemoveConnectionId(id uint32) {
	return
}

func (s *Server) GetConnManager() *xnet.ConnManager {
	return s.connManager
}

func (s *Server) GetPool() *sync.Pool {
	return s.pool
}

// 设置路由管理
func (s *Server) SetDistpatcher(dispatch xnet.Dispatcher) {
	s.dispatch = dispatch
}

// 获取路由管理
func (s *Server) GetDispatcher() xnet.Dispatcher {
	return s.dispatch
}

func (s *Server) SetPoolSize(size uint32) {
	s.opts.PoolSize = size
}

func (s *Server) GetStates() xnet.States {
	return s.states
}

// 根据用户id 获取单个链接id
func (s *Server) GetConnectionByID(id string) (xnet.Connection, bool) {
	if len(id) == 0 {
		return nil, false
	}
	conn, err := s.connManager.Get(id)
	if err != nil {
		return conn, false
	}
	return conn, true
}

func (s *Server) ChatWith(uid string, msgID uint32, data []byte) error {
	conn, ok := s.GetConnectionByID(uid)
	if !ok || conn == nil {
		return ErrConnNotFound
	}

	return conn.SendMsg(msgID, data)
}

// id管理
type idManager struct {
	maxsize int // 当偏移量大出这个范围后充值map
	used    map[uint32]struct{}
	has     map[uint32]struct{}
	m       sync.Mutex
}

func newIdManager(size int) *idManager {
	if size < 1 {
		panic("max connection must >= 1")
	}
	ids := idManager{
		used: make(map[uint32]struct{}),
		has:  make(map[uint32]struct{}),
	}

	for i := 1; i <= size; i++ {
		ids.has[uint32(i)] = struct{}{}
	}
	return &ids
}

// 获取可以用的id
func (i *idManager) GetId() uint32 {
	var has uint32
	i.m.Lock()
	defer i.m.Unlock()
	for id := range i.has {
		has = id
		delete(i.has, has)
		i.used[has] = empty{}
		break
	}
	return has
}

func (i *idManager) RemoveId(id uint32) {
	i.m.Lock()
	defer i.m.Unlock()
	if _, ok := i.used[id]; ok {
		delete(i.used, id)
		i.has[id] = empty{}
		return
	}
}

func (s *Server) GetID() int64 {
	return s.id
}

func (s *Server) SetAcceptFunc(accept func([]byte) (bool, xnet.IdCreater)) {
	s.acceptFucn = accept
}

func (s *Server) GetAcceptFunc() func([]byte) (bool, xnet.IdCreater) {
	return s.acceptFucn
}
