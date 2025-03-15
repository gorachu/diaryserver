DROP TABLE IF EXISTS sets;
DROP TABLE IF EXISTS workout_exercises;
DROP TABLE IF EXISTS workouts;
DROP TABLE IF EXISTS allowed_exercises;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS blacklisted_tokens;

DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_workouts_user_id;
DROP INDEX IF EXISTS idx_workout_exercises_workout_id;
DROP INDEX IF EXISTS idx_sets_workout_exercise_id;
DROP INDEX IF EXISTS idx_blacklisted_tokens_token;
DROP INDEX IF EXISTS idx_blacklisted_tokens_expiration;