package utils

import (
	"net/http"
	"task-manager-api/models"

	"github.com/gin-gonic/gin"
)

type ProblemDetails struct {
	Type     string      `json:"type,omitempty"`
	Title    string      `json:"title"`
	Status   int         `json:"status"`
	Detail   string      `json:"detail,omitempty"`
	Instance string      `json:"instance,omitempty"`
	Extras   interface{} `json:"extras,omitempty"`
}

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, models.BaseResponse{
		Status:  "success",
		Message: "Request was successful",
		Data:    data,
	})
}

func Error(c *gin.Context, status int, title, detail string, extras ...interface{}) {
	problem := ProblemDetails{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: c.Request.URL.Path,
	}
	if len(extras) > 0 {
		problem.Extras = extras[0]
	}

	c.JSON(status, problem)
}

func Paging(c *gin.Context, data interface{}, page, limit, total int) {
	c.JSON(http.StatusOK, models.PagingResponse{
		Status:  "success",
		Message: "Request was successful",
		Data:    data,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: (total + limit - 1) / limit, // Calculate total pages
		},
	})
}
