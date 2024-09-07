package push

import (
	"github.com/mitchellh/mapstructure"
	"llb-chat/apps/im/ws/internal/svc"
	"llb-chat/apps/im/ws/websocket"
	"llb-chat/apps/im/ws/ws"
	"llb-chat/pkg/constants"
)

// 是一个 WebSocket 处理函数，用于将消息推送到相应的用户或用户组。

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		// 将接收到的消息数据解码到 ws.Push 结构体中
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			// 如果解码失败，发送错误消息
			srv.Send(websocket.NewErrMessage(err))
			return
		}
		// 根据消息的聊天类型，决定消息的发送目标
		switch data.ChatType {
		case constants.SingleChatType:
			// 单聊，调用 single 函数处理
			single(srv, &data, data.RecvId)
		case constants.GroupChatType:
			// 群聊，调用 group 函数处理
			group(srv, &data)
		}
	}
}

// 处理单聊消息，将消息发送给指定的接收者
func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	// 获取接收者的 WebSocket 连接
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// 如果接收者不在线，记录离线信息

		// todo: 目标离线

		return nil
	}

	// 打印推送消息的日志
	srv.Infof("push msg %v", data)

	// 创建并推送消息
	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			ReadRecords: data.ReadRecords,
			MsgId:       data.MsgId,
			MType:       data.MType,
			Content:     data.Content,
		},
	}), rconn)
}

// 处理群聊消息，将消息发送给群组中的每个接收者
func group(srv *websocket.Server, data *ws.Push) error {
	// 遍历接收者 ID 列表
	for _, id := range data.RecvIds {
		// 使用匿名函数和 Schedule 方法并发处理每个接收者的消息发送
		func(id string) {
			srv.Schedule(func() {
				// 调用 single 函数处理单个接收者的消息
				single(srv, data, id)
			})
		}(id)
	}
	return nil
}
