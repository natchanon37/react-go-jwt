package database

import (
	"react-go-jwt/server/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(){
	connect, err := gorm.Open(mysql.Open("root:non256199@/jwt-go-react"),&gorm.Config{})
	if err != nil {
		panic("could not connect to database")
	}

	DB = connect

	connect.AutoMigrate(&models.User{})

}