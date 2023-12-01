package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

type Size uint8

const (
	UnknownSize Size = iota
	Small
	Medium
	Large
	Galactic
	Cosmic
)

var sizeNames = []string{"unknown", "small", "medium", "large", "galactic", "cosmic"}

type Faction uint8

const (
	UnknownFaction Faction = iota
	Supremacy
	Harmony
	Purity
)

var factionNames = []string{"unknown", "supremacy", "harmony", "purity"}

type Characteristic uint8

const (
	UnknownCharacteristic Characteristic = iota
	Benevolent
	Progressive
	Humanitarian
	Charismatic
	Indusrialist
	Diplomat
	Warrior
	Economist
)

var characteristicNames = []string{"unknown", "benevolent", "progressive", "humanitarian", "charismatic", "industrialist", "diplomat", "warrior", "economist"}

type StarClass uint8

const (
	UnknownStarClass StarClass = iota
	Os
	Bs
	As
	Fs
	Gs
	Ks
	Ms
)

var starClassNames = []string{"UNKNOWN", "O", "B", "A", "F", "G", "K", "M"}

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

//=====================================================================================
// Stringer interface
//=====================================================================================

func (s Size) String() string {
	return sizeNames[s]
}

func (f Faction) String() string {
	return factionNames[f]
}

func (c Characteristic) String() string {
	return characteristicNames[c]
}

func (s StarClass) String() string {
	return starClassNames[s]
}

func (p PlanetClass) String() string {
	return planetClassNames[p]
}

//=====================================================================================
// Valuer interface
//=====================================================================================

func (s Size) Value() (driver.Value, error) {
	return sizeNames[s], nil
}

func (f Faction) Value() (driver.Value, error) {
	return factionNames[f], nil
}

func (c Characteristic) Value() (driver.Value, error) {
	return characteristicNames[c], nil
}

func (s StarClass) Value() (driver.Value, error) {
	return starClassNames[s], nil
}

func (p PlanetClass) Value() (driver.Value, error) {
	return planetClassNames[p], nil
}

//=====================================================================================
// Scanner interface
//=====================================================================================

var (
	ErrScanSize           = errors.New("failed to scan size enum")
	ErrScanFaction        = errors.New("failed to scan faction enum")
	ErrScanCharacteristic = errors.New("failed to scan characteristic enum")
	ErrScanStarClass      = errors.New("failed to scan star class enum")
	ErrScanPlanetClass    = errors.New("failed to scan planet class enum")
)

func (s *Size) Scan(value interface{}) error {
	// If value is nil set size to unknown
	if value == nil {
		*s = UnknownSize
		return nil
	}

	// Convert the value to a string
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			// Parse the value of v
			switch v {
			case "unknown":
				*s = UnknownSize
			case "small":
				*s = Small
			case "medium":
				*s = Medium
			case "large":
				*s = Large
			case "galactic":
				*s = Galactic
			case "cosmic":
				*s = Cosmic
			default:
				return ErrScanSize
			}
			return nil
		}
	}

	return ErrScanSize
}

func (f *Faction) Scan(value interface{}) error {
	// If value is nil set size to unknown
	if value == nil {
		*f = UnknownFaction
		return nil
	}

	// Convert the value to a string
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			// Parse the value of v
			switch v {
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

func (c *Characteristic) Scan(value interface{}) error {
	// If value is nil set size to unknown
	if value == nil {
		*c = UnknownCharacteristic
		return nil
	}

	// Convert the value to a string
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			// Parse the value of v
			switch v {
			case "unknown":
				*c = UnknownCharacteristic
			case "benevolent":
				*c = Benevolent
			case "progressive":
				*c = Progressive
			case "humanitarian":
				*c = Humanitarian
			case "charismatic":
				*c = Charismatic
			case "industrialist":
				*c = Indusrialist
			case "diplomat":
				*c = Diplomat
			case "warrior":
				*c = Warrior
			case "economist":
				*c = Economist
			default:
				return ErrScanCharacteristic
			}
			return nil
		}
	}

	return ErrScanCharacteristic
}

func (s *StarClass) Scan(value interface{}) error {
	// If value is nil set size to unknown
	if value == nil {
		*s = UnknownStarClass
		return nil
	}

	// Convert the value to a string
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			// Parse the value of v
			switch v {
			case "UNKNOWN":
				*s = UnknownStarClass
			case "O":
				*s = Os
			case "B":
				*s = Bs
			case "A":
				*s = As
			case "F":
				*s = Fs
			case "G":
				*s = Gs
			case "K":
				*s = Ks
			case "M":
				*s = Ms
			default:
				return ErrScanStarClass
			}
			return nil
		}
	}

	return ErrScanStarClass
}

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

func (s Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (f Faction) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (c Characteristic) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (s StarClass) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (p PlanetClass) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

//=====================================================================================
// JSON Unmarshaler interface
//=====================================================================================

func (s *Size) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToLower(strings.TrimSpace(sv))
	return s.Scan(sv)
}

func (f *Faction) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToLower(strings.TrimSpace(sv))
	return f.Scan(sv)
}

func (c *Characteristic) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToLower(strings.TrimSpace(sv))
	return c.Scan(sv)
}

func (s *StarClass) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToUpper(strings.TrimSpace(sv))
	return s.Scan(sv)
}

func (p *PlanetClass) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToUpper(strings.TrimSpace(sv))
	return p.Scan(sv)
}
