package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	// 嵌入式结构体 rest.RestConf，包含 REST 服务的相关配置
	Database string
	// 数据库连接字符串，用于连接数据库

	Redisx redis.RedisConf
	// Redis 配置，类型为 redis.RedisConf，包含 Redis 连接所需的配置信息

	UserRpc zrpc.RpcClientConf
	// UserRpc 服务的 RPC 客户端配置，类型为 zrpc.RpcClientConf，包含 RPC 连接所需的配置信息

	JwtAuth struct {
		AccessSecret string
		// 用于 JWT 认证的密钥

		AccessExpire int64
		// JWT 令牌的过期时间，以秒为单位，已注释掉
	}
}
