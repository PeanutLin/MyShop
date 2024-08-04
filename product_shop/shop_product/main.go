package main

import (
	"log"
	"net"
	product "productshop/kitex_gen/shop/product/productservice"
	"productshop/product_shop/middleware/jeager"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/redis"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	consulapi "github.com/hashicorp/consul/api"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	consul "github.com/kitex-contrib/registry-consul"
	internal_opentracing "github.com/kitex-contrib/tracer-opentracing"
)

func main() {
	// init consul
	r, err := consul.NewConsulRegister("127.0.0.1:8500", consul.WithCheck(&consulapi.AgentServiceCheck{
		Interval:                       "7s",
		Timeout:                        "5s",
		DeregisterCriticalServiceAfter: "1m",
	}))
	if err != nil {
		log.Fatal(err)
	}
	// init jeager
	_, closer := jeager.Init("shop.product")
	defer closer.Close()
	// init log
	defer logs.Init().Sync()
	// init redis
	redis.Init()

	addr, err := net.ResolveTCPAddr("tcp", ":8084")
	if err != nil {
		logs.Error(err.Error())
	}

	tracer := internal_opentracing.NewDefaultServerSuite()
	svr := product.NewServer(
		NewProductServiceImpl(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "shop.product",
		}),
		server.WithSuite(tracer),
		server.WithTracer(prometheus.NewServerTracer(":9095", "/kitexserver")),
	)

	err = svr.Run()
	if err != nil {
		logs.Error(err.Error())
	}
}
