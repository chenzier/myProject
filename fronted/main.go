package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"product/common"
	"product/fronted/middlerware"
	"product/fronted/web/controllers"
	"product/rabbitmq"
	"product/repositories"
	"product/services"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	tmplate := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tmplate)
	//4.设置模板目标
	app.HandleDir("/public", iris.Dir("./fronted/web/public"))
	app.HandleDir("/html", "./fronted/web/htmlProductShow")

	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("view2.html")
		ctx.View("shared/error.html")
	})

	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {
		//log.Error(err)
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //注意这里在最后将所有协程关闭

	//sess := sessions.New(sessions.Config{
	//	Cookie:  "helloworld",
	//	Expires: 60 * time.Minute,
	//})

	//5.注册控制器
	user := repositories.NewUserRepository("user", db)
	userService := services.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	//userPro.Register(userService, ctx, sess.Start)
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	rabbitmq := rabbitmq.NewRabbitMQSimple("imoocProduct")

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderMangerRepository("order", db)
	orderService := services.NewOrderService(order)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middlerware.AuthConProduct)
	//pro.Register(productService, orderService, sess.Start)
	pro.Register(productService, orderService, rabbitmq)
	pro.Handle(new(controllers.ProductController)) //调用 userPrro.Handle(new(controllers.UserController)) 后，框架会将 UserController 中的各个方法与对应的 HTTP 请求方式（如 GET、POST 等）进行绑定，以便在收到请求时调用相应的方法来处理请求。例如，如果 UserController 中有一个 Get 方法，则该方法会处理 /user 路由的 GET 请求；如果有一个 Post 方法，则处理 POST 请求，以此类推。

	app.Get("/", func(ctx iris.Context) {
		// 渲染模板时传递了一个包含数据的 Map 对象
		ctx.View("index.html", iris.Map{"title": "Hello, Iris!"})
	})
	//6.启动服务
	app.Run(
		iris.Addr("localhost:8083"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
