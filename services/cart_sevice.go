package services

import (
	"product/repositories"
)

type ICartService interface {
}

func NewCartService(repository repositories.ICartRepository) ICartService {
	return &CartService{CartRepository: repository}
}

type CartService struct {
	CartRepository repositories.ICartRepository
}

func LoadCartFromMysql() {
	//新建redis列表
}
