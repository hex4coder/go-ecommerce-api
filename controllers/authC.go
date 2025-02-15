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

func (c *JWTConfig) GetClaimsFromToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(c.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to casting token claims")
	}

	cl := new(MyClaims)
	cl.Email = claims["email"].(string)
	cl.Id = int(claims["id"].(float64))
	cl.Role = int(claims["role"].(float64))

	return cl, nil
}

func (c *JWTConfig) VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(c.JwtSecret), nil
	})

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

type RegisterRequest struct {
	// biodata user
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`

	// Alamat
	Kodepos   string `json:"kodepos" validate:"required"`
	Provinsi  string `json:"provinsi" validate:"required"`
	Kota      string `json:"kota" validate:"required"`
	Kecamatan string `json:"kecamatan" validate:"required"`
	Desa      string `json:"desa" validate:"required"`
	Dusun     string `json:"dusun" validate:"required"`
	Jalan     string `json:"jalan" validate:"required"`
	Nomorhp   string `json:"nomorhp" validate:"required"`
}

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

func (a *AuthAPI) GetJWTConfig() *JWTConfig {
	return a.JWT
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

	// buat temporary
	var emailCount int64

	// validasi email
	a.DB.Table("users").Where(&models.User{Email: r.Email}).Count(&emailCount)
	if emailCount > 0 {
		return fmt.Errorf("email %s sudah terdaftar", r.Email)
	}

	// buat model baru
	user := new(models.User)

	// create password hash
	hashed, err := a.HashPassword(r.Password)
	if err != nil {
		return fmt.Errorf("failed to hashing the password %s", err.Error())
	}

	// assign data to user
	user.Email = r.Email
	user.Name = r.Name
	user.Password = hashed
	user.Role = 1 // customer

	// insert user to database
	insert := a.DB.Table("users").Create(user)
	if insert.Error != nil {
		return fmt.Errorf("error create new user %s", insert.Error.Error())
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
		return fmt.Errorf("error create new address user %s", insertAddr.Error.Error())
	}

	// no error
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
