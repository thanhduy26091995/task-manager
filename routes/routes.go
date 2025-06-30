package routes

import (
	"task-manager-api/controllers"
	"task-manager-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(r *gin.Engine) {
	r.POST("/api/v1/login", controllers.Login)
	r.POST("/api/v1/register", controllers.Register)
	r.POST("/api/v1/refresh", controllers.Register)

	protected := r.Group("/api/v1`")
	protected.Use(middleware.JWTAuthMiddleware(), middleware.RateLimitMiddleware())
	{
		protected.GET("/tasks", controllers.GetTasks)
		protected.GET("/tasks/:id", controllers.GetTask)
		protected.POST("/tasks", controllers.CreateTask)
		protected.PUT("/tasks/:id", controllers.UpdateTask)
		protected.DELETE("/tasks/:id", controllers.DeleteTask)
		protected.POST("/webhook/apple", controllers.HandleAppleWebhook)
	}
}
