package datamodels

import "sync"

type CartItem struct {
	ID         int64 `sql:"ID""`
	UserID     int64 `sql:"userID"`
	ProductID  int64 `sql:"productID"`
	ProductNum int64 `sql:"productNum"`
}

type Cart struct {
	UserID  int64
	CartMap *sync.Map
}
