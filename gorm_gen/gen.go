package main

import (
	"fmt"
	"productshop/product_shop/common"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../product_shop/middleware/mysql/gen/query",
		ModelPkgPath: "../product_shop/middleware/mysql/gen/model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", common.MySQLUserName, common.MySQLPassWord, common.MySQLHost, common.MySQLPort, common.MySQLDbName)
	gormdb, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	g.UseDB(gormdb) // reuse your gorm db

	g.ApplyBasic(g.GenerateModel("users"))
	g.ApplyBasic(g.GenerateModel("orders"))
	g.ApplyBasic(g.GenerateModel("products"))

	// Generate the code
	g.Execute()
}
