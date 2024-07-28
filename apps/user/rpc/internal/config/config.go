package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	Mysql struct {
		DataSource string
	}

	Cache cache.CacheConf

	Redisx redis.RedisConf

	Jwt struct {
		AccessSecret string // 签署JWT所用的密钥，确保生成的令牌是可信的且未被篡改
		AccessExpire int64  // 令牌的有效期（以秒为单位），定义JWT的过期时间以增强安全性
	}
}
