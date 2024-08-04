package service

import (
	"context"
	"productshop/kitex_gen/shop/product"
	"productshop/kitex_gen/shop/sale"
	"productshop/kitex_gen/shop/validate"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/mq"
	"productshop/product_shop/middleware/mq_content"
	"productshop/product_shop/shop_sale/rpc"
)

type SaleService struct {
}

func NewSaleService() *SaleService {
	return &SaleService{}
}

func (s *SaleService) BuyProduct(ctx context.Context, req *sale.GetSaleReq) (resp *sale.GetSaleResp, err error) {
	userID := req.GetUserID()
	productID := req.GetProductID()
	productNum := req.GetProductNum()
	userCookie := req.GetUserCookie()

	// 1. 权限认证
	validateCli := rpc.GetValidateClient()

	validateReq := &validate.GetValidateReq{
		UserID:     userID,
		UserCookie: userCookie,
	}

	validateResp, err := validateCli.GetValidate(ctx, validateReq)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	if !validateResp.IsSuccess {
		return &sale.GetSaleResp{
			IsSuccess: false,
		}, nil
	}

	// 2. 商品购买
	productCli := rpc.GetProductCliernt()

	productReq := &product.GetProductReq{
		UserID:     userID,
		ProductID:  productID,
		ProductNum: 1,
	}

	productResp, err := productCli.GetProduct(ctx, productReq)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}

	// 判断数量控制接口请求状态
	if productResp.IsSuccess {
		message := mq_content.NewMessage(userID, productID, productNum)
		err = mq.SendOrderMessage(ctx, message)
		if err != nil {
			logs.Error("[SendOrderMessage] error", logs.String("error msg", err.Error()))
			return nil, err
		}
		return &sale.GetSaleResp{
			IsSuccess: true,
		}, nil
	}

	return &sale.GetSaleResp{
		IsSuccess: false,
	}, nil
}
