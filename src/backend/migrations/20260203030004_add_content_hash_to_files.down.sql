-- Remove content_hash column and index

DROP INDEX IF EXISTS idx_files_user_id_content_hash;
ALTER TABLE files DROP COLUMN IF EXISTS content_hash;
