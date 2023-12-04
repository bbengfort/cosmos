package models

import (
	"time"

	"github.com/bbengfort/cosmos/pkg/enums"
)

type Planet struct {
	ID          int64
	SystemID    int64
	Name        string
	PlanetClass enums.PlanetClass
	IsHomeworld bool
	Orbit       int16
	Labs        int16
	Tech        int64
	Mines       int16
	Metals      int64
	Reactors    int16
	Energy      int64
	Cities      int16
	Credits     int64
	Farms       int16
	Food        int64
	Created     time.Time
	Modified    time.Time
}
