package handlers

import (
	"diaryserver/internal/storage/sqlite"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func (h *Handler) LoadCalendar(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling load user requeset")
	user_ID_Object, successful := c.Get("user_id")
	if !successful {
		logger.Error("User id not found")
		c.JSON(400, gin.H{"error": "User id not found"})
		return
	}
	user_ID, ok := user_ID_Object.(int64)
	if !ok {
		logger.Error("User id is not integer")
		c.JSON(400, gin.H{"error": "User id is not integer"})
		return
	}
	workoutsInfo, err := h.storage.GetAllWorkouts(user_ID)
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	type WorkoutTimings struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	}
	type Day struct {
		Date            string           `json:"date"`
		WorkoutsTimings []WorkoutTimings `json:"workoutsTimings"`
	}
	type Request struct {
		Days []Day `json:"days"`
	}
	var request Request
	workoutsGrouped := make(map[string][]WorkoutTimings)
	for _, workout := range workoutsInfo {
		workoutsGrouped[workout.Date] = append(workoutsGrouped[workout.Date], WorkoutTimings{
			StartTime: workout.StartTime,
			EndTime:   workout.EndTime,
		})
	}

	for date, workoutList := range workoutsGrouped {
		request.Days = append(request.Days, Day{
			Date:            date,
			WorkoutsTimings: workoutList,
		})
	}

	c.JSON(200, request)
}

func (h *Handler) LoadTrainings(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling load trainings requeset")
	user_ID_Object, successful := c.Get("user_id")
	if !successful {
		logger.Error("User id not found")
		c.JSON(400, gin.H{"error": "User id not found"})
		return
	}
	user_ID, ok := user_ID_Object.(int64)
	if !ok {
		logger.Error("User id is not integer")
		c.JSON(400, gin.H{"error": "User id is not integer"})
		return
	}
	date := c.Param("date")
	if date == "" {
		logger.Error("date is required")
		c.JSON(400, gin.H{"error": "date is required"})
		return

	}
	type WorkoutExercises struct {
		WorkoutExerciseID int64  `json:"workoutExerciseId"`
		ExerciseName      string `json:"exerciseName"`
	}
	type TrainingsInDay struct {
		Date             string             `json:"date"`
		StartTime        string             `json:"timeStart"`
		EndTime          string             `json:"timeEnd"`
		WorkoutExercises []WorkoutExercises `json:"workoutExercises"`
	}
	type Request struct {
		TrainingsInDay []TrainingsInDay `json:"trainingsInDay"`
	}
	workoutsInfo, err := h.storage.GetWorkoutsFromDate(user_ID, date)
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	var request Request
	for _, workout := range workoutsInfo {
		training := TrainingsInDay{
			Date:      workout.Date,
			StartTime: workout.StartTime,
			EndTime:   workout.EndTime,
		}
		exercises, err := h.storage.GetWorkoutExercises(workout.WorkoutID)
		if err != nil {
			logger.Error("Internal server error while accessing the DB", "error", err)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		for _, exercise := range exercises {
			allowedExercise, err := h.storage.GetAllowedExercise(exercise.ExerciseID)
			if err != nil {
				logger.Error("Internal server error while accessing the DB", "error", err)
				c.JSON(500, gin.H{"error": "Internal server error"})
				return
			}
			training.WorkoutExercises = append(training.WorkoutExercises, WorkoutExercises{
				WorkoutExerciseID: exercise.WorkoutExerciseID,
				ExerciseName:      allowedExercise.Name,
			})
		}
		request.TrainingsInDay = append(request.TrainingsInDay, training)
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
	user_ID, ok := user_ID_Object.(int64)
	if !ok {
		logger.Error("User id is not integer")
		c.JSON(400, gin.H{"error": "User id is not integer"})
		return
	}
	type Workout struct {
		Date  string `json:"date"`
		Note  string `json:"note"`
		Photo string `json:"photo"`
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
		Photo:  workout_temp.Photo,
	}
	if err := h.storage.AddWorkout(workout); err != nil {
		logger.Error("failed to create workout", "error", err)
		c.JSON(500, gin.H{"error": "failed to create workout"})
		return
	}
	c.JSON(201, workout)
}
