package rpc

import (
	"productshop/kitex_gen/shop/product/productservice"
	"productshop/kitex_gen/shop/validate/validateservice"

	"github.com/cloudwego/kitex/client"
	consulapi "github.com/hashicorp/consul/api"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	consul "github.com/kitex-contrib/registry-consul"
	internal_opentracing "github.com/kitex-contrib/tracer-opentracing"
)

var (
	ValidateClient validateservice.Client
	ProductCliernt productservice.Client
)

func MustInit() {
	var err error
	tracer := internal_opentracing.NewDefaultClientSuite()

	// kitex 内置 prometheus 监控
	prometheusCli := prometheus.NewClientTracer(":9093", "/kitexclient")

	consulConfig := consulapi.Config{
		Address: "127.0.0.1:8500",
	}
	r, err := consul.NewConsulResolverWithConfig(&consulConfig)
	if err != nil {
		panic(err)
	}

	ValidateClient, err = validateservice.NewClient(
		"shop.validate",
		// 服务发现
		client.WithResolver(r),
		// 服务追踪
		client.WithSuite(tracer),
		// 服务监控
		client.WithTracer(prometheusCli),
	)
	if err != nil {
		panic(err)
	}

	ProductCliernt, err = productservice.NewClient(
		"shop.product",
		// 服务发现
		client.WithResolver(r),
		// 服务追踪
		client.WithSuite(tracer),
		// 服务监控
		client.WithTracer(prometheusCli),
	)
	if err != nil {
		panic(err)
	}
}

func GetValidateClient() validateservice.Client {
	return ValidateClient
}

func GetProductCliernt() productservice.Client {
	return ProductCliernt
}
