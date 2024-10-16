// Code generated by Kitex v0.10.3. DO NOT EDIT.

package validateservice

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	validate "productshop/kitex_gen/shop/validate"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"GetValidate": kitex.NewMethodInfo(
		getValidateHandler,
		newValidateServiceGetValidateArgs,
		newValidateServiceGetValidateResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	validateServiceServiceInfo                = NewServiceInfo()
	validateServiceServiceInfoForClient       = NewServiceInfoForClient()
	validateServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return validateServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return validateServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return validateServiceServiceInfoForClient
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
	serviceName := "ValidateService"
	handlerType := (*validate.ValidateService)(nil)
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
		"PackageName": "validate",
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

func getValidateHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*validate.ValidateServiceGetValidateArgs)
	realResult := result.(*validate.ValidateServiceGetValidateResult)
	success, err := handler.(validate.ValidateService).GetValidate(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newValidateServiceGetValidateArgs() interface{} {
	return validate.NewValidateServiceGetValidateArgs()
}

func newValidateServiceGetValidateResult() interface{} {
	return validate.NewValidateServiceGetValidateResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) GetValidate(ctx context.Context, req *validate.GetValidateReq) (r *validate.GetValidateResp, err error) {
	var _args validate.ValidateServiceGetValidateArgs
	_args.Req = req
	var _result validate.ValidateServiceGetValidateResult
	if err = p.c.Call(ctx, "GetValidate", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
