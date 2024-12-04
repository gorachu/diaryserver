package service

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (AuthService *AuthService) Login(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling login requeset")

	var request struct {
		Username     string `json:"username"`
		PasswordHash string `json:"password_hash"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("invalid request body", "error", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	user, err := AuthService.ValidateUser(request.Username, request.PasswordHash)
	if err != nil {
		logger.Error("failed to validate user", "error", err)
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := AuthService.GenerateAccessToken(user.UserID)
	if err != nil {
		logger.Error("failed to generate access token", "error", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	refreshToken, err := AuthService.GenerateRefreshToken(user.UserID)
	if err != nil {
		logger.Error("failed to generate refresh token", "error", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
