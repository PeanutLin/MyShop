package controllers

import (
	service "productshop/product_shop/base_service"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService service.IOrderService
}

// func (o *OrderController) GetAll() mvc.View {
// 	orderMap, err := o.OrderService.GetAllOrderInfo()
// 	if err != nil {
// 		logs.Error(err.Error())
// 	}

// 	return mvc.View{
// 		Name: "order/view.html",
// 		Data: iris.Map{
// 			"order": orderMap,
// 		},
// 	}
// }

func (o *OrderController) GetIndex() mvc.View {
	return mvc.View{
		Name: "order/index.html",
	}
}
