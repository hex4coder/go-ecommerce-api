package main

import (
	"github.com/hex4coder/go-ecommerce-api/controllers"
	"github.com/hex4coder/go-ecommerce-api/models"
)

type UserInterface interface {
	GetUserById(id int) (models.User, error)
	GetUserAddressById(id int) (models.Address, error)
}

type AuthInterface interface {
	Login(*controllers.LoginRequest) (string, error) // jwt, error
	Register(*controllers.RegisterRequest) error
	Logout() error
	HashPassword(string) (string, error)
	CheckPassword(string, string) bool
	GetJWTConfig() *controllers.JWTConfig
}

type KategoriInterface interface {
	GetAll() ([]*models.Kategori, error)
	GetById(int) (*models.Kategori, error)
}

type BrandInterface interface {
	GetAll() ([]*models.Brand, error)
	GetById(int) (*models.Brand, error)
}

type ProductInterface interface {
	// order by desc created_at / yang terbaru paling atas
	GetAllProducts() ([]*models.Product, error)
	GetProductsByCategoryID(int) ([]*models.Product, error)
	GetProductsByBrandID(int) ([]*models.Product, error)
	GetDetailProduct(int) (*models.Product, error)
	GetProductPhotosByID(int) ([]*models.PhotoProducts, error)
	GetUkuranProdukByID(int) ([]*models.UkuranProduks, error)
	GetPopularProducts(int) ([]*models.Product, error)
}

type OrderInterface interface {
	GetMyOrders(int) ([]*models.Order, error)
}
