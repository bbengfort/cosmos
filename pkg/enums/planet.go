package enums

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

type PlanetClass uint8

const (
	UnknownPlanetClass PlanetClass = iota
	Ap
	Bp
	Cp
	Dp
	Ep
	Fp
	Gp
	Hp
	Ip
	Jp
	Kp
	Lp
	Mp
	Np
	Op
	Pp
	Qp
	Rp
	Sp
	Up
	Xp
	Yp
)

var planetClassNames = []string{"UNKNOWN", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "U", "X", "Y"}

// =====================================================================================
// Stringer interface
// =====================================================================================

func (p PlanetClass) String() string {
	return planetClassNames[p]
}

//=====================================================================================
// Valuer interface
//=====================================================================================

func (p PlanetClass) Value() (driver.Value, error) {
	return planetClassNames[p], nil
}

//=====================================================================================
// Scanner interface
//=====================================================================================

func (p *PlanetClass) Scan(value interface{}) error {
	// If value is nil set size to unknown
	if value == nil {
		*p = UnknownPlanetClass
		return nil
	}

	// Convert the value to a string
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			// Parse the value of v
			switch v {
			case "UNKNOWN":
				*p = UnknownPlanetClass
			case "A":
				*p = Ap
			case "B":
				*p = Bp
			case "C":
				*p = Cp
			case "D":
				*p = Dp
			case "E":
				*p = Ep
			case "F":
				*p = Fp
			case "G":
				*p = Gp
			case "H":
				*p = Hp
			case "I":
				*p = Ip
			case "J":
				*p = Jp
			case "K":
				*p = Kp
			case "L":
				*p = Lp
			case "M":
				*p = Mp
			case "N":
				*p = Np
			case "O":
				*p = Op
			case "P":
				*p = Pp
			case "Q":
				*p = Qp
			case "R":
				*p = Rp
			case "S":
				*p = Sp
			case "U":
				*p = Up
			case "X":
				*p = Xp
			case "Y":
				*p = Yp
			default:
				return ErrScanPlanetClass
			}
			return nil
		}
	}

	return ErrScanPlanetClass
}

//=====================================================================================
// JSON Marshaler interface
//=====================================================================================

func (p PlanetClass) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

//=====================================================================================
// JSON Unmarshaler interface
//=====================================================================================

func (p *PlanetClass) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToUpper(strings.TrimSpace(sv))
	return p.Scan(sv)
}

//=====================================================================================
// Nullable Type
//=====================================================================================

type NullPlanetClass struct {
	PlanetClass PlanetClass
	Valid       bool // Valid is true if PlanetClass is not NULL
}

func (p *NullPlanetClass) Scan(value any) (err error) {
	if value == nil {
		p.PlanetClass, p.Valid = UnknownPlanetClass, false
		return nil
	}

	p.Valid = true
	return p.PlanetClass.Scan(value)
}

func (p NullPlanetClass) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return p.PlanetClass.Value()
}
