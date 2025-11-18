-- приклад для Postgres
ALTER TABLE users ADD COLUMN email_norm text NOT NULL;
UPDATE users SET email_norm = lower(trim(email));
-- ALTER TABLE users ALTER COLUMN email_norm SET NOT NULL;
-- CREATE UNIQUE INDEX users_email_norm_uidx ON users (email_norm);