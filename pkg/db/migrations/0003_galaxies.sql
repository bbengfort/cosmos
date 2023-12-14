-- Galaxies defines a single game in cosmos along with its associated players and map.
BEGIN;

/*
 * Tables
 */

-- Describes the size of the galaxy for 2, 10, 20, 50, and 100 players respectively.
CREATE TYPE SIZE AS ENUM ('small', 'medium', 'large', 'galactic', 'cosmic');

-- Describes the current state of the game.
CREATE TYPE GAME_STATE AS ENUM ('pending', 'playing', 'completed');

-- Galaxies is the list of all games that are currently running.
CREATE TABLE IF NOT EXISTS galaxies (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL DEFAULT '',
    turn        INTEGER NOT NULL DEFAULT 0,
    size        SIZE NOT NULL DEFAULT 'medium',
    max_players SMALLINT NOT NULL DEFAULT 10,
    max_turns   INTEGER NOT NULL DEFAULT 1000,
    join_code   VARCHAR(16) NOT NULL UNIQUE,
    game_state  GAME_STATE NOT NULL DEFAULT 'pending',
    created     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT max_players_nonnegative CHECK (max_players > 0),
    CONSTRAINT max_turns_nonnegative CHECK (max_turns > 0)
);

-- Factions describe the affinity of the player
CREATE TYPE FACTION AS ENUM ('supremacy', 'harmony', 'purity');

-- Characteristics describe special bonuses the player will receive
CREATE TYPE CHARACTERISTIC AS ENUM ('benevolent', 'progressive', 'humanitarian', 'charismatic', 'industrialist', 'diplomat', 'warrior', 'economist');

-- Maps players to the galaxies that they are associated with.
CREATE TABLE IF NOT EXISTS players (
    galaxy_id       INTEGER NOT NULL,
    player_id       INTEGER NOT NULL,
    role_id         INTEGER NOT NULL DEFAULT 2,
    home_system_id  INTEGER DEFAULT NULL,
    name            VARCHAR(255) NOT NULL DEFAULT '',
    faction         FACTION NOT NULL DEFAULT 'harmony',
    character       CHARACTERISTIC NOT NULL DEFAULT 'warrior',
    created         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (galaxy_id, player_id)
);

-- Star class determine the kind of star the system has
CREATE TYPE STAR_CLASS AS ENUM ('O', 'B', 'A', 'F', 'G', 'K', 'M');

-- Systems describe the map of the galaxy and the properties of the game as it runs.
CREATE TABLE IF NOT EXISTS systems (
    id              SERIAL PRIMARY KEY,
    galaxy_id       INTEGER NOT NULL,
    name            VARCHAR(255) NOT NULL,
    is_home_system  BOOLEAN NOT NULL DEFAULT 'f',
    star_class      STAR_CLASS NOT NULL DEFAULT 'G',
    system_radius   SMALLINT NOT NULL DEFAULT 64,
    warp_gate       SMALLINT NOT NULL DEFAULT 0,
    shipyard        SMALLINT NOT NULL DEFAULT 0,
    created         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT system_radius_minimum CHECK (system_radius >= 32),
    CONSTRAINT system_radius_maximum CHECK (system_radius <= 512),
    CONSTRAINT warp_gate_nonnegative CHECK (warp_gate >= 0),
    CONSTRAINT warp_gate_maximum CHECK (warp_gate < 9),
    CONSTRAINT shipyard_nonnegative CHECK (shipyard >= 0),
    CONSTRAINT shipyard_maximum CHECK (shipyard < 9)
);

