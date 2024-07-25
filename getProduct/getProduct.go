package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"productshop/common"

	"github.com/go-redis/redis/v8"
)

// lua 脚本
var script = `
	-- redis-cli --eval product.lua product-1 3
	-- 获取商品库存信息
	local counts = redis.call("HMGET", KEYS[1], "total", "ordered")
	-- 将总库存转换为数值
	local total = tonumber(counts[1])
	-- 将已被秒杀的库存转换为数值
	local ordered = tonumber(counts[2])
	-- 获得订单购买数量
	local delta = KEYS[2]
	-- 如果当前请求的库存量加上已被秒杀的库存量仍然小于总库存量，就可以更新库存
	if ordered + delta <= total then
		-- 更新已秒杀的库存量
		redis.call("HINCRBY", KEYS[1], "ordered", delta)
		return delta
	end
	return 0
`
// redis 全局句柄
var rdb *redis.Client;


// 获取秒杀商品
func GetProductFromRedis(productName string, productNum string) bool {
	ctx := context.Background()
	fmt.Println("productName: ", productName)
	fmt.Println("productNum: ", productNum)
	result, err := rdb.Eval(ctx, script, []string{productName, productNum}).Result()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	if result == productNum {
		fmt.Println("success")
		return true
	} else {
		fmt.Println("failed")
		return false
	}
}

// 
func GetProductHandler(w http.ResponseWriter, req *http.Request) {
	// 从 URL 查询字符串中获取参数
	queryParams := req.URL.Query()

	// 获取单个参数
	productName := queryParams.Get("productName")
	productNum := queryParams.Get("productNum")
	if GetProductFromRedis(productName, productNum) {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

// http 服务
func StartHTTPServer() {
	
	http.HandleFunc("/getProduct", GetProductHandler)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("Err: ", err)
	}
}

func main() {
	// 启动 redis 客户端
	rdb = redis.NewClient(&redis.Options{
		Addr:     common.RedisHost + ":" + common.RedisPort,
		Password: common.RedisPassword,
		DB:       0,  // use default DB
	})

	StartHTTPServer()

}