package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	dir := filepath.Dir(storagePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("%s: failed to create directory: %w", op, err)
	}
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) DropAllTables() error {
	const op = "storage.sqlite.DropAllTables"
	queries := []string{
		`DROP TABLE IF EXISTS sets;`,
		`DROP TABLE IF EXISTS workout_exercises;`,
		`DROP TABLE IF EXISTS workouts;`,
		`DROP TABLE IF EXISTS allowed_exercises;`,
		`DROP TABLE IF EXISTS users;`,
		`DROP TABLE IF EXISTS blacklisted_tokens;`,

		`DROP INDEX IF EXISTS idx_users_username;`,
		`DROP INDEX IF EXISTS idx_users_email;`,
		`DROP INDEX IF EXISTS idx_workouts_user_id;`,
		`DROP INDEX IF EXISTS idx_workout_exercises_workout_id;`,
		`DROP INDEX IF EXISTS idx_sets_workout_exercise_id;`,
		`DROP INDEX IF EXISTS idx_blacklisted_tokens_token;`,
		`DROP INDEX IF EXISTS idx_blacklisted_tokens_expiration;`,
	}
	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
