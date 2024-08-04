package common

import (
	"productshop/product_shop/middleware/logs"

	"github.com/kataras/iris/v12"
)

func AuthProduct(ctx iris.Context) {
	uid := ctx.GetCookie("user_id")
	if uid == "" {
		logs.Info("must log in first")
		ctx.Redirect("/user/login")
		return
	}
	logs.Info("already logged in")
	ctx.Next()
}
