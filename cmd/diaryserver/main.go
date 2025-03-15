package main

import (
	"diaryserver/internal/config"
	"diaryserver/internal/router"
	"diaryserver/internal/service"
	"diaryserver/internal/storage/sqlite"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := SetupLogger(cfg.Env)
	log.Info("starting diary server", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage := InitStorage(cfg, log)
	_ = storage.AddAllowedExercises([]sqlite.AllowedExercise{
		{
			Name:        "Жим лежа",
			Description: "Нужно лечь и жать",
		},
		{
			Name:        "Присед со штангой",
			Description: "Нужно приседать вместе со штангой",
		},
		{
			Name:        "Становая тяга",
			Description: "Поднять штангу с пола",
		},
	})
	var passwords []string
	n := 10
	for i := 1; i <= n; i++ {
		passwords = append(passwords, fmt.Sprintf("user%d@gmail.com", i))
	}
	var hashedPasswords []string
	for _, password := range passwords {
		hashedPassword, _ := service.HashPassword(password)
		hashedPasswords = append(hashedPasswords, hashedPassword)
	}
	var users []sqlite.User
	for i, password := range passwords {
		users = append(users, sqlite.User{
			Username:     fmt.Sprintf("user%d", i+1),
			Email:        password,
			PasswordHash: hashedPasswords[i],
		})
	}
	_ = storage.AddUsers(users)
	_ = storage.AddWorkoutExercises([]sqlite.WorkoutExercise{
		{
			WorkoutID:  1,
			ExerciseID: 2,
		},
		{
			WorkoutID:  1,
			ExerciseID: 3,
		},
		{
			WorkoutID:  2,
			ExerciseID: 1,
		},
	})
	error := storage.AddWorkouts([]sqlite.Workout{
		{
			UserID:    1,
			Date:      "2025-01-01",
			StartTime: "08:00:04",
			EndTime:   "09:00:24",
			Notes:     "Тренировка на силу",
			Photo:     "",
		},
		{
			UserID:    1,
			Date:      "2025-01-01",
			StartTime: "10:00",
			EndTime:   "11:00",
			Notes:     "Кардио тренировка",
			Photo:     "",
		},
	})
	if error != nil {
		log.Debug("Error from inserting", "error", error)
		return
	}
	// TODO: init router
	router := router.SetupRouter(storage, log, cfg)
	// TODO: run server
	log.Info("stating server", slog.String("address", cfg.HTTPServer.Address))
	if err := router.Run(cfg.HTTPServer.Address); err != nil {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
func InitStorage(cfg *config.Config, log *slog.Logger) *sqlite.Storage {
	absolutePath, err := filepath.Abs(cfg.StoragePath)
	if err != nil {
		log.Debug("error in getting absolute path", "error", err)
	} else {
		log.Debug(cfg.StoragePath, "absolute path", absolutePath)
	}
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", "error", err)
		os.Exit(1)
	}
	return storage
}
