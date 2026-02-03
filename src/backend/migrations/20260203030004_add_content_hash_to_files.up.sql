-- Add content_hash column for duplicate file detection
-- SHA-256 hash is 64 hex characters

ALTER TABLE files ADD COLUMN content_hash VARCHAR(64);

-- Index for efficient duplicate detection queries
-- Using (user_id, content_hash) allows checking for duplicates per-user
CREATE INDEX idx_files_user_id_content_hash ON files(user_id, content_hash);
