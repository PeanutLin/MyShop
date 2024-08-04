package redis

import (
	"productshop/product_shop/common"

	"github.com/go-redis/redis"
)

// redis 全局句柄
var rdb *redis.Client

func GetRedisClient() *redis.Client {
	return rdb
}

func Init() {
	// 启动 redis 客户端
	rdb = redis.NewClient(&redis.Options{
		Addr:     common.RedisHost + ":" + common.RedisPort,
		Password: common.RedisPassword,
		DB:       0, // use default DB
	})
}
