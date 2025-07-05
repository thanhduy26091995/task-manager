package main

import (
	"task-manager-api/config"
	"task-manager-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDatabase()
	config.ConnectWithRedis()
	config.InitLogger()

	r := gin.Default()
	routes.RegisterRoute(r)
	config.Log.Info("Starting Task Manager API...")

	if err := r.Run(":8080"); err != nil {
		config.Log.Sugar().Fatalf("Failed to start server: %v", err)
	}
	config.Log.Info("Task Manager API is running on port 8080")
}
