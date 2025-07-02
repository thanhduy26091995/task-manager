package middleware

import (
	"context"
	"net/http"
	"strconv"
	"task-manager-api/config"
	"task-manager-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	limiterRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

var rateLimit = 10
var rateLimitPeriod = 1 * time.Minute
var redisKeyPrefix = "rate:limit"

func RateLimitMiddleware() gin.HandlerFunc {
	// Defind rate limit parameters
	// For example, allow 100 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  100,
	}

	// Create store with Redis backend
	store, err := limiterRedis.NewStoreWithOptions(config.RedisClient, limiter.StoreOptions{
		Prefix:   "rate_limit",
		MaxRetry: 3,
	})

	if err != nil {
		panic(err) // Handle error appropriately in production code
	}
	return ginlimiter.NewMiddleware(limiter.New(store, rate))
}

func RateLimitPerIPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := redisKeyPrefix + ip

		rdb := config.RedisClient
		// Increase and get the current count
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "Redis error", err.Error())
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, rateLimitPeriod) // Set expiration for the key
		}

		if count > int64(rateLimit) {
			ttl, _ := rdb.TTL(ctx, key).Result()
			c.Header("X-RateLimit-Limit", strconv.Itoa(rateLimit))
			c.Header("X-Retry-After", strconv.Itoa(int(ttl.Seconds())))
			c.Abort()
			utils.Error(c, http.StatusTooManyRequests, "Rate limit exceeded. Try again later.", "You have exceeded the rate limit of "+strconv.Itoa(rateLimit)+" requests per "+rateLimitPeriod.String()+". Please try again after "+ttl.String()+".")
			return
		}

		c.Next()
	}
}
