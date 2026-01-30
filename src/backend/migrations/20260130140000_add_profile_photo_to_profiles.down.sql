-- Remove profile photo file reference from profiles table

ALTER TABLE profiles
    DROP COLUMN IF EXISTS profile_photo_file_id;
