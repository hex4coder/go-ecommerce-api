package controllers

import (
	"fmt"
	"testing"
)

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
