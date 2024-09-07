package ws

import "llb-chat/pkg/constants"

// 消息体定义

type (
	// Msg 结构体表示一个消息的基本信息。
	Msg struct {
		// 消息ID，用于唯一标识一条消息。
		MsgId string `mapstructure:"msgId"`

		// 已读记录，记录各个用户对该消息的阅读状态。
		ReadRecords map[string]string `mapstructure:"readRecords"`

		// 消息类型，表示消息的种类（如文本、图片等），继承自 constants.MType。
		constants.MType `mapstructure:"mType"`

		// 消息内容，存储实际的消息文本或内容。
		Content string `mapstructure:"content"`
	}

	// Chat 结构体表示一个聊天记录。
	Chat struct {
		// 聊天类型，表示聊天的种类（如单聊、群聊），继承自 constants.ChatType。
		constants.ChatType `mapstructure:"chatType"`

		// 嵌套的消息体，包含消息的基本信息。
		Msg `mapstructure:"msg"`

		// 对话ID，标识一个对话的唯一标识符。
		ConversationId string `mapstructure:"conversationId"`

		// 发送者ID，标识消息的发送者。
		SendId string `mapstructure:"sendId"`

		// 接收者ID，标识消息的接收者。
		RecvId string `mapstructure:"recvId"`

		// 发送时间，消息发送的时间戳（以毫秒为单位）。
		SendTime int64 `mapstructure:"sendTime"`
	}

	// Push 结构体表示推送消息的详细信息。
	Push struct {
		// 聊天类型，表示聊天的种类（如单聊、群聊），继承自 constants.ChatType。
		constants.ChatType `mapstructure:"chatType"`

		// 消息类型，表示消息的种类（如文本、图片等），继承自 constants.MType。
		constants.MType `mapstructure:"mType"`

		// 对话ID，标识一个对话的唯一标识符。
		ConversationId string `mapstructure:"conversationId"`

		// 发送者ID，标识消息的发送者。
		SendId string `mapstructure:"sendId"`

		// 接收者ID，标识消息的接收者。
		RecvId string `mapstructure:"recvId"`

		// 接收者ID列表，标识消息的所有接收者（用于群聊）。
		RecvIds []string `mapstructure:"recvIds"`

		// 发送时间，消息发送的时间戳（以毫秒为单位）。
		SendTime int64 `mapstructure:"sendTime"`

		// 消息ID，用于唯一标识一条消息。
		MsgId string `mapstructure:"msgId"`

		// 已读记录，记录各个用户对该消息的阅读状态。
		ReadRecords map[string]string `mapstructure:"readRecords"`

		// 内容类型，表示消息内容的类型（如文本、图片等），继承自 constants.ContentType。
		ContentType constants.ContentType `mapstructure:"contentType"`

		// 消息内容，存储实际的消息文本或内容。
		Content string `mapstructure:"content"`
	}

	// MarkRead 结构体表示标记消息已读的请求。
	MarkRead struct {
		// 聊天类型，表示聊天的种类（如单聊、群聊），继承自 constants.ChatType。
		constants.ChatType `mapstructure:"chatType"`

		// 接收者ID，标识消息的接收者。
		RecvId string `mapstructure:"recvId"`

		// 对话ID，标识一个对话的唯一标识符。
		ConversationId string `mapstructure:"conversationId"`

		// 消息ID列表，标识要标记为已读的消息ID。
		MsgIds []string `mapstructure:"msgIds"`
	}
)
