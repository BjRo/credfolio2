-- Create authors table to store testimonial author information as a proper entity
-- This enables LinkedIn linking, editing, and deduplication across testimonials

CREATE TABLE authors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    title TEXT,
    company TEXT,
    linkedin_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for efficient lookup by profile
CREATE INDEX idx_authors_profile_id ON authors(profile_id);

-- Prevent obvious duplicates: same name+company for same profile
-- Using a unique index with COALESCE to handle NULL company values
CREATE UNIQUE INDEX idx_authors_profile_name_company
    ON authors(profile_id, name, COALESCE(company, ''));

-- Add author_id foreign key to testimonials
-- Keep existing author columns temporarily for backward compatibility during migration
ALTER TABLE testimonials
    ADD COLUMN author_id UUID REFERENCES authors(id) ON DELETE SET NULL;

-- Make author_name nullable since we'll use author_id going forward
ALTER TABLE testimonials
    ALTER COLUMN author_name DROP NOT NULL;

-- Index for looking up testimonials by author
CREATE INDEX idx_testimonials_author_id ON testimonials(author_id);

-- Migrate existing testimonial data to authors table
-- For each unique author (profile_id, name, company), create an author record
INSERT INTO authors (profile_id, name, title, company)
SELECT DISTINCT ON (t.profile_id, t.author_name, COALESCE(t.author_company, ''))
    t.profile_id,
    t.author_name,
    t.author_title,
    t.author_company
FROM testimonials t
WHERE t.author_name IS NOT NULL;

-- Link existing testimonials to their newly created authors
UPDATE testimonials t
SET author_id = a.id
FROM authors a
WHERE t.profile_id = a.profile_id
  AND t.author_name = a.name
  AND COALESCE(t.author_company, '') = COALESCE(a.company, '');
