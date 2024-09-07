package ctxdata

import "github.com/golang-jwt/jwt"

// Identify 是一个常量，代表 JWT 声明中的一个键，其值为 "llb"
const Identify = "llb"

// GetJwtToken 生成一个 JWT 令牌
// secretKey: 用于签名 JWT 的密钥
// iat: JWT 的签发时间（UNIX 时间戳）
// seconds: JWT 的有效期（从签发时间起的秒数）
// uid: 用户标识符，作为 JWT 的一个声明
// 返回生成的 JWT 令牌字符串和可能的错误
func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	// 创建一个 JWT 声明对象 (MapClaims 是一个类型，表示一个键值对的声明集合)
	claims := make(jwt.MapClaims)

	// 设置 JWT 的过期时间 (exp) 为签发时间 (iat) 加上有效期 (seconds)
	claims["exp"] = iat + seconds

	// 设置 JWT 的签发时间 (iat)
	claims["iat"] = iat

	// 设置一个自定义声明，使用常量 Identify 作为键，将用户标识符 uid 作为值
	claims[Identify] = uid

	// 创建一个新的 JWT 对象，指定签名方法为 HMAC SHA-256
	token := jwt.New(jwt.SigningMethodHS256)

	// 将声明赋值给 JWT 对象
	token.Claims = claims

	// 使用指定的密钥对 JWT 进行签名，并返回生成的 JWT 字符串
	// SignedString 方法返回生成的 JWT 字符串和可能的错误
	return token.SignedString([]byte(secretKey))
}
