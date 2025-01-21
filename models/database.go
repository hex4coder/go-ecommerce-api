package models

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	DatabaseName string
	Username     string
	Password     string
	ParseTime    bool
	Host         string
	Port         int
}

func LoadDatabaseConfigFromEnv(parseTime bool) *DatabaseConfig {

	// read env file for database config

	dbname := os.Getenv("DATABASE_NAME")
	if len(dbname) == 0 {
		dbname = "ecommercebalanipa"
	}

	dbhost := os.Getenv("DATABASE_HOST")
	if len(dbhost) == 0 {
		dbhost = "172.28.208.1"
	}

	dbportstr := os.Getenv("DATABASE_PORT")
	if len(dbportstr) == 0 {
		dbportstr = "3306"
	}
	dbport, err := strconv.Atoi(dbportstr)
	if err != nil {
		fmt.Println("Failed to convert database port from env file ", err)
		panic(err)
	}

	dbusername := os.Getenv("DATABASE_USERNAME")
	if len(dbusername) == 0 {
		dbusername = "aaa"
	}

	dbpassword := os.Getenv("DATABASE_PASSWORD")
	if len(dbpassword) == 0 {
		dbpassword = "anu123"
	}

	return &DatabaseConfig{
		ParseTime:    parseTime,
		DatabaseName: dbname,
		Host:         dbhost,
		Port:         dbport,
		Username:     dbusername,
		Password:     dbpassword,
	}
}

func (d *DatabaseConfig) GetConnectionString() string {
	connString := ""
	connString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.Username, d.Password, d.Host, d.Port, d.DatabaseName)

	if d.ParseTime {
		connString = fmt.Sprintf("%s?parseTime=true", connString)
	}
	return connString
}

func ConnectDB() (*gorm.DB, error) {
	// database, err := gorm.Open(mysql.Open("root@tcp(127.0.0.1:3306)/ecommercebalanipa?parseTime=true"))
	dbConfig := LoadDatabaseConfigFromEnv(true)
	// panic(dbConfig.GetConnectionString())
	database, err := gorm.Open(mysql.Open(dbConfig.GetConnectionString()))
	if err != nil {
		return nil, fmt.Errorf("database connection failed %s", dbConfig.GetConnectionString())
	}

	return database, nil
}
