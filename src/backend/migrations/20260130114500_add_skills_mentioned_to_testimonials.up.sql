-- Add skills_mentioned column to testimonials table
-- This stores the skill names mentioned in each specific testimonial quote
-- Used to filter validatedSkills to only show skills relevant to that quote

ALTER TABLE testimonials ADD COLUMN skills_mentioned TEXT[];
