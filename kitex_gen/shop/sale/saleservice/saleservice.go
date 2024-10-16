// Code generated by Kitex v0.10.3. DO NOT EDIT.

package saleservice

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	sale "productshop/kitex_gen/shop/sale"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"GetSale": kitex.NewMethodInfo(
		getSaleHandler,
		newSaleServiceGetSaleArgs,
		newSaleServiceGetSaleResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	saleServiceServiceInfo                = NewServiceInfo()
	saleServiceServiceInfoForClient       = NewServiceInfoForClient()
	saleServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return saleServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return saleServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return saleServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "SaleService"
	handlerType := (*sale.SaleService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "sale",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.10.3",
		Extra:           extra,
	}
	return svcInfo
}

func getSaleHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*sale.SaleServiceGetSaleArgs)
	realResult := result.(*sale.SaleServiceGetSaleResult)
	success, err := handler.(sale.SaleService).GetSale(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSaleServiceGetSaleArgs() interface{} {
	return sale.NewSaleServiceGetSaleArgs()
}

func newSaleServiceGetSaleResult() interface{} {
	return sale.NewSaleServiceGetSaleResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) GetSale(ctx context.Context, req *sale.GetSaleReq) (r *sale.GetSaleResp, err error) {
	var _args sale.SaleServiceGetSaleArgs
	_args.Req = req
	var _result sale.SaleServiceGetSaleResult
	if err = p.c.Call(ctx, "GetSale", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
