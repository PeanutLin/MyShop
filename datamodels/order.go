package datamodels

// 订单
type Order struct {
	// 订单 ID，主键自增
	ID int64
	// 用户 ID
	UserId int64 `gorm:"column:userID"`
	// 商品 ID
	ProductId int64 `gorm:"column:productID"`
	// 订单状态
	OrderStatus int64 `gorm:"column:orderStatus"`
}

const (
	OrderWait    = iota
	OrderSuccess // 1
	OrderFailed  // 2
)
