package sqlite

import (
	"database/sql"
	"fmt"
)

type Set struct {
	WorkoutExerciseID int64
	Repetitions       int
	Weight            float64
}
type SetInfo struct {
	SetID             int64   `json:"setId"`
	WorkoutExerciseID int64   `json:"workoutExerciseId"`
	Repetitions       int     `json:"repetitions"`
	Weight            float64 `json:"weight"`
}

func (s *Storage) AddSet(set Set) error {
	const op = "storage.sqlite.AddSet"

	if set.WorkoutExerciseID == 0 {
		return fmt.Errorf("%s: workout exercise ID is required", op)
	}

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
		if set.WorkoutExerciseID == 0 {
			return fmt.Errorf("%s: workout exercise ID is required", op)
		}
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
func (s *Storage) DeleteSet(setID int64) error {
	const op = "storage.sqlite.DeleteSet"
	query := `DELETE FROM sets WHERE set_id = ?`

	_, err := s.db.Exec(query, setID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteSets(setIDs []int64) error {
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
func (s *Storage) GetSet(setID int64) (*SetInfo, error) {
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

func (s *Storage) GetSets(workoutExerciseID int64) ([]SetInfo, error) {
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
func (s *Storage) ReplaceSets(workoutExerciseID int64, sets []SetInfo) error {
	const op = "storage.sqlite.ReplaceSets"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	queryDelete := `DELETE FROM sets WHERE workout_exercise_id = ? AND set_id IS NOT NULL`
	_, err = tx.Exec(queryDelete, workoutExerciseID)
	if err != nil {
		return fmt.Errorf("%s: failed to delete existing sets: %w", op, err)
	}

	queryInsert := `INSERT INTO sets (workout_exercise_id, repetitions, weight) VALUES (?, ?, ?)`
	stmt, err := tx.Prepare(queryInsert)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare insert statement: %w", op, err)
	}
	defer stmt.Close()

	for _, set := range sets {
		_, err := stmt.Exec(workoutExerciseID, set.Repetitions, set.Weight)
		if err != nil {
			return fmt.Errorf("%s: failed to insert set: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
