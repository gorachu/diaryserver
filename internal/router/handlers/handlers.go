package handlers

import (
	"diaryserver/internal/storage/sqlite"
	"log/slog"
)

type Handler struct {
	storage *sqlite.Storage
	log     *slog.Logger
}

func NewHandlers(storage *sqlite.Storage, log *slog.Logger) *Handler {
	return &Handler{storage: storage, log: log}
}
