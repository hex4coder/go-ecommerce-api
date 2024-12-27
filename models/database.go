package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	database, err := gorm.Open(mysql.Open("ecom:ecom123;@tcp(127.0.0.1:3306)/ecommercebalanipa?parseTime=true"))
	if err != nil {
		return nil, fmt.Errorf("database connection failed")
	}

	return database, nil
}
