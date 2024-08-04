package db

import (
	"database/sql"
	"fmt"
	"productshop/common"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMysqlConn() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", common.MySQLUserName, common.MySQLPassWord, common.MySQLHost, common.MySQLPort, common.MySQLDbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 获取返回值，获取一条
func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				//fmt.Println(reflect.TypeOf(col))
				record[columns[i]] = string(v)
			}
		}
	}
	return record
}

// 获取所有
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	//返回所有列
	columns, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(columns))
	//这里表示一行填充数据
	scans := make([]interface{}, len(columns))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := columns[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result
}
