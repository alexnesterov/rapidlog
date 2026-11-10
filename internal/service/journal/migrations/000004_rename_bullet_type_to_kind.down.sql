ALTER TYPE kind RENAME TO bullet_type;
ALTER TABLE bullets RENAME COLUMN kind TO type;