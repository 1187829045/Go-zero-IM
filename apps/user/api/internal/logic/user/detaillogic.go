package user

import (
	"context"
	"github.com/jinzhu/copier"
	"llb-chat/apps/user/rpc/user"
	"llb-chat/pkg/ctxdata"

	"llb-chat/apps/user/api/internal/svc"
	"llb-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// 定义 DetailLogic 结构体
// 包含 Logger、context.Context 和 *svc.ServiceContext

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDetailLogic 是一个工厂函数，用于创建 DetailLogic 实例
// 接受 context.Context 和 *svc.ServiceContext 参数，并返回初始化的 DetailLogic 实例
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//通过用户id查询到用户信息，并返回

func (l *DetailLogic) Detail(req *types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line
	// 从上下文中获取用户 ID
	uid := ctxdata.GetUId(l.ctx)
	// 调用用户客户端的 GetUserInfo 方法，传入用户 ID，获取用户信息
	userInfoResp, err := l.svcCtx.User.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}

	var res types.User
	copier.Copy(&res, userInfoResp.User)
	return &types.UserInfoResp{Info: res}, nil
}
