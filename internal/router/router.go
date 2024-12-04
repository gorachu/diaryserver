package router

import (
	"diaryserver/internal/router/handlers"
	"diaryserver/internal/storage/sqlite"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(storage *sqlite.Storage, log *slog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		reqLogger := log.With(
			slog.String("path", c.Request.URL.Path),
			slog.String("method", c.Request.Method),
			slog.String("client_ip", c.ClientIP()),
		)
		c.Set("logger", reqLogger)
		start := time.Now()
		c.Next()
		reqLogger.Info("request completed",
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.Int("errors", len(c.Errors)),
		)
	})
	users := r.Group("/users")
	{
		users.GET("", handlers.NewHandlers(storage, log).GetUsers)
		users.GET("/:username", handlers.NewHandlers(storage, log).GetUser)
		users.POST("", handlers.NewHandlers(storage, log).CreateUser)
	}

	return r
}