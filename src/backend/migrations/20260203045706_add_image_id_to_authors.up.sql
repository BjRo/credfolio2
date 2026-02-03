-- Add image_id column to authors table for profile pictures
ALTER TABLE authors
ADD COLUMN image_id UUID REFERENCES files(id) ON DELETE SET NULL;

-- Create index for faster lookups
CREATE INDEX idx_authors_image_id ON authors(image_id) WHERE image_id IS NOT NULL;
