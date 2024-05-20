package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"product/datamodels"
	"product/encrypt"
	"product/services"
	"product/tool"
	"strconv"
)

type UserController struct {
	Ctx     iris.Context
	Service services.IUserService
	Session *sessions.Session
}

func (u *UserController) GetRegister() mvc.View {
	fmt.Println("访问register")
	return mvc.View{
		Name: "user/register.html",
	}
}

func (c *UserController) PostRegister() {
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)
	user := &datamodels.User{
		UserName:     userName,
		NickName:     nickName,
		HashPassword: password,
	}
	_, err := c.Service.AddUser(user)
	if err != nil {
		c.Ctx.Redirect("/user/error")
	}
	c.Ctx.Redirect("/user/login")
	fmt.Println("新用户注册成功")
	return
}

func (c *UserController) GetLogin() mvc.View {
	fmt.Println("进入登陆页面")
	return mvc.View{
		Name: "user/login.html",
	}
}

// 3.写入用户ID到cookie中
func (c *UserController) PostLogin() mvc.Response {
	//1.获取用户提交的表单信息
	var (
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)
	//2.验证账户密码正确
	user, isOk := c.Service.IsPwdSuccess(userName, password)
	if !isOk {
		return mvc.Response{
			Path: "/user/login",
		}
	}
	userID := strconv.FormatInt(user.ID, 10)
	//	c.Session.Set("userID", strconv.FormatInt(user.ID, 10))
	encryptedUserID, err := encrypt.EnPwdCode([]byte(userID))
	if err != nil {
		fmt.Println("加密用户ID失败:", err)
		return mvc.Response{
			Path: "/user/login",
		}
	}
	//一旦设置了 Cookie，每次用户发送请求时，浏览器都会将相应的 Cookie 信息发送给服务器。
	//在 Iris 框架中，可以使用 ctx.GetCookie(name string) 方法来获取特定名称的 Cookie 的值。
	//例如，可以通过 uid := ctx.GetCookie("uid") 来获取名为 "uid" 的 Cookie 的值。

	//c.Session.Set("userID", strconv.FormatInt(user.ID, 10))
	tool.GlobalCookie(c.Ctx, "uid", userID)
	tool.GlobalCookie(c.Ctx, "sign", encryptedUserID)
	//fmt.Println("用户加密成功")
	fmt.Println("用户登陆成功")

	return mvc.Response{
		Path: "/product/",
	}
}

//func (c *UserController) PostLogin() mvc.Response {
//
//	var (
//		userName = c.Ctx.FormValue("userName")
//		password = c.Ctx.FormValue("password")
//	)
//	user, isOk := c.Service.IsPwdSuccess(userName, password)
//	fmt.Println("用户登陆成功")
//	if !isOk {
//		return mvc.Response{
//			Path: "/user/login",
//		}
//	}
//	//一旦设置了 Cookie，每次用户发送请求时，浏览器都会将相应的 Cookie 信息发送给服务器。
//	//在 Iris 框架中，可以使用 ctx.GetCookie(name string) 方法来获取特定名称的 Cookie 的值。
//	//例如，可以通过 uid := ctx.GetCookie("uid") 来获取名为 "uid" 的 Cookie 的值。
//	tool.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))
//	c.Session.Set("userID", strconv.FormatInt(user.ID, 10))
//	return mvc.Response{
//		Path: "/product/",
//	}
//}
