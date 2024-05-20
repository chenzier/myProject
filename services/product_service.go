package services

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"product/datamodels"
	"product/repositories"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProduct() ([]*datamodels.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
	SubNumberOne(productID int64) error
	GetProductByID1(ctx iris.Context, productID int64) (*datamodels.Product, error)
}

type ProductService struct {
	productRepository repositories.IProduct
}

// 初始化函数
func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{repository}
}

// 根据id查商品
func (p *ProductService) GetProductByID(productID int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(productID)
}
func (p *ProductService) GetProductByID1(ctx iris.Context, productID int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey1(ctx, productID)
}

// 获取所有商品信息
func (p *ProductService) GetAllProduct() ([]*datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

// 删
func (p *ProductService) DeleteProductByID(productID int64) bool {
	return p.productRepository.Delete(productID)
}

// 增
func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	id, error := p.productRepository.Insert(product)
	fmt.Println("添加商品为", id, error)
	return id, error
}

// 改
func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productRepository.Update(product)
}

func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productRepository.SubProductNum(productID)
}
