package models

type BaseResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

type PagingResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	TotalPages int `json:"total_pages"`
}