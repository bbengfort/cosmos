package models

import "time"

type Player struct {
	GalaxyID     int64
	PlayerID     int64
	RoleID       int64
	HomeSystemID int64
	Name         string
	Faction      Faction
	Character    Characteristic
	Created      time.Time
	Modified     time.Time
}
