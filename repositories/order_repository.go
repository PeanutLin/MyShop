package repositories

import (
	"log"
	"productshop/datamodels"
	"productshop/db"

	"gorm.io/gorm"
)

type IOrderRepository interface {
	Conn() error
	Insert(order *datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(order *datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManagerRepository struct {
	sqlConn *gorm.DB
}

// 新建订单管理库
func NewOrderManagerRepository(sql *gorm.DB) IOrderRepository {
	return &OrderManagerRepository{
		sqlConn: sql,
	}
}

// 检查 mysql 连接
func (o *OrderManagerRepository) Conn() error {
	if o.sqlConn == nil {
		sqlConn, err := db.NewMysqlConn()
		if err != nil {
			return err
		}
		o.sqlConn = sqlConn
	}
	return nil
}

// 插入订单
func (o *OrderManagerRepository) Insert(order *datamodels.Order) (orderId int64, err error) {
	if err = o.Conn(); err != nil {
		return -1, err
	}

	err = o.sqlConn.Create(order).Error
	if err != nil {
		log.Println("insert fail : ", err)
	}

	return orderId, nil
}

// 根据订单 ID 删除订单
func (o *OrderManagerRepository) Delete(orderId int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}

	err := o.sqlConn.Where("id=?", orderId).Delete(&datamodels.Order{})
	if err != nil {
		log.Println("delete order fail : ", err)
		return false
	}

	return true
}

// 更新订单信息
func (o *OrderManagerRepository) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}

	err := o.sqlConn.Save(order).Error
	if err != nil {
		log.Println("update order fail", err)
		return err
	}

	return nil
}

// 根据主键查找订单
func (o *OrderManagerRepository) SelectByKey(id int64) (*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	order := &datamodels.Order{}
	err := o.sqlConn.First(&order, id).Error
	if err != nil {
		log.Println("Select By Key fail", err)
		return nil, err
	}
	return order, nil
}

// 查找所有订单
func (o *OrderManagerRepository) SelectAll() (order []*datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}

	err = o.sqlConn.Find(order).Error
	if err != nil {
		log.Println("Select All fail", err)
		return nil, err
	}
	return order, nil
}

// 查找所有订单
func (o *OrderManagerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}

	rows, err := o.sqlConn.Model(&datamodels.Order{}).Select("orders.ID, users.userName, products.productName, orders.orderStatus").Joins("left join products on products.ID = orders.productID").Joins("left join users on orders.userID = users.ID").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := db.GetResultRows(rows)
	return result, nil
}
