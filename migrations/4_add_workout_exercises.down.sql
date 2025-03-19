DELETE FROM workout_exercises 
WHERE (workout_id, exercise_id) IN (
    (2, 4),
    (1, 2),
    (2, 1),
    (1, 5),
    (1, 1),
    (2, 2),
    (2, 5)
);                                 