package logic

import (
	"context"
	"github.com/pkg/errors"
	"llb-chat/pkg/xerr"

	"llb-chat/apps/im/rpc/im"
	"llb-chat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话记录
func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// TODO: 添加具体业务逻辑代码

	// 根据消息 ID 查找聊天记录
	if in.MsgId != "" {
		// 查找指定消息 ID 的聊天记录
		chatlog, err := l.svcCtx.ChatLogModel.FindOne(l.ctx, in.MsgId)
		if err != nil {
			// 如果查找时发生错误，返回一个包装了错误信息的错误对象，并附带详细的上下文信息
			return nil, errors.Wrapf(xerr.NewDBErr(), "find chatLog by msgId err %v, req %v", err, in.MsgId)
		}

		// 返回包含查找到的聊天记录的响应对象
		return &im.GetChatLogResp{
			List: []*im.ChatLog{{
				Id:             chatlog.ID.Hex(),       // 聊天记录的唯一 ID
				ConversationId: chatlog.ConversationId, // 对话 ID
				SendId:         chatlog.SendId,
				RecvId:         chatlog.RecvId,
				MsgType:        int32(chatlog.MsgType),
				MsgContent:     chatlog.MsgContent,
				ChatType:       int32(chatlog.ChatType), // 对话类型
				SendTime:       chatlog.SendTime,
				ReadRecords:    chatlog.ReadRecords,
			}},
		}, nil
	}

	// 根据发送时间范围查找聊天记录
	data, err := l.svcCtx.ChatLogModel.ListBySendTime(l.ctx, in.ConversationId, in.StartSendTime, in.EndSendTime, in.Count)
	if err != nil {
		// 如果查找时发生错误，返回一个包装了错误信息的错误对象，并附带详细的上下文信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "find chatLog list by SendTime err %v, req %v", err, in)
	}

	// 初始化响应对象中的聊天记录列表
	res := make([]*im.ChatLog, 0, len(data))
	// 将查找到的聊天记录转换为响应对象中的格式
	for _, datum := range data {
		res = append(res, &im.ChatLog{
			Id:             datum.ID.Hex(),        // 聊天记录的唯一 ID
			ConversationId: datum.ConversationId,  // 对话 ID
			SendId:         datum.SendId,          // 发送者 ID
			RecvId:         datum.RecvId,          // 接收者 ID
			MsgType:        int32(datum.MsgType),  // 消息类型
			MsgContent:     datum.MsgContent,      // 消息内容
			ChatType:       int32(datum.ChatType), // 对话类型
			SendTime:       datum.SendTime,        // 发送时间
			ReadRecords:    datum.ReadRecords,     // 阅读记录
		})
	}

	// 返回包含查找到的聊天记录列表的响应对象
	return &im.GetChatLogResp{
		List: res,
	}, nil
}
