package datamodels

type Order struct {
	ID          int64 `sql:"ID""`
	UserID      int64 `sql:"userID"`
	ProductID   int64 `sql:"productID"`
	OrderStatus int   `sql:"orderStatus"`
}

const (
	OrderWait = iota
	OrderSuccess
	OrderFailed
)
