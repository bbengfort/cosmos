package enums

import "errors"

var (
	ErrScanSize           = errors.New("failed to parse size enum")
	ErrScanFaction        = errors.New("failed to parse faction enum")
	ErrScanCharacteristic = errors.New("failed to parse characteristic enum")
	ErrScanStarClass      = errors.New("failed to parse star class enum")
	ErrScanPlanetClass    = errors.New("failed to parse planet class enum")
	ErrScanGameState      = errors.New("failed to parse game state enum")
)
