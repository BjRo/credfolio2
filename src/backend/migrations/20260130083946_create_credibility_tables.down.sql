-- Rollback credibility system tables

DROP TABLE IF EXISTS experience_validations;
DROP TABLE IF EXISTS skill_validations;
DROP TABLE IF EXISTS testimonials;
ALTER TABLE reference_letters DROP COLUMN IF EXISTS error_message;
