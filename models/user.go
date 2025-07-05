package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"-"`

	Tassks []Task `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
}