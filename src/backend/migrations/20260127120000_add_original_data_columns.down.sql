DROP INDEX IF EXISTS idx_profile_education_source_resume_id;
DROP INDEX IF EXISTS idx_profile_experiences_source_resume_id;
ALTER TABLE profile_education DROP COLUMN IF EXISTS original_data;
ALTER TABLE profile_experiences DROP COLUMN IF EXISTS original_data;
