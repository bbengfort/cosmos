BEGIN;

-- drop tables in order of foreign key dependencies (to prevent errors)
DROP TABLE IF EXISTS users;

-- drop triggers after all table dependencies have been removed
DROP FUNCTION trigger_set_modified_timestamp;

COMMIT;