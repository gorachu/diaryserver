package handlers

import (
	"diaryserver/internal/storage/sqlite"
	"log/slog"
	"strconv"

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
		WorkoutId        int64              `json:"workoutId"`
		Date             string             `json:"date"`
		StartTime        string             `json:"timeStart"`
		EndTime          string             `json:"timeEnd"`
		Notes            string             `json:"notes"`
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
			WorkoutId: workout.WorkoutID,
			Date:      workout.Date,
			StartTime: workout.StartTime,
			EndTime:   workout.EndTime,
			Notes:     workout.Notes,
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
	var setInfo Workout
	if err := c.ShouldBindJSON(&setInfo); err != nil {
		logger.Error("invalid request body", "error", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	var workout = sqlite.Workout{
		UserID: user_ID,
		Date:   setInfo.Date,
		Notes:  setInfo.Note,
		Photo:  setInfo.Photo,
	}
	if err := h.storage.AddWorkout(workout); err != nil {
		logger.Error("failed to create workout", "error", err)
		c.JSON(500, gin.H{"error": "failed to create workout"})
		return
	}
	c.JSON(201, workout)
}
func (h *Handler) LoadTrainingSingle(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling load single training requeset")
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
	workoutIdStr := c.Param("workoutId")
	if workoutIdStr == "" {
		logger.Error("workoutId is required")
		c.JSON(400, gin.H{"error": "WorkoutId is required"})
		return
	}
	workoutId, err := strconv.ParseInt(workoutIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid workoutId, must be an integer"})
		return
	}
	type WorkoutExercises struct {
		WorkoutExerciseID int64            `json:"workoutExerciseId"`
		ExerciseName      string           `json:"exerciseName"`
		Sets              []sqlite.SetInfo `json:"sets"`
	}
	type Training struct {
		StartTime        string                       `json:"timeStart"`
		EndTime          string                       `json:"timeEnd"`
		Notes            string                       `json:"notes"`
		WorkoutExercises []WorkoutExercises           `json:"workoutExercises"`
		ListOfExercises  []sqlite.AllowedExerciseInfo `json:"listOfExercises"`
	}
	workoutInfo, err := h.storage.GetWorkoutFromID(workoutId)
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if workoutInfo.UserID != user_ID {
		logger.Error("Access Denied")
		c.JSON(403, gin.H{"error": "Access Denied"})
		return
	}
	listOfExercises, err := h.storage.GetAllowedExercises()
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	training := Training{
		StartTime:       workoutInfo.StartTime,
		EndTime:         workoutInfo.EndTime,
		Notes:           workoutInfo.Notes,
		ListOfExercises: listOfExercises,
	}
	exercises, err := h.storage.GetWorkoutExercises(workoutInfo.WorkoutID)
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

		sets, err := h.storage.GetSets(exercise.WorkoutExerciseID)
		if err != nil {
			logger.Error("Internal server error while accessing the DB", "error", err)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		training.WorkoutExercises = append(training.WorkoutExercises, WorkoutExercises{
			WorkoutExerciseID: exercise.WorkoutExerciseID,
			ExerciseName:      allowedExercise.Name,
			Sets:              sets,
		})
	}
	c.JSON(200, training)
}
func (h *Handler) CreateSets(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling creating set requeset")
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
	workoutIdStr := c.Param("workoutId")
	if workoutIdStr == "" {
		logger.Error("workoutId is required")
		c.JSON(400, gin.H{"error": "WorkoutId is required"})
		return
	}
	workoutId, err := strconv.ParseInt(workoutIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid workoutId, must be an integer"})
		return
	}
	workoutInfo, err := h.storage.GetWorkoutFromID(workoutId)
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if workoutInfo.UserID != user_ID {
		logger.Error("Access Denied")
		c.JSON(403, gin.H{"error": "Access Denied"})
		return
	}
	type SingleSet struct {
		Weight      float64 `json:"weight"`
		Repetitions int     `json:"repetitions"`
	}
	type Sets struct {
		AllowedExercise sqlite.AllowedExerciseInfo `json:"allowedExercise"`
		Sets            []SingleSet                `json:"sets"`
	}
	var setsInfo Sets
	if err := c.ShouldBindJSON(&setsInfo); err != nil {
		logger.Error("invalid request body", "error", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	workoutExerciseId, err := h.storage.AddWorkoutExercise(sqlite.WorkoutExercise{
		WorkoutID:  workoutId,
		ExerciseID: setsInfo.AllowedExercise.AllowedExerciseId,
	})
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	var sets []sqlite.Set
	for _, set := range setsInfo.Sets {
		sets = append(sets, sqlite.Set{WorkoutExerciseID: workoutExerciseId, Repetitions: set.Repetitions, Weight: set.Weight})
	}
	if err := h.storage.AddSets(sets); err != nil {
		logger.Error("failed to create sets", "error", err)
		c.JSON(500, gin.H{"error": "failed to create set"})
		return
	}
	c.JSON(201, sets)
}
func (h *Handler) DeleteExercise(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling DeleteExercise")
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
	workoutIdStr := c.Param("workoutId")
	if workoutIdStr == "" {
		logger.Error("workoutId is required")
		c.JSON(400, gin.H{"error": "WorkoutId is required"})
		return
	}
	workoutId, err := strconv.ParseInt(workoutIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid workoutId, must be an integer"})
		return
	}
	workoutInfo, err := h.storage.GetWorkoutFromID(workoutId)
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if workoutInfo.UserID != user_ID {
		logger.Error("Access Denied")
		c.JSON(403, gin.H{"error": "Access Denied"})
		return
	}
	type Exercise struct {
		Sets              []sqlite.SetInfo `json:"sets"`
		WorkoutExerciseId int64            `json:"workoutExerciseId"`
	}
	var exerciseInfo Exercise
	if err := c.ShouldBindJSON(&exerciseInfo); err != nil {
		logger.Error("invalid request body", "error", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.storage.DeleteWorkoutExercise(exerciseInfo.WorkoutExerciseId); err != nil {
		logger.Error("Internal server error", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, exerciseInfo)
}
func (h *Handler) ChangeExercise(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling ChangeExercise")
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
	workoutIdStr := c.Param("workoutId")
	if workoutIdStr == "" {
		logger.Error("workoutId is required")
		c.JSON(400, gin.H{"error": "WorkoutId is required"})
		return
	}
	workoutId, err := strconv.ParseInt(workoutIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid workoutId, must be an integer"})
		return
	}
	workoutInfo, err := h.storage.GetWorkoutFromID(workoutId)
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if workoutInfo.UserID != user_ID {
		logger.Error("Access Denied")
		c.JSON(403, gin.H{"error": "Access Denied"})
		return
	}
	exerciseIdStr := c.Param("exerciseId")
	if workoutIdStr == "" {
		logger.Error("workoutId is required")
		c.JSON(400, gin.H{"error": "WorkoutId is required"})
		return
	}
	exerciseId, err := strconv.ParseInt(exerciseIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid workoutId, must be an integer"})
		return
	}
	type Exercise struct {
		Sets []sqlite.SetInfo `json:"sets"`
	}
	var exerciseInfo Exercise
	if err := c.ShouldBindJSON(&exerciseInfo); err != nil {
		logger.Error("invalid request body", "error", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.storage.ReplaceSets(exerciseId, exerciseInfo.Sets); err != nil {
		logger.Error("Internal server error", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	type RequestInfo struct {
		ExeExercise Exercise `json:"ExeExercise"`
		ExerciseId  int64    `json:"exerciseId"`
	}
	request := RequestInfo{
		ExeExercise: exerciseInfo,
		ExerciseId:  exerciseId,
	}
	c.JSON(200, request)
}
func (h *Handler) ChangeWorkoutInfo(c *gin.Context) {
	logger := c.MustGet("logger").(*slog.Logger)
	logger.Debug("handling ChangeWorkoutInfo")
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
	workoutIdStr := c.Param("workoutId")
	if workoutIdStr == "" {
		logger.Error("workoutId is required")
		c.JSON(400, gin.H{"error": "WorkoutId is required"})
		return
	}
	workoutId, err := strconv.ParseInt(workoutIdStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid workoutId, must be an integer"})
		return
	}
	workoutInfo, err := h.storage.GetWorkoutFromID(workoutId)
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if workoutInfo.UserID != user_ID {
		logger.Error("Access Denied")
		c.JSON(403, gin.H{"error": "Access Denied"})
		return
	}
	var updateData sqlite.WorkoutInfo
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Error("Invalid request body", "error", err)
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	updateData.WorkoutID = workoutId
	updateData.UserID = user_ID
	updatedFields := getChangedFields(*workoutInfo, updateData)
	if err := h.storage.PartialUpdateWorkout(workoutId, updatedFields); err != nil {
		logger.Error("Failed to update workout", "error", err)
		c.JSON(500, gin.H{"error": "Failed to update workout"})
		return
	}
	workoutInfo, err = h.storage.GetWorkoutFromID(workoutId)
	if err != nil {
		logger.Error("Internal server error while accessing the DB", "error", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, workoutInfo)
}
func getChangedFields(original, updated sqlite.WorkoutInfo) map[string]interface{} {
	changes := make(map[string]interface{})
	if updated.Date != "" && updated.Date != original.Date {
		changes["date"] = updated.Date
	}
	if updated.StartTime != "" && updated.StartTime != original.StartTime {
		changes["time_start"] = updated.StartTime
	}
	if updated.EndTime != "" && updated.EndTime != original.EndTime {
		changes["time_end"] = updated.EndTime
	}
	if updated.Notes != original.Notes {
		changes["notes"] = updated.Notes
	}
	if updated.Photo != original.Photo {
		changes["photo"] = updated.Photo
	}
	return changes
}
