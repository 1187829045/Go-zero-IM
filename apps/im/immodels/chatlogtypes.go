package immodels

import (
	"llb-chat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DefaultChatLogLimit int64 = 100

// ChatLog 代表一条聊天记录
type ChatLog struct {
	// ID 是 MongoDB 自动生成的唯一标识符
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// ConversationId 是会话的唯一标识符
	ConversationId string `bson:"conversationId"`
	// SendId 是发送者的唯一标识符
	SendId string `bson:"sendId"`
	// RecvId 是接收者的唯一标识符
	RecvId string `bson:"recvId"`
	// MsgFrom 表示消息的来源（例如，用户或系统）
	MsgFrom    int                `bson:"msgFrom"`
	ChatType   constants.ChatType `bson:"chatType"`
	MsgType    constants.MType    `bson:"msgType"`
	MsgContent string             `bson:"msgContent"`
	SendTime   int64              `bson:"sendTime"`
	// Status 表示消息的状态（例如，已发送、已接收）
	Status int `bson:"status"`
	// ReadRecords 是消息的已读记录
	ReadRecords []byte `bson:"readRecords"`
	// TODO: 填写你自己的字段
	// UpdateAt 是记录的最后更新时间
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	// CreateAt 是记录的创建时间
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
