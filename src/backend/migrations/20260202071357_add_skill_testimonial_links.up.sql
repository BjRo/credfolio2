-- Add testimonial_id to skill_validations for granular attribution
-- This enables linking a skill validation to the specific testimonial that validates it
ALTER TABLE skill_validations ADD COLUMN IF NOT EXISTS testimonial_id UUID REFERENCES testimonials(id) ON DELETE SET NULL;

-- Add source_reference_letter_id to profile_skills for tracking discovered skills
-- This enables knowing which reference letter a skill was discovered from
ALTER TABLE profile_skills ADD COLUMN IF NOT EXISTS source_reference_letter_id UUID REFERENCES reference_letters(id) ON DELETE SET NULL;

-- Create indexes for efficient queries
CREATE INDEX IF NOT EXISTS idx_skill_validations_testimonial_id ON skill_validations(testimonial_id);
CREATE INDEX IF NOT EXISTS idx_profile_skills_source_reference_letter_id ON profile_skills(source_reference_letter_id);
