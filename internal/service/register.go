package service

import (
	"crypto/tls"
	"log/slog"

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
	type RegisterRequest struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("invalid request", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	logger.Debug("received registration data",
		"username", request.Username,
		"email", request.Email)

	c.JSON(200, gin.H{"message": "Registration successful"})

}
