package controllers

import (
	"net/http"
	"strconv"
	"task-manager-api/config"
	"task-manager-api/models"
	"task-manager-api/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	var total int64

	userID := utils.GetUserID(c);
	if userID == 0 {
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	completed := c.Query("completed")

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit
	query := config.DB.Model(&models.Task{}).Where("user_id = ?", userID)

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
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&tasks).Error
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to retrieve tasks", err.Error())
		return
	}

	utils.Paging(c, tasks, page, limit, int(total))
}

func GetTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")
	userID := utils.GetUserID(c)
	if err := config.DB.First(&task, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Task not found", err.Error())
		return
	}

	if task.UserID != userID {
		utils.Error(c, http.StatusForbidden, "You do not have permission to access this task", "Task does not belong to the user")
		return	
	}
	
	utils.Success(c, http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		utils.Error(c, http.StatusNotFound, "Invalid request data", err.Error())
		return
	}

	task.UserID = utils.GetUserID(c)
	if err := config.DB.Create(&task).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Failed to create task", err.Error())
		return
	}

	// Logging
	go func(t models.Task) {
		defer func() {
			if r := recover(); r != nil {
				config.Log.Error("Recovered from panic in CreateTask", zap.Any("error", r))
			}
		}()

		config.Log.Info("Task created", zap.Int("ID", int(t.ID)), zap.String("Title", t.Title), zap.Bool("Completed", t.Completed))
	}(task)

	utils.Success(c, http.StatusCreated, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	userID := utils.GetUserID(c)

	// Check if the task exists
	if err := config.DB.First(&task, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Task not found", err.Error())
		return
	}

	// Ensure the task belongs to the user
	if task.UserID != userID {
		utils.Error(c, http.StatusForbidden, "You do not have permission to update this task", "Task does not belong to the user")
		return
	}

	// Check if the request body is valid
	if err := c.ShouldBindJSON(&task); err != nil {
		utils.Error(c, http.StatusNotFound, "Invalid request data", err.Error())
		return
	}

	// Update the task
	if err := config.DB.Save(&task).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Failed to update task", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	userID := utils.GetUserID(c)

	// Check if the task exists
	if err := config.DB.First(&task, id).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Task not found", err.Error())
		return
	}

	if task.UserID != userID {
		utils.Error(c, http.StatusForbidden, "You do not have permission to delete this task", "Task does not belong to the user")
		return
	}

	// Delete the task
	if err := config.DB.Delete(&task).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Failed to delete task", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, task)
}
