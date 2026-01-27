-- Create profile_education table for education entries in user profiles
-- Mirrors profile_experiences structure adapted for education data

CREATE TABLE profile_education (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    institution VARCHAR(255) NOT NULL,
    degree VARCHAR(255) NOT NULL,
    field VARCHAR(255),
    start_date VARCHAR(50),
    end_date VARCHAR(50),
    is_current BOOLEAN NOT NULL DEFAULT false,
    description TEXT,
    gpa VARCHAR(20),
    display_order INT NOT NULL DEFAULT 0,
    source VARCHAR(50) NOT NULL DEFAULT 'manual',
    source_resume_id UUID REFERENCES resumes(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for looking up education by profile
CREATE INDEX idx_profile_education_profile_id ON profile_education(profile_id);

-- Index for ordering education entries
CREATE INDEX idx_profile_education_display_order ON profile_education(profile_id, display_order);

-- Apply update trigger to profile_education table
CREATE TRIGGER update_profile_education_updated_at
    BEFORE UPDATE ON profile_education
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
