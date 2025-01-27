package sqlite

import (
	"database/sql"
	"fmt"
)

type Workout struct {
	UserID int
	Date   string
	Notes  string
}
type WorkoutInfo struct {
	WorkoutID int
	UserID    int
	Date      string
	Notes     string
}

func (s *Storage) AddWorkout(workout Workout) error {
	const op = "storage.sqlite.AddWorkout"
	if workout.UserID == 0 {
		return fmt.Errorf("%s: user ID is required", op)
	}
	query := `INSERT INTO workouts (user_id, date, notes) VALUES (?, ?, ?)`

	_, err := s.db.Exec(query, workout.UserID, workout.Date, workout.Notes)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddWorkouts(workouts []Workout) error {
	const op = "storage.sqlite.AddWorkouts"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `INSERT INTO workouts (user_id, date, notes) VALUES (?, ?, ?)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, workout := range workouts {
		if workout.UserID == 0 {
			return fmt.Errorf("%s: user ID is required", op)
		}
		_, err := stmt.Exec(workout.UserID, workout.Date, workout.Notes)
		if err != nil {
			return fmt.Errorf("%s: failed to add workout for user %d: %w", op, workout.UserID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
func (s *Storage) DeleteWorkout(workoutID int) error {
	const op = "storage.sqlite.DeleteWorkout"
	query := `DELETE FROM workouts WHERE workout_id = ?`

	_, err := s.db.Exec(query, workoutID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteWorkouts(workoutIDs []int) error {
	const op = "storage.sqlite.DeleteWorkouts"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `DELETE FROM workouts WHERE workout_id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}
	defer stmt.Close()

	for _, id := range workoutIDs {
		_, err := stmt.Exec(id)
		if err != nil {
			return fmt.Errorf("%s: failed to delete workout with ID %d: %w", op, id, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}
func (s *Storage) GetWorkout(user_id int, date string) (*WorkoutInfo, error) {
	const op = "storage.sqlite.GetWorkout"
	query := `SELECT workout_id, user_id, date, notes 
			 FROM workouts WHERE user_id = ? AND date = ?`

	workout := &WorkoutInfo{}
	err := s.db.QueryRow(query, user_id, date).Scan(
		&workout.WorkoutID,
		&workout.UserID,
		&workout.Date,
		&workout.Notes,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: workout not found", op)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return workout, nil
}
func (s *Storage) GetWorkoutFromID(workout_ID int) (*WorkoutInfo, error) {
	const op = "storage.sqlite.GetWorkout"
	query := `SELECT workout_id, user_id, date, notes 
			 FROM workouts WHERE workout_id = ?`

	workout := &WorkoutInfo{}
	err := s.db.QueryRow(query, workout_ID).Scan(
		&workout.WorkoutID,
		&workout.UserID,
		&workout.Date,
		&workout.Notes,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: workout not found", op)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return workout, nil
}

func (s *Storage) GetWorkouts(userID int) ([]WorkoutInfo, error) {
	const op = "storage.sqlite.GetWorkouts"
	query := `SELECT workout_id, user_id, date, notes 
			 FROM workouts WHERE user_id = ?
			 ORDER BY workout_id`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var workouts []WorkoutInfo
	for rows.Next() {
		var workout WorkoutInfo
		err := rows.Scan(
			&workout.WorkoutID,
			&workout.UserID,
			&workout.Date,
			&workout.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		workouts = append(workouts, workout)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return workouts, nil
}
