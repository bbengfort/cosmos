package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/bbengfort/cosmos/pkg/enums"
	"github.com/bbengfort/cosmos/pkg/jcode"
	"github.com/jmoiron/sqlx"
)

type Galaxy struct {
	ID         int64           `db:"id"`
	Name       string          `db:"name"`
	Turn       int64           `db:"turn"`
	Size       enums.Size      `db:"size"`
	MaxPlayers int16           `db:"max_players"`
	MaxTurns   int64           `db:"max_turns"`
	JoinCode   jcode.JoinCode  `db:"join_code"`
	GameState  enums.GameState `db:"game_state"`
	Created    time.Time       `db:"created"`
	Modified   time.Time       `db:"modified"`
}

const (
	createGalaxySQL = "INSERT INTO galaxies (name, turn, size, max_players, max_turns, join_code, created, modified) VALUES (:name, :turn, :size, :max_players, :max_turns, :join_code, :created, :modified) RETURNING ID;"
)

func CreateGalaxy(ctx context.Context, galaxy *Galaxy) (err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false}); err != nil {
		return err
	}
	defer tx.Rollback()

	// Set created and updated timestamps
	galaxy.Created = time.Now()
	galaxy.Modified = galaxy.Created

	var (
		query string
		args  []interface{}
	)

	if query, args, err = tx.BindNamed(createGalaxySQL, galaxy); err != nil {
		return err
	}

	if err = tx.Get(galaxy, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}

const (
	listGalaxiesSQL = "SELECT g.* FROM galaxies g JOIN players p on g.id=p.galaxy_id WHERE p.player_id=$1 AND g.game_state=$2"
)

func ListGalaxies(ctx context.Context, userID int64, state enums.GameState) (galaxies []*Galaxy, err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	galaxies = make([]*Galaxy, 0)
	if err = tx.Select(&galaxies, listGalaxiesSQL, userID, state); err != nil {
		return nil, err
	}

	tx.Commit()
	return galaxies, nil
}

const (
	listActiveGalaxiesSQL = "SELECT g.* FROM galaxies g JOIN players p on g.id=p.galaxy_id WHERE p.player_id=$1 AND g.game_state!='completed'"
)

func ListActiveGalaxies(ctx context.Context, userID int64) (galaxies []*Galaxy, err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	galaxies = make([]*Galaxy, 0)
	if err = tx.Select(&galaxies, listActiveGalaxiesSQL, userID); err != nil {
		return nil, err
	}

	tx.Commit()
	return galaxies, nil
}
