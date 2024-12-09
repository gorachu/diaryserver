package sqlite

import (
	"fmt"
	"time"
)

func (s *Storage) AddBlacklistedToken(token string, expirationTime time.Time) error {
	const op = "storage.sqlite.AddBlacklistedToken"

	query := `
		INSERT INTO blacklisted_tokens (token, expiration_time)
		VALUES (?, ?)`
	_, err := s.db.Exec(query, token, expirationTime)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) IsTokenBlacklisted(token string) (bool, error) {
	const op = "storage.sqlite.IsTokenBlacklisted"

	query := `
		SELECT COUNT(*) 
		FROM blacklisted_tokens 
		WHERE token = ? AND expiration_time > datetime('now')`

	var count int
	err := s.db.QueryRow(query, token).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return count > 0, nil
}

func (s *Storage) RemoveExpiredTokens() error {
	const op = "storage.sqlite.RemoveExpiredTokens"

	query := `
		DELETE FROM blacklisted_tokens 
		WHERE expiration_time <= datetime('now')`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
