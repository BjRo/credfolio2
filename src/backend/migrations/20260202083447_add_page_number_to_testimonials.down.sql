-- Remove page_number column from testimonials table

ALTER TABLE testimonials
DROP COLUMN IF EXISTS page_number;
