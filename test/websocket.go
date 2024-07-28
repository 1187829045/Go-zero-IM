package main

import (
	"fmt"                          // 导入用于格式化 I/O 的标准库
	"github.com/gorilla/websocket" // 导入 Gorilla WebSocket 包
	"log"                          // 导入日志包
	"net/http"                     // 导入 HTTP 包
)

// 定义一个 WebSocket 升级器，使用默认选项
var upgrade = websocket.Upgrader{}

// 定义处理 WebSocket 连接的函数
func serverWs(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		// 如果升级失败，记录错误并返回
		log.Print("upgrade:", err)
		return
	}
	// 确保在函数结束时关闭连接
	defer conn.Close()
	for {
		// 读取 WebSocket 消息
		//消息类型 (mt)，消息内容 (message)，以及一个错误对象 (err)。
		mt, message, err := conn.ReadMessage()
		if err != nil {
			// 如果读取消息失败，记录错误并退出循环
			log.Println("read:", err)
			break
		}
		// 打印接收到的消息
		log.Printf("recv: %s", message)
		// 将接收到的消息原样返回
		err = conn.WriteMessage(mt, message)
		if err != nil {
			// 如果写入消息失败，记录错误并退出循环
			log.Println("write:", err)
			break
		}
	}
}

// 定义主函数
func main() {
	// 将 "/ws" 路径映射到处理 WebSocket 连接的函数
	http.HandleFunc("/ws", serverWs)
	// 打印启动 WebSocket 服务器的消息
	fmt.Println("启动websocket")
	// 启动 HTTP 服务器，监听 0.0.0.0:1234 地址，日志记录任何错误
	log.Fatal(http.ListenAndServe("0.0.0.0:1234", nil))
}
