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
			utils.Error(c, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		claims, err := utils.ParseAccessToken(tokenString)
		if err != nil {
			c.Abort()
			utils.Error(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
