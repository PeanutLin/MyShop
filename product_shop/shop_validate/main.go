package main

import (
	"log"
	"net"
	validate "productshop/kitex_gen/shop/validate/validateservice"
	"productshop/product_shop/middleware/jeager"
	"productshop/product_shop/middleware/logs"

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
	_, closer := jeager.Init("shop.validate")
	defer closer.Close()
	// init log
	defer logs.Init().Sync()

	addr, err := net.ResolveTCPAddr("tcp", ":8083")
	if err != nil {
		logs.Error(err.Error())
	}

	tracer := internal_opentracing.NewDefaultServerSuite()
	svr := validate.NewServer(
		NewValidateServiceImpl(),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "shop.validate",
		}),
		server.WithSuite(tracer),
		server.WithTracer(prometheus.NewServerTracer(":9094", "/kitexserver")),
	)

	err = svr.Run()
	if err != nil {
		logs.Error(err.Error())
	}
}
