package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/liujingkaiai/x-socket/xnet/websocket"
)

func main() {
	wsserver := websocket.Default()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		temp := template.Must(template.ParseFiles("./ws.html"))
		temp.Execute(w, nil)
	})
	http.HandleFunc("/ws", wsserver.ServeWs)
	fmt.Println("server starting at 7772")
	log.Fatal(http.ListenAndServe(":7772", nil))
}
