package middleware

import (
	"diaryserver/internal/service"
	"diaryserver/internal/storage/sqlite"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(storage *sqlite.Storage, accessSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.MustGet("logger").(*slog.Logger)
		logger.Debug("checking authentication")

		accessToken, err := c.Cookie("access_token")
		logger.Debug("accessToken from cookie", "accessToken", accessToken)
		if err != nil || accessToken == "" {
			logger.Warn("no access token cookie")
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		isBlacklisted, err := storage.IsTokenBlacklisted(accessToken)
		if err != nil {
			logger.Error("failed to check token blacklist", "error", err)
			c.JSON(500, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		if isBlacklisted {
			logger.Warn("token is blacklisted")
			c.JSON(401, gin.H{"error": "token is blacklisted"})
			c.Abort()
			return
		}

		authService := service.NewAuthService(
			storage,
			accessSecret,
			"", // refresh secret не нужен для валидации
			0,  // TTL не нужны для валидации
			0,
		)

		userID, err := authService.ValidateAccessToken(accessToken)
		if err != nil {
			logger.Warn("invalid access token", "error", err)
			c.JSON(401, gin.H{"error": "invalid access token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
