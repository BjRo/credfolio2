-- Rollback initial schema
-- Drop in reverse order of creation due to foreign key dependencies

-- Drop triggers first
DROP TRIGGER IF EXISTS update_reference_letters_updated_at ON reference_letters;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (respecting foreign key dependencies)
DROP TABLE IF EXISTS reference_letters;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS users;

-- Note: We don't drop the uuid-ossp extension as other databases/schemas might use it
