package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf" // 引入 go-zero 的配置加载库
	"llb-chat/apps/im/ws/internal/config"
	"llb-chat/apps/im/ws/internal/handler"
	"llb-chat/apps/im/ws/internal/svc"
	"llb-chat/apps/im/ws/websocket"
	"time"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse() // 解析命令行标志

	// 定义一个 config.Config 类型的变量 c，用于保存配置
	var c config.Config
	// 使用 go-zero 的 conf.MustLoad 函数加载配置文件，并将内容填充到 c 变量中
	conf.MustLoad(*configFile, &c)

	// 调用配置对象的 SetUp 方法，设置go-zero日志和监听的一些处理，如果出错则 panic
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	// 创建服务上下文 ctx，传入配置对象 c
	ctx := svc.NewServiceContext(c)

	// 创建一个新的 WebSocket 服务器 srv
	// 使用各种选项进行配置：
	// - WebSocket 监听地址
	// - 身份验证处理器
	// - 消息确认方式设置为 RigorAck
	// - 最大连接空闲时间设置为 10 秒
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerAck(websocket.RigorAck),
		websocket.WithServerMaxConnectionIdle(10*time.Second),
	)

	defer srv.Stop() // 确保在 main 函数结束时停止服务器

	// 注册处理器
	handler.RegisterHandlers(srv, ctx)

	// 输出服务器启动的消息
	fmt.Println("start websocket server at ", c.ListenOn, " ..... ")

	// 启动 WebSocket 服务器
	srv.Start()
}
