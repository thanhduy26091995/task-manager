package utils

import (
	"net/http"
	"task-manager-api/models"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, models.BaseResponse{
		Status:  "success",
		Message: "Request was successful",
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, models.BaseResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	})
}

func Paging(c *gin.Context, data interface{}, page, limit, total int) {
	c.JSON(http.StatusOK, models.PagingResponse {
		Status:  "success",
		Message: "Request was successful",
		Data:    data,
		Pagination: models.Pagination{
			Page:         page,
			Limit:       limit,
			Total:       total,
			TotalPages:  (total + limit - 1) / limit, // Calculate total pages
		},
	})
}