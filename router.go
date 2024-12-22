package main

import (
	"fmt"
	"net/http"
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

func APISuccessResponse(statusCode int, message string, data any, c *gin.Context) {
	c.IndentedJSON(statusCode, map[string]any{
		"message": message,
		"status":  "success",
		"data":    data,
	})
}

func (app *App) RegisterRoutes() {
	// custom validators
	cv := NewCustomValidator()

	// ---------------------------------AUTH API-----------------------------
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
		APISuccessResponse(http.StatusOK, "login function", map[string]any{
			"jwt": jwt,
		}, c)
	})
	app.router.POST("/register", func(c *gin.Context) {})

	// ---------------------------------USERS-------------------------------

	ar := app.router.Group("/", AuthMiddleware(app))
	ar.POST("/logout", func(c *gin.Context) {
		c.SetCookie("jwt", "", 0, "/", app.url, true, true)
		err := app.auth.Logout()
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}
		APISuccessResponse(http.StatusOK, "logout", nil, c)
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

		APISuccessResponse(http.StatusOK, fmt.Sprintf("user with id %d", id), user, c)

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

		APISuccessResponse(http.StatusOK, fmt.Sprintf("address with user id %d", id), address, c)
	})
}
