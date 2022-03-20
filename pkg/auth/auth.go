package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bbengfort/cosmos/pkg/db"
	"github.com/rs/zerolog/log"
)

var (
	ErrUserNotFound    = errors.New("username or email address not found")
	ErrPasswordInvalid = errors.New("incorrect password")
	ErrNoUserID        = errors.New("cannot load a user without an ID")
)

type User struct {
	ID       uint64
	Username string
	Email    string
	LastSeen sql.NullTime
	Created  time.Time
	Modified time.Time
}

const (
	userPasswordSQL = "SELECT id, password FROM users WHERE username=$1 OR email=$1;"
	getUserDataSQL  = "SELECT username, email, last_seen, created, modified FROM users WHERE id=$1"
)

// Authenticate a user with a username (or email) and password.
func Authenticate(ctx context.Context, username, password string) (user *User, err error) {
	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true}); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create the user struct
	user = &User{}

	// Lookup the password derived key from the database
	var dk string
	if err = tx.QueryRow(userPasswordSQL, username).Scan(&user.ID, &dk); err != nil {
		if err != sql.ErrNoRows {
			log.Error().Err(err).Msg("could not query users table")
		}
		return nil, ErrUserNotFound
	}

	// Check the derived key with the password
	var valid bool
	if valid, err = VerifyDerivedKey(dk, password); err != nil {
		log.Error().Err(err).Msg("could not verify user password")
	}

	if !valid {
		return nil, ErrPasswordInvalid
	}

	// If the user is authenticated read the user from the database
	if err = user.load(tx); err != nil {
		log.Error().Err(err).Msg("could not load user from the database")
		return nil, ErrUserNotFound
	}

	tx.Commit()
	return user, nil
}

func (u *User) load(tx *sql.Tx) error {
	if u.ID == 0 {
		return ErrNoUserID
	}

	if err := tx.QueryRow(getUserDataSQL, u.ID).Scan(&u.Username, &u.Email, &u.LastSeen, &u.Created, &u.Modified); err != nil {
		return err
	}
	return nil
}

const (
	insertUserSQL = "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id;"
)

// Register a user inserting their user details into the database
func Register(ctx context.Context, username, email, password string) (user *User, err error) {
	// Create the derived key for the password
	var dk string
	if dk, err = CreateDerivedKey(password); err != nil {
		return nil, err
	}

	var tx *sql.Tx
	if tx, err = db.BeginTx(ctx, nil); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create the user struct
	user = &User{}
	if err = tx.QueryRow(insertUserSQL, username, email, dk).Scan(&user.ID); err != nil {
		return nil, err
	}

	if err = user.load(tx); err != nil {
		return nil, err
	}

	tx.Commit()
	return user, nil
}
