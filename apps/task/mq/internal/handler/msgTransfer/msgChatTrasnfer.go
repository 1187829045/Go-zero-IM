package msgTransfer

import (
	"context"                                    // 导入上下文包，用于管理请求的生命周期
	"encoding/json"                              // 导入 JSON 包，用于处理 JSON 数据的编解码
	"fmt"                                        // 导入 fmt 包，用于格式化输出
	"go.mongodb.org/mongo-driver/bson/primitive" // 导入 MongoDB BSON 库，用于生成 ObjectID
	"llb-chat/apps/im/immodels"                  // 导入聊天模型
	"llb-chat/apps/im/ws/ws"                     // 导入 WebSocket 包
	"llb-chat/apps/task/mq/internal/svc"         // 导入服务上下文包
	"llb-chat/apps/task/mq/mq"                   // 导入消息队列相关包
	"llb-chat/pkg/bitmap"                        // 导入位图包，用于处理消息的已读状态
)

// MsgChatTransfer 结构体，嵌套 baseMsgTransfer，继承其功能
type MsgChatTransfer struct {
	*baseMsgTransfer // 嵌套 baseMsgTransfer 以继承其功能
}

// NewMsgChatTransfer 构造函数，用于创建 MsgChatTransfer 实例
func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	// 返回 MsgChatTransfer 的新实例
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc), // 初始化 baseMsgTransfer
	}
}

//方法处理消费的消息
//mq存储消息后发送给websocket

func (m *MsgChatTransfer) Consume(key, value string) error {
	// 打印消息的键和值
	fmt.Println("key : ", key, " value : ", value)

	var (
		data  mq.MsgChatTransfer        // 消息数据
		ctx   = context.Background()    // 创建上下文
		msgId = primitive.NewObjectID() // 生成一个新的 ObjectID 作为消息 ID
	)
	// 解码 JSON 数据到 MsgChatTransfer 结构体
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err // 如果解码失败，返回错误
	}

	// 记录聊天数据
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err // 如果记录失败，返回错误
	}

	// 转发消息到 WebSocket
	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId, // 会话 ID
		ChatType:       data.ChatType,       // 聊天类型
		SendId:         data.SendId,         // 发送者 ID
		RecvId:         data.RecvId,         // 接收者 ID
		RecvIds:        data.RecvIds,        // 多个接收者 ID
		SendTime:       data.SendTime,       // 发送时间
		MType:          data.MType,          // 消息类型
		MsgId:          msgId.Hex(),         // 消息 ID，转换为十六进制字符串
		Content:        data.Content,        // 消息内容
	})
}

// addChatLog 方法记录聊天信息
func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	// 创建聊天消息
	chatLog := immodels.ChatLog{
		ID:             msgId,               // 消息 ID
		ConversationId: data.ConversationId, // 会话 ID
		SendId:         data.SendId,         // 发送者 ID
		RecvId:         data.RecvId,         // 接收者 ID
		ChatType:       data.ChatType,       // 聊天类型
		MsgFrom:        0,                   // 消息来源，0 可能表示默认值或系统消息
		MsgType:        data.MType,          // 消息类型
		MsgContent:     data.Content,        // 消息内容
		SendTime:       data.SendTime,       // 发送时间
	}

	// 创建并初始化位图，用于跟踪已读记录
	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)            // 将发送者 ID 标记为已读
	chatLog.ReadRecords = readRecords.Export() // 导出位图为可存储格式

	// 插入聊天消息到数据库
	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err // 如果插入失败，返回错误
	}

	// 更新会话中的消息记录
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
