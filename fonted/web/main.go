package main

import (
	"context"
	"fmt"
	"productshop/common"
	"productshop/db"
	"productshop/fonted/middleware"
	"productshop/fonted/web/controllers"
	"productshop/repositories"
	"productshop/services"

	_ "github.com/go-sql-driver/mysql"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	// 1.创建 iris 实例 -------------------------------------
	app := iris.New()

	// 设置重定向 / 至 /product
	app.Get("/", func(ctx iris.Context) {
		ctx.Redirect("/user/")
	})

	// 2.设置错误模式 -------------------------------------
	app.Logger().SetLevel("debug")

	// 3.注册模板 -------------------------------------
	tempalte := iris.HTML("./fonted/web/views", ".html").Reload(true)
	app.RegisterView(tempalte)

	// 静态页面
	// app.HandleDir("/html", "./fonted/web/htmlProductShow")

	// 资源目录 -------------------------------------
	app.HandleDir("/assets", "./assets")

	// 出现异常跳转指定页面 -------------------------------------
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewLayout("")
		ctx.ViewData("Message", ctx.Values().GetStringDefault("Message", "visit error"))
		ctx.View("shared/error.html")
	})

	// 连接数据库 -------------------------------------
	db, err := db.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 5.注册控制器
	// 注册 user 控制器 -------------------------------------
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))

	// 注册 product 控制器 -------------------------------------
	productRepository := repositories.NewProductManager(db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")

	// 验证中间件
	productParty.Use(middleware.AuthProduct)

	// order 服务
	orderRepository := repositories.NewOrderManagerRepository(db)
	orderService := services.NewOrderService(orderRepository)

	product := mvc.New(productParty)
	product.Register(ctx, productService, orderService)
	product.Handle(new(controllers.ProductController))

	// 6.启动服务 -------------------------------------
	app.Run(
		iris.Addr(common.FontedHost+":"+common.FontedPort),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
