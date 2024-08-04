package main

import (
	"context"
	"fmt"
	"productshop/backend/web/controllers"
	"productshop/common"
	"productshop/db"
	"productshop/repositories"
	"productshop/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	// 1.创建 iris 实例
	app := iris.New()

	// 2.设置错误模式
	app.Logger().SetLevel("debug")

	// 设置重定向 / 至 /order/index
	app.Get("/", func(ctx iris.Context) {
		ctx.Redirect("/order/index")
	})

	// 3.注册模板
	tempalte := iris.HTML("./backend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(tempalte)

	// 4.设置资源目录
	app.HandleDir("/assets", "./assets")

	// 出现异常跳转指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewLayout("")
		ctx.ViewData("Message", ctx.Values().GetStringDefault("Message", "visit error"))
		ctx.View("shared/error.html")
	})

	// 连接数据库
	db, err := db.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 5.注册控制器
	productRepository := repositories.NewProductManager(db)
	productService := services.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	orderRepository := repositories.NewOrderManagerRepository(db)
	orderService := services.NewOrderService(orderRepository)
	orderParty := app.Party("/order")
	order := mvc.New(orderParty)
	order.Register(ctx, orderService)
	order.Handle(new(controllers.OrderController))

	// 6.启动服务
	app.Run(
		iris.Addr(common.BackendHost+":"+common.BackendPort),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
