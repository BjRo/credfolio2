-- Add original_data column to store the raw extracted JSON for each item
ALTER TABLE profile_experiences ADD COLUMN original_data JSONB;
ALTER TABLE profile_education ADD COLUMN original_data JSONB;

-- Add indexes on source_resume_id for efficient lookups during re-processing
CREATE INDEX idx_profile_experiences_source_resume_id ON profile_experiences(source_resume_id) WHERE source_resume_id IS NOT NULL;
CREATE INDEX idx_profile_education_source_resume_id ON profile_education(source_resume_id) WHERE source_resume_id IS NOT NULL;
