-- Remove header fields from profiles table

ALTER TABLE profiles
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS location,
    DROP COLUMN IF EXISTS summary;
