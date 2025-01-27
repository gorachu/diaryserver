package middleware

import (
	"diaryserver/internal/config"
	"diaryserver/internal/service"
	"diaryserver/internal/storage/sqlite"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(storage *sqlite.Storage, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.MustGet("logger").(*slog.Logger)
		logger.Debug("checking authentication")
		accessToken, err := c.Cookie("access_token")
		logger.Debug("accessToken from cookie", "accessToken", accessToken)
		authService := service.NewAuthService(
			storage,
			cfg.JWT.AccessSecret,
			cfg.JWT.RefreshSecret,
			cfg.JWT.AccessTokenTTL,
			cfg.JWT.RefreshTokenTTL,
		)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				logger.Warn("No access token cookie")
			} else {
				logger.Warn("Couldn't get access token cookie")
				c.JSON(500, gin.H{"error": "getting access token cookie error"})
				c.Abort()
				return
			}
		} else if accessToken != " " {
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
			_, err = authService.ValidateAccessToken(accessToken)
			if err != nil {
				logger.Warn("Access Token validation error")
				c.JSON(401, gin.H{"error": "Access token validation error"})
				c.Abort()
				return
			}
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err != nil || refreshToken == "" {
			logger.Warn("no refresh token cookie")
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		isBlacklisted, err := storage.IsTokenBlacklisted(refreshToken)
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

		newAccessToken, newRefreshToken, err := authService.RefreshTokens(refreshToken)
		if err != nil {
			logger.Warn("invalid refresh token", "error", err)
			c.JSON(401, gin.H{"error": "invalid refresh token"})
			c.Abort()
			return
		}

		c.SetCookie("access_token", newAccessToken, int(cfg.JWT.AccessTokenTTL.Seconds()), "/", "", false, true)
		c.SetCookie("refresh_token", newRefreshToken, int(cfg.JWT.RefreshTokenTTL.Seconds()), "/", "", false, true)

		userID, err := authService.ValidateAccessToken(newAccessToken)
		if err != nil {
			logger.Error("failed to validate new access token", "error", err)
			c.JSON(500, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
