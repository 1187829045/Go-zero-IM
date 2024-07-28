package user

import (
	"context"
	"github.com/jinzhu/copier"
	"llb-chat/apps/user/rpc/user"
	"llb-chat/pkg/constants"

	"llb-chat/apps/user/api/internal/svc"
	"llb-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 创建 Login 方法，接收一个 LoginReq 类型的请求参数，返回 LoginResp 类型的响应和一个错误
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {

	// 调用 svcCtx 中的 User 服务的 Login 方法进行用户登录，传入用户请求的 Phone 和 Password
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	var res types.LoginResp
	copier.Copy(&res, loginResp)

	// 处理登录后的业务逻辑，将用户登录状态存储到 Redis 中
	// 使用 Redis 的 HsetCtx 方法，在 constants.REDIS_ONLINE_USER 这个 hash 中，设置 loginResp.Id 为 "1"
	l.svcCtx.Redis.HsetCtx(l.ctx, constants.REDIS_ONLINE_USER, loginResp.Id, "1")

	// 返回登录响应 res 和 nil 错误信息
	return &res, nil
}
