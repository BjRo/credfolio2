-- Create credibility system tables for reference letter validation
-- Adds: error_message to reference_letters, testimonials, skill_validations, experience_validations

-- Add error_message column to reference_letters for storing failure reasons
ALTER TABLE reference_letters ADD COLUMN IF NOT EXISTS error_message TEXT;

-- Testimonials table: stores full quotes from reference letters for display
CREATE TABLE testimonials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    reference_letter_id UUID NOT NULL REFERENCES reference_letters(id) ON DELETE CASCADE,
    quote TEXT NOT NULL,
    author_name TEXT NOT NULL,
    author_title TEXT,
    author_company TEXT,
    relationship TEXT NOT NULL DEFAULT 'other',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for testimonials
CREATE INDEX idx_testimonials_profile_id ON testimonials(profile_id);
CREATE INDEX idx_testimonials_reference_letter_id ON testimonials(reference_letter_id);

-- Apply update trigger to testimonials table
CREATE TRIGGER update_testimonials_updated_at
    BEFORE UPDATE ON testimonials
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Skill validations table: links profile skills to reference letters that validate them
CREATE TABLE skill_validations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_skill_id UUID NOT NULL REFERENCES profile_skills(id) ON DELETE CASCADE,
    reference_letter_id UUID NOT NULL REFERENCES reference_letters(id) ON DELETE CASCADE,
    quote_snippet TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (profile_skill_id, reference_letter_id)
);

-- Indexes for skill_validations
CREATE INDEX idx_skill_validations_profile_skill_id ON skill_validations(profile_skill_id);
CREATE INDEX idx_skill_validations_reference_letter_id ON skill_validations(reference_letter_id);

-- Experience validations table: links profile experiences to reference letters that validate them
CREATE TABLE experience_validations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_experience_id UUID NOT NULL REFERENCES profile_experiences(id) ON DELETE CASCADE,
    reference_letter_id UUID NOT NULL REFERENCES reference_letters(id) ON DELETE CASCADE,
    quote_snippet TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (profile_experience_id, reference_letter_id)
);

-- Indexes for experience_validations
CREATE INDEX idx_experience_validations_profile_experience_id ON experience_validations(profile_experience_id);
CREATE INDEX idx_experience_validations_reference_letter_id ON experience_validations(reference_letter_id);
