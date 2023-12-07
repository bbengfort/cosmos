package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/jmoiron/sqlx"
)

const defaultRole = "DefaultRole"

type Role struct {
	ID          int64          `db:"id"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	IsDefault   bool           `db:"is_default"`
	Created     time.Time      `db:"created"`
	Modified    time.Time      `db:"modified"`
	permissions []*Permission
}

type Permission struct {
	ID          int64          `db:"id"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	Created     time.Time      `db:"created"`
	Modified    time.Time      `db:"modified"`
}

// Get role by ID (int64) or by title (string).
func GetRole(ctx context.Context, nameOrID any) (role *Role, err error) {
	var tx *sqlx.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if role, err = getRole(tx, nameOrID); err != nil {
		return nil, err
	}

	tx.Commit()
	return role, err
}

// Get role by ID (int64) or by title (string) from any transaction.
func getRole(tx *sqlx.Tx, nameOrID any) (role *Role, err error) {
	var query string
	params := []interface{}{nameOrID}

	switch t := nameOrID.(type) {
	case int64, sql.NullInt64:
		query = "SELECT * FROM roles WHERE id=$1"
	case string:
		if t == defaultRole {
			query = "SELECT * FROM roles WHERE is_default IS true LIMIT 1"
			params = []interface{}{}
		} else {
			query = "SELECT * FROM roles WHERE title=$1"
		}
	default:
		return nil, fmt.Errorf("unknown role id type %T", nameOrID)
	}

	// Fetch the role
	role = &Role{}
	if err = tx.QueryRow(query, params...).Scan(&role.ID, &role.Title, &role.Description, &role.IsDefault, &role.Created, &role.Modified); err != nil {
		return nil, err
	}

	// Fetch the role's permissions
	if err = role.getPermissions(tx); err != nil {
		return nil, err
	}

	return role, err
}

// Fetch the roles permissions; if they're already cached on the role they're returned
// directly, otherwise the database is queried to populate the permissions.
func (r *Role) Permissions(ctx context.Context) (_ []*Permission, err error) {
	if r.permissions == nil {
		var tx *sqlx.Tx
		if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
			return nil, err
		}
		defer tx.Rollback()

		if err = r.getPermissions(tx); err != nil {
			return nil, err
		}
		tx.Commit()
	}
	return r.permissions, nil
}

const (
	getRolePerms = "SELECT p.id, p.title, p.description, p.created, p.modified FROM role_permissions rp JOIN permissions p on rp.permission_id=p.id WHERE rp.role_id=$1"
)

func (r *Role) getPermissions(tx *sqlx.Tx) (err error) {
	r.permissions = make([]*Permission, 0, 4)

	var rows *sql.Rows
	if rows, err = tx.Query(getRolePerms, r.ID); err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		p := &Permission{}
		if err = rows.Scan(&p.ID, &p.Title, &p.Description, &p.Created, &p.Modified); err != nil {
			return err
		}
		r.permissions = append(r.permissions, p)
	}

	return rows.Err()
}
