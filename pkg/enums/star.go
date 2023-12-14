package enums

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

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

var starClassNames = [8]string{"UNKNOWN", "O", "B", "A", "F", "G", "K", "M"}

//=====================================================================================
// Stringer interface
//=====================================================================================

func (s StarClass) String() string {
	return starClassNames[s]
}

//=====================================================================================
// Valuer interface
//=====================================================================================

func (s StarClass) Value() (driver.Value, error) {
	return starClassNames[s], nil
}

//=====================================================================================
// Scanner interface
//=====================================================================================

func (s *StarClass) Scan(value interface{}) error {
	// If value is nil set size to unknown
	if value == nil {
		*s = UnknownStarClass
		return nil
	}

	// Convert the value to a string
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.([]byte); ok {
			// Parse the value of v
			switch string(v) {
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

//=====================================================================================
// JSON Marshaler interface
//=====================================================================================

func (s StarClass) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

//=====================================================================================
// JSON Unmarshaler interface
//=====================================================================================

func (s *StarClass) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToUpper(strings.TrimSpace(sv))
	return s.Scan(sv)
}

//=====================================================================================
// Nullable Type
//=====================================================================================

type NullStarClass struct {
	StarClass StarClass
	Valid     bool // Valid is true if StarClass is not NULL
}

func (p *NullStarClass) Scan(value any) (err error) {
	if value == nil {
		p.StarClass, p.Valid = UnknownStarClass, false
		return nil
	}

	p.Valid = true
	return p.StarClass.Scan(value)
}

func (p NullStarClass) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return p.StarClass.Value()
}
