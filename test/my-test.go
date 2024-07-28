package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// 这个文件，用于重写一些函数
var upgrader = websocket.Upgrader{}

func ServerWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(string(message))
		err = conn.WriteMessage(mt, message)
	}
}

func main() {
	// 将 "/ws" 路径映射到处理 WebSocket 连接的函数
	http.HandleFunc("/ws", ServerWs)
	// 打印启动 WebSocket 服务器的消息
	fmt.Println("启动websocket")
	// 启动 HTTP 服务器，监听 0.0.0.0:1234 地址，日志记录任何错误
	log.Fatal(http.ListenAndServe("0.0.0.0:1234", nil))
}
