package test

import (
	"fmt"
	"testing"

	"github.com/hex4coder/go-ecommerce-api/controllers"
	"github.com/hex4coder/go-ecommerce-api/models"
)

func TestEmailCheck(t *testing.T) {

	var r = new(controllers.RegisterRequest)
	r.Email = "customer1@gmail.com"
	r.Name = "Customer1"
	r.Password = "12345678"

	r.Jalan = "Jalan Poros Mambu"
	r.Kodepos = "91353"
	r.Nomorhp = "083238999909"
	r.Dusun = "Malise"
	r.Desa = "Baru"
	r.Kecamatan = "Luyo"
	r.Kota = "Polewali Mandar"
	r.Provinsi = "Sulawesi Barat"

	// buat temporary
	var emailCount int64

	db, err := models.ConnectDB()

	if err != nil {
		t.Fatalf("ada kesalahan pada koneksi database %v", err)
	}

	a := controllers.NewAuthAPI(db)

	// validasi email
	a.DB.Table("users").Where(&models.User{Email: r.Email}).Count(&emailCount)
	if emailCount > 0 {
		t.Fatalf("email %s sudah terdaftar", r.Email)
	} else {
		fmt.Println("email belum ada")
	}

	// buat model baru
	user := new(models.User)

	// create password hash
	hashed, err := a.HashPassword(r.Password)
	if err != nil {
		t.Errorf("failed to hashing the password %s", err.Error())
	}

	// assign data to user
	user.Email = r.Email
	user.Name = r.Name
	user.Password = hashed
	user.Role = 1 // customer

	// insert user to database
	insert := a.DB.Table("users").Create(user)
	if insert.Error != nil {
		t.Errorf("error create new user %s", insert.Error.Error())
	}

	// get user id
	userId := user.Id

	// create address model
	address := new(models.Address)
	address.UserID = userId
	address.Kodepos = r.Kodepos
	address.Provinsi = r.Provinsi
	address.Kota = r.Kota
	address.Kecamatan = r.Kecamatan
	address.Desa = r.Desa
	address.Dusun = r.Dusun
	address.Jalan = r.Jalan
	address.Nomorhp = r.Nomorhp

	// insert address to database
	insertAddr := a.DB.Table("addresses").Create(address)
	if insertAddr.Error != nil {
		t.Errorf("error create new address user %s", insertAddr.Error.Error())
	}

	fmt.Println("User berhasil dibuat")
}

func TestPasswordHash(t *testing.T) {
	hashedPassword := "$2y$12$lksGE6.lHRcG5BW3MNBu/eO7jRJxu8QcU6wRdGUkewvnEI9K28fw2"
	password := "12345678"

	api := controllers.NewAuthAPI(nil)

	if ok := api.CheckPassword(password, hashedPassword); !ok {
		t.Fatal("password tidak valid")
	}
}

func TestHashingPassword(t *testing.T) {
	password := "12345678"
	api := controllers.NewAuthAPI(nil)

	hashed, err := api.HashPassword(password)
	if err != nil {
		t.Fatalf("%v", err)
	}

	fmt.Printf("Hasil : %s", hashed)
}
