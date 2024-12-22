package main

import (
	"github.com/hex4coder/go-ecommerce-api/controllers"
	"github.com/hex4coder/go-ecommerce-api/models"
)

type UserAPIInterface interface {
	GetUserById(id int) (models.User, error)
	GetUserAddressById(id int) (models.Addresses, error)
}

type AuthInterface interface {
	Login(*controllers.LoginRequest) (string, error) // jwt, error
	Register(*controllers.RegisterRequest) error
	Logout() error
	HashPassword(string) (string, error)
	CheckPassword(string, string) bool
	GetJWTConfig() *controllers.JWTConfig
}
