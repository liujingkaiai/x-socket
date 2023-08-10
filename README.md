### 基于gorilla/websocket库进行二次封装 
- 自定义路由 
- 自定义鉴权，自定义描述符id目前支持 string
- 链接管理器
- 限制连接数 
- 单聊
- wooker pool 
- task queue 
- 自定义用户tag 
- 发布订阅模式，向某种tag标签的用户发送消息


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


### 路由

消息格式 消息id + 消息内容   |uint32|消息内容二进制
```go 

const (
	Hello uint32 iota
)

//定义路由
func HelloHandler(req xnet.Request) {
	//  向连接写入消息  id:0  写入内容 world
	req.SendMsg(Hello , []byte("world"))
}


// 实例化 webscoket server 
s := xnet.NewDefault() 
// 设置路由 id:0 处理方法  HelloWorld
s.AddRouter(Hello, HelloHandler) 

```

