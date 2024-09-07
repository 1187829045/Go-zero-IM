package mq

import "llb-chat/pkg/constants"

//MType，ChatTyped都是int类型

// 结构体表示一条聊天消息的传输信息

type MsgChatTransfer struct {
	MsgId string `mapstructure:"msgId"`

	ConversationId     string `json:"conversationId"`
	constants.ChatType `json:"chatType"`
	SendId             string   `json:"sendId"`
	RecvId             string   `json:"recvId"`
	RecvIds            []string `json:"recvIds"` // 接收者ID列表，表示消息接收方的用户ID列表（群聊）
	SendTime           int64    `json:"sendTime"`

	constants.MType `json:"mType"` // 消息类型，使用常量定义，可能是文本、图片、视频等
	Content         string         `json:"content"` // 消息内容，表示消息的实际内容
}

// 结构体表示已读消息的标记信息

type MsgMarkRead struct {
	constants.ChatType `json:"chatType"`
	ConversationId     string   `json:"conversationId"`
	SendId             string   `json:"sendId"`
	RecvId             string   `json:"recvId"`
	MsgIds             []string `json:"msgIds"`
}
