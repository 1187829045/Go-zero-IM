package main

import (
	"flag"
	"fmt"
	"llb-chat/pkg/interceptor/rpcserver"

	"llb-chat/apps/user/rpc/internal/config"
	"llb-chat/apps/user/rpc/internal/server"
	"llb-chat/apps/user/rpc/internal/svc"
	"llb-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 定义一个命令行标志 -f，用于指定配置文件路径，默认值是 "etc/dev/user.yaml"，描述是 "the config file"
var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

func main() {
	// 解析命令行标志
	flag.Parse()

	// 创建一个 Config 类型的变量 c
	var c config.Config

	// 加载配置文件到 c 变量中，使用 conf.MustLoad 函数来加载配置文件，失败时会引发 panic
	conf.MustLoad(*configFile, &c)

	// 使用加载的配置创建一个服务上下文 ctx
	ctx := svc.NewServiceContext(c)

	// 设置根令牌，如果设置失败，则引发 panic
	if err := ctx.SetRootToken(); err != nil {
		panic(err)
	}

	// 创建一个新的 gRPC 服务器，配置来自 c.RpcServerConf
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// 注册 User 服务到 gRPC 服务器中
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		// 如果模式是开发模式或测试模式，则注册 gRPC 反射服务
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	// 添加一个 Unary 拦截器，用于日志记录
	s.AddUnaryInterceptors(rpcserver.LogInterceptor)

	// 确保在 main 函数结束时停止 gRPC 服务器
	defer s.Stop()

	// 打印服务器启动的地址
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)

	// 启动 gRPC 服务器
	s.Start()
}
