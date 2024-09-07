package logic

import (
	"context"
	"llb-chat/apps/im/rpc/im"
	"llb-chat/apps/social/rpc/socialclient"
	"llb-chat/apps/user/rpc/user"
	"llb-chat/pkg/bitmap"
	"llb-chat/pkg/constants"

	"llb-chat/apps/im/api/internal/svc"
	"llb-chat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetChatLogReadRecordsLogic 是一个包含日志记录器、上下文和服务上下文的结构体
type GetChatLogReadRecordsLogic struct {
	logx.Logger                     // 嵌入日志记录器
	ctx         context.Context     // 上下文
	svcCtx      *svc.ServiceContext // 服务上下文
}

// NewGetChatLogReadRecordsLogic 创建一个新的 GetChatLogReadRecordsLogic 实例
func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx), // 通过上下文创建日志记录器
		ctx:    ctx,                   // 初始化上下文
		svcCtx: svcCtx,                // 初始化服务上下文
	}
}

// 用于获取聊天记录的已读未读信息
func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
	// todo: add your logic here and delete this line

	// 调用服务上下文的 Im.GetChatLog 方法获取聊天记录
	chatlogs, err := l.svcCtx.Im.GetChatLog(l.ctx, &im.GetChatLogReq{
		MsgId: req.MsgId,
	})
	if err != nil || len(chatlogs.List) == 0 {
		// 如果发生错误或没有聊天记录，返回错误
		return nil, err
	}

	// 初始化一些变量
	var (
		chatlog = chatlogs.List[0]         // 取第一个聊天记录
		reads   = []string{chatlog.SendId} // 已读的用户列表，初始包含发送者ID
		unreads []string                   // 未读的用户列表
		ids     []string                   // 所有相关用户的ID列表
	)

	// 根据聊天类型分别设置已读和未读列表
	switch constants.ChatType(chatlog.ChatType) {
	case constants.SingleChatType:
		// 单聊类型
		if len(chatlog.ReadRecords) == 0 || chatlog.ReadRecords[0] == 0 {
			unreads = []string{chatlog.RecvId} // 没有已读记录或已读记录为0，接收者未读
		} else {
			reads = append(reads, chatlog.RecvId) // 有已读记录，接收者已读
		}
		ids = []string{chatlog.RecvId, chatlog.SendId} // 添加接收者和发送者的ID到ID列表
	case constants.GroupChatType:
		//群聊类型
		groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
			GroupId: chatlog.RecvId, // 设置群组ID
		})
		if err != nil {
			// 如果发生错误，返回错误
			return nil, err
		}

		// 加载已读记录的位图
		bitmaps := bitmap.Load(chatlog.ReadRecords)
		for _, members := range groupUsers.List {
			ids = append(ids, members.UserId) // 添加群组成员的用户ID到ID列表

			if members.UserId == chatlog.SendId {
				// 发送者跳过
				continue
			}

			if bitmaps.IsSet(members.UserId) {
				reads = append(reads, members.UserId) // 成员已读
			} else {
				unreads = append(unreads, members.UserId) // 成员未读
			}
		}
	}

	// 查找所有相关用户的信息
	userEntitys, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{
		Ids: ids, // 设置要查找的用户ID列表
	})
	if err != nil {
		// 如果发生错误，返回错误
		return nil, err
	}
	userEntitySet := make(map[string]*user.UserEntity, len(userEntitys.User))
	for i, entity := range userEntitys.User {
		userEntitySet[entity.Id] = userEntitys.User[i] // 将用户信息存入映射表
	}

	// 设置已读和未读用户的手机号码
	for i, read := range reads {
		if u := userEntitySet[read]; u != nil {
			reads[i] = u.Phone // 设置已读用户的手机号码
		}
	}
	for i, unread := range unreads {
		if u := userEntitySet[unread]; u != nil {
			unreads[i] = u.Phone // 设置未读用户的手机号码
		}
	}

	// 返回已读和未读用户的手机号码
	return &types.GetChatLogReadRecordsResp{
		Reads:   reads,
		UnReads: unreads,
	}, nil
}
