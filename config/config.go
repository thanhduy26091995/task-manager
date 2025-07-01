package config

import (
	"log"
	"os"
	"strconv"
	"task-manager-api/models"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
}

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = db // Assign the DB variable to the global DB variable
	log.Println("Database connection established")
	migrationErr := DB.AutoMigrate(models.Task{}, models.User{}) // Assuming Task is the model you want to migrate
	if migrationErr != nil {
		log.Fatalf("Failed to migrate database: %v", migrationErr)
	}
	log.Println("Database migration completed")
	log.Println("Database is ready to use")
}

func ConnectWithRedis() {
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisAddr := os.Getenv("REDIS_SERVER")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value, must be an integer: %v", err)
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	log.Println("Redis connection established")
}
