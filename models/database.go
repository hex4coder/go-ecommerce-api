package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	database, err := gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/ecommercebalanipa"))
	if err != nil {
		return nil, fmt.Errorf("database connection failed")
	}

	return database, nil
}
