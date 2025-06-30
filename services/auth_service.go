package services

import (
	"errors"
	"task-manager-api/models"
	"task-manager-api/utils"

	"gorm.io/gorm"
)

func Authenticate(db *gorm.DB, email, password string) (uint, error) {
	var user models.User

	// Validate email format
	if !utils.IsValidEmail(email) {
		return 0, errors.New("invalid email format")
	}
	// Validate password format
	if !utils.IsValidPassword(password) {
		return 0, errors.New("password must be at least 6 characters long")
	}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return 0, errors.New("user not found")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return 0, errors.New("invalid password")
	}

	return user.ID, nil
}

func Register(db *gorm.DB, email, password string) error {
	var existing models.User

	// Validate email format
	if !utils.IsValidEmail(email) {
		return errors.New("invalid email format")
	}
	// Validate password format
	if !utils.IsValidPassword(password) {
		return errors.New("password must be at least 6 characters long")
	}

	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    email,
		Password: hashedPassword,
	}

	return db.Create(&user).Error
}
