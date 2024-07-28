package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"llb-chat/apps/user/api/internal/config"
	"llb-chat/apps/user/rpc/userclient"
)

// 定义了一个重试策略的 JSON 字符串
// "methodConfig" 配置了服务名称和重试策略
// "service" 指定了服务名称 "user.User"
// "waitForReady" 设置为 true 表示等待服务准备就绪
// "retryPolicy" 包含最大尝试次数5、初始回退时间0.001、最大回退时间0.002、回退倍数1.0和可重试状态码UNKNOWN

var retryPolicy = `{
	"methodConfig" : [{
		"name": [{
			"service": "user.User"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN"]
		}
	}]
}`

// 定义了一个 ServiceContext 结构体
// 包含配置项 Config，Redis 客户端指针 *redis.Redis，以及用户客户端 userclient.User

type ServiceContext struct {
	Config config.Config

	*redis.Redis
	userclient.User
}

// NewServiceContext 是一个工厂函数，用于创建新的 ServiceContext 实例
// 接受一个配置参数 c，并返回一个初始化好的 ServiceContext 实例
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		Redis: redis.MustNewRedis(c.Redisx),
		// 初始化 Redis 客户端

		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithDialOption(grpc.WithDefaultServiceConfig(
			retryPolicy)))),
		// 初始化用户客户端 应用重试策略 retryPolicy
	}
}
