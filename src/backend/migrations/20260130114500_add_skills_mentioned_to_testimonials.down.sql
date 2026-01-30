-- Remove skills_mentioned column from testimonials table

ALTER TABLE testimonials DROP COLUMN IF EXISTS skills_mentioned;
