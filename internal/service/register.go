package service

import (
	"crypto/tls"
	"log/slog"

	"diaryserver/internal/storage/sqlite"

	"github.com/gin-gonic/gin"
)

func (AuthService *AuthService) Register(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling register request")

	if c.Request.TLS == nil {
		logger.Error("non-secure connection rejected")
		c.JSON(400, gin.H{"error": "HTTPS required"})
		return
	}

	// Проверка версии TLS
	tlsVersion := c.Request.TLS.Version
	if tlsVersion < tls.VersionTLS12 {
		logger.Error("outdated TLS version", "version", tlsVersion)
		c.JSON(400, gin.H{"error": "TLS 1.2 or higher required"})
		return
	}
	type UserData struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	type RegisterRequest struct {
		User UserData `json:"user" binding:"required"`
	}
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("invalid request", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	logger.Debug("received registration data",
		"username", request.User.Username,
		"email", request.User.Email)
	hashedPassword, err := HashPassword(request.User.Password)
	if err != nil {
		logger.Error("failed to hash password", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	err = c.MustGet("db").(*sqlite.Storage).AddUser(sqlite.User{
		Username:     request.User.Username,
		Email:        request.User.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		logger.Error("failed to create user", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "Registration successful"})

}
