package main

import (
	"flag" // 导入 flag 包，用于解析命令行参数
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx" // 用于处理 HTTP 请求
	"llb-chat/pkg/resultx"

	"llb-chat/apps/user/api/internal/config"
	"llb-chat/apps/user/api/internal/handler"
	"llb-chat/apps/user/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf" // 导入 go-zero 的配置包，用于加载配置
	"github.com/zeromicro/go-zero/rest"      // 导入 go-zero 的 REST 包，用于创建 REST 服务器
)

// 定义 configFile 变量，用于存储配置文件的路径，默认值为 "etc/dev/user.yaml"
var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

func main() {
	// 解析命令行参数
	flag.Parse()

	var c config.Config
	// 加载配置文件到 c 结构体中
	conf.MustLoad(*configFile, &c)

	// 创建一个新的 REST 服务器，使用配置中的 RestConf 部分
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 创建一个新的服务上下文，传入配置 c
	ctx := svc.NewServiceContext(c)
	// 注册处理器，将请求路径与处理函数绑定
	handler.RegisterHandlers(server, ctx)

	// 设置自定义错误处理器
	httpx.SetErrorHandlerCtx(resultx.ErrHandler(c.Name))
	// 设置自定义成功处理器
	httpx.SetOkHandler(resultx.OkHandler)

	// 打印服务器启动信息
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	// 启动服务器
	server.Start()
}
