package handler

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"llb-chat/apps/task/mq/internal/handler/msgTransfer"
	"llb-chat/apps/task/mq/internal/svc"
)

type Listen struct {
	svc *svc.ServiceContext
}

func NewListen(svc *svc.ServiceContext) *Listen {
	return &Listen{svc: svc}
}

//有很多消费者所以要进行监听

func (l *Listen) Services() []service.Service {
	return []service.Service{
		kq.MustNewQueue(l.svc.Config.MsgReadTransfer, msgTransfer.NewMsgReadTransfer(l.svc)),
		// todo: 此处可以加载多个消费者
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer, msgTransfer.NewMsgChatTransfer(l.svc)),
	}
}
