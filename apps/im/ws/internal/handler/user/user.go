package user

import (
	"llb-chat/apps/im/ws/internal/svc"
	"llb-chat/apps/im/ws/websocket"
)

//获取所有在线用户

//返回一个处理 WebSocket 消息的函数，用于获取所有在线用户

func OnLine(svc *svc.ServiceContext) websocket.HandlerFunc {
	// 返回一个处理 WebSocket 消息的函数
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 获取服务器上所有在线用户的 UID
		uids := srv.GetUsers()

		// 获取当前连接的用户信息
		u := srv.GetUsers(conn)

		// 发送包含所有在线用户的 UID 的消息到当前连接的用户
		err := srv.Send(websocket.NewMessage(u[0], uids), conn)

		// 打印发送消息时的错误信息，如果有的话
		srv.Info("err ", err)
	}
}
