package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	//嵌套了service.ServiceConf 的配置项，通常包含服务相关的基本配置
	service.ServiceConf

	// 指定服务监听的地址和端口
	ListenOn string

	// 消息聊天转发的配置
	MsgChatTransfer kq.KqConf

	// 消息阅读转发的配置
	MsgReadTransfer kq.KqConf

	// Redis 的配置
	Redisx redis.RedisConf

	// MongoDB 的配置
	Mongo struct {
		// MongoDB 的连接 URL
		Url string
		// MongoDB 数据库名称
		Db string
	}

	// 消息阅读处理的相关配置
	MsgReadHandler struct {
		// 群组消息阅读处理的 ID
		GroupMsgReadHandler int
		// 群组消息阅读记录的延迟时间，单位为秒
		GroupMsgReadRecordDelayTime int64
		// 群组消息阅读记录的延迟次数
		GroupMsgReadRecordDelayCount int
	}

	// 社交 RPC 客户端的配置
	SocialRpc zrpc.RpcClientConf

	// WebSocket 服务器的配置
	Ws struct {
		// WebSocket 服务器的主机地址
		Host string
	}
}
