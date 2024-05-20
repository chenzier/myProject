package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"product/backend/web/controllers"
	"product/common"
	"product/repositories"
	"product/services"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")

	//3.注册模板
	tmplate := iris.HTML("./backend/web/views", ".html").Layout(
		"shared/layout.html").Reload(true)
	app.RegisterView(tmplate)

	//4.设置模板目标
	app.HandleDir("/assets", iris.Dir("./backend/web/assets"))
	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("view1")
		ctx.View("shared/error.html")
	})

	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		//log.Error(err)
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //当main函数结束时，cancel()函数执行，然后上下文会通知每个函数(或协程)终止执行，防止函数(或函数)继续占有内存等资源
	//和协程并发 有关

	//5.注册控制器
	productRepository := repositories.NewProductManager("product", db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))
	//6.启动服务
	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
