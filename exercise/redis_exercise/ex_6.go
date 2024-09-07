package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// 展示了如何使用 Redis 的 TxPipeline 和 TxPipelined 方法来批量执行事务命令
var rdb4 *redis.Client

// 初始化 Redis 客户端
func initRedisClient4() {
	rdb4 = redis.NewClient(&redis.Options{
		Addr:     "192.168.128.128:6379", // Redis 服务器地址
		Password: "llb-easy-chat",        // 密码
		DB:       0,                      // 数据库索引
		PoolSize: 20,                     // 连接池大小
	})
}

// txPipelineDemo 使用事务管道示例
func txPipelineDemo() {
	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // 确保在函数结束时取消上下文

	// 创建一个事务管道
	pipe := rdb4.TxPipeline()

	// 增加计数器值
	incr := pipe.Incr(ctx, "tx_pipeline_counter")
	// 设置计数器过期时间为1小时
	pipe.Expire(ctx, "tx_pipeline_counter", time.Hour)

	// 执行管道中的命令
	_, err := pipe.Exec(ctx)
	if err != nil {
		// 如果执行过程中出现错误，panic抛出错误
		panic(err)
	}

	// 打印增量后的值和执行结果
	fmt.Println("TxPipeline Counter Value:", incr.Val(), "Error:", err)
}

// txPipelinedDemo 使用 TxPipelined 方法示例
func txPipelinedDemo() {
	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // 确保在函数结束时取消上下文

	// 定义一个指针变量来存储增量命令的返回结果
	var incr2 *redis.IntCmd

	// 使用 TxPipelined 方法执行事务管道
	_, err := rdb4.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		// 增加计数器值
		incr2 = pipe.Incr(ctx, "tx_pipeline_counter")
		// 设置计数器过期时间为1小时
		pipe.Expire(ctx, "tx_pipeline_counter", time.Hour)
		return nil
	})
	if err != nil {
		// 如果执行过程中出现错误，panic抛出错误
		panic(err)
	}

	// 打印增量后的值和执行结果
	fmt.Println("TxPipelined Counter Value:", incr2.Val(), "Error:", err)
}

func main() {
	// 初始化 Redis 客户端
	initRedisClient4()
	// 调用 txPipelineDemo 函数
	txPipelineDemo()
	// 调用 txPipelinedDemo 函数
	txPipelinedDemo()
}
