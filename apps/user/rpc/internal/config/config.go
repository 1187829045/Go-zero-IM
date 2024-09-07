package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

//zrpc.RpcServerConf结构体字段

//type RpcServerConf struct {
//	ListenOn string // RPC 服务器监听的地址和端口
//	Etcd     struct {
//		Hosts []string // Etcd 集群的地址
//		Key   string   // 服务在 Etcd 中的注册键
//	}
//}
//type RedisConf struct {
//	Addr     string // Redis 服务器的地址，例如 "localhost:6379"
//	Password string // Redis 认证密码（如果设置了）
//	DB       int    // Redis 数据库索引，默认是 0
//	PoolSize int    // Redis 连接池大小
//}
//type CacheConf struct {
//	// 示例配置字段
//	Type   string // 缓存类型，例如 "memory" 或 "redis"
//	Expiry int64  // 缓存数据的过期时间（秒）
//}

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
