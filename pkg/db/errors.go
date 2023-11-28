package db

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrNotFound      = errors.New("object not found in database")
	ErrAlreadyExists = errors.New("object already exists in database")
)

func Check(err error) error {
	if err == nil {
		return nil
	}

	if pgErr, ok := err.(*pq.Error); ok {
		if pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
