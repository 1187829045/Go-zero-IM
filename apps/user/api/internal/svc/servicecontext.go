package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"llb-chat/apps/user/api/internal/config"
	"llb-chat/apps/user/rpc/userclient"
	// N * client =》 别名
)

type ServiceContext struct {
	Config config.Config

	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
