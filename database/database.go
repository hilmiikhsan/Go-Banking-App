package database

import (
	"banking-app/helpers"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	database, err := gorm.Open("mysql", "root:@/bankingapp?charset=utf8&parseTime=True")
	helpers.HandleErr(err)

	database.DB().SetMaxIdleConns(20)
	database.DB().SetMaxOpenConns(200)

	DB = database
}