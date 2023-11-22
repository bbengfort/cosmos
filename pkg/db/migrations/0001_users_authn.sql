-- Schema for users, roles, and authentication.
BEGIN;

/*
 * Tables
 */

-- Primary authentication table that holds usernames and hashed passwords.
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(512) DEFAULT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL UNIQUE,
    role_id     INTEGER DEFAULT NULL,
    last_login  TIMESTAMPTZ DEFAULT NULL,
    created     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Roles are collections of permissions that can be quickly assigned to a user
CREATE TABLE IF NOT EXISTS roles (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(512) DEFAULT NULL,
    created     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Permissions (or scopes) authorize the user to perform actions on the API
CREATE TABLE IF NOT EXISTS permissions (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(512) DEFAULT NULL,
    created     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Maps the default permissions to a role
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

/*
 * Foreign Key Relationships
 */

ALTER TABLE users ADD CONSTRAINT fk_users_role
    FOREIGN KEY (role_id) REFERENCES roles (id)
    ON DELETE RESTRICT;

ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_role
    FOREIGN KEY (role_id) REFERENCES roles (id)
    ON DELETE CASCADE;

ALTER TABLE role_permissions ADD CONSTRAINT fk_role_permissions_permission
    FOREIGN KEY (permission_id) REFERENCES permissions (id)
    ON DELETE CASCADE;

/*
 * Views
 */

-- Allows the easy selection of all permissions for a user based on their role
CREATE OR REPLACE VIEW user_permissions AS
    SELECT u.id as user_id, p.title as permission
        FROM users u
        JOIN role_permissions rp ON rp.role_id = u.role_id
        JOIN permissions p ON p.id = rp.permission_id
;

/*
 * Automatically update modified timestamps
 */

-- Users modified timestamp
CREATE TRIGGER set_users_modified
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

-- Roles modified timestamp
CREATE TRIGGER set_roles_modified
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

-- Permissions modified timestamp
CREATE TRIGGER set_permissions_modified
BEFORE UPDATE ON permissions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

-- Role Permissions modified timestamp
CREATE TRIGGER set_role_permissions_modified
BEFORE UPDATE ON role_permissions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

COMMIT;