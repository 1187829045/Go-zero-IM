package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"llb-chat/apps/social/socialmodels"
	"llb-chat/pkg/constants"
	"llb-chat/pkg/wuid"
	"llb-chat/pkg/xerr"

	"llb-chat/apps/social/rpc/internal/svc"
	"llb-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 群要求
func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// todo: add your logic here and delete this line

	groups := &socialmodels.Groups{
		Id:         wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Name:       in.Name,
		Icon:       in.Icon,
		CreatorUid: in.CreatorUid,
		//IsVerify:   true,
		IsVerify: false,
	}

	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.GroupsModel.Insert(l.ctx, session, groups)

		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group err %v req %v", err, in)
		}

		_, err = l.svcCtx.GroupMembersModel.Insert(l.ctx, session, &socialmodels.GroupMembers{
			GroupId:   groups.Id,
			UserId:    in.CreatorUid,
			RoleLevel: int(constants.CreatorGroupRoleLevel),
		})
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "insert group member err %v req %v", err, in)
		}
		return nil
	})

	return &social.GroupCreateResp{}, err
}
