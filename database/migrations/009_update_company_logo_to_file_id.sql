-- Migration 009: Update company logo field to store file ID instead of base64 data

-- Add new logo_file_id column
ALTER TABLE companies ADD COLUMN logo_file_id INTEGER;

-- Create index for logo_file_id
CREATE INDEX IF NOT EXISTS idx_companies_logo_file_id ON companies(logo_file_id);

-- Note: The old logo column will be kept for backward compatibility during transition
-- It can be removed in a future migration once all data is migrated