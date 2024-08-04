package main

import (
	"log"
	"net"
	sale "productshop/kitex_gen/shop/sale/saleservice"
	"productshop/product_shop/middleware/jeager"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/mq"
	"productshop/product_shop/shop_sale/rpc"

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
	// 初始化 Jeager
	_, closer := jeager.Init("shop.sale")
	defer closer.Close()
	// init log
	defer logs.Init().Sync()
	// init rpc client
	rpc.MustInit()
	// init mq
	mq.MustInitBookProductProducer()

	addr, err := net.ResolveTCPAddr("tcp", ":8082")
	if err != nil {
		logs.Error(err.Error())
	}

	tracer := internal_opentracing.NewDefaultServerSuite()
	svr := sale.NewServer(
		NewSaleServiceImpl(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "shop.sale",
		}),
		server.WithSuite(tracer),
		server.WithTracer(prometheus.NewServerTracer(":9092", "/kitexserver")),
	)
	err = svr.Run()

	if err != nil {
		logs.Error(err.Error())
	}
}
