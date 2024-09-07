package websocket

import "time"

// FrameType 定义帧的类型
type FrameType uint8

// 定义帧类型的常量
const (
	FrameData  FrameType = 0x0 // 数据帧
	FramePing  FrameType = 0x1 // Ping 帧，用于检查连接是否仍然活跃
	FrameAck   FrameType = 0x2 // 确认帧，表示接收到消息
	FrameNoAck FrameType = 0x3 // 不确认帧，表示消息未被确认
	FrameErr   FrameType = 0x9 // 错误帧，表示发生了错误

	//FrameHeaders      FrameType = 0x1
	//FramePriority     FrameType = 0x2
	//FrameRSTStream    FrameType = 0x3
	//FrameSettings     FrameType = 0x4
	//FramePushPromise  FrameType = 0x5
	//FrameGoAway       FrameType = 0x7
	//FrameWindowUpdate FrameType = 0x8
	//FrameContinuation FrameType = 0x9
)

// Message 结构体表示一个 WebSocket 消息
type Message struct {
	FrameType `json:"frameType"` // 帧类型
	Id        string             `json:"id"`     // 消息 ID
	AckSeq    int                `json:"ackSeq"` // 确认序列号
	ackTime   time.Time
	errCount  int
	Method    string      `json:"method"` // 请求方法名
	FormId    string      `json:"formId"` // 消息请求的来源，接受可以不用，已经和websocket登陆了，所以在系统种可以获得ID，用于客户方发给接收方
	Data      interface{} `json:"data"`
}

// NewMessage 创建一个新的数据帧消息
func NewMessage(formId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData, // 设置帧类型为数据帧
		FormId:    formId,
		Data:      data,
	}
}

// NewErrMessage 创建一个新的错误帧消息
func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,    // 设置帧类型为错误帧
		Data:      err.Error(), // 设置错误信息
	}
}
