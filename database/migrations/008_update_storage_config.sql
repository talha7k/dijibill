-- Update system_settings table for separate file database configuration
-- Replace old S3 fields with new flexible storage and file database fields

-- Add new storage configuration columns
ALTER TABLE system_settings ADD COLUMN storage_base_path TEXT DEFAULT './storage';

-- Add file database configuration columns
ALTER TABLE system_settings ADD COLUMN file_db_path TEXT DEFAULT './data/files.db';
ALTER TABLE system_settings ADD COLUMN file_db_sync_url TEXT DEFAULT '';
ALTER TABLE system_settings ADD COLUMN file_db_sync_token TEXT DEFAULT '';
ALTER TABLE system_settings ADD COLUMN file_db_auto_sync BOOLEAN DEFAULT 0;
ALTER TABLE system_settings ADD COLUMN file_db_sync_interval INTEGER DEFAULT 60;

-- Update existing records with default values
UPDATE system_settings SET 
    storage_type = 'local',
    storage_base_path = COALESCE(storage_local_path, './storage'),
    file_db_path = './data/files.db',
    file_db_sync_url = '',
    file_db_sync_token = '',
    file_db_auto_sync = 0,
    file_db_sync_interval = 60
WHERE storage_base_path IS NULL;

-- Remove old S3-specific columns (if they exist)
-- Note: SQLite doesn't support DROP COLUMN directly, so we'll leave them for backward compatibility
-- They can be ignored in the application code