package middleware

import (
	"context"
	"net/http"
	"task-manager-api/config"
	"task-manager-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request has an Idempotency-Key header
		idempotencyKey := c.GetHeader("Idempotency-Key")
		if idempotencyKey == "" {
			c.Abort()
			utils.Error(c, http.StatusBadRequest, "Idempotency-Key header is required")
			return
		}

		rdb := config.RedisClient
		ctx := context.Background()
		// Check if the key already exists in Redis
		val, err := rdb.Get(ctx, idempotencyKey).Result()
		if err == nil && val != "" {
			// If the key exists, it means the request has already been processed
			c.Abort()
			utils.Error(c, http.StatusConflict, "Request has already been processed")
			return
		}
		
		rdb.Set(ctx, idempotencyKey, "processed", time.Hour)
		c.Next()
	}
}