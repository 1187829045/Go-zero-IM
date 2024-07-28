package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"llb-chat/apps/im/immodels"
	"llb-chat/pkg/xerr"

	"llb-chat/apps/im/rpc/im"
	"llb-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// 根据用户 ID 查询用户会话列表
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		// 如果未找到会话列表，返回空响应
		if err == immodels.ErrNotFound {
			return &im.GetConversationsResp{}, nil
		}
		// 查询过程中出现其他错误，返回包装后的错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindByUserId err %v, req %v", err, in.UserId)
	}

	// 创建一个响应对象
	var res im.GetConversationsResp

	// 使用 copier 复制数据到响应对象中
	copier.Copy(&res, &data)

	// 从会话列表中提取会话 ID 列表
	ids := make([]string, 0, len(data.ConversationList))
	for _, conversation := range data.ConversationList {
		ids = append(ids, conversation.ConversationId)
	}

	// 根据会话 ID 列表查询具体的会话信息
	conversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		// 查询过程中出现错误，返回包装后的错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.ListByConversationIds err %v, req %v", err, ids)
	}

	// 遍历查询到的具体会话信息，计算是否存在未读消息
	for _, conversation := range conversations {
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			// 如果会话列表中没有该会话，继续下一个
			continue
		}

		// 用户已读取的消息量
		total := res.ConversationList[conversation.ConversationId].Total

		if total < int32(conversation.Total) {
			// 如果用户已读取的消息量小于会话中的总消息量，则有新的未读消息
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			// 计算未读消息量
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total
			// 将当前会话设置为显示状态
			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}

	// 返回响应对象
	return &res, nil
}
