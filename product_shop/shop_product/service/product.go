package service

import (
	"context"
	"productshop/kitex_gen/shop/product"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/redis"
	"strconv"

	"github.com/pkg/errors"
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

type ProductService struct {
}

func NewProductService() *ProductService {
	return &ProductService{}
}

// Redis 缓存处理订单
func (p *ProductService) SolveProductFromRedis(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	productRedisID := "product-" + strconv.FormatInt(req.GetProductID(), 10)
	productNum := strconv.FormatInt(req.GetProductNum(), 10)

	redisCli := redis.GetRedisClient()
	result, err := redisCli.Eval(script, []string{productRedisID, productNum}).Result()

	if err != nil {
		logs.Error("redis eval error", logs.String("error msg", err.Error()))
		return nil, errors.Wrap(err, "redis failed")
	}

	if result == productNum {
		logs.Info("Buy Product Success")
		return &product.GetProductResp{
			IsSuccess: true,
		}, nil
	} else {
		logs.Info("Buy Product Failed")
		return &product.GetProductResp{
			IsSuccess: false,
		}, nil
	}
}
