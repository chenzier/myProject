package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"os"
	"path/filepath"
	"product/datamodels"
	"product/rabbitmq"
	"product/services"
	"strconv"
	"text/template"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Rabbitmq       *rabbitmq.RabbitMQ
	Session        *sessions.Session
}

var (
	htmlOutPath  = "./fronted/web/htmlProductShow" //生成的Html保存目录
	templatePath = "./fronted/web/views/template/"
)

func (p *ProductController) GetGenerateHtml() {
	productString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	//1.获取模版
	contenstTmp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//2.获取html生成路径
	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")

	//3.获取模版渲染数据
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	//4.生成静态文件
	generateStaticHtml(p.Ctx, contenstTmp, fileName, product)
}

// 生成html静态文件
func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
	fmt.Println("1223")
	//1.判断静态文件是否存在
	if exist(fileName) {
		err := os.Remove(fileName) // 如果文件已经存在，则会调用 os.Remove 函数将其删除。这一步是为了确保生成的静态文件是最新的，而不是覆盖原有的内容。
		if err != nil {
			ctx.Application().Logger().Error(err)
		}
		fmt.Println("存在并删除成功")
	}
	//2.生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	//接下来，函数会调用 os.OpenFile 函数打开指定的文件
	//如果文件不存在，则会创建它
	//然后，使用 os.O_CREATE|os.O_WRONLY 标志以只写模式打开文件
	//并设置文件权限为 os.ModePerm。这意味着函数将以写入模式打开文件

	if err != nil {
		ctx.Application().Logger().Error(err)
	}
	fmt.Println("生成文件", fileName)
	defer file.Close()
	template.Execute(file, &product)
	//函数会调用 template.Execute 方法，将模板中的内容填充到打开的文件中
	//这里的模板是一个 template.Template 类型的变量，通常是一个预定义的HTML模板文件，其中包含了一些占位符，需要根据具体的数据进行替换。
	//在这个例子中，模板会根据传入的 product 参数，替换模板中的占位符，并将结果写入到打开的文件中。
}

// 判断文件是否存在
func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

//func (p *ProductController) GetGenerateHtml() {
//	productString := p.Ctx.URLParam("productID")
//	productID, err := strconv.Atoi(productString)
//	if err != nil {
//		p.Ctx.Application().Logger().Debug(err)
//	}
//	//1.获取模板文件地址
//	contentTmp, err := template.ParseFiles(filepath.Join(templatePath), "product.html")
//	if err != nil {
//		p.Ctx.Application().Logger().Debug(err)
//	}
//	//2.获取html生成路径
//	fileName := filepath.Join(htmlOutPath, "htmlProduct.html")
//	//3.获取模板渲染数据
//	product, err := p.ProductService.GetProductByID(int64(productID))
//	if err != nil {
//		p.Ctx.Application().Logger().Debug(err)
//	}
//	//4.生成静态文件
//	generateStaticHtml(p.Ctx, contentTmp, fileName, product)
//}
//
//// 生成html静态文件
//func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
//	//判断静态文件是否存在
//	if exist(fileName) {
//		err := os.Remove(fileName)
//		if err != nil {
//			ctx.Application().Logger().Error(err)
//		}
//	}
//	//2.生成静态文件
//	file, err := os.OpenFile(fileName, os.O_CREATE, os.ModePerm)
//	if err != nil {
//		ctx.Application().Logger().Error(err)
//	}
//	defer file.Close()
//	template.Execute(file, &product)
//}
//
//// 判断文件是否存在
//func exist(fileName string) bool {
//	_, err := os.Stat(fileName)
//	return err == nil || os.IsExist(err)
//}

// 控制器定义了两个方法 GetDetail和GetOrder，访问对应的网页会调用get方法

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductByID1(p.Ctx, 1)
	if err != nil {
		p.Ctx.Application().Logger().Error(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

//func (p *ProductController) GetOrder() mvc.View {
//	productString := p.Ctx.URLParam("productID")
//	userString := p.Ctx.GetCookie("uid")
//	productID, err := strconv.Atoi(productString)
//	if err != nil {
//		p.Ctx.Application().Logger().Debug(err)
//	}
//	product, err := p.ProductService.GetProductByID(int64(productID))
//	if err != nil {
//		p.Ctx.Application().Logger().Debug(err)
//	}
//	var orderID int64
//	showMessage := "抢购失败！"
//	//判断商品数量是否满足需求
//	if product.ProductNum > 0 {
//		//扣除商品数量
//		product.ProductNum -= 1
//		err := p.ProductService.UpdateProduct(product)
//		if err != nil {
//			p.Ctx.Application().Logger().Debug(err)
//		}
//		//创建订单
//		userID, err := strconv.Atoi(userString)
//		//fmt.Println(userID, product.ProductName)
//		if err != nil {
//			p.Ctx.Application().Logger().Debug(err)
//		}
//
//		order := &datamodels.Order{
//			ID:          100,
//			UserID:      int64(userID),
//			ProductID:   int64(productID),
//			OrderStatus: datamodels.OrderWait,
//		}
//		//新建订单
//		orderID, err = p.OrderService.InsertOrder(order)
//		fmt.Println("orderID", orderID)
//		if err != nil {
//			p.Ctx.Application().Logger().Debug(err)
//		} else {
//			showMessage = "抢购成功！"
//		}
//	}
//
//	return mvc.View{
//		Layout: "shared/productLayout.html",
//		Name:   "product/result.html",
//		Data: iris.Map{
//			"orderID":     orderID,
//			"showMessage": showMessage,
//		},
//	}
//
//}

func (p *ProductController) GetOrder() []byte {
	//http://localhost:8083/product/ order?productID=1
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.ParseInt(productString, 10, 64)
	//fmt.Println(productString, userString, productID, err)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)

	}
	userID, err := strconv.ParseInt(userString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)

	}

	//创建消息体
	message := datamodels.NewMessage(userID, productID)
	//类型转化
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	err = p.Rabbitmq.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return []byte("true")

}
