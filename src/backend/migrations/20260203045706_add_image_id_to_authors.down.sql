-- Remove image_id column from authors table
DROP INDEX IF EXISTS idx_authors_image_id;
ALTER TABLE authors DROP COLUMN IF EXISTS image_id;
