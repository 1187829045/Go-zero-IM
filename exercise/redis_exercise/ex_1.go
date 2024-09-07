package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8" // 引入 go-redis 包，用于操作 Redis
	"time"
)

// doCommand go-redis基本使用示例
func doCommand1() {
	// 创建一个带有超时的上下文，这里设置超时时间为 500 毫秒
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.128.128:6379", // Redis 服务器的地址和端口
		Password: "llb-easy-chat",        // Redis 服务器的密码
		DB:       0,                      // 使用的数据库索引（默认数据库 0）
		PoolSize: 20,                     // Redis 连接池的大小
	})

	// 执行 GET 命令从 Redis 获取键 "key" 的值
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		// 如果发生错误，打印错误信息
		fmt.Println("Error getting key:", err)
	} else {
		// 如果没有错误，打印获取到的值
		fmt.Println("Value:", val)
	}

	// 先获取到命令对象
	cmder := rdb.Get(ctx, "key")
	if cmder.Err() != nil {
		// 如果发生错误，打印错误信息
		fmt.Println("Error getting key (using cmder):", cmder.Err())
	} else {
		// 如果没有错误，打印获取到的值
		fmt.Println("Value (using cmder):", cmder.Val())
	}

	// 执行 SET 命令，将键 "key" 的值设置为 10，并指定过期时间为 1 小时
	err = rdb.Set(ctx, "key", 10, time.Hour).Err()
	if err != nil {
		// 如果发生错误，打印错误信息
		fmt.Println("Error setting key:", err)
	}

	// 再次执行 GET 命令从 Redis 获取键 "key" 的新值
	value := rdb.Get(ctx, "key").Val()
	// 打印获取到的新值
	fmt.Println("New Value:", value)
}

// doDemo rdb.Do 方法使用示例
func doDemo() {
	// 创建一个带有超时的上下文，这里设置超时时间为 500 毫秒
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // 在函数退出时调用 cancel 以释放资源

	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.128.128:6379", // Redis 服务器的地址和端口
		Password: "llb-easy-chat",        // Redis 服务器的密码
		DB:       0,                      // 使用的数据库索引（默认数据库 0）
		PoolSize: 20,                     // Redis 连接池的大小
	})

	// 直接执行 SET 命令，设置键 "key" 的值为 10，并设置过期时间为 3600 秒
	err := rdb.Do(ctx, "set", "key", 10, "EX", 3600).Err()
	if err != nil {
		// 如果发生错误，打印错误信息
		fmt.Println("Error setting key:", err)
	}

	// 执行 GET 命令获取键 "key" 的值
	val, err := rdb.Do(ctx, "get", "key").Result()
	if err != nil {
		// 如果发生错误，打印错误信息
		fmt.Println("Error getting key:", err)
	} else {
		// 如果没有错误，打印获取到的值
		fmt.Println("Value:", val)
	}
}
func main() {
	// 调用 doCommand 函数以执行 Redis 操作示例
	doCommand1()
	doDemo()
}
