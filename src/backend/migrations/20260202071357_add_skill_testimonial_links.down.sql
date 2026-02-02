-- Remove indexes
DROP INDEX IF EXISTS idx_skill_validations_testimonial_id;
DROP INDEX IF EXISTS idx_profile_skills_source_reference_letter_id;

-- Remove columns
ALTER TABLE skill_validations DROP COLUMN IF EXISTS testimonial_id;
ALTER TABLE profile_skills DROP COLUMN IF EXISTS source_reference_letter_id;
