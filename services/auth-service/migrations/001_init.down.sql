-- Drop trigger first
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (refresh_tokens first due to foreign key)
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
