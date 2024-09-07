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

// Login 方法是用于处理用户登录的业务逻辑。
// 它接收一个 LoginReq 类型的请求参数，并返回 LoginResp 类型的响应和一个错误信息。
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {

	// 调用服务上下文 svcCtx 中的 User 服务的 Login 方法进行用户登录验证。
	// 传入的参数是用户请求中的手机号码和密码，Phone 和 Password。
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})

	// 如果在调用 User 服务的 Login 方法时发生错误（如登录失败、服务不可用等），
	// 立即返回错误信息 err，并终止后续的处理。
	if err != nil {
		return nil, err
	}

	// 定义一个 LoginResp 类型的变量 res，用于存储转换后的登录响应数据。
	var res types.LoginResp

	// 使用 copier.Copy 方法将 loginResp 中的内容复制到 res 中。
	// copier 是一个可以轻松地将一个结构体复制到另一个结构体中的库。
	copier.Copy(&res, loginResp)

	// 登录成功后，处理后续业务逻辑。
	// 在 Redis 中存储用户的在线状态。
	// 使用 Redis 的 HsetCtx 方法，在 constants.REDIS_ONLINE_USER 这个 hash 表中，
	// 将登录响应 loginResp 中的用户 ID 设置为 "1"，表示该用户在线。
	l.svcCtx.Redis.HsetCtx(l.ctx, constants.REDIS_ONLINE_USER, loginResp.Id, "1")

	// 最后，返回登录响应 res 和 nil 错误信息，表示登录操作成功。
	return &res, nil
}
