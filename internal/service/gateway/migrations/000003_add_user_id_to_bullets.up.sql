DELETE FROM bullets;
ALTER TABLE bullets ADD COLUMN user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE;
CREATE INDEX idx_bullets_user_id ON bullets(user_id);