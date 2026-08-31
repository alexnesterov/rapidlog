DROP INDEX IF EXISTS idx_bullets_user_id;
ALTER TABLE bullets DROP COLUMN IF EXISTS user_id;