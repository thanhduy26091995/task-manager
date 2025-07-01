package main

import (
	"log"
	"task-manager-api/config"
	"task-manager-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDatabase()
	config.ConnectWithRedis()

	r := gin.Default()
	routes.RegisterRoute(r)
	log.Println(("Starting Task Manager API..."))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Println("Task Manager API is running on port 8080")
}
