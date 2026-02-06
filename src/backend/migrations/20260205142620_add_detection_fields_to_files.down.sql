ALTER TABLE files
  DROP COLUMN IF EXISTS detection_error,
  DROP COLUMN IF EXISTS detection_result,
  DROP COLUMN IF EXISTS detection_status;
