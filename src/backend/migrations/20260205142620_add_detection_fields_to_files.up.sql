ALTER TABLE files
  ADD COLUMN detection_status VARCHAR(50),
  ADD COLUMN detection_result JSONB,
  ADD COLUMN detection_error TEXT;
