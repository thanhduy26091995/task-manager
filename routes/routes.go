package routes

import (
	"task-manager-api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.GET("/tasks", controllers.GetTasks)
		v1.GET("/tasks/:id", controllers.GetTask)
		v1.POST("/tasks", controllers.CreateTask)
		v1.PUT("/tasks/:id", controllers.UpdateTask)
		v1.DELETE("/tasks/:id", controllers.DeleteTask)
		v1.POST("/webhook/apple", controllers.HandleAppleWebhook)
	}
}
