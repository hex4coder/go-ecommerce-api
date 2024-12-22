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
}

func NewApp(db *gorm.DB) *App {
	return &App{
		db: db,
		auth: &controllers.AuthAPI{
			DB: db,
		},
		user: &controllers.UserAPI{
			DB: db,
		},
		router: gin.Default(),
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

	// release gin
	gmode := os.Getenv("GIN_MODE")
	gin.SetMode(gmode)

	// buat koneksi
	db, err := models.ConnectDB()
	if err != nil {
		panic(err)
	}

	// buat app
	app := NewApp(db)

	// register routes
	app.RegisterRoutes()

	// run the router
	app.Run()
}
