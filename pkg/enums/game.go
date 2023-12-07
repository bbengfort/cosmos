package enums

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

type GameState uint8

const (
	UnknownGameState GameState = iota
	Pending
	Playing
	Completed
)

var gameStateNames = []string{"unknown", "pending", "playing", "completed"}

//=====================================================================================
// Stringer interface
//=====================================================================================

func (s GameState) String() string {
	return gameStateNames[s]
}

//=====================================================================================
// Valuer interface
//=====================================================================================

func (s GameState) Value() (driver.Value, error) {
	return gameStateNames[s], nil
}

//=====================================================================================
// Scanner interface
//=====================================================================================

func (s *GameState) Scan(value interface{}) error {
	// If value is nil set size to unknown
	if value == nil {
		*s = UnknownGameState
		return nil
	}

	// Convert the value to a string
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.([]byte); ok {
			// Parse the value of v
			switch string(v) {
			case "unknown":
				*s = UnknownGameState
			case "pending":
				*s = Pending
			case "playing":
				*s = Playing
			case "completed":
				*s = Completed
			default:
				return ErrScanGameState
			}
			return nil
		}
	}

	return ErrScanGameState
}

//=====================================================================================
// JSON Marshaler interface
//=====================================================================================

func (s GameState) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

//=====================================================================================
// JSON Unmarshaler interface
//=====================================================================================

func (s *GameState) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToLower(strings.TrimSpace(sv))
	return s.Scan(sv)
}

//=====================================================================================
// Nullable Type
//=====================================================================================

type NullGameState struct {
	GameState GameState
	Valid     bool // Valid is true if GameState is not NULL
}

func (g *NullGameState) Scan(value any) (err error) {
	if value == nil {
		g.GameState, g.Valid = UnknownGameState, false
		return nil
	}

	g.Valid = true
	return g.GameState.Scan(value)
}

func (g NullGameState) Value() (driver.Value, error) {
	if !g.Valid {
		return nil, nil
	}
	return g.GameState.Value()
}
