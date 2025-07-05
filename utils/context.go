package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) uint {
	val, exists := c.Get("user_id")
	if !exists {
		Error(c, http.StatusBadRequest, "User ID not found in context", "The user ID is required for this operation.")
		c.Abort()
		return 0
	}

	userID, ok := val.(uint)
	if !ok {
		Error(c, http.StatusBadRequest, "Invalid user ID type", "The user ID must be a valid unsigned integer.")
		c.Abort()
		return 0
	}

	return userID
}
