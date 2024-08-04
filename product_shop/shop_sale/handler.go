package main

import (
	"context"
	"fmt"
	"productshop/kitex_gen/shop/sale"
	"productshop/product_shop/shop_sale/service"
)

// SaleServiceImpl implements the last service interface defined in the IDL.
type SaleServiceImpl struct {
	SaleService *service.SaleService
}

func NewSaleServiceImpl() *SaleServiceImpl {
	return &SaleServiceImpl{
		SaleService: service.NewSaleService(),
	}
}

// GetSale implements the SaleServiceImpl interface.
func (s *SaleServiceImpl) GetSale(ctx context.Context, req *sale.GetSaleReq) (resp *sale.GetSaleResp, err error) {
	fmt.Print("get dddd")
	return s.SaleService.BuyProduct(ctx, req)
}
