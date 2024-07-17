package group

import (
	"context"
	"github.com/jinzhu/copier"
	"llb-chat/apps/social/rpc/socialclient"
	"llb-chat/pkg/ctxdata"

	"llb-chat/apps/social/api/internal/svc"
	"llb-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupListLogic) GroupList(req *types.GroupListRep) (resp *types.GroupListResp, err error) {
	// todo: add your logic here and delete this line

	uid := ctxdata.GetUId(l.ctx)
	list, err := l.svcCtx.Social.GroupList(l.ctx, &socialclient.GroupListReq{
		UserId: uid,
	})

	if err != nil {
		return nil, err
	}

	var respList []*types.Groups
	copier.Copy(&respList, list.List)

	return &types.GroupListResp{List: respList}, nil
}
