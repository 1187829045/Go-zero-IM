package svc

import (
	"llb-chat/apps/im/immodels"
	"llb-chat/apps/im/ws/internal/config"
	"llb-chat/apps/task/mq/mqclient"
)

//整个服务的上下文和全局属性

type ServiceContext struct {
	Config config.Config

	immodels.ChatLogModel
	mqclient.MsgChatTransferClient //消息推送到 Kafka
	mqclient.MsgReadTransferClient //推送已读消息Kafka
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                c,
		MsgChatTransferClient: mqclient.NewMsgChatTransferClient(c.MsgChatTransfer.Addrs, c.MsgChatTransfer.Topic),
		MsgReadTransferClient: mqclient.NewMsgReadTransferClient(c.MsgReadTransfer.Addrs, c.MsgReadTransfer.Topic),
		ChatLogModel:          immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
	}
}
