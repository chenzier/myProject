package main

import (
	"fmt"
	"product/common"
	"product/rabbitmq"
	"product/repositories"
	"product/services"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	//创建product数据库操作实例
	product := repositories.NewProductManager("product", db)
	//创建product service
	productService := services.NewProductService(product)
	//创建Order数据库实例
	order := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("imoocProduct")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService, 3)
}
