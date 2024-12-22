package controllers

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct{}

type AuthAPI struct {
	DB *gorm.DB
}

func NewAuthAPI(db *gorm.DB) *AuthAPI {
	return &AuthAPI{
		DB: db,
	}
}

func (a *AuthAPI) Login(r *LoginRequest) (string, error) {
	return "", nil
}
func (a *AuthAPI) Register(r *RegisterRequest) error {
	return nil
}
func (a *AuthAPI) Logout() error {
	return nil
}
func (a *AuthAPI) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func (a *AuthAPI) CheckPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
