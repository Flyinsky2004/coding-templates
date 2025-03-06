/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package config

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RedisClient 全局 Redis 客户端
var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() {
	db, _ := strconv.Atoi(fmt.Sprintf("%v", Config.Redis.Database))
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Config.Redis.Host, Config.Redis.Port),
		Password: Config.Redis.Password,
		DB:       db,
	})

	// 测试连接
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}
