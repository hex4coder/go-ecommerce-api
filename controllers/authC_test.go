package controllers

import (
	"fmt"
	"testing"

	"github.com/hex4coder/go-ecommerce-api/models"
)

func TestEmailCheck(t *testing.T) {
	// buat temporary
	var emailCount int64

	db, err := models.ConnectDB()

	if err != nil {
		t.Fatalf("ada kesalahan pada koneksi database %v", err)
	}

	a := NewAuthAPI(db)

	email := "customer1@gmail.com"

	// validasi email
	a.DB.Table("users").Where(&models.User{Email: email}).Count(&emailCount)
	if emailCount > 0 {
		t.Fatalf("email %s sudah terdaftar", email)
	} else {
		fmt.Println("email belum ada")
	}

}

func TestPasswordHash(t *testing.T) {
	hashedPassword := "$2y$12$lksGE6.lHRcG5BW3MNBu/eO7jRJxu8QcU6wRdGUkewvnEI9K28fw2"
	password := "12345678"

	api := NewAuthAPI(nil)

	if ok := api.CheckPassword(password, hashedPassword); !ok {
		t.Fatal("password tidak valid")
	}
}

func TestHashingPassword(t *testing.T) {
	password := "12345678"
	api := NewAuthAPI(nil)

	hashed, err := api.HashPassword(password)
	if err != nil {
		t.Fatalf("%v", err)
	}

	fmt.Printf("Hasil : %s", hashed)
}
