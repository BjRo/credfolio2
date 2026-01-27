-- Rollback: Drop profile_education table

DROP TRIGGER IF EXISTS update_profile_education_updated_at ON profile_education;
DROP INDEX IF EXISTS idx_profile_education_display_order;
DROP INDEX IF EXISTS idx_profile_education_profile_id;
DROP TABLE IF EXISTS profile_education;
