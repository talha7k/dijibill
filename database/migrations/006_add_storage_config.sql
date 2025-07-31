-- Migration: Add storage configuration to system_settings table
-- This migration adds storage configuration fields to support both local and S3 storage

-- Add storage configuration columns to system_settings table
ALTER TABLE system_settings ADD COLUMN storage_type TEXT DEFAULT 'local';
ALTER TABLE system_settings ADD COLUMN storage_local_path TEXT DEFAULT './app_images';
ALTER TABLE system_settings ADD COLUMN storage_s3_bucket TEXT DEFAULT '';
ALTER TABLE system_settings ADD COLUMN storage_s3_region TEXT DEFAULT '';
ALTER TABLE system_settings ADD COLUMN storage_s3_prefix TEXT DEFAULT '';

-- Update existing records with default values
UPDATE system_settings SET 
    storage_type = 'local',
    storage_local_path = './app_images',
    storage_s3_bucket = '',
    storage_s3_region = '',
    storage_s3_prefix = ''
WHERE storage_type IS NULL OR storage_type = '';