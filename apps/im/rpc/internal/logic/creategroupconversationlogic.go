package logic

import (
	"context"
	"github.com/pkg/errors"
	"llb-chat/apps/im/immodels"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/xerr"

	"llb-chat/apps/im/rpc/im"
	"llb-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	// TODO: 添加具体业务逻辑代码

	res := &im.CreateGroupConversationResp{}

	// 尝试在数据库中查找是否已存在指定群组 ID 的对话记录
	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	if err == nil {
		return res, nil
	}
	// 如果发生了除“记录未找到”之外的其他错误
	if err != immodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.FindOne err %v, req %v", err, in.GroupId)
	}

	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})
	if err != nil {
		// 返回一个包装了错误信息的错误对象，并附带详细的上下文信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "ConversationModel.Insert err %v", err)
	}

	// 设置用户群组对话信息
	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return res, err
}
