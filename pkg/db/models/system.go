package models

import "time"

type System struct {
	ID           int64
	GalaxyID     int64
	Name         string
	IsHomeSystem bool
	StarClass    StarClass
	SystemRadius int16
	WarpGate     int16
	Shipyard     int16
	Created      time.Time
	Modified     time.Time
}
