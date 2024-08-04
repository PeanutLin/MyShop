package main

import (
	"context"
	"productshop/kitex_gen/shop/validate"
	"productshop/product_shop/shop_validate/service"
)

// ValidateServiceImpl implements the last service interface defined in the IDL.
type ValidateServiceImpl struct {
	valiedateService *service.ValiedateService
}

func NewValidateServiceImpl() *ValidateServiceImpl {
	return &ValidateServiceImpl{
		valiedateService: service.NewValiedateService(),
	}
}

// GetValidate implements the ValidateServiceImpl interface.
func (s *ValidateServiceImpl) GetValidate(ctx context.Context, req *validate.GetValidateReq) (resp *validate.GetValidateResp, err error) {
	return s.valiedateService.Validate(ctx, req)
}
