-- Add page_number column to testimonials table
-- Stores the page number in the source PDF where the testimonial quote appears
-- This enables deep linking to the specific page when viewing the source document

ALTER TABLE testimonials
ADD COLUMN page_number INTEGER;

COMMENT ON COLUMN testimonials.page_number IS 'Page number in source PDF where quote appears (1-indexed)';
