package handlers

import (
	"diaryserver/internal/storage/sqlite"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (h *Handler) LoadCalendar(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling load user requeset")
	type CustomData struct {
		Key       string   `json:"key"`
		Trainings []string `json:"trainings"`
	}
	type Request struct {
		Data       []string     `json:"data"`
		CustomData []CustomData `json:"customData"`
	}
	request := Request{
		Data: []string{"2025-01-01", "2025-01-02"},
		CustomData: []CustomData{
			{
				Key:       "2025-01-01",
				Trainings: []string{"09:00 - 10:30", "16:00 - 17:00"},
			},
			{
				Key:       "2025-01-02",
				Trainings: []string{"09:01 - 10:31", "16:01 - 17:01"},
			},
		},
	}
	c.JSON(200, request)
}

func (h *Handler) LoadTrainings(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling load trainings requeset")

	date := c.Param("date")
	if date == "" {
		logger.Error("date is required")
		c.JSON(400, gin.H{"error": "date is required"})
		return
	}
	type TrainingsInDay struct {
		Time     string   `json:"time"`
		Duration string   `json:"duration"`
		Sets     []string `json:"sets"`
	}
	type Request struct {
		TrainingsInDay []TrainingsInDay `json:"trainingsInDay"`
	}
	request := Request{
		TrainingsInDay: []TrainingsInDay{
			{
				Time:     "7:00",
				Duration: "1h2m",
				Sets:     []string{"1 set", "2 set"},
			},
			{
				Time:     "18:00",
				Duration: "1h30m",
				Sets:     []string{"1 set", "2 set"},
			},
		},
	}
	c.JSON(200, request)
}
func (h *Handler) CreateTraining(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling creating training requeset")
	user_ID_Object, successful := c.Get("user_id")
	if !successful {
		logger.Error("User id not found")
		c.JSON(400, gin.H{"error": "User id not found"})
		return
	}

	user_ID, ok := user_ID_Object.(int)
	if !ok {
		logger.Error("User id is not integer")
		c.JSON(400, gin.H{"error": "User id is not integer"})
		return
	}
	type Workout struct {
		Date string `json:"date"`
		Note string `json:"note"`
	}
	var workout_temp Workout
	if err := c.ShouldBindJSON(&workout_temp); err != nil {
		logger.Error("invalid request body", "error", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	var workout = sqlite.Workout{
		UserID: user_ID,
		Date:   workout_temp.Date,
		Notes:  workout_temp.Note,
	}
	if err := h.storage.AddWorkout(workout); err != nil {
		logger.Error("failed to create workout", "error", err)
		c.JSON(500, gin.H{"error": "failed to create workout"})
		return
	}
	c.JSON(201, workout)
}
