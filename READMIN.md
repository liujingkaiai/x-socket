### 基于gorilla/websocket/blob库进行二次封装 
- 自定义路由 
- 自定义鉴权，自定义描述符id目前支持 int  int8  int32  int64  uint8  uint32  uint64  string
- 链接管理器
- 限制连接数 
- 中间件
- wooker pool 
- task queue 
- 自定义用户tag 
- 向某种tag标签的用户发送消息


### todo 
- context上下文
- 数据校验
- 负载均衡
- 跨服务通信
- 日志收集
- 链路追踪

### 使用方法 

```go 
package main

import (
	"net/http"

	"github.com/liujingkaiai/x-websocekt/xnet"
)

func main() {
	s := xnet.NewDefault()  //设置httpServer
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})
	http.HandleFunc("/ws", s.Shake)
}

```


###路由

消息格式 消息id + 消息内容   |uint32|消息内容二进制
```go 

const (
	Hello uint32 iota
)

// 使用方法 
s := xnet.NewDefault()  //设置httpServer
s.AddRouter(Hello, func(iface.Connection , data)) 

```