-- Drop triggers
DROP TRIGGER IF EXISTS update_user_preferences_updated_at ON user_preferences;
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS otps;
DROP TABLE IF EXISTS email_verifications;
DROP TABLE IF EXISTS password_resets;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

-- Drop ENUM types
DROP TYPE IF EXISTS auth_provider;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;
