package middleware

import (
	"net/http"
	"task-manager-api/utils"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Abort()
			utils.Error(c, http.StatusUnauthorized, "Authorization header is required", "Please provide a valid JWT token in the Authorization header.")
			return
		}

		claims, err := utils.ParseAccessToken(tokenString)
		if err != nil {
			c.Abort()
			utils.Error(c, http.StatusUnauthorized, "Invalid token", err.Error())
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
