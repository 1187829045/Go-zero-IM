package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var rdb1 *redis.Client

func initRedisClient() {
	// 初始化 Redis 客户端
	rdb1 = redis.NewClient(&redis.Options{
		Addr:     "192.168.128.128:6379", // Redis 服务器的地址和端口
		Password: "llb-easy-chat",        // Redis 服务器的密码
		DB:       0,                      // 使用的数据库索引（默认数据库 0）
		PoolSize: 20,                     // Redis 连接池的大小
	})
}

// getValueFromRedis redis.Nil判断
func getValueFromRedis(key, defaultValue string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // 在函数退出时调用 cancel 以释放资源

	val, err := rdb1.Get(ctx, key).Result()
	if err != nil {
		// 如果返回的错误是key不存在
		if errors.Is(err, redis.Nil) {
			return defaultValue, nil
		}
		// 出其他错了
		return "", err
	}
	return val, nil
}

func main() {
	// 初始化 Redis 客户端
	initRedisClient()

	// 测试 getValueFromRedis 函数
	key := "key"
	defaultValue := "default"
	value, err := getValueFromRedis(key, defaultValue)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", value)
	}
}
