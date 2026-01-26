-- Drop resumes table and related objects

DROP TRIGGER IF EXISTS update_resumes_updated_at ON resumes;
DROP INDEX IF EXISTS idx_resumes_status;
DROP INDEX IF EXISTS idx_resumes_user_id;
DROP TABLE IF EXISTS resumes;
