-- Add header fields to profiles table for user-editable profile information
-- These fields allow users to override/edit the data extracted from their resume

ALTER TABLE profiles
    ADD COLUMN name VARCHAR(255),
    ADD COLUMN email VARCHAR(255),
    ADD COLUMN phone VARCHAR(50),
    ADD COLUMN location VARCHAR(255),
    ADD COLUMN summary TEXT;

-- Add comment for documentation
COMMENT ON COLUMN profiles.name IS 'User-edited name (overrides resume extraction if set)';
COMMENT ON COLUMN profiles.email IS 'User-edited email (overrides resume extraction if set)';
COMMENT ON COLUMN profiles.phone IS 'User-edited phone (overrides resume extraction if set)';
COMMENT ON COLUMN profiles.location IS 'User-edited location (overrides resume extraction if set)';
COMMENT ON COLUMN profiles.summary IS 'User-edited professional summary (overrides resume extraction if set)';
