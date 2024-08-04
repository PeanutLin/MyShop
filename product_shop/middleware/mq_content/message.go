package mq_content

type Message struct {
	UserID     int64
	ProductID  int64
	ProductNum int64
}

const (
	OrderWait    = iota
	OrderSuccess // 1
	OrderFailed  // 2
)

func NewMessage(userID int64, productID int64, productNum int64) *Message {
	return &Message{
		UserID:     userID,
		ProductID:  productID,
		ProductNum: productNum,
	}
}
