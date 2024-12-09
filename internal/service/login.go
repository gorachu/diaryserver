package service

import (
	"crypto/tls"
	"diaryserver/internal/config"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (AuthService *AuthService) Login(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling login request")

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

	logger.Debug("TLS connection details",
		"version", c.Request.TLS.Version,
		"cipher_suite", c.Request.TLS.CipherSuite,
		"server_name", c.Request.TLS.ServerName)

	type LoginCredentials struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
		Password string `json:"password" binding:"required,min=8"`
	}

	type LoginRequest struct {
		User LoginCredentials `json:"user" binding:"required"`
	}

	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("invalid request body", "error", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := AuthService.ValidateUser(request.User.Username, request.User.Password)
	if err != nil {
		logger.Error("failed to validate user", "error", err)
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, refreshToken, err := AuthService.GenerateTokens(user.UserID)
	if err != nil {
		logger.Error("failed to generate tokens", "error", err)
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		int(c.MustGet("cfg").(*config.Config).JWT.AccessTokenTTL.Seconds()), // время жизни в секундах
		"/",
		"",   // домен
		true, // secure flag (только HTTPS)
		true, // httpOnly flag
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(c.MustGet("cfg").(*config.Config).JWT.RefreshTokenTTL.Seconds()), // время жизни в секундах
		"/",
		"",   // домен
		true, // secure flag (только HTTPS)
		true, // httpOnly flag
	)

	c.JSON(200, gin.H{
		"message": "login successful",
	})
}
