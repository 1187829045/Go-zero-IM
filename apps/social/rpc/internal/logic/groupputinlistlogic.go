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

type GroupPutinListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinListLogic {
	return &GroupPutinListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 群申请记录
func (l *GroupPutinListLogic) GroupPutinList(in *social.GroupPutinListReq) (*social.GroupPutinListResp, error) {
	// todo: add your logic here and delete this line

	// 调用 ListNoHandler 方法，获取指定群组中所有未处理的群组申请列表
	groupReqs, err := l.svcCtx.GroupRequestsModel.ListNoHandler(l.ctx, in.GroupId)
	if err != nil {
		// 如果获取过程中出现错误，返回带错误信息的响应
		return nil, errors.Wrapf(xerr.NewDBErr(), "list group req err %v req %v", err, in.GroupId)
	}

	// 定义一个空的响应列表变量，用于存储转换后的群组申请记录
	var respList []*social.GroupRequests
	// 使用 copier 工具，将数据库查询结果复制到响应列表中
	copier.Copy(&respList, groupReqs)

	// 返回包含群组申请记录的响应
	return &social.GroupPutinListResp{
		List: respList,
	}, nil
}
