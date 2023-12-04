package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/bbengfort/cosmos/pkg/enums"
)

type Player struct {
	GalaxyID     int64
	PlayerID     int64
	RoleID       int64
	HomeSystemID int64
	Name         string
	Faction      enums.Faction
	Character    enums.Characteristic
	Created      time.Time
	Modified     time.Time
	role         *Role
}

const (
	createPlayerSQL = "INSERT INTO players (galaxy_id, player_id, role_id, home_system_id, name, faction, character, created, modified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
)

func CreatePlayer(ctx context.Context, player *Player) (err error) {
	var tx *sql.Tx
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

	if _, err = tx.Exec(createPlayerSQL,
		player.GalaxyID,
		player.PlayerID,
		player.RoleID,
		player.HomeSystemID,
		player.Name,
		player.Faction,
		player.Character,
		player.Created,
		player.Modified,
	); err != nil {
		return err
	}
	return tx.Commit()
}
