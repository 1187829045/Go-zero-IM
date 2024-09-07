package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"llb-chat/apps/social/socialmodels"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/xerr"
	"time"

	"llb-chat/apps/social/rpc/internal/svc"
	"llb-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutIn 方法处理好友申请的逻辑
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: 在这里添加逻辑并删除此行

	// 检查申请人是否已经是目标用户的好友
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && err != socialmodels.ErrNotFound {
		// 如果查询好友关系时出错，并且错误不是记录未找到的错误，返回包装后的数据库错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friends by uid and fid err %v req %v ", err, in)
	}
	if friends != nil {
		// 如果找到了好友记录，说明申请人已经是目标用户的好友，直接返回成功响应
		return &social.FriendPutInResp{}, err
	}

	//检查是否已经存在未成功的申请记录
	friendReqs, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && err != socialmodels.ErrNotFound {
		// 如果查询申请记录时出错，并且错误不是记录未找到的错误，返回包装后的数据库错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by rid and uid err %v req %v ", err, in)
	}
	if friendReqs != nil {
		// 如果找到已有的申请记录，说明申请记录已经存在，直接返回成功响应
		return &social.FriendPutInResp{}, err
	}

	// 创建新的好友申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId, // 申请人用户ID
		ReqUid: in.ReqUid, // 被申请人用户ID
		ReqMsg: sql.NullString{
			Valid:  true,      // 标记请求消息字段有效
			String: in.ReqMsg, // 设置请求消息内容
		},
		ReqTime: time.Unix(in.ReqTime, 0), // 设置申请时间
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult), // 设置处理结果为“未处理”
			Valid: true,                             // 标记处理结果字段有效
		},
	})

	if err != nil {
		// 如果插入申请记录时出错，返回包装后的数据库错误
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert friendRequest err %v req %v ", err, in)
	}

	// 成功处理好友申请，返回成功响应
	return &social.FriendPutInResp{}, nil
}
