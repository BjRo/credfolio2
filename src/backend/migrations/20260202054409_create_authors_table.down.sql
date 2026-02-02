-- Revert: remove authors table and restore testimonials to original state

-- First, copy author data back to testimonials (in case it was modified)
UPDATE testimonials t
SET
    author_name = COALESCE(t.author_name, a.name),
    author_title = COALESCE(t.author_title, a.title),
    author_company = COALESCE(t.author_company, a.company)
FROM authors a
WHERE t.author_id = a.id;

-- Drop the foreign key index
DROP INDEX IF EXISTS idx_testimonials_author_id;

-- Remove author_id from testimonials
ALTER TABLE testimonials DROP COLUMN IF EXISTS author_id;

-- Restore NOT NULL constraint on author_name
ALTER TABLE testimonials ALTER COLUMN author_name SET NOT NULL;

-- Drop the unique index and authors table
DROP INDEX IF EXISTS idx_authors_profile_name_company;
DROP INDEX IF EXISTS idx_authors_profile_id;
DROP TABLE IF EXISTS authors;
