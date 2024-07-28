package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"llb-chat/apps/im/immodels"
	"llb-chat/apps/im/rpc/im"
	"llb-chat/apps/im/rpc/internal/svc"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/wuid"
	"llb-chat/pkg/xerr"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 建立会话: 群聊, 私聊
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// 创建一个 SetUpUserConversationResp 响应对象
	var res im.SetUpUserConversationResp

	// 根据 ChatType 进行逻辑处理
	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		// 生成会话的 ID
		conversationId := wuid.CombineId(in.SendId, in.RecvId)

		// 检查是否已存在该会话
		conversationRes, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			// 如果未找到该会话，建立新会话
			if err == immodels.ErrNotFound {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ConversationId: conversationId,           // 会话 ID
					ChatType:       constants.SingleChatType, // 聊天类型
				})

				if err != nil {
					// 插入会话失败，返回包装后的错误
					return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v", err)
				}
			} else {
				// 查找会话失败，返回包装后的错误
				return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, conversationId)
			}
		} else if conversationRes != nil {
			// 如果会话已存在，直接返回响应
			return &res, nil
		}

		// 为发送方和接收方分别建立会话
		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, true)
		if err != nil {
			return nil, err
		}
		err = l.setUpUserConversation(conversationId, in.RecvId, in.SendId, constants.SingleChatType, false)
		if err != nil {
			return nil, err
		}

	case constants.GroupChatType:
		// 为群聊建立会话
		err := l.setUpUserConversation(in.RecvId, in.SendId, in.RecvId, constants.GroupChatType, true)
		if err != nil {
			return nil, err
		}
	}

	// 返回响应对象
	return &res, nil
}

func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string,
	chatType constants.ChatType, isShow bool) error {
	// 查找用户的会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == immodels.ErrNotFound {
			// 如果未找到用户会话列表，创建一个新的
			conversations = &immodels.Conversations{
				ID:               primitive.NewObjectID(),                 // 新的 ObjectID
				UserId:           userId,                                  // 用户 ID
				ConversationList: make(map[string]*immodels.Conversation), // 初始化会话列表
			}
		} else {
			// 查找用户会话列表失败，返回包装后的错误
			return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, userId)
		}
	}

	// 如果会话列表中已有该会话记录，则不做更新
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}

	// 添加新的会话记录
	conversations.ConversationList[conversationId] = &immodels.Conversation{
		ConversationId: conversationId, // 会话 ID
		ChatType:       chatType,       // 聊天类型
		IsShow:         isShow,         // 是否显示
	}

	// 更新会话列表
	err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		// 更新失败，返回包装后的错误
		return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v, req %v", err, conversations)
	}
	return nil
}
