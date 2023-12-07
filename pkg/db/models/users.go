package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        int64          `db:"id"`
	Name      sql.NullString `db:"name"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	RoleID    sql.NullInt64  `db:"role_id"`
	LastLogin sql.NullTime   `db:"last_login"`
	Created   time.Time      `db:"created"`
	Modified  time.Time      `db:"modified"`
	role      *Role
}

const (
	createUserSQL     = "INSERT INTO users (name, email, password, role_id, last_login) VALUES (:name, :email, :password, :role_id, :last_login);"
	popCreatedUserSQL = "SELECT id, created, modified FROM users WHERE email=$1"
)

// Create a new user in the database.
func CreateUser(ctx context.Context, user *User) (err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false}); err != nil {
		return err
	}
	defer tx.Rollback()

	// Assign the user to the default role in the database
	if user.role, err = getRole(tx, defaultRole); err != nil {
		return fmt.Errorf("could not get default role: %w", err)
	}

	// Populate fields not set by API calls
	user.LastLogin = sql.NullTime{}
	user.RoleID = sql.NullInt64{Valid: true, Int64: user.role.ID}

	// Execute the insert query
	if _, err = tx.NamedExec(createUserSQL, user); err != nil {
		return err
	}

	// Populate the final fields for creating the user
	if err = tx.QueryRowx(popCreatedUserSQL, user.Email).StructScan(user); err != nil {
		return err
	}

	return tx.Commit()
}

// Get user by ID (int64) or by email (string).
func GetUser(ctx context.Context, id any) (u *User, err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var query string
	switch id.(type) {
	case int64:
		query = "SELECT * FROM users WHERE id=$1"
	case string:
		query = "SELECT * FROM users WHERE email=$1"
	default:
		return nil, fmt.Errorf("unknown user id type %T", id)
	}

	u = &User{}
	if err = tx.QueryRowx(query, id).StructScan(u); err != nil {
		return nil, err
	}

	// Fetch the role and permissions of the user
	if u.role, err = getRole(tx, u.RoleID); err != nil {
		return nil, err
	}

	tx.Commit()
	return u, nil
}

func (u *User) Role(ctx context.Context) (_ *Role, err error) {
	if u.role == nil {
		var tx *sqlx.Tx
		if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
			return nil, err
		}
		defer tx.Rollback()

		if u.role, err = getRole(tx, u.RoleID); err != nil {
			return nil, err
		}
		tx.Commit()
	}
	return u.role, nil
}

func (u *User) Permissions(ctx context.Context) (_ []*Permission, err error) {
	var role *Role
	if role, err = u.Role(ctx); err != nil {
		return nil, err
	}
	return role.Permissions(ctx)
}

const updateLastLoginSQL = "UPDATE users SET last_login=:last_login WHERE id=:id"

func (u *User) LoggedIn(ctx context.Context) (err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: false}); err != nil {
		return err
	}
	defer tx.Rollback()

	u.LastLogin = sql.NullTime{Valid: true, Time: time.Now()}
	if _, err = tx.NamedExec(updateLastLoginSQL, u); err != nil {
		return err
	}
	return tx.Commit()
}
