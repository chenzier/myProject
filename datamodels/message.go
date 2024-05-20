package datamodels

// 简单消息体
type Message struct {
	ProductID int64
	UserID    int64
}

// 创建结构体
func NewMessage(userId int64, ProductId int64) *Message {
	return &Message{ProductID: ProductId, UserID: userId}
}
