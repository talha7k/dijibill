-- Add intro_viewed column to users table
ALTER TABLE users ADD COLUMN intro_viewed BOOLEAN DEFAULT FALSE;

-- Update existing users to have intro_viewed as false
UPDATE users SET intro_viewed = FALSE WHERE intro_viewed IS NULL;