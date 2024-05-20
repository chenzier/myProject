package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"product/services"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray, err := o.OrderService.GetAllOrderInfo()
	//fmt.Println(orderArray)
	if err != nil {
		o.Ctx.Application().Logger().Debug("订单信息失败")
	}
	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}
