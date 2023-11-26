package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bbengfort/cosmos/pkg/db"
)

type User struct {
	ID        int64
	Name      sql.NullString
	Email     string
	Password  string
	RoleID    sql.NullInt64
	LastLogin sql.NullTime
	Created   time.Time
	Modified  time.Time
}

// Get user by ID (int64) or by email (string).
func GetUser(ctx context.Context, id any) (u *User, err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var query string
	switch id.(type) {
	case int64:
		query = "SELECT * FROM users WHERE id=$1"
	case string:
		query = "SELECT * F?ROM users WHERE email=$1"
	default:
		return nil, fmt.Errorf("unknown user id type %T", id)
	}

	u = &User{}
	if err = tx.QueryRow(query, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.RoleID,
		&u.LastLogin, &u.Created, &u.Modified,
	); err != nil {
		return nil, err
	}

	tx.Commit()
	return u, nil
}
