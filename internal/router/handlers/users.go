package handlers

import (
	"log/slog"

	"diaryserver/internal/storage/sqlite"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUser(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling users requeset")

	username := c.Param("username")
	if username == "" {
		logger.Error("username is required")
		c.JSON(400, gin.H{"error": "username is required"})
		return
	}
	user, err := h.storage.GetUser(username)
	if err != nil {
		logger.Error("failed to get user", "error", err)
		c.JSON(500, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(200, user)
}
func (h *Handler) GetUsers(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling users requeset")

	users, err := h.storage.GetUsers()
	if err != nil {
		logger.Error("failed to get users", "error", err)
		c.JSON(500, gin.H{"error": "failed to get users"})
		return
	}

	c.JSON(200, users)
}

func (h *Handler) CreateUser(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling create user request")

	var user sqlite.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("invalid request body", "error", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.storage.AddUser(user); err != nil {
		logger.Error("failed to create user", "error", err)
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(201, user)
}
