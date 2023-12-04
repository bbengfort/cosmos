package enums

import (
	"database/sql/driver"
	"encoding/json"
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

//=====================================================================================
// Stringer interface
//=====================================================================================

func (s Size) String() string {
	return sizeNames[s]
}

//=====================================================================================
// Valuer interface
//=====================================================================================

func (s Size) Value() (driver.Value, error) {
	return sizeNames[s], nil
}

//=====================================================================================
// Scanner interface
//=====================================================================================

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

//=====================================================================================
// JSON Marshaler interface
//=====================================================================================

func (s Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
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

//=====================================================================================
// Nullable Type
//=====================================================================================

type NullSize struct {
	Size  Size
	Valid bool // Valid is true if Size is not NULL
}

func (p *NullSize) Scan(value any) (err error) {
	if value == nil {
		p.Size, p.Valid = UnknownSize, false
		return nil
	}

	p.Valid = true
	return p.Size.Scan(value)
}

func (p NullSize) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return p.Size.Value()
}
