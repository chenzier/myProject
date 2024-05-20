package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"product/common"
	"product/datamodels"
	"product/services"
	"strconv"
)

func restoreIpAddresses(s string) []string {
	res := []string{}
	path := []string{}
	var dfs func(start int, count int)
	dfs = func(start int, count int) {
		// fmt.Println(start,count,path)
		if count == 4 {
			if start >= len(s) {
				temp := make([]string, len(path))
				copy(temp, path)
				str := ""
				for i := 0; i < len(temp)-1; i++ {
					str += temp[i]
					str += "."
				}
				str += temp[len(temp)-1]
				res = append(res, str)
			}
			return
		}
		if start >= len(s) {
			return
		}
		i := start
		fmt.Println(i, i+2 < len(s)-1, path)
		if s[i] == '0' {
			path = append(path, string(s[i]))
			dfs(i+1, count+1) //不回溯
		} else {
			if i+2 < len(s)-1 {
				a, _ := strconv.Atoi(s[i : i+3])
				if a <= 255 {
					path = append(path, s[i:i+3])
					dfs(i+3, count+1)
					path = path[:len(path)-1]
				}
			}
		}
		if i+1 < len(s)-1 {
			path = append(path, s[i:i+2])
			dfs(i+2, count+1)
			path = append(path, s[i:i+2])

		}
		path = append(path, s[i:i+1])
		dfs(i+1, count+1)
		path = append(path, s[i:i+1])
	}

	dfs(0, 0)
	return res
}

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
}

func (p *ProductController) GetAll() mvc.View { //Http的Get请求
	productArray, _ := p.ProductService.GetAllProduct()
	//fmt.Println("访问all")
	//fmt.Println(productArray[0].ProductName)
	// 尝试创建模板视图
	view := mvc.View{
		Name: "product/view.html", //设置模版名称，即设置使用"product/view.html"这个模板
		Data: iris.Map{
			"productArray": productArray,
			//iris.Map 类型是 Iris 框架提供的一种键值对类型，用于在视图中传递数据。
			//这里将商品数组 productArray 存储在了键名为 "productArray" 的键下。
		},
	}
	//模版视图会根据所使用的模板"product/view.html"以及Data(模板设置了在某处需要某数据)自动渲染
	// 返回创建的模板视图
	return view
}

// 修改商品
func (p *ProductController) PostUpdate() { //Post请求
	fmt.Println("访问update")
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	//p.Ctx.Request().ParseForm() 是一个方法调用，用于解析 HTTP 请求的表单数据。
	//在 Web 应用中，当客户端向服务器发送 POST 请求时，通常会携带表单数据。这些数据可能来自 HTML 表单中的输入字段，例如文本框、复选框、下拉列表等。服务器需要解析这些表单数据以便进行处理。
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"}) //根据表单标签设置decoder
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		//decoder根据表单标签 将表单信息 映射到product{}结构体中中
		//p.Ctx.Request().Form 是一个字段，表示 HTTP 请求中包含的表单数据。
		//在 Go 语言的标准库 net/http 中，Request 结构体的 Form 字段存储了请求中的表单数据。

		//举个例子，如果表单中有一个字段叫做 username，可以通过 p.Ctx.Request().Form["username"] 或 p.Ctx.Request().Form.Get("username") 来获取该字段的值

		p.Ctx.Application().Logger().Debug(err) //这是一个日志记录操作，用于在应用程序中输出调试级别的日志信息。
	}
	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetAdd() mvc.View {
	fmt.Println("访问add")
	return mvc.View{
		Name: "product/add.html",
	}
}

func (p *ProductController) PostAdd() {
	product := &datamodels.Product{}
	p.Ctx.Request().ParseForm()
	//p.Ctx.Request() 方法，你可以在控制器中获取当前 HTTP 请求的所有信息，并对请求进行处理和响应。例如，你可以获取请求的参数、读取请求体、设置响应头等。
	//当前请求的 HTTP 请求对象" 指的是当前客户端发送到服务器的 HTTP 请求。
	//	在 Web 应用程序中，客户端（通常是浏览器）向服务器发送请求，服务器接收并处理请求，然后发送响应给客户端。
	//HTTP 请求对象（http.Request）是 Go 语言标准库中 net/http 包中的一个类型。
	//它代表了一个 HTTP 请求，包含了客户端发送的所有请求信息，如 URL、方法、请求头、请求体等。
	//每当客户端发送一个 HTTP 请求时，服务器都会创建一个对应的 http.Request 对象来表示该请求。
	dec := common.NewDecoder(&common.DecoderOptions{TagName: "imooc"})
	if err := dec.Decode(p.Ctx.Request().Form, product); err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	_, err := p.ProductService.InsertProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

func (p *ProductController) GetManager() mvc.View {
	idString := p.Ctx.URLParam("id") //通过上下文获取id
	//通过id查询商品 十进制 转为 十六进制
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	product, err := p.ProductService.GetProductByID(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetDelete() {
	idString := p.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	isOk := p.ProductService.DeleteProductByID(id)
	if isOk {
		p.Ctx.Application().Logger().Debug("删除商品成功，ID为：" + idString)
	} else {
		p.Ctx.Application().Logger().Debug("删除商品失败，ID为：" + idString)
	}
	p.Ctx.Redirect("/product/all")
}
