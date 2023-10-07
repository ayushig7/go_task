package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB
var urlDSN = "root:ayushigupta7@@tcp(localhost:3306)/company"
var err error

func DataMigration() {
	Database, err = gorm.Open(mysql.Open(urlDSN), &gorm.Config{})
	if err != nil {
		fmt.Print(err.Error())
		panic("connection failed")
	}
	Database.AutoMigrate(&Employee{})
}
