package base_service

import (
	"context"
	"productshop/product_shop/base_repository"
	"productshop/product_shop/middleware/mq_content"
	"productshop/product_shop/middleware/mysql/gen/model"
	"strconv"
)

type IOrderService interface {
	// GetOrderById(int64) (*model.Order, error)
	// DeleteOrderById(int64) bool
	// UpdateOrder(order *model.Order) error
	InsertOrder(ctx context.Context, order *model.Order) (int64, error)
	// GetAllOrder() ([]*model.Order, error)
	// GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(ctx context.Context, message *mq_content.Message) (int64, error)
}

type OrderService struct {
	OrderRepository base_repository.IOrderRepository
}

func NewOrderService(repository base_repository.IOrderRepository) IOrderService {
	return &OrderService{
		OrderRepository: repository,
	}
}

// func (o *OrderService) GetOrderById(orderId int64) (*model.Order, error) {
// 	return o.OrderRepository.SelectByKey(orderId)
// }

// func (o *OrderService) DeleteOrderById(orderId int64) bool {
// 	return o.OrderRepository.Delete(orderId)
// }

// func (o *OrderService) UpdateOrder(order *model.Order) error {
// 	return o.OrderRepository.Update(order)
// }

func (o *OrderService) InsertOrder(ctx context.Context, order *model.Order) (int64, error) {
	return o.OrderRepository.Insert(ctx, order)
}

// func (o *OrderService) GetAllOrder() ([]*model.Order, error) {
// 	return o.OrderRepository.SelectAll()
// }

// func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
// 	return o.OrderRepository.SelectAllWithInfo()
// }

// 根据消息创建订单
func (o *OrderService) InsertOrderByMessage(ctx context.Context, message *mq_content.Message) (int64, error) {
	order := &model.Order{
		UserID:      strconv.Itoa(int(message.UserID)),
		ProductID:   strconv.Itoa(int(message.ProductID)),
		OrderStatus: mq_content.OrderSuccess,
	}
	return o.InsertOrder(ctx, order)
}
