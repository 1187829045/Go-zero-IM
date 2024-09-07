package mqclient

import (
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"llb-chat/apps/task/mq/mq"
)

type MsgChatTransferClient interface {
	Push(msg *mq.MsgChatTransfer) error
}

type msgChatTransferClient struct {
	pusher *kq.Pusher
}

// addr 是 Kafka 集群的地址列表，topic 是消息的主题

func NewMsgChatTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgChatTransferClient {
	return &msgChatTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

// Push 将 mq.MsgChatTransfer 消息推送到 Kafka服务端
func (c *msgChatTransferClient) Push(msg *mq.MsgChatTransfer) error {
	// 将消息序列化为 JSON 格式
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 使用 pusher 将消息推送到 Kafka
	return c.pusher.Push(string(body))
}

// 用于推送已读消息
type MsgReadTransferClient interface {
	Push(msg *mq.MsgMarkRead) error
}

// msgReadTransferClient 是 MsgReadTransferClient 接口的具体实现，包含一个 pusher
type msgReadTransferClient struct {
	pusher *kq.Pusher
}

// NewMsgReadTransferClient 创建并返回一个新的 msgReadTransferClient 实例
// addr 是 Kafka 集群的地址列表，topic 是消息的主题
func NewMsgReadTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgReadTransferClient {
	return &msgReadTransferClient{
		pusher: kq.NewPusher(addr, topic),
	}
}

// Push 将 mq.MsgMarkRead 消息推送到 Kafka
func (c *msgReadTransferClient) Push(msg *mq.MsgMarkRead) error {
	// 将消息序列化为 JSON 格式
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 使用 pusher 将消息推送到 Kafka
	return c.pusher.Push(string(body))
}
