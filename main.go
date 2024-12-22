package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hex4coder/go-ecommerce-api/controllers"
	"github.com/hex4coder/go-ecommerce-api/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	auth   AuthInterface
	user   UserAPIInterface
	router *gin.Engine
	db     *gorm.DB
	url    string
}

func NewApp(db *gorm.DB, url string) *App {
	return &App{
		db:     db,
		auth:   controllers.NewAuthAPI(db),
		user:   controllers.NewUserAPI(db),
		router: gin.Default(),
		url:    url,
	}
}

func (app *App) Run() {
	port := 3000
	app.router.Run(fmt.Sprintf(":%d", port))
}

func main() {

	// load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// app name
	appUrl := os.Getenv("APP_NAME")

	// release gin
	gmode := os.Getenv("GIN_MODE")
	gin.SetMode(gmode)

	// buat koneksi
	db, err := models.ConnectDB()
	if err != nil {
		panic(err)
	}

	// buat app
	app := NewApp(db, appUrl)

	// register routes
	app.RegisterRoutes()

	// run the router
	app.Run()
}
