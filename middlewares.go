package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)



func AuthMiddleware(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		token := ""

		// get cookie from request
		cookie, err := c.Cookie("jwt")
		if err == http.ErrNoCookie {
			// no cookie, perform header bearer token authorization
			// get headers bearer token
			tokenString := c.GetHeader("Authorization")
			if tokenString == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
					"status":  "error",
					"message": "token empty",
				})
				return
			}

			// get the token from the string Authorization: Bearer .....token.....
			token = tokenString[len("Bearer:"):]
		} else {
			// set token to cookie
			token = cookie
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
				"status":  "error",
				"message": "token empty",
			})
			return
		}

		err = app.auth.GetJWTConfig().VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.Next()
	}
}
