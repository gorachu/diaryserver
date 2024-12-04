package sqlite

import (
	"database/sql"
	"fmt"
)

type Set struct {
	WorkoutExerciseID int
	Repetitions       int
	Weight            float64
}
type SetInfo struct {
	SetID             int
	WorkoutExerciseID int
	Repetitions       int
	Weight            float64
}

func (s *Storage) AddSet(set Set) error {
	const op = "storage.sqlite.AddSet"

	query := `INSERT INTO sets (workout_exercise_id, repetitions, weight) VALUES (?, ?, ?)`

	_, err := s.db.Exec(query, set.WorkoutExerciseID, set.Repetitions, set.Weight)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddSets(sets []Set) error {
	const op = "storage.sqlite.AddSets"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `INSERT INTO sets (workout_exercise_id, repetitions, weight) VALUES (?, ?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, set := range sets {
		_, err := stmt.Exec(set.WorkoutExerciseID, set.Repetitions, set.Weight)
		if err != nil {
			return fmt.Errorf("%s: failed to add set: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
func (s *Storage) DeleteSet(setID int) error {
	const op = "storage.sqlite.DeleteSet"
	query := `DELETE FROM sets WHERE set_id = ?`

	_, err := s.db.Exec(query, setID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteSets(setIDs []int) error {
	const op = "storage.sqlite.DeleteSets"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `DELETE FROM sets WHERE set_id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, setID := range setIDs {
		_, err := stmt.Exec(setID)
		if err != nil {
			return fmt.Errorf("%s: failed to delete set with ID %d: %w", op, setID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
func (s *Storage) GetSet(setID int) (*SetInfo, error) {
	const op = "storage.sqlite.GetSet"
	query := `SELECT set_id, workout_exercise_id, repetitions, weight 
			 FROM sets WHERE set_id = ?`

	set := &SetInfo{}
	err := s.db.QueryRow(query, setID).Scan(
		&set.SetID,
		&set.WorkoutExerciseID,
		&set.Repetitions,
		&set.Weight,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: set not found", op)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return set, nil
}

func (s *Storage) GetSets(workoutExerciseID int) ([]SetInfo, error) {
	const op = "storage.sqlite.GetSets"
	query := `SELECT set_id, workout_exercise_id, repetitions, weight 
			 FROM sets WHERE workout_exercise_id = ?
			 ORDER BY set_id`

	rows, err := s.db.Query(query, workoutExerciseID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var sets []SetInfo
	for rows.Next() {
		var set SetInfo
		err := rows.Scan(
			&set.SetID,
			&set.WorkoutExerciseID,
			&set.Repetitions,
			&set.Weight,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		sets = append(sets, set)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sets, nil
}
