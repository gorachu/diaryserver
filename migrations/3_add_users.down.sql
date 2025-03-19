DELETE FROM users WHERE user_id IN (
    SELECT user_id FROM users WHERE username LIKE 'user%' ORDER BY user_id LIMIT 10
);