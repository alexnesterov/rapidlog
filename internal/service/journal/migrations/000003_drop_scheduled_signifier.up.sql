DELETE FROM bullets WHERE signifier = 'scheduled';

CREATE TYPE signifier_new AS ENUM ('open', 'completed', 'cancelled');

ALTER TABLE bullets
  ALTER COLUMN signifier TYPE signifier_new USING signifier::text::signifier_new;

DROP TYPE signifier;
ALTER TYPE signifier_new RENAME TO signifier;