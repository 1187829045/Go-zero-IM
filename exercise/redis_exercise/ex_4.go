package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// 使用Scan命令来遍历所有符合要求的 key
var rdb2 *redis.Client

// 初始化 Redis 客户端
func initRedisClient2() {
	rdb2 = redis.NewClient(&redis.Options{
		Addr:     "192.168.128.128:6379", // Redis 服务器地址
		Password: "llb-easy-chat",        // 密码
		DB:       0,                      // 数据库索引
		PoolSize: 20,                     // 连接池大小
	})
}

// scanKeysDemo1 按前缀查找所有key示例
func scanKeysDemo1() {
	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // 确保在函数结束时取消上下文

	// 初始化游标为0
	var cursor uint64
	for {
		var keys []string
		var err error

		// 扫描以 "prefix:" 开头的键，获取键名和新的游标
		keys, cursor, err = rdb2.Scan(ctx, cursor, "prefix:*", 0).Result()
		if err != nil {
			// 如果出现错误，则panic抛出错误
			panic(err)
		}

		// 打印扫描到的键名
		for _, key := range keys {
			fmt.Println("key", key)
		}

		// 如果游标为0，表示没有更多的键需要扫描，跳出循环
		if cursor == 0 {
			break
		}
	}
}

// scanKeysDemo2 按前缀扫描key示例
func scanKeysDemo2() {
	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // 确保在函数结束时取消上下文

	// 按前缀扫描key，初始化一个新的迭代器
	iter := rdb2.Scan(ctx, 0, "prefix:*", 0).Iterator()
	for iter.Next(ctx) { // 遍历迭代器，获取扫描到的键
		fmt.Println("keys", iter.Val()) // 打印当前键的值
	}

	// 检查迭代器是否出现错误
	if err := iter.Err(); err != nil {
		// 如果有错误，则panic抛出错误
		panic(err)
	}
}
func main() {
	// 初始化 Redis 客户端
	initRedisClient2()
	// 调用 scanKeysDemo1 函数
	scanKeysDemo1()
	scanKeysDemo2()

}
