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

type FriendPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutInList 方法处理获取用户未处理的好友申请列表
func (l *FriendPutInListLogic) FriendPutInList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	// todo: 在这里添加逻辑并删除此行

	// 调用服务上下文中的 FriendRequestsModel 获取未处理的好友申请列表
	friendReqList, err := l.svcCtx.FriendRequestsModel.ListNoHandler(l.ctx, in.UserId)
	if err != nil {
		// 如果查询好友申请列表出错，包装错误并返回
		return nil, errors.Wrapf(xerr.NewDBErr(), "find list friend req err %v req %v", err, in.UserId)
	}

	var resp []*social.FriendRequests // 定义返回的好友申请列表变量
	// 使用 copier.Copy 复制从模型获取的数据到响应变量
	copier.Copy(&resp, &friendReqList)

	// 返回好友申请列表响应
	return &social.FriendPutInListResp{
		List: resp,
	}, nil
}
