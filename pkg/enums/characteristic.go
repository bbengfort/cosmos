package enums

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

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

//=====================================================================================
// Stringer interface
//=====================================================================================

func (c Characteristic) String() string {
	return characteristicNames[c]
}

//=====================================================================================
// Valuer interface
//=====================================================================================

func (c Characteristic) Value() (driver.Value, error) {
	return characteristicNames[c], nil
}

//=====================================================================================
// Scanner interface
//=====================================================================================

func (c *Characteristic) Scan(value interface{}) error {
	// If value is nil set size to unknown
	if value == nil {
		*c = UnknownCharacteristic
		return nil
	}

	// Convert the value to a string
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.([]byte); ok {
			// Parse the value of v
			switch string(v) {
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

//=====================================================================================
// JSON Marshaler interface
//=====================================================================================

func (c Characteristic) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

//=====================================================================================
// JSON Unmarshaler interface
//=====================================================================================

func (c *Characteristic) UnmarshalJSON(data []byte) (err error) {
	var sv string
	if err = json.Unmarshal(data, &sv); err != nil {
		return err
	}

	sv = strings.ToLower(strings.TrimSpace(sv))
	return c.Scan(sv)
}

//=====================================================================================
// Nullable Type
//=====================================================================================

type NullCharacteristic struct {
	Characteristic Characteristic
	Valid          bool // Valid is true if Characteristic is not NULL
}

func (p *NullCharacteristic) Scan(value any) (err error) {
	if value == nil {
		p.Characteristic, p.Valid = UnknownCharacteristic, false
		return nil
	}

	p.Valid = true
	return p.Characteristic.Scan(value)
}

func (p NullCharacteristic) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}
	return p.Characteristic.Value()
}
