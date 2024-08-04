package controllers

import (
	"context"
	"productshop/product_shop/base_service"
	"productshop/product_shop/common"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/mysql/gen/model"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

type UserController struct {
	Ctx         iris.Context
	UserService base_service.IUserService
	Session     *sessions.Session
}

func (c *UserController) Get() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/welcome.html",
	}
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/register.html",
	}
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Layout: "",
		Name:   "user/login.html",
	}
}

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
	user := &model.User{
		NickName: nickName,
		UserName: userName,
		Password: pwd,
	}
	_, err := c.UserService.AddUser(context.Background(), user)
	if err != nil {
		logs.Error(err.Error())
	}

	if err != nil {
		c.Ctx.Redirect("/user/error")
		return
	}
	c.Ctx.Redirect("/user/login")
}

func (c *UserController) PostLogin() mvc.Response {
	ctx := context.Background()
	// 1.获取用户提交的表单信息
	var (
		userName = c.Ctx.FormValue("userName")
		password = c.Ctx.FormValue("password")
	)
	// 2.验证用户账号密码是否正确
	user, isOK, err := c.UserService.IsLoginSuccess(ctx, userName, password)
	if err != nil {
		logs.Error("login error", logs.String("error msg", err.Error()))
	}

	// Login Failed
	if !isOK {
		return mvc.Response{
			Path: "loginerr",
		}
	}

	// 写入用户 ID 到 Cookie 中
	common.GlobalCookie(c.Ctx, "user_id", strconv.FormatInt(int64(user.ID), 10))
	uidByte := strconv.FormatInt(int64(user.ID), 10)
	uidString, err := common.EnPwdCode([]byte(uidByte))
	if err != nil {
		logs.Error("[common.EnPwdCode] error")
	}
	common.GlobalCookie(c.Ctx, "sign", uidString)

	return mvc.Response{
		Path: "/product",
	}
}
