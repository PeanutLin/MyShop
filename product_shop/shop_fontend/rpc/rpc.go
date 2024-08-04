package rpc

import (
	"productshop/kitex_gen/shop/sale/saleservice"

	"github.com/cloudwego/kitex/client"
	consulapi "github.com/hashicorp/consul/api"

	prometheus "github.com/kitex-contrib/monitor-prometheus"
	consul "github.com/kitex-contrib/registry-consul"
	internal_opentracing "github.com/kitex-contrib/tracer-opentracing"
)

var (
	SaleClient saleservice.Client
)

func MustInit() {
	var err error

	// init consul
	consulConfig := consulapi.Config{
		Address: "127.0.0.1:8500",
	}
	r, err := consul.NewConsulResolverWithConfig(&consulConfig)
	if err != nil {
		panic(err)
	}

	// init OpenTracing
	tracer := internal_opentracing.NewDefaultClientSuite()
	shopServiceNmae := "shop.sale"
	SaleClient, err = saleservice.NewClient(
		// 服务名称
		shopServiceNmae,
		// 服务发现
		client.WithResolver(r),
		// 服务追踪
		client.WithSuite(tracer),
		// 服务监控
		client.WithTracer(prometheus.NewClientTracer(":9091", "/kitexclient")),
	)
	if err != nil {
		panic(err)
	}
}

func GetSaleClient() saleservice.Client {
	return SaleClient
}
