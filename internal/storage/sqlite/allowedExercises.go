package sqlite

import "fmt"

type AllowedExercise struct {
	Name        string
	Description string
}
type AllowedExerciseInfo struct {
	AllowedExerciseId int64  `json:"allowedExerciseId"`
	Name              string `json:"name"`
	Description       string `json:"description"`
}

func (s *Storage) AddAllowedExercise(exercise AllowedExercise) error {
	const op = "storage.sqlite.AddAllowedExercise"
	if exercise.Name == "" {
		return fmt.Errorf("%s: name is required", op)
	}
	query := `INSERT INTO allowed_exercises (name, description) VALUES (?, ?)`

	_, err := s.db.Exec(query, exercise.Name, exercise.Description)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (s *Storage) AddAllowedExercises(exercises []AllowedExercise) error {
	const op = "storage.sqlite.AddAllowedExercises"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `INSERT INTO allowed_exercises (name, description) VALUES (?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, exercise := range exercises {
		if exercise.Name == "" {
			return fmt.Errorf("%s: name is required", op)
		}
		_, err := stmt.Exec(exercise.Name, exercise.Description)
		if err != nil {
			return fmt.Errorf("%s: failed to add exercise %s: %w", op, exercise.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
func (s *Storage) DeleteAllowedExercise(name string) error {
	const op = "storage.sqlite.DeleteAllowedExercise"
	query := `DELETE FROM allowed_exercises WHERE name = ?`

	_, err := s.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteAllowedExercises(names []string) error {
	const op = "storage.sqlite.DeleteAllowedExercises"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `DELETE FROM allowed_exercises WHERE name = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, name := range names {
		_, err := stmt.Exec(name)
		if err != nil {
			return fmt.Errorf("%s: failed to delete exercise %s: %w", op, name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAllowedExercise(id int64) (AllowedExerciseInfo, error) {
	const op = "storage.sqlite.GetAllowedExercise"

	query := `SELECT exercise_id, name, description FROM allowed_exercises WHERE exercise_id = ?`

	var exercise AllowedExerciseInfo
	err := s.db.QueryRow(query, id).Scan(&exercise.AllowedExerciseId, &exercise.Name, &exercise.Description)
	if err != nil {
		return AllowedExerciseInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	return exercise, nil
}

func (s *Storage) GetAllowedExercises() ([]AllowedExerciseInfo, error) {
	const op = "storage.sqlite.GetAllowedExercises"

	query := `SELECT exercise_id, name, description FROM allowed_exercises`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var exercises []AllowedExerciseInfo
	for rows.Next() {
		var exercise AllowedExerciseInfo
		if err := rows.Scan(&exercise.AllowedExerciseId, &exercise.Name, &exercise.Description); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return exercises, nil
}
