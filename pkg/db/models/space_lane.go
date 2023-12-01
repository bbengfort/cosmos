package models

import "time"

type SpaceLane struct {
	OriginID int64
	TargetID int64
	Distance int16
	Hazards  int16
	Created  time.Time
	Modified time.Time
}
