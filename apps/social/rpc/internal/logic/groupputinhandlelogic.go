package logic

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"llb-chat/apps/social/socialmodels"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/xerr"

	"llb-chat/apps/social/rpc/internal/svc"
	"llb-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrGroupReqBeforePass   = xerr.NewMsg("请求已通过")
	ErrGroupReqBeforeRefuse = xerr.NewMsg("请求已拒绝")
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 群申请处理
func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	// 查找群组申请请求记录
	groupReq, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, int64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group req err %v req %v", err, in.GroupReqId)
	}

	// 检查申请处理结果状态
	switch constants.HandlerResult(groupReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforePass)
	case constants.RefuseHandlerResult:
		return nil, errors.WithStack(ErrGroupReqBeforeRefuse)
	}
	// 更新申请处理结果状态
	groupReq.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}
	// 开启事务处理
	err = l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {

		// 更新群组申请记录
		if err := l.svcCtx.GroupRequestsModel.Update(l.ctx, session, groupReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "update group req err %v req %v", err, groupReq)
		}

		// 如果申请处理结果不是通过，则结束事务
		if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
			return nil
		}

		// 如果申请通过，添加用户到群组成员中
		groupMember := &socialmodels.GroupMembers{
			GroupId:     groupReq.GroupId,
			UserId:      groupReq.ReqId,
			RoleLevel:   int(constants.AtLargeGroupRoleLevel),
			OperatorUid: in.HandleUid,
		}
		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, groupMember)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v req %v", err, groupMember)
		}

		return nil
	})

	// 如果申请处理结果不是通过，返回空响应
	if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
		return &social.GroupPutInHandleResp{}, err
	}

	// 返回群组ID
	return &social.GroupPutInHandleResp{
		GroupId: groupReq.GroupId,
	}, err
}
