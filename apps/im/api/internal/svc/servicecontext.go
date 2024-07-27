package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"llb-chat/apps/im/api/internal/config"
	"llb-chat/apps/im/rpc/imclient"
	"llb-chat/apps/social/rpc/socialclient"
	"llb-chat/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config

	imclient.Im
	userclient.User
	socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
}
