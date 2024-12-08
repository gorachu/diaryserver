package main

import (
	"diaryserver/internal/config"
	"diaryserver/internal/router"
	"diaryserver/internal/storage/sqlite"
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
