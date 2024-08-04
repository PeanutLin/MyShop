package mysql

import (
	"fmt"
	"productshop/product_shop/common"
	"productshop/product_shop/middleware/logs"
	"productshop/product_shop/middleware/mysql/gen/query"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db      *gorm.DB
	QueryDB *query.Query
)

func Init() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", common.MySQLUserName, common.MySQLPassWord, common.MySQLHost, common.MySQLPort, common.MySQLDbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}

	InitQuery()
	logs.Info("MySQL Init Success")
}

func InitQuery() {
	QueryDB = query.Use(db)
}

func GetQuery() *query.Query {
	return query.Use(db)
}
