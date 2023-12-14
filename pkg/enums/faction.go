package enums

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

type Faction uint8

const (
	UnknownFaction Faction = iota
	Supremacy
	Harmony
	Purity
)

var factionNames = [4]string{"unknown", "supremacy", "harmony", "purity"}

//=====================================================================================
// Stringer interface
//=====================================================================================

func (f Faction) String() string {
	return factionNames[f]
}

//=====================================================================================
// Valuer interface
//=====================================================================================

func (f Faction) Value() (driver.Value, error) {
	return factionNames[f], nil
}

//=====================================================================================
// Scanner interface
//=====================================================================================

func (f *Faction) Scan(value interface{}) error {
	// If value is nil set size to unknown
	if value == nil {
		*f = UnknownFaction
		return nil
	}

	// Convert the value to a string
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.([]byte); ok {
			// Parse the value of v
			switch string(v) {
			case "unknown":
				*f = UnknownFaction
			case "supremacy":
				*f = Supremacy
			case "harmony":
				*f = Harmony
			case "purity":
				*f = Purity
			default:
				return ErrScanFaction
			}
			return nil
		}
	}

	return ErrScanFaction
}

//=====================================================================================
// JSON Marshaler interface
//=====================================================================================

func (f Faction) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

//=====================================================================================
// JSON Unmarshaler interface
//=====================================================================================

func (f *Faction) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToLower(strings.TrimSpace(sv))
	return f.Scan(sv)
}

//=====================================================================================
// Nullable Type
//=====================================================================================

type NullFaction struct {
	Faction Faction
	Valid   bool // Valid is true if Faction is not NULL
}

func (p *NullFaction) Scan(value any) (err error) {
	if value == nil {
		p.Faction, p.Valid = UnknownFaction, false
		return nil
	}

	p.Valid = true
	return p.Faction.Scan(value)
}

func (p NullFaction) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return p.Faction.Value()
}
