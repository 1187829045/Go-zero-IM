package handler

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"llb-chat/apps/im/ws/internal/svc"
	"llb-chat/pkg/ctxdata"
	"net/http"
)

// 结构体用于处理 JWT 认证
type JwtAuth struct {
	svc    *svc.ServiceContext
	parser *token.TokenParser // 用于解析 JWT 的工具
	logx.Logger
}

// 创建一个新的 JwtAuth 实例
func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

// 进行 JWT 认证，检查请求中的 JWT 是否有效
func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	// 从请求中解析 JWT
	tok, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		// 解析 JWT 失败，记录错误日志
		j.Errorf("parse token err %v ", err)
		return false
	}

	// 检查 JWT 是否有效
	if !tok.Valid {
		return false
	}

	// 将 JWT 中的声明解析为 MapClaims
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// 将用户识别信息添加到请求上下文中
	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))

	return true
}

// UserId 从请求中提取用户 ID
func (j *JwtAuth) UserId(r *http.Request) string {
	// 从请求上下文中获取用户 ID
	return ctxdata.GetUId(r.Context())
}
