package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// 操作有序集合
var rdb *redis.Client

func initRedisClient1() {
	// 初始化 Redis 客户端
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.128.128:6379", // Redis 服务器的地址和端口
		Password: "llb-easy-chat",        // Redis 服务器的密码
		DB:       0,                      // 使用的数据库索引（默认数据库 0）
		PoolSize: 20,                     // Redis 连接池的大小
	})
}

// zsetDemo 操作zset示例
func zsetDemo() {
	// 定义zset的key
	zsetKey := "language_rank"
	// 定义zset的value
	languages := []*redis.Z{
		{Score: 90.0, Member: "Golang"},
		{Score: 98.0, Member: "Java"},
		{Score: 95.0, Member: "Python"},
		{Score: 97.0, Member: "JavaScript"},
		{Score: 99.0, Member: "C/C++"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// ZADD命令添加元素到zset中
	err := rdb.ZAdd(ctx, zsetKey, languages...).Err()
	if err != nil {
		fmt.Printf("zadd failed, err:%v\n", err)
		return
	}
	fmt.Println("zadd success")

	// ZINCRBY命令将Golang的分数增加10
	newScore, err := rdb.ZIncrBy(ctx, zsetKey, 10.0, "Golang").Result()
	if err != nil {
		fmt.Printf("zincrby failed, err:%v\n", err)
		return
	}
	fmt.Printf("Golang's score is %f now.\n", newScore)

	// ZREVRANGE命令获取分数最高的3个元素
	ret := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Val()
	fmt.Println("Top 3 languages:")
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}

	// ZRANGEBYSCORE命令获取分数在95到100之间的元素
	op := &redis.ZRangeBy{
		Min: "95",
		Max: "100",
	}
	ret, err = rdb.ZRangeByScoreWithScores(ctx, zsetKey, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}
	fmt.Println("Languages with scores between 95 and 100:")
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
}

func main() {
	// 初始化 Redis 客户端
	initRedisClient1()
	// 调用zsetDemo函数
	zsetDemo()
}
