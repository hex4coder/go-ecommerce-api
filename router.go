package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hex4coder/go-ecommerce-api/controllers"
)

func APIErrorResponse(statusCode int, message string, c *gin.Context) {
	c.IndentedJSON(statusCode, map[string]any{
		"message": message,
		"status":  "error",
	})
}

func APISuccessResponse(message string, data any, c *gin.Context) {
	c.IndentedJSON(http.StatusOK, map[string]any{
		"message": message,
		"status":  "success",
		"data":    data,
	})
}

func (app *App) RegisterRoutes() {
	// custom validators
	cv := NewCustomValidator()

	// ---------------------------------AUTH API------------------------------------------------
	app.router.GET("/", func(c *gin.Context) {
		APISuccessResponse("we are online and ready", map[string]any{
			"programmed_by": "Ardan, S.Kom",
			"design_by":     "Ardan, S.Kom",
		}, c)
	})

	app.router.POST("/login", func(c *gin.Context) {
		loginReq := new(controllers.LoginRequest)

		// request binding
		if err := c.BindJSON(loginReq); err != nil {
			APIErrorResponse(http.StatusBadRequest, err.Error(), c)
			return
		}

		// validator
		if err := cv.Validate(loginReq); err != nil {
			APIErrorResponse(http.StatusBadRequest, err.Error(), c)
			return
		}

		// login function
		jwt, err := app.auth.Login(loginReq)
		if err != nil {
			APIErrorResponse(http.StatusBadRequest, err.Error(), c)
			return
		}

		// set cookie
		c.SetCookie("jwt", jwt, int(time.Now().Add(app.auth.GetJWTConfig().ExpiredDuration).Unix()), "/", app.url, true, true)

		// jwt account
		APISuccessResponse("login function", map[string]any{
			"jwt": jwt,
		}, c)
	})

	// register function
	app.router.POST("/register", func(c *gin.Context) {
		registerReq := new(controllers.RegisterRequest)

		// request binding
		if err := c.BindJSON(registerReq); err != nil {
			APIErrorResponse(http.StatusBadRequest, err.Error(), c)
			return
		}

		// validator
		if err := cv.Validate(registerReq); err != nil {
			APIErrorResponse(http.StatusBadRequest, err.Error(), c)
			return
		}

		// create new user
		if err := app.auth.Register(registerReq); err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		// success
		APISuccessResponse("new user created", registerReq, c)
	})

	// ambil list kategori
	app.router.GET("/categories", func(c *gin.Context) {
		data, err := app.kategori.GetAll()

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse("list kategori", data, c)
	})

	// ambil list brand
	app.router.GET("/brands", func(c *gin.Context) {
		data, err := app.brand.GetAll()

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse("list brand", data, c)
	})

	// ambil list produk berdasarkan waktu terbaru
	app.router.GET("/products", func(c *gin.Context) {
		data, err := app.product.GetAllProducts()

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse("list of products", data, c)
	})

	// get detail product by id
	app.router.GET("/product/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("invalid product_id given %s", idStr), c)
			return
		}

		// convert string to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("error conver id %s => %v", idStr, err), c)
			return
		}

		// buat temporary data
		data := make(map[string]any)

		// find products by categori id
		product, err := app.product.GetDetailProduct(id)
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		// assign product to data
		data["product"] = product

		// cari ukuran product
		ukuran, err := app.product.GetUkuranProdukByID(id)
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		// assign ukuran to map
		data["ukuran"] = ukuran

		// get list photos of product
		photos, err := app.product.GetProductPhotosByID(id)
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		// assign photos to map
		data["photos"] = photos

		// cari brand dari product tersebut
		brand, err := app.brand.GetById(product.BrandID)
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		// assign brand to data
		data["brand"] = brand

		// cari kategori detail dari product
		cat, err := app.kategori.GetById(product.KategoriID)
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		// assign category to data
		data["kategori"] = cat

		// return success and data
		APISuccessResponse("get detail product", data, c)
	})

	// ambil list product berdasarkan merek
	app.router.GET("/kategori/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("invalid kategori_id given %s", idStr), c)
			return
		}

		// convert string to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("error conver id %s => %v", idStr, err), c)
			return
		}

		// find products by categori id
		data, err := app.product.GetProductsByCategoryID(id)

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse("list of products", data, c)
	})

	// ambil list products berdasarkan brand id
	app.router.GET("/brand/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("invalid brand_id given %s", idStr), c)
			return
		}

		// convert string to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("error conver id %s => %v", idStr, err), c)
			return
		}

		// find products by categori id
		data, err := app.product.GetProductsByBrandID(id)

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse("list of products", data, c)
	})

	// ambil list foto dari produk
	app.router.GET("/photos/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("invalid produk_id given %s", idStr), c)
			return
		}

		// convert string to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("error conver id %s => %v", idStr, err), c)
			return
		}

		data, err := app.product.GetProductPhotosByID(id)

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse("list of photo products", data, c)
	})

	// ambil list ukuran dari produk
	app.router.GET("/ukuran/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("invalid produk_id given %s", idStr), c)
			return
		}

		// convert string to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("error conver id %s => %v", idStr, err), c)
			return
		}

		data, err := app.product.GetUkuranProdukByID(id)

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse("list of photo products", data, c)
	})

	// get popular products
	app.router.GET("/popular-products/:limit", func(c *gin.Context) {
		limitStr := c.Param("limit")
		if limitStr == "" {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("invalid limit given %s", limitStr), c)
			return
		}

		// convert string to int
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			APIErrorResponse(http.StatusBadRequest, fmt.Sprintf("error conver limit %s => %v", limitStr, err), c)
			return
		}

		// find products by categori id
		data, err := app.product.GetPopularProducts(limit)

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse("list of popular products", data, c)
	})

	// ---------------------------------USERS-------------------------------

	ar := app.router.Group("/", AuthMiddleware(app))
	ar.POST("/logout", func(c *gin.Context) {
		c.SetCookie("jwt", "", 0, "/", app.url, true, true)
		err := app.auth.Logout()
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}
		APISuccessResponse("logout", nil, c)
	})

	// get user by id
	ar.GET("/user", func(c *gin.Context) {
		claims, e := c.Get("claims")
		if !e {
			APIErrorResponse(http.StatusInternalServerError, "failed to set context with claims", c)
			return
		}
		cl := claims.(*controllers.MyClaims)

		id := cl.Id
		user, err := app.user.GetUserById(id)
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, fmt.Sprintf("failed to get user with %d error: %s", id, err.Error()), c)
			return
		}

		APISuccessResponse(fmt.Sprintf("user with id %d", id), user, c)

	})

	// get users address with id
	ar.GET("/user-address", func(c *gin.Context) {
		claims, e := c.Get("claims")
		if !e {
			APIErrorResponse(http.StatusInternalServerError, "failed to set context with claims", c)
			return
		}

		cl := claims.(*controllers.MyClaims)
		id := cl.Id

		address, err := app.user.GetUserAddressById(id)
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, fmt.Sprintf("failed to get user with %d error: %s", id, err.Error()), c)
			return
		}

		APISuccessResponse(fmt.Sprintf("address with user id %d", id), address, c)
	})

	// promo code check
	ar.GET("/check-promo/:code", func(c *gin.Context) {

		// get code parameter
		code := c.Param("code")

		// if no data in code
		if len(code) == 0 {
			// error with no data
			APIErrorResponse(http.StatusBadRequest, "no code is passed", c)
			// return
			return
		}

		// find the code
		promo, err := app.promo.CheckPromo(code)

		// check error
		if err != nil {
			// return with error
			APIErrorResponse(http.StatusBadRequest, err.Error(), c)
			return
		}

		// return data if success
		APISuccessResponse(fmt.Sprintf("promo with code %s", code), promo, c)
	})

}
