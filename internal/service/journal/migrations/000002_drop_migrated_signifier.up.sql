DELETE FROM bullets WHERE signifier = 'migrated';

CREATE TYPE signifier_new AS ENUM ('open', 'completed', 'scheduled', 'cancelled');

ALTER TABLE bullets
  ALTER COLUMN signifier TYPE signifier_new USING signifier::text::signifier_new;

DROP TYPE signifier;
ALTER TYPE signifier_new RENAME TO signifier;
