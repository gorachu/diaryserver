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
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE TABLE IF NOT EXISTS allowed_exercises (
			exercise_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS workouts (
			workout_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATETIME DEFAULT CURRENT_TIMESTAMP,
			notes TEXT,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS workout_exercises (
			workout_exercise_id INTEGER PRIMARY KEY AUTOINCREMENT,
			workout_id INTEGER NOT NULL,
			exercise_id INTEGER NOT NULL,
			FOREIGN KEY (workout_id) REFERENCES workouts(workout_id) ON DELETE CASCADE,
			FOREIGN KEY (exercise_id) REFERENCES allowed_exercises(exercise_id)
		);`,

		`CREATE TABLE IF NOT EXISTS sets (
			set_id INTEGER PRIMARY KEY AUTOINCREMENT,
			workout_exercise_id INTEGER NOT NULL,
			repetitions INTEGER NOT NULL,
			weight REAL NOT NULL,
			FOREIGN KEY (workout_exercise_id) REFERENCES workout_exercises(workout_exercise_id) ON DELETE CASCADE
		);`,

		`CREATE TABLE IF NOT EXISTS blacklisted_tokens (
			token_id INTEGER PRIMARY KEY AUTOINCREMENT,
			token TEXT NOT NULL UNIQUE,
			expiration_time DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,

		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`,
		`CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_workout_exercises_workout_id ON workout_exercises(workout_id);`,
		`CREATE INDEX IF NOT EXISTS idx_sets_workout_exercise_id ON sets(workout_exercise_id);`,
		`CREATE INDEX IF NOT EXISTS idx_blacklisted_tokens_token ON blacklisted_tokens(token);`,
		`CREATE INDEX IF NOT EXISTS idx_blacklisted_tokens_expiration ON blacklisted_tokens(expiration_time);`,
	}
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
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
