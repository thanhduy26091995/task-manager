package config

import (
	"log"
	"os"
	"task-manager-api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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
	migrationErr := DB.AutoMigrate(models.Task{}) // Assuming Task is the model you want to migrate
	if migrationErr != nil {
		log.Fatalf("Failed to migrate database: %v", migrationErr)
	}
	log.Println("Database migration completed")
	log.Println("Database is ready to use")
}
