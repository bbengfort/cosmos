package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/bbengfort/cosmos/pkg/enums"
	"github.com/jmoiron/sqlx"
)

type Player struct {
	GalaxyID     int64                `db:"galaxy_id"`
	PlayerID     int64                `db:"player_id"`
	RoleID       int64                `db:"role_id"`
	HomeSystemID sql.NullInt64        `db:"home_system_id"`
	Name         string               `db:"name"`
	Faction      enums.Faction        `db:"faction"`
	Character    enums.Characteristic `db:"character"`
	Created      time.Time            `db:"created"`
	Modified     time.Time            `db:"modified"`
	role         *Role
}

const (
	createPlayerSQL = "INSERT INTO players (galaxy_id, player_id, role_id, home_system_id, name, faction, character, created, modified) VALUES (:galaxy_id, :player_id, :role_id, :home_system_id, :name, :faction, :character, :created, :modified)"
)

func CreatePlayer(ctx context.Context, player *Player) (err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false}); err != nil {
		return err
	}
	defer tx.Rollback()

	// Assign the default role to the player if one isn't on the player.
	if player.RoleID == 0 {
		if player.role, err = getRole(tx, defaultRole); err != nil {
			return fmt.Errorf("could not get default role: %w", err)
		}
		player.RoleID = player.role.ID
	}

	// Set created and updated timestamps
	player.Created = time.Now()
	player.Modified = player.Created

	if _, err = tx.NamedExec(createPlayerSQL, player); err != nil {
		return err
	}
	return tx.Commit()
}
