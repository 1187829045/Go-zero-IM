package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/proc"
	"llb-chat/pkg/configserver"
	"sync"

	"llb-chat/apps/im/api/internal/config"
	"llb-chat/apps/im/api/internal/handler"
	"llb-chat/apps/im/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

// 定义一个等待组，用于同步多个 goroutine
var wg sync.WaitGroup

func main() {
	// 解析命令行标志
	flag.Parse()

	// 定义配置结构体变量 c
	var c config.Config

	// 创建配置服务器并加载配置文件
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "192.168.128.128:3379",             // etcd 服务器地址
		ProjectKey:     "2f5bb7747efda0546636fb385a3fa593", // 项目密钥
		Namespace:      "im",                               // 命名空间
		Configs:        "im-api.yaml",                      // 配置文件名
		ConfigFilePath: "./etc/conf",                       // 配置文件路径
		LogLevel:       "DEBUG",                            // 日志级别
	})).MustLoad(&c, func(bytes []byte) error {
		// 加载配置文件内容并解析成配置结构体
		var c config.Config
		configserver.LoadFromJsonBytes(bytes, &c)

		// 包装过程
		proc.WrapUp()

		// 增加一个等待组计数
		wg.Add(1)
		go func(c config.Config) {
			defer wg.Done() // 减少等待组计数

			Run(c) // 运行服务
		}(c)
		return nil
	})
	if err != nil {
		panic(err) // 如果加载配置文件出错，则程序崩溃
	}

	// 增加一个等待组计数
	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done() // 减少等待组计数

		Run(c) // 运行服务
	}(c)

	// 等待所有等待组计数归零
	wg.Wait()
}

// Run 函数启动 REST 服务器
func Run(c config.Config) {
	// 创建并启动 REST 服务器
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop() // 在函数退出时停止服务器

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)
	// 注册处理程序
	handler.RegisterHandlers(server, ctx)

	// 打印服务器启动信息
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	// 启动服务器
	server.Start()
}
