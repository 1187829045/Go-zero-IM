package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"llb-chat/apps/user/models"
	"llb-chat/apps/user/rpc/internal/svc"
	"llb-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	// todo: add your logic here and delete this line

	var (
		userEntitys []*models.Users
		err         error
	)
	if in.Phone != "" {
		userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntitys = append(userEntitys, userEntity)
		}
	} else if in.Name != "" {
		userEntitys, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		// 否则，如果 Ids 列表长度大于0，通过用户ID列表查找用户
		userEntitys, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}
	if err != nil {
		return nil, err
	}
	var resp []*user.UserEntity
	// 将 userEntitys 复制到 resp 中
	copier.Copy(&resp, &userEntitys)

	// 返回查找结果
	return &user.FindUserResp{
		User: resp,
	}, nil
}
