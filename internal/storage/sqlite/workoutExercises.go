package sqlite

import (
	"database/sql"
	"fmt"
)

type WorkoutExercise struct {
	WorkoutID  int64
	ExerciseID int64
}

type WorkoutExerciseInfo struct {
	WorkoutExerciseID int64
	WorkoutID         int64
	ExerciseID        int64
}

func (s *Storage) AddWorkoutExercise(workoutExercise WorkoutExercise) (int64, error) {
	const op = "storage.sqlite.AddWorkoutExercise"
	if workoutExercise.WorkoutID == 0 || workoutExercise.ExerciseID == 0 {
		return 0, fmt.Errorf("%s: workout ID or exercise ID is required", op)
	}
	query := `INSERT INTO workout_exercises (workout_id, exercise_id) VALUES (?, ?)`

	result, err := s.db.Exec(query, workoutExercise.WorkoutID, workoutExercise.ExerciseID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	workoutExerciseID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert ID: %w", op, err)
	}

	return workoutExerciseID, nil
}

func (s *Storage) AddWorkoutExercises(workoutExercises []WorkoutExercise) error {
	const op = "storage.sqlite.AddWorkoutExercises"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `INSERT INTO workout_exercises (workout_id, exercise_id) VALUES (?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, workoutExercise := range workoutExercises {
		if workoutExercise.WorkoutID == 0 || workoutExercise.ExerciseID == 0 {
			return fmt.Errorf("%s: workout ID or exercise ID is required", op)
		}
		_, err := stmt.Exec(workoutExercise.WorkoutID, workoutExercise.ExerciseID)
		if err != nil {
			return fmt.Errorf("%s: failed to add workout exercise: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteWorkoutExercise(workoutExerciseID int64) error {
	const op = "storage.sqlite.DeleteWorkoutExercise"
	query := `DELETE FROM workout_exercises WHERE workout_exercise_id = ?`

	_, err := s.db.Exec(query, workoutExerciseID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteWorkoutExercises(workoutExerciseIDs []int64) error {
	const op = "storage.sqlite.DeleteWorkoutExercises"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `DELETE FROM workout_exercises WHERE workout_exercise_id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, id := range workoutExerciseIDs {
		_, err := stmt.Exec(id)
		if err != nil {
			return fmt.Errorf("%s: failed to delete workout exercise with ID %d: %w", op, id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) GetWorkoutExercise(workoutExerciseID int64) (*WorkoutExerciseInfo, error) {
	const op = "storage.sqlite.GetWorkoutExercise"
	query := `SELECT workout_exercise_id, workout_id, exercise_id 
			 FROM workout_exercises WHERE workout_exercise_id = ?`

	we := &WorkoutExerciseInfo{}
	err := s.db.QueryRow(query, workoutExerciseID).Scan(
		&we.WorkoutExerciseID,
		&we.WorkoutID,
		&we.ExerciseID,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: workout exercise not found", op)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return we, nil
}

func (s *Storage) GetWorkoutExercises(workoutID int64) ([]WorkoutExerciseInfo, error) {
	const op = "storage.sqlite.GetWorkoutExercises"
	query := `SELECT workout_exercise_id, workout_id, exercise_id 
			 FROM workout_exercises WHERE workout_id = ?
			 ORDER BY workout_exercise_id`

	rows, err := s.db.Query(query, workoutID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var workoutExercises []WorkoutExerciseInfo
	for rows.Next() {
		var we WorkoutExerciseInfo
		err := rows.Scan(
			&we.WorkoutExerciseID,
			&we.WorkoutID,
			&we.ExerciseID,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		workoutExercises = append(workoutExercises, we)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return workoutExercises, nil
}
