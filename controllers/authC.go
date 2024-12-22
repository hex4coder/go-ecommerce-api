package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hex4coder/go-ecommerce-api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --------------------------------- JWT AREA----------------------------------------------
type JWTConfig struct {
	AppName         string
	JwtSecret       string
	ExpiredDuration time.Duration
}

type MyClaims struct {
	Id    int
	Email string
	Role  int
}

func NewClaimsFromUserModel(user *models.User) *MyClaims {
	return &MyClaims{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}
}

func loadJWTConfig() *JWTConfig {
	config := new(JWTConfig)

	config.AppName = "Go Ecommerce API"
	config.JwtSecret = "terlalu-secret"
	config.ExpiredDuration = time.Hour * 1

	// load config from .env
	secret := os.Getenv("JWT_SECRET")
	if len(secret) > 1 {
		config.JwtSecret = secret
	}

	return config
}

func (c *JWTConfig) GenerateToken(claim *MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    claim.Id,
		"email": claim.Email,
		"role":  claim.Role,
		"exp":   time.Now().Add(c.ExpiredDuration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(c.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (c *JWTConfig) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) { return c.JwtSecret, nil })

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// -----------------------------------AUTH API---------------------------------------------
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct{}

type AuthAPI struct {
	DB  *gorm.DB
	JWT *JWTConfig
}

func NewAuthAPI(db *gorm.DB) *AuthAPI {
	jwtconfig := loadJWTConfig()
	return &AuthAPI{
		DB:  db,
		JWT: jwtconfig,
	}
}

func (a *AuthAPI) Login(r *LoginRequest) (string, error) {

	// check email and password
	user := new(models.User)
	res := a.DB.Preload("Address").Where(&models.User{Email: r.Email}).First(user)

	// no account
	if res.RowsAffected < 1 {
		return "", fmt.Errorf("no user with email %s", r.Email)
		// return "", fmt.Errorf("creadential tidak valid")
	}

	// check pass
	hashedPassword := user.Password
	if !a.CheckPassword(r.Password, hashedPassword) {
		return "", fmt.Errorf("password tidak valid")
		// return "", fmt.Errorf("creadential tidak valid")
	}

	// create jwt token
	cl := NewClaimsFromUserModel(user)
	token, err := a.JWT.GenerateToken(cl)
	if err != nil {
		return "", err
	}

	return token, nil
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
