-- Create profile_skills table for managing user skills (manual + extraction-sourced)
CREATE TABLE profile_skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    normalized_name TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT 'TECHNICAL',
    display_order INT NOT NULL DEFAULT 0,
    source TEXT NOT NULL DEFAULT 'manual',
    source_resume_id UUID REFERENCES resumes(id) ON DELETE SET NULL,
    original_data JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for efficient querying
CREATE INDEX idx_profile_skills_profile_id ON profile_skills(profile_id);
CREATE INDEX idx_profile_skills_display_order ON profile_skills(profile_id, display_order);
CREATE INDEX idx_profile_skills_source_resume_id ON profile_skills(source_resume_id);

-- Unique constraint on normalized_name per profile to prevent duplicate skills
CREATE UNIQUE INDEX idx_profile_skills_unique_name ON profile_skills(profile_id, normalized_name);

-- Auto-update updated_at timestamp
CREATE TRIGGER set_profile_skills_updated_at
    BEFORE UPDATE ON profile_skills
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
