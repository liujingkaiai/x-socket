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

消息格式 uint32 + 消息内容   |uint32|消息内容  
```go 
// 使用方法 
s := xnet.NewDefault()  //设置httpServer
s.AddRouter(1, func(xnet.Request)) 


```

