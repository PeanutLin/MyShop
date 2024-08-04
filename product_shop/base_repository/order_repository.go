package base_repository

import (
	"context"
	"productshop/product_shop/middleware/mysql"
	"productshop/product_shop/middleware/mysql/gen/model"
	"productshop/product_shop/middleware/mysql/gen/query"
)

type IOrderRepository interface {
	Insert(ctx context.Context, order *model.Order) (orderID int64, err error)
	Delete(ctx context.Context, orderID int64) bool
	SelectByID(ctx context.Context, orderID int64) (*model.Order, error)
	SelectAll(ctx context.Context) ([]*model.Order, error)
	// SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderRepository struct {
	Q *query.Query
}

// 新建订单管理库
func NewOrderRepository() IOrderRepository {
	return &OrderRepository{
		Q: mysql.QueryDB,
	}
}

// 插入订单
func (o *OrderRepository) Insert(ctx context.Context, order *model.Order) (orderID int64, err error) {
	err = o.Q.Order.WithContext(ctx).Save(order)
	if err != nil {
		return -1, err
	}

	return orderID, nil
}

// 根据订单 ID 删除订单
func (o *OrderRepository) Delete(ctx context.Context, orderID int64) bool {
	panic("not implement")
}

// 根据主键查找订单
func (o *OrderRepository) SelectByID(ctx context.Context, orderID int64) (order *model.Order, err error) {
	orderPO := o.Q.Order
	order, err = o.Q.Order.Where(orderPO.ID.Eq(int32(orderID))).First()
	if err != nil {
		return nil, nil
	}
	return order, nil
}

// 查找所有订单
func (o *OrderRepository) SelectAll(ctx context.Context) (orders []*model.Order, err error) {
	orders, err = o.Q.Order.Find()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// // 查找所有订单
// func (o *OrderRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
// 	if errConn := o.Conn(); errConn != nil {
// 		return nil, errConn
// 	}

// 	rows, err := o.sqlConn.Model(&model.Order{}).Select("orders.ID, users.userName, products.productName, orders.orderStatus").Joins("left join products on products.ID = orders.productID").Joins("left join users on orders.userID = users.ID").Rows()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	result := db.GetResultRows(rows)
// 	return result, nil
// }
