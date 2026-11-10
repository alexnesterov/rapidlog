CREATE TYPE signifier_old AS ENUM ('open', 'completed', 'scheduled', 'cancelled');

ALTER TABLE bullets
  ALTER COLUMN signifier TYPE signifier_old USING signifier::text::signifier_old;

DROP TYPE signifier;
ALTER TYPE signifier_old RENAME TO signifier;