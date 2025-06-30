package middleware

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	limiterRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimitMiddleware() gin.HandlerFunc {
	// Defind rate limit parameters
	// For example, allow 100 requests per minute
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  100,
	}

	// Setup Redis store and middleware
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisAddr := os.Getenv("REDIS_SERVER")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		panic("Invalid REDIS_DB value, must be an integer")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr, // Adjust as necessary
		Password: redisPassword,
		DB:       redisDB, // Use default DB
	})

	// Create store with Redis backend
	store, err := limiterRedis.NewStoreWithOptions(rdb, limiter.StoreOptions{
		Prefix:   "rate_limit",
		MaxRetry: 3,
	})

	if err != nil {
		panic(err) // Handle error appropriately in production code
	}
	return ginlimiter.NewMiddleware(limiter.New(store, rate))
}
