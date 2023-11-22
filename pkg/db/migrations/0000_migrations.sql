-- Migrations for cosmos data storage.
-- These migrations target an postgresql database that is used to coordinate the cosmos
-- deployment and scale. The migrations table allows a booting node to determine which
-- version its schema is at so that it can quickly make changes to its data store when
-- the node starts or during runtime.
BEGIN;

-- The migrations table stores the migrations applied to arrive at the current schema
-- of the database. The cosmos application checks this table for the version the
-- db is at and applies any later migrations as needed.
CREATE TABLE IF NOT EXISTS migrations (
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(512) NOT NULL,
    version VARCHAR(255) NOT NULL,
    created TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Automatically update modified timestamps
CREATE OR REPLACE FUNCTION trigger_set_modified_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.modified = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMIT;