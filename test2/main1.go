package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// 注册模板引擎，指定模板文件路径为 "./views" 目录，扩展名为 ".html"
	tmpl := iris.HTML("./test2/views", ".html").Layout("layout.html")
	app.RegisterView(tmpl)

	// 设置路由，当访问根路径时，渲染 "index.html" 模板
	app.Get("/", func(ctx iris.Context) {
		// 渲染模板时传递了一个包含数据的 Map 对象
		//ctx.View("login.html", iris.Map{"message": "Hello, Iris!"})
		//ctx.ViewData("message", "Hello world!")
		////ctx.View("index.html", iris.Map{"title": "Page Title", "message": "Hello, Iris!"})
		//ctx.View("index.html")
		ctx.View("index.html", iris.Map{"title": "My Website - Home"})

	})

	// 运行应用，监听 8080 端口
	app.Run(iris.Addr(":8080"))
}
