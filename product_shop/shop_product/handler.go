package main

import (
	"context"
	product "productshop/kitex_gen/shop/product"
	"productshop/product_shop/shop_product/service"
)

// ProductServiceImpl implements the last service interface defined in the IDL.
type ProductServiceImpl struct {
	productService *service.ProductService
}

func NewProductServiceImpl() *ProductServiceImpl {
	return &ProductServiceImpl{
		productService: service.NewProductService(),
	}
}

// GetProduct implements the ProductServiceImpl interface.
func (p *ProductServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	return p.productService.SolveProductFromRedis(ctx, req)

}
