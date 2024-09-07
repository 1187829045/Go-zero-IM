//调用mq客户端把消息对送到mq服务端

package conversation

import (
	"github.com/mitchellh/mapstructure"
	"llb-chat/apps/im/ws/internal/svc"
	"llb-chat/apps/im/ws/websocket"
	"llb-chat/apps/im/ws/ws"
	"llb-chat/apps/task/mq/mq"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/wuid"
	"time"
)

// 处理私聊消息的 WebSocket 处理函数
func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// TODO: 私聊逻辑

		// 定义一个 ws.Chat 类型的变量，用于存储解码后的消息数据
		var data ws.Chat

		// 将消息数据解码到 data 变量中，如果解码出错，发送错误消息并返回
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		// 如果会话 ID 为空，根据聊天类型生成会话 ID
		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.SingleChatType:
				// 如果是单聊，生成一个由发送者和接收者 ID 组合的会话 ID
				data.ConversationId = wuid.CombineId(conn.Uid, data.RecvId)
			case constants.GroupChatType:
				// 如果是群聊，使用接收者 ID 作为会话 ID
				data.ConversationId = data.RecvId
			}
		}

		// 将聊天消息推送到消息传输到kafka服务端

		err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			SendTime:       time.Now().UnixMilli(),
			MType:          data.Msg.MType,
			Content:        data.Msg.Content,
			MsgId:          msg.Id,
		})

		// 如果推送出错，发送错误消息并返回
		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}

//处理已读未读消息的 WebSocket 处理函数

func MarkRead(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// TODO: 已读未读处理

		// 定义一个 ws.MarkRead 类型的变量，用于存储解码后的消息数据
		var data ws.MarkRead

		// 将消息数据解码到 data 变量中，如果解码出错，发送错误消息并返回
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}

		// 将已读消息推送到 Kafka
		err := svc.MsgReadTransferClient.Push(&mq.MsgMarkRead{
			ChatType:       data.ChatType,
			ConversationId: data.ConversationId,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			MsgIds:         data.MsgIds,
		})

		// 如果推送出错，发送错误消息并返回
		if err != nil {
			srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}
