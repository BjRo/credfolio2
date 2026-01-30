-- Add profile photo file reference to profiles table
-- Allows users to upload and display a profile photo

ALTER TABLE profiles
    ADD COLUMN profile_photo_file_id UUID REFERENCES files(id) ON DELETE SET NULL;
