package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTokenFromCookiesOrBearer(c *gin.Context) (string, error) {
	c.Header("Content-Type", "application/json")

	token := ""

	// get cookie from request
	cookie, err := c.Cookie("jwt")
	if err != http.ErrNoCookie {
		// set token to cookie
		token = cookie
	}

	// no cookie, perform header bearer token authorization
	// get headers bearer token
	tokenString := c.GetHeader("Authorization")

	// get the token from the string Authorization: Bearer .....token.....
	if tokenString != "" {
		token = tokenString[len("Bearer "):]
	}

	if token == "" {
		return "", fmt.Errorf("token empty")
	}

	return token, nil
}

func AuthMiddleware(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := GetTokenFromCookiesOrBearer(c)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		cl, err := app.auth.GetJWTConfig().GetClaimsFromToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.Set("claims", cl)

		c.Next()
	}
}
