package main

import (
	"context"
	"productshop/product_shop/base_repository"
	"productshop/product_shop/base_service"
	"productshop/product_shop/common"
	"productshop/product_shop/middleware/jeager"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/mysql"
	"productshop/product_shop/shop_fontend/controllers"
	"productshop/product_shop/shop_fontend/rpc"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	// 初始化 log
	defer logs.Init().Sync()

	// 连接数据库 -------------------------------------
	mysql.Init()

	// 初始化 Jeager
	_, closer := jeager.Init("shop.fontend")
	defer closer.Close()

	// rpc 客户端
	rpc.MustInit()

	// 1.创建 iris 实例 -------------------------------------
	app := iris.New()

	// 设置重定向 / 至 /product
	app.Get("/", func(ctx iris.Context) {
		ctx.Redirect("/user/")
	})

	// 2.设置错误模式 -------------------------------------
	app.Logger().SetLevel("debug")

	// 3.注册模板 -------------------------------------
	tempalte := iris.HTML("./views", ".html").Reload(true)
	app.RegisterView(tempalte)

	// 静态页面
	// app.HandleDir("/html", "./fonted/web/htmlProductShow")

	// 资源目录 -------------------------------------
	app.HandleDir("/assets", "../assets")

	// 出现异常跳转指定页面 -------------------------------------
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewLayout("")
		ctx.ViewData("Message", ctx.Values().GetStringDefault("Message", "visit error"))
		ctx.View("shared/error.html")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 5.注册控制器
	// 注册 user 控制器 -------------------------------------
	userRepository := base_repository.NewUserRepository()
	userService := base_service.NewUserService(userRepository)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))

	// 注册 product 控制器 -------------------------------------
	productRepository := base_repository.NewProductRepository()
	productService := base_service.NewProductService(productRepository)
	productParty := app.Party("/product")

	// 添加验证中间件
	productParty.Use(common.AuthProduct)

	// order 服务
	orderRepository := base_repository.NewOrderRepository()
	orderService := base_service.NewOrderService(orderRepository)

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