-- Plant class describes the kind of planet it is
-- See: https://uss-theurgy.com/w/index.php?title=Planetary_Classification
CREATE TYPE PLANET_CLASS AS ENUM ('A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'U', 'X', 'Y');

-- Planets are the economic building blocks of systems that generate credits, energy, research, etc.
CREATE TABLE IF NOT EXISTS planets (
    id              SERIAL PRIMARY KEY,
    system_id       INTEGER NOT NULL,
    name            VARCHAR(255) NOT NULL,
    planet_class    PLANET_CLASS NOT NULL DEFAULT 'M',
    is_homeworld    BOOLEAN NOT NULL DEFAULT 'f',
    orbit           SMALLINT NOT NULL DEFAULT 8,
    orbital_speed   FLOAT NOT NULL DEFAULT 3.0,  -- turns per revolution
    labs            SMALLINT NOT NULL DEFAULT 0,
    tech            INTEGER NOT NULL DEFAULT 0,
    mines           SMALLINT NOT NULL DEFAULT 0,
    metals          INTEGER NOT NULL DEFAULT 0,
    reactors        SMALLINT NOT NULL DEFAULT 0,
    energy          INTEGER NOT NULL DEFAULT 0,
    cities          SMALLINT NOT NULL DEFAULT 0,
    credits         INTEGER NOT NULL DEFAULT 0,
    farms           SMALLINT NOT NULL DEFAULT 0,
    food            INTEGER NOT NULL DEFAULT 0,
    created         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_planets_system_orbit UNIQUE(system_id, orbit),
    CONSTRAINT orbit_minimum CHECK (orbit >= 8),
    CONSTRAINT orbit_maximum CHECK (orbit < 512),
    CONSTRAINT labs_nonnegative CHECK (labs >= 0),
    CONSTRAINT labs_maximum CHECK (labs < 512),
    CONSTRAINT mines_nonnegative CHECK (mines >= 0),
    CONSTRAINT mines_maximum CHECK (mines < 512),
    CONSTRAINT reactors_nonnegative CHECK (reactors >= 0),
    CONSTRAINT reactors_maximum CHECK (reactors < 512),
    CONSTRAINT cities_nonnegative CHECK (cities >= 0),
    CONSTRAINT cities_maximum CHECK (cities < 512),
    CONSTRAINT farms_nonnegative CHECK (farms >= 0),
    CONSTRAINT farms_maximum CHECK (farms < 512)
);

-- Astroid belts are objects in solar systems that provide cover and block space battles
CREATE TABLE IF NOT EXISTS asteroids (
    id              SERIAL PRIMARY KEY,
    system_id       INTEGER NOT NULL,
    orbit           SMALLINT NOT NULL DEFAULT 32,
    density         FLOAT NOT NULL DEFAULT 0.5,
    created         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_asteroids_system_orbit UNIQUE(system_id, orbit),
    CONSTRAINT orbit_minimum CHECK (orbit >= 1),
    CONSTRAINT orbit_maximum CHECK (orbit < 512),
    CONSTRAINT density_minimum CHECK (orbit >= 0.0),
    CONSTRAINT density_maximum CHECK (orbit <= 1.0)
);

-- Space Lanes describe the connections between star systems that fleets can travel on.
-- Note that the data definition of space lanes is directional.
CREATE TABLE IF NOT EXISTS space_lanes (
    origin_id   INTEGER NOT NULL,
    target_id   INTEGER NOT NULL,
    distance    SMALLINT NOT NULL,
    hazards     SMALLINT NOT NULL,
    created     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (origin_id, target_id),
    CONSTRAINT distance_nonnegative CHECK (distance > 0),
    CONSTRAINT hazards_nonnegative CHECK (hazards >= 0),
    CONSTRAINT hazards_maximum CHECK (hazards < 512)
);


/*
 * Foreign Key Relationships
 */

ALTER TABLE players ADD CONSTRAINT fk_players_galaxy
    FOREIGN KEY (galaxy_id) REFERENCES galaxies (id)
    ON DELETE CASCADE;

ALTER TABLE players ADD CONSTRAINT fk_players_player
    FOREIGN KEY (player_id) REFERENCES users (id)
    ON DELETE CASCADE;

ALTER TABLE players ADD CONSTRAINT fk_players_role
    FOREIGN KEY (role_id) REFERENCES roles (id)
    ON DELETE RESTRICT;

ALTER TABLE players ADD CONSTRAINT fk_players_home_system
    FOREIGN KEY (home_system_id) REFERENCES systems (id)
    ON DELETE RESTRICT;

-- Only one player can occupy a home system
ALTER TABLE players ADD CONSTRAINT unique_players_home_system UNIQUE(home_system_id);

ALTER TABLE systems ADD CONSTRAINT fk_systems_galaxy
    FOREIGN KEY (galaxy_id) REFERENCES galaxies (id)
    ON DELETE CASCADE;

ALTER TABLE planets ADD CONSTRAINT fk_planets_system
    FOREIGN KEY (system_id) REFERENCES systems (id)
    ON DELETE CASCADE;

ALTER TABLE asteroids ADD CONSTRAINT fk_asteroids_system
    FOREIGN KEY (system_id) REFERENCES systems (id)
    ON DELETE CASCADE;

ALTER TABLE space_lanes ADD CONSTRAINT fk_space_lanes_origin
    FOREIGN KEY (origin_id) REFERENCES systems (id)
    ON DELETE CASCADE;

ALTER TABLE space_lanes ADD CONSTRAINT fk_space_lanes_target
    FOREIGN KEY (target_id) REFERENCES systems (id)
    ON DELETE CASCADE;

/*
 * Views
 */

/*
 * Automatically update modified timestamps
 */

-- Galaxies modified timestamp
CREATE TRIGGER set_galaxies_modified
BEFORE UPDATE ON galaxies
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

-- Players modified timestamp
CREATE TRIGGER set_players_modified
BEFORE UPDATE ON players
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

-- systems modified timestamp
CREATE TRIGGER set_systems_modified
BEFORE UPDATE ON systems
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

-- Planets modified timestamps
CREATE TRIGGER set_planets_modified
BEFORE UPDATE ON planets
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

-- Astroids modified timestamps
CREATE TRIGGER set_asteroids_modified
BEFORE UPDATE ON asteroids
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

-- Space Lanes modified timestamps
CREATE TRIGGER set_space_lanes_modified
BEFORE UPDATE ON space_lanes
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modified_timestamp();

COMMIT;