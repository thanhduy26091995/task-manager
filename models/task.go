package models

import "time"

type Task struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Completed bool `gorm:"default:false" json:"completed"`
	DueDate *time.Time `json:"due_date,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"` 
}