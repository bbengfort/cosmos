/*
 * Initial Cosmos database schema for a PostgreSQL database.
 */
BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id        SERIAL PRIMARY KEY,
    email     VARCHAR(255) UNIQUE NOT NULL,
    username  VARCHAR(255) UNIQUE NOT NULL,
    password  VARCHAR(255) NOT NULL,
    last_seen TIMESTAMPTZ,
    created   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Automatically update modified at timestamps
CREATE OR REPLACE FUNCTION trigger_set_modified_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.modified = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Users modified timestamp
CREATE TRIGGER set_users_modified
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

COMMIT;