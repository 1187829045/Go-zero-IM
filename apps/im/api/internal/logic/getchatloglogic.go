package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"llb-chat/apps/im/rpc/imclient"

	"llb-chat/apps/im/api/internal/svc"
	"llb-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetChatLogLogic 是一个包含日志记录器、上下文和服务上下文的结构体
type GetChatLogLogic struct {
	logx.Logger // 嵌入日志记录器
	ctx         context.Context
	svcCtx      *svc.ServiceContext // 服务上下文
}

// NewGetChatLogLogic 创建一个新的 GetChatLogLogic 实例
func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		Logger: logx.WithContext(ctx), // 通过上下文创建日志记录器
		ctx:    ctx,                   // 初始化上下文
		svcCtx: svcCtx,                // 初始化服务上下文
	}
}

// GetChatLog 是 GetChatLogLogic 的方法，用于获取聊天记录
func (l *GetChatLogLogic) GetChatLog(req *types.ChatLogReq) (resp *types.ChatLogResp, err error) {
	// todo: add your logic here and delete this line

	// 调用服务上下文的 GetChatLog 方法获取聊天记录
	data, err := l.svcCtx.GetChatLog(l.ctx, &imclient.GetChatLogReq{
		ConversationId: req.ConversationId, // 设置会话ID
		StartSendTime:  req.StartSendTime,  // 设置开始发送时间
		EndSendTime:    req.EndSendTime,    // 设置结束发送时间
		Count:          req.Count,          // 设置记录数
	})
	if err != nil {
		// 如果发生错误，返回错误
		return nil, err
	}

	// 定义返回的响应结构体
	var res types.ChatLogResp
	// 使用 copier 复制 data 到 res
	copier.Copy(&res, data)

	// 返回响应结果和错误
	return &res, err
}
