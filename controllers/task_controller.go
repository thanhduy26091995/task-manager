package controllers

import (
	"net/http"
	"strconv"
	"task-manager-api/config"
	"task-manager-api/models"
	"task-manager-api/utils"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	var total int64

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	completed := c.Query("completed")

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit
	query := config.DB.Model(&models.Task{})

	// Apply filter
	if completed != "" {
		if completed == "true" || completed == "1" {
			query = query.Where("completed = ?", true)
		} else if completed == "false" || completed == "0" {
			query = query.Where("completed = ?", false)
		}
	}

	// Count total tasks
	query.Count(&total)

	// Apply pagination
	err := query.Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		utils.Error(c, 500, "Failed to retrieve tasks")
		return
	}

	utils.Paging(c, tasks, page, limit, int(total))
}

func GetTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")
	if err := config.DB.First(&task, id).Error; err != nil {
		utils.Error(c, 404, "Task not found")
		return
	}

	utils.Success(c, http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		utils.Error(c, 400, "Invalid request data")
		return
	}

	if err := config.DB.Create(&task).Error; err != nil {
		utils.Error(c, 400, "Failed to create task")
		return
	}

	utils.Success(c, http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	// Check if the task exists
	if err := config.DB.First(&task, id).Error; err != nil {
		utils.Error(c, 404, "Task not found")
		return
	}

	// Check if the request body is valid
	if err := c.ShouldBindJSON(&task); err != nil {
		utils.Error(c, 400, "Invalid request data")
		return
	}

	// Update the task
	if err := config.DB.Save(&task).Error; err != nil {
		utils.Error(c, 400, "Failed to update task")
		return
	}

	utils.Success(c, http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	// Check if the task exists
	if err := config.DB.First(&task, id).Error; err != nil {
		utils.Error(c, 404, "Task not found")
		return
	}

	// Delete the task
	if err := config.DB.Delete(&task).Error; err != nil {
		utils.Error(c, 400, "Failed to delete task")
		return
	}

	utils.Success(c, http.StatusOK, task)
}
