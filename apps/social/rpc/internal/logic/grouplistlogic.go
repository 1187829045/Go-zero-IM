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

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取群列表
func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	//获取用户的群列表
	userGroup, err := l.svcCtx.GroupMembersModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		// 如果查询出错，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group member err %v req %v", err, in.UserId)
	}
	// 如果用户没有加入任何群组，返回空响应
	if len(userGroup) == 0 {
		return &social.GroupListResp{}, nil
	}

	// 创建一个存储群组ID的切片
	ids := make([]string, 0, len(userGroup))
	// 将用户加入的所有群组ID添加到切片中
	for _, v := range userGroup {
		ids = append(ids, v.GroupId)
	}

	// 根据群组ID列表获取群组信息
	groups, err := l.svcCtx.GroupsModel.ListByGroupIds(l.ctx, ids)
	if err != nil {
		// 如果查询出错，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group err %v req %v", err, ids)
	}

	// 定义响应的群组列表
	var respList []*social.Groups
	// 使用copier复制群组信息到响应列表中
	copier.Copy(&respList, &groups)

	// 返回群组列表响应
	return &social.GroupListResp{
		List: respList,
	}, nil
}
