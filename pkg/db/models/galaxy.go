package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/bbengfort/cosmos/pkg/enums"
	"github.com/bbengfort/cosmos/pkg/jcode"
)

type Galaxy struct {
	ID         int64
	Name       string
	Turn       int64
	Size       enums.Size
	MaxPlayers int16
	MaxTurns   int64
	JoinCode   jcode.JoinCode
	GameState  enums.GameState
	Created    time.Time
	Modified   time.Time
}

const (
	createGalaxySQL = "INSERT INTO galaxies (name, turn, size, max_players, max_turns, join_code, created, modified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING ID"
)

func CreateGalaxy(ctx context.Context, galaxy *Galaxy) (err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false}); err != nil {
		return err
	}
	defer tx.Rollback()

	// Set created and updated timestamps
	galaxy.Created = time.Now()
	galaxy.Modified = galaxy.Created

	if err = tx.QueryRow(createGalaxySQL,
		galaxy.Name,
		galaxy.Turn,
		galaxy.Size,
		galaxy.MaxPlayers,
		galaxy.MaxTurns,
		galaxy.JoinCode,
		galaxy.Created,
		galaxy.Modified,
	).Scan(&galaxy.ID); err != nil {
		return err
	}
	return tx.Commit()
}

const (
	listGalaxiesSQL = "SELECT g.* FROM galaxies g JOIN players p on g.id=p.galaxy_id WHERE p.player_id=$1 AND g.game_state=$2"
)

func ListGalaxies(ctx context.Context, userID int64, state enums.GameState) (galaxies []*Galaxy, err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var rows *sql.Rows
	if rows, err = tx.Query(listGalaxiesSQL, userID, state); err != nil {
		return nil, err
	}
	defer rows.Close()

	galaxies = make([]*Galaxy, 0)
	for rows.Next() {
		g := &Galaxy{}
		if err = rows.Scan(&g.ID, &g.Name, &g.Turn, &g.Size, &g.MaxPlayers, &g.MaxTurns, &g.JoinCode, &g.GameState, &g.Created, &g.Modified); err != nil {
			return nil, err
		}
		galaxies = append(galaxies, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	tx.Commit()
	return galaxies, nil
}
