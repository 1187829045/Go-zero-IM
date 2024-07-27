package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"llb-chat/apps/user/api/internal/config"
	"llb-chat/apps/user/rpc/userclient"
	// N * client =》 别名
)

type ServiceContext struct {
	Config config.Config

	*redis.Redis
	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		Redis: redis.MustNewRedis(c.Redisx),
		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithDialOption(grpc.WithDefaultServiceConfig(
			retryPolicy)))),
	}
}
