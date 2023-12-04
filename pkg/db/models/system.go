package models

import (
	"time"

	"github.com/bbengfort/cosmos/pkg/enums"
)

type System struct {
	ID           int64
	GalaxyID     int64
	Name         string
	IsHomeSystem bool
	StarClass    enums.StarClass
	SystemRadius int16
	WarpGate     int16
	Shipyard     int16
	Created      time.Time
	Modified     time.Time
}
