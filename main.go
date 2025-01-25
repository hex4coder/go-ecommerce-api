package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hex4coder/go-ecommerce-api/controllers"
	"github.com/hex4coder/go-ecommerce-api/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type App struct {
	// controllers
	auth     AuthInterface
	user     UserInterface
	kategori KategoriInterface
	brand    BrandInterface
	product  ProductInterface
	promo    PromoCodeInterface
	order    OrderInterface

	// app base
	router *gin.Engine
	db     *gorm.DB
	url    string
}

func NewApp(db *gorm.DB, url string, router *gin.Engine) *App {
	return &App{
		auth:     controllers.NewAuthAPI(db),
		user:     controllers.NewUserAPI(db),
		kategori: controllers.NewKategoriAPI(db),
		brand:    controllers.NewBrandAPI(db),
		product:  controllers.NewProductAPI(db),
		promo:    controllers.NewPromoCodeAPI(db),
		order:    controllers.NewOrderAPI(db),

		db:     db,
		router: router,
		url:    url,
	}
}

func (app *App) Run() {
	port := 3000
	log.Printf("%s - [STARTED] - Server started at port %d\n", time.Now(), port)

	err := app.router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed start server : %v", err)
	}
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

	// build router using logger middleware
	router := gin.Default()
	router.Use(gin.Logger())

	// buat koneksi
	db, err := models.ConnectDB()
	if err != nil {
		panic(err)
	}

	// buat app
	app := NewApp(db, appUrl, router)

	// register routes
	app.RegisterRoutes()

	// run the router
	app.Run()
}
