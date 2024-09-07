package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"llb-chat/apps/social/socialmodels"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/xerr"

	"llb-chat/apps/social/rpc/internal/svc"
	"llb-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

// 错误定义
var (
	ErrFriendReqBeforePass   = xerr.NewMsg("好友申请并已经通过")
	ErrFriendReqBeforeRefuse = xerr.NewMsg("好友申请已经被拒绝")
)

// FriendPutInHandleLogic 结构体包含上下文、服务上下文和日志记录器
type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewFriendPutInHandleLogic 函数用于创建一个新的 FriendPutInHandleLogic 实例
func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx), // 使用上下文创建日志记录器
	}
}

// 处理好友申请

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// todo: 添加你的逻辑代码并删除此行

	// 获取好友申请记录
	firendReq, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, int64(in.FriendReqId))
	if err != nil {
		// 如果查询好友申请记录出错，返回错误信息
		return nil, errors.Wrapf(xerr.NewDBErr(), "find friendsRequest by friendReqid err %v req %v ", err, in.FriendReqId)
	}

	// 验证好友申请是否已被处理
	switch constants.HandlerResult(firendReq.HandleResult.Int64) {
	case constants.PassHandlerResult: // 如果已通过
		return nil, errors.WithStack(ErrFriendReqBeforePass)
	case constants.RefuseHandlerResult: // 如果已拒绝
		return nil, errors.WithStack(ErrFriendReqBeforeRefuse)
	}

	// 设置新的处理结果
	firendReq.HandleResult.Int64 = int64(in.HandleResult)

	// 修改申请结果并处理好友关系 - 使用事务处理
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 更新好友申请记录
		if err := l.svcCtx.FriendRequestsModel.Update(l.ctx, session, firendReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update friend request err %v, req %v", err, firendReq)
		}

		// 如果处理结果不是通过，则无需建立好友关系记录
		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		// 创建两条新的好友关系记录
		friends := []*socialmodels.Friends{
			{
				UserId:    firendReq.UserId,
				FriendUid: firendReq.ReqUid,
			}, {
				UserId:    firendReq.ReqUid,
				FriendUid: firendReq.UserId,
			},
		}

		// 批量插入好友关系记录
		_, err = l.svcCtx.FriendsModel.Inserts(l.ctx, session, friends...)
		if err != nil {
			// 如果插入好友关系记录出错，返回错误信息
			return errors.Wrapf(xerr.NewDBErr(), "friends inserts err %v, req %v", err, friends)
		}
		return nil
	})

	// 返回处理结果和错误信息
	return &social.FriendPutInHandleResp{}, err
}
