package controllers

import (
	"net/http"
	"task-manager-api/config"
	"task-manager-api/services"
	"task-manager-api/utils"

	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	userId, err := services.Authenticate(config.DB, input.Email, input.Password)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Generate JWT token
	accessToken, _ := utils.GenerateAccessToken(userId)
	refreshToken, _ := utils.GenerateRefreshToken(userId)

	utils.Success(c, http.StatusOK, gin.H{"token": accessToken, "refresh_token": refreshToken})
}

func Register(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	if err := services.Register(config.DB, input.Email, input.Password); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(c, http.StatusCreated, "User registered successfully")
}

func Refresh(c *gin.Context) {
	var refreshInput RefreshInput
	if err := c.ShouldBindJSON(&refreshInput); err != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input format")
		return
	}

	claims, err := utils.ParseRefreshToken(refreshInput.RefreshToken)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	userId := claims["user_id"].(uint)
	accessToken, err := utils.GenerateAccessToken(userId)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate access token")
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"access_token": accessToken})
}
