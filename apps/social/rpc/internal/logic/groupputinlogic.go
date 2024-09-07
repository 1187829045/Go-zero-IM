package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"llb-chat/apps/social/rpc/internal/svc"
	"llb-chat/apps/social/rpc/social"
	"llb-chat/apps/social/socialmodels"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/xerr"
)

type GroupPutinLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinLogic {
	return &GroupPutinLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {
	// 处理用户加入群组的请求

	var (
		inviteGroupMember *socialmodels.GroupMembers // 邀请人的群组成员信息
		userGroupMember   *socialmodels.GroupMembers // 用户的群组成员信息
		groupInfo         *socialmodels.Groups       // 群组信息

		err error
	)

	// 查找用户是否已经是群组成员
	userGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.ReqId, in.GroupId)
	if err != nil && err != socialmodels.ErrNotFound {
		// 如果查找过程中出现错误，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and  req id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if userGroupMember != nil {
		// 如果用户已经是群组成员，返回成功响应
		return &social.GroupPutinResp{}, nil
	}

	// 查找群组申请记录是否已经存在
	groupReq, err := l.svcCtx.GroupRequestsModel.FindByGroupIdAndReqId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && err != socialmodels.ErrNotFound {
		// 如果查找过程中出现错误，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group req by groud id and user id err %v, req %v, %v", err,
			in.GroupId, in.ReqId)
	}
	if groupReq != nil {
		// 如果申请记录已存在，返回成功响应
		return &social.GroupPutinResp{}, nil
	}

	// 创建新的群组申请记录
	groupReq = &socialmodels.GroupRequests{
		ReqId:   in.ReqId,
		GroupId: in.GroupId,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: sql.NullTime{
			Time:  time.Unix(in.ReqTime, 0),
			Valid: true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	}

	// 定义用于创建群组成员的函数
	createGroupMember := func() {
		if err != nil {
			return
		}
		// 调用创建群组成员的函数
		err = l.createGroupMember(in)
	}

	// 查找群组信息
	groupInfo, err = l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		// 如果查找过程中出现错误，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group by groud id err %v, req %v", err, in.GroupId)
	}

	// 根据群组的验证方式处理申请
	if !groupInfo.IsVerify {
		// 群组不需要验证
		defer createGroupMember() // 在函数退出时创建群组成员

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}

		// 创建群组申请记录并返回
		return l.createGroupReq(groupReq, true)
	}

	// 验证群组加入方式
	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		// 如果是通过申请加入的
		return l.createGroupReq(groupReq, false)
	}

	// 查找邀请人的群组成员信息
	inviteGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.InviterUid, in.GroupId)
	if err != nil {
		// 如果查找过程中出现错误，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group member by groud id and user id err %v, req %v",
			in.InviterUid, in.GroupId)
	}

	if constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.CreatorGroupRoleLevel || constants.
		GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.ManagerGroupRoleLevel {
		// 如果是管理者或创建者邀请
		defer createGroupMember() // 在函数退出时创建群组成员

		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		groupReq.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		// 创建群组申请记录并返回
		return l.createGroupReq(groupReq, true)
	}

	// 其他情况
	return l.createGroupReq(groupReq, false)
}

func (l *GroupPutinLogic) createGroupReq(groupReq *socialmodels.GroupRequests, isPass bool) (*social.GroupPutinResp, error) {
	// 插入群组申请记录
	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, groupReq)
	if err != nil {
		// 如果插入过程中出现错误，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert group req err %v req %v", err, groupReq)
	}

	// 如果申请通过，返回包含群组 ID 的响应
	if isPass {
		return &social.GroupPutinResp{GroupId: groupReq.GroupId}, nil
	}

	// 否则返回空响应
	return &social.GroupPutinResp{}, nil
}

func (l *GroupPutinLogic) createGroupMember(in *social.GroupPutinReq) error {
	// 创建群组成员记录
	groupMember := &socialmodels.GroupMembers{
		GroupId:     in.GroupId,
		UserId:      in.ReqId,
		RoleLevel:   int(constants.AtLargeGroupRoleLevel),
		OperatorUid: in.InviterUid,
	}
	_, err := l.svcCtx.GroupMembersModel.Insert(l.ctx, nil, groupMember)
	if err != nil {
		// 如果插入过程中出现错误，返回错误信息
		return errors.Wrapf(xerr.NewDBErr(), "insert friend err %v req %v", err, groupMember)
	}

	return nil
}
