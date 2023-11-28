-- Populate the database with initial roles and permissions data
BEGIN;

INSERT INTO roles (id, title, description, is_default) VALUES
    (1, 'Admin', 'Able to manage the cosmos API and other users on the server', 'f'),
    (2, 'Player', 'Able to create games and play and manage games they belong to', 't'),
    (3, 'Observer', 'Only has read only access to the cosmos server', 'f')
;

INSERT INTO permissions (id, title, description) VALUES
    (1, 'users:manage', 'Can manage users, roles, and permissions'),
    (2, 'games:read', 'Can view the games that are on the server'),
    (3, 'games:create', 'Can create a new game on the server'),
    (4, 'games:manage', 'Can manage all games that are on the server')
;

INSERT INTO role_permissions (role_id, permission_id) VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (1, 4),
    (2, 2),
    (2, 3),
    (3, 2)
;

COMMIT;