package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// 使用pipe
var rdb3 *redis.Client

// 初始化 Redis 客户端
func initRedisClient3() {
	rdb3 = redis.NewClient(&redis.Options{
		Addr:     "192.168.128.128:6379", // Redis 服务器地址
		Password: "llb-easy-chat",        // 密码
		DB:       0,                      // 数据库索引
		PoolSize: 20,                     // 连接池大小
	})
}

// pipelineDemo 使用管道示例
func pipelineDemo() {
	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // 确保在函数结束时取消上下文

	// 创建一个管道
	pipe := rdb3.Pipeline()

	// 增加计数器值
	incr := pipe.Incr(ctx, "pipeline_counter")
	// 设置计数器过期时间为1小时
	pipe.Expire(ctx, "pipeline_counter", time.Hour)

	// 执行管道中的命令
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		// 如果执行过程中出现错误，panic抛出错误
		panic(err)
	}

	// 打印执行管道命令后的返回结果
	fmt.Println(cmds)

	// 在执行pipe.Exec之后才能获取到结果，获取并打印增量后的值
	fmt.Println("Pipeline Counter Value:", incr.Val())
}

func main() {
	// 初始化 Redis 客户端
	initRedisClient3()
	// 调用 pipelineDemo 函数
	pipelineDemo()
}
