package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"llb-chat/pkg/xerr"

	"llb-chat/apps/social/rpc/internal/svc"
	"llb-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

// 结构体包含上下文、服务上下文和日志记录器
type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// 函数用于创建一个新的 FriendListLogic 实例
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 方法处理好友列表请求，返回好友列表响应
func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	// todo: 添加你的逻辑代码并删除此行

	// 从服务上下文中的 FriendsModel 获取指定用户 ID 的好友列表
	friendsList, err := l.svcCtx.FriendsModel.ListByUserid(l.ctx, in.UserId)
	if err != nil {
		// 如果获取好友列表出错，返回错误信息并记录日志
		return nil, errors.Wrapf(xerr.NewDBErr(), "list friend by uid err %v req %v ", err, in.UserId)
	}

	// 定义响应列表变量
	var respList []*social.Friends
	// 复制好友列表到响应列表
	copier.Copy(&respList, &friendsList)

	// 返回好友列表响应
	return &social.FriendListResp{
		List: respList, // 好友列表
	}, nil
}
