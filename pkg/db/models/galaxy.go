package models

import "time"

type Galaxy struct {
	ID         int64
	Name       string
	Turn       int64
	Size       Size
	MaxPlayers int16
	MaxTurns   int64
	Created    time.Time
	Modified   time.Time
}
