package main

import (
	"fmt"
	"net/http"
	"strconv"

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

		// jwt account
		APISuccessResponse(http.StatusOK, "login function", map[string]any{
			"jwt": jwt,
		}, c)
	})
	app.router.POST("/register", func(c *gin.Context) {})
	app.router.POST("/logout", func(c *gin.Context) {})

	// ---------------------------------USERS-------------------------------

	// get list of users
	app.router.GET("/users", func(c *gin.Context) {
		users, err := app.user.GetUsers()
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse(http.StatusOK, "list of users", users, c)
	})

	//  list of users with address
	app.router.GET("/users-with-address", func(c *gin.Context) {
		users, err := app.user.GetUsersWithAddress()
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, err.Error(), c)
			return
		}

		APISuccessResponse(http.StatusOK, "list of users with address", users, c)
	})

	// get user by id
	app.router.GET("/user/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, fmt.Sprintf("failed to convert id params %v", idStr), c)
			return
		}

		user, err := app.user.GetUserById(id)
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, fmt.Sprintf("failed to get user with %d error: %s", id, err.Error()), c)
			return
		}

		APISuccessResponse(http.StatusOK, fmt.Sprintf("user with id %d", id), user, c)
	})

	// get users address with id
	app.router.GET("/user-address/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, fmt.Sprintf("failed to convert id params %v", idStr), c)
			return
		}

		address, err := app.user.GetUserAddressById(id)
		if err != nil {
			APIErrorResponse(http.StatusInternalServerError, fmt.Sprintf("failed to get user with %d error: %s", id, err.Error()), c)
			return
		}

		APISuccessResponse(http.StatusOK, fmt.Sprintf("address with user id %d", id), address, c)
	})
}