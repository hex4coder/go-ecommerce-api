package test

import (
	"fmt"
	"testing"

	"github.com/hex4coder/go-ecommerce-api/controllers"
	"github.com/hex4coder/go-ecommerce-api/models"
	"github.com/joho/godotenv"
)

func TestMyOrder(t *testing.T) {
	fmt.Println("testing my order")

	// load env
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("error loading .env file : %s", err)
	}

	// connect to database
	// open database
	db, err := models.ConnectDB()
	if err != nil {
		t.Fatalf("%v", err)
	}

	// get user id
	userId := 2

	// buat order api
	orderapi := controllers.NewOrderAPI(db)

	// get data
	orders, err := orderapi.GetMyOrders(userId)

	// check error
	if err != nil {
		t.Fatalf("error in get myorders : %s", err)
	}

	// print orders
	fmt.Printf("%v", orders)
}
