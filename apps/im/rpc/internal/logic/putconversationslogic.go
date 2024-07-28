package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"llb-chat/apps/im/immodels"
	"llb-chat/apps/im/rpc/im"
	"llb-chat/apps/im/rpc/internal/svc"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/xerr"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	// TODO: 添加具体业务逻辑代码

	// 根据用户 ID 查找会话数据
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		// 如果查找时发生错误，返回一个包装了错误信息的错误对象，并附带详细的上下文信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindByUserId err %v, req %v", err, in.UserId)
	}

	// 如果会话列表为空，则初始化一个空的会话列表
	if data.ConversationList == nil {
		data.ConversationList = make(map[string]*immodels.Conversation)
	}

	// 遍历输入中的会话列表
	for s, conversation := range in.ConversationList {
		var oldTotal int
		// 如果旧的会话列表中已有该会话，则获取其原有的总数
		if data.ConversationList[s] != nil {
			oldTotal = data.ConversationList[s].Total
		}

		// 更新或添加会话信息到会话列表中
		data.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversation.ConversationId,               // 会话 ID
			ChatType:       constants.ChatType(conversation.ChatType), // 聊天类型
			IsShow:         conversation.IsShow,                       // 是否显示
			Total:          int(conversation.Read) + oldTotal,         // 消息总数
			Seq:            conversation.Seq,                          // 序列号
		}
	}

	// 更新会话数据
	err = l.svcCtx.ConversationsModel.Update(l.ctx, data)
	if err != nil {
		// 如果更新时发生错误，返回一个包装了错误信息的错误对象，并附带详细的上下文信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Update err %v, req %v", err, data)
	}

	// 返回空的响应对象表示操作成功
	return &im.PutConversationsResp{}, nil
}
