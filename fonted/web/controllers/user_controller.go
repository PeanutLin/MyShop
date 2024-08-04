package controllers

import (
	"fmt"
	"productshop/common"
	"productshop/datamodels"
	"productshop/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.IUserService
	Session     *sessions.Session
}

func (c *UserController) Get() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/welcome.html",
	}
}

// http://111.230.70.68:8080/user/register
func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/register.html",
	}
}

// http://111.230.70.68:8080/user/login
func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/login.html",
	}
}

// http://111.230.70.68:8080/user/loginerr
func (c *UserController) GetLoginerr() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/login.html",
		Data: iris.Map{
			"showMessage": "提示：用户名或密码错误，请重试",
		},
	}
}

func (c *UserController) PostRegister() {
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		pwd      = c.Ctx.FormValue("password")
	)
	user := &datamodels.User{
		NickName:     nickName,
		UserName:     userName,
		HashPassword: pwd,
	}
	_, err := c.UserService.AddUser(user)
	fmt.Println(err)
	if err != nil {
		c.Ctx.Redirect("/user/error")
		return
	}
	c.Ctx.Redirect("/user/login")
}

func (c *UserController) PostLogin() mvc.Response {
	// 1.获取用户提交的表单信息
	var (
		userName = c.Ctx.FormValue("userName")
		pwd      = c.Ctx.FormValue("password")
	)

	// 2.验证用户账号密码是否正确
	user, isOk := c.UserService.IsLoginSuccess(userName, pwd)

	// Login Failed
	if !isOk {
		return mvc.Response{
			Path: "loginerr",
		}
	}

	// 写入用户 ID 到 Cookie 中
	common.GlobalCookie(c.Ctx, "uid", strconv.FormatInt(user.ID, 10))
	uidByte := strconv.FormatInt(user.ID, 10)
	uidString, err := common.EnPwdCode([]byte(uidByte))
	if err != nil {
		fmt.Println(err)
	}
	common.GlobalCookie(c.Ctx, "sign", uidString)

	return mvc.Response{
		Path: "/product",
	}
}
