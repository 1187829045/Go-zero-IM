package logic

import (
	"context"
	"database/sql"
	"errors"
	"llb-chat/apps/user/models"
	"llb-chat/pkg/ctxdata"
	"llb-chat/pkg/encrypt"
	"llb-chat/pkg/wuid"
	"time"

	"llb-chat/apps/user/rpc/internal/svc"
	"llb-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegister = errors.New("手机号已经注册过")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line

	// 1. 验证用户是否注册，根据手机号码验证
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil && err != models.ErrNotFound {
		// 如果发生错误且错误不是用户未找到，返回错误
		return nil, err
	}

	if userEntity != nil {
		// 如果用户实体不为空，说明用户已经注册，返回 ErrPhoneIsRegister 错误
		return nil, ErrPhoneIsRegister
	}

	// 2. 定义用户数据
	userEntity = &models.Users{
		// 生成用户 ID
		Id: wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		// 设置头像
		Avatar: in.Avatar,
		// 设置昵称
		Nickname: in.Nickname,
		// 设置手机号
		Phone: in.Phone,
		// 设置性别
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	// 3. 如果密码长度大于 0，生成密码哈希
	if len(in.Password) > 0 {
		genPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			// 如果生成密码哈希出错，返回错误
			return nil, err
		}
		userEntity.Password = sql.NullString{
			String: string(genPassword),
			Valid:  true,
		}
	}

	// 4. 插入用户数据到数据库
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		// 如果插入数据出错，返回错误
		return nil, err
	}

	// 5. 生成 JWT token
	now := time.Now().Unix() // 获取当前时间的 Unix 时间戳
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		// 如果生成 token 出错，返回错误
		return nil, err
	}

	// 返回注册响应，包括 token 和过期时间
	return &user.RegisterResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
