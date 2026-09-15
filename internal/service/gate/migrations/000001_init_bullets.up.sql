CREATE TYPE bullet_type AS ENUM ('task', 'event', 'note');
CREATE TYPE signifier AS ENUM ('open', 'completed', 'migrated', 'scheduled', 'cancelled');

CREATE TABLE bullets (
  id UUID PRIMARY KEY,
  type bullet_type NOT NULL,
  signifier signifier NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_bullets_created_at ON bullets(created_at);