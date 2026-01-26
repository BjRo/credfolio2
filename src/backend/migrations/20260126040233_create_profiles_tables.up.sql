-- Create profiles and profile_experiences tables for manual profile editing
-- Profiles store user profile data separate from resume extractions
-- Profile experiences are manually entered or edited work history

-- Profiles table: stores user profile data
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Index for looking up profile by user
CREATE INDEX idx_profiles_user_id ON profiles(user_id);

-- Profile experiences table: stores work experience entries
CREATE TABLE profile_experiences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    company VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    start_date VARCHAR(50),
    end_date VARCHAR(50),
    is_current BOOLEAN NOT NULL DEFAULT false,
    description TEXT,
    highlights TEXT[],
    display_order INT NOT NULL DEFAULT 0,
    source VARCHAR(50) NOT NULL DEFAULT 'manual',
    source_resume_id UUID REFERENCES resumes(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for looking up experiences by profile
CREATE INDEX idx_profile_experiences_profile_id ON profile_experiences(profile_id);

-- Index for ordering experiences
CREATE INDEX idx_profile_experiences_display_order ON profile_experiences(profile_id, display_order);

-- Apply update trigger to profiles table
CREATE TRIGGER update_profiles_updated_at
    BEFORE UPDATE ON profiles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply update trigger to profile_experiences table
CREATE TRIGGER update_profile_experiences_updated_at
    BEFORE UPDATE ON profile_experiences
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
