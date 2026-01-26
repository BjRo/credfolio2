-- Rollback: Drop profiles and profile_experiences tables

DROP TRIGGER IF EXISTS update_profile_experiences_updated_at ON profile_experiences;
DROP TRIGGER IF EXISTS update_profiles_updated_at ON profiles;
DROP INDEX IF EXISTS idx_profile_experiences_display_order;
DROP INDEX IF EXISTS idx_profile_experiences_profile_id;
DROP TABLE IF EXISTS profile_experiences;
DROP INDEX IF EXISTS idx_profiles_user_id;
DROP TABLE IF EXISTS profiles;
