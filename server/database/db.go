package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var database *gorm.DB

//打开数据库
func Open(dialect string, args ...interface{}) *gorm.DB {
	db, err := gorm.Open(dialect, args)
	if err != nil {
		panic(err)
	}
	database = db
	return database
}

func OpenDefault() *gorm.DB {
	database, err := gorm.Open("mysql", "root:2209931449Qq@(rm-bp136wu99t27ypvykxo.mysql.rds.aliyuncs.com:3306)/fantastic_chat?"+
		"charset=utf8mb4")
	if err != nil {
		panic(err)
	}
	return database
}
