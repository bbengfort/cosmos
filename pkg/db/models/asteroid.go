package models

import "time"

type Asteroid struct {
	ID       int64
	SystemID int64
	Orbit    int16
	Density  float64
	Created  time.Time
	Modified time.Time
}
