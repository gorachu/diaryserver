package sqlite

import (
	"fmt"
)

type User struct {
	Username     string
	Email        string
	PasswordHash string
}

type UserInfo struct {
	UserID       int64
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    string
}

func (s *Storage) AddUser(user User) error {
	const op = "storage.sqlite.AddUser"
	if user.Username == "" || user.Email == "" || user.PasswordHash == "" {
		return fmt.Errorf("%s: username, email and password_hash are required", op)
	}
	query := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`

	_, err := s.db.Exec(query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddUsers(users []User) error {
	const op = "storage.sqlite.AddUsers"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, user := range users {
		if user.Username == "" || user.Email == "" || user.PasswordHash == "" {
			return fmt.Errorf("%s: username, email and password_hash are required", op)
		}
		_, err := stmt.Exec(user.Username, user.Email, user.PasswordHash)
		if err != nil {
			return fmt.Errorf("%s: failed to add user %s: %w", op, user.Username, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUser(username string) error {
	const op = "storage.sqlite.DeleteUser"
	query := `DELETE FROM users WHERE username = ?`
	if _, err := s.db.Exec(query, username); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteUsers(usernames []string) error {
	const op = "storage.sqlite.DeleteUsers"
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()
	query := `DELETE FROM users WHERE username = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()
	for _, username := range usernames {
		_, err := stmt.Exec(username)
		if err != nil {
			return fmt.Errorf("%s: failed to delete user %s: %w", op, username, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}
	return nil
}
func (s *Storage) DeleteAllUsers() error {
	const op = "storage.sqlite.DeleteAllUsers"
	query := `DELETE FROM users`
	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetUser(username string) (*UserInfo, error) {
	const op = "storage.sqlite.GetUser"

	query := `SELECT user_id, username, email, password_hash FROM users WHERE username = ?`

	row := s.db.QueryRow(query, username)

	var user UserInfo
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) GetUsers() ([]UserInfo, error) {
	const op = "storage.sqlite.GetUsers"

	query := `SELECT user_id, username, email, password_hash, created_at FROM users`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users []UserInfo
	for rows.Next() {
		var user UserInfo
		err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}
