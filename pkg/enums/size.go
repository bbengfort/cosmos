package enums

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

// Size indicates the game size, which defines the maximum number of players, the
// maximum number of turns, as well as the the number of systems that are generated in
// the galaxy. The size of a galaxy cannot be changed once it is created.
type Size uint8

const (
	UnknownSize Size = iota
	Small
	Medium
	Large
	Galactic
	Cosmic
)

var sizeNames = [6]string{"unknown", "small", "medium", "large", "galactic", "cosmic"}

//=====================================================================================
// Size Specific Methods
//=====================================================================================

// MaxPlayers returns a constant number of players based on the size.
func (s Size) MaxPlayers() int16 {
	switch s {
	case UnknownSize:
		return 0
	case Small:
		return 2
	case Medium:
		return 10
	case Large:
		return 20
	case Galactic:
		return 50
	case Cosmic:
		return 100
	default:
		panic(fmt.Errorf("unknown size %v", s))
	}
}

var (
	minSystems = [6]int{0, 20, 100, 200, 500, 1000}
	maxSystems = [6]int{0, 40, 200, 400, 1000, 2000}
)

// NumSystems returns a random number of systems based on a range given by the size.
// Multiple calls to this function will return different numbers of systems bounded by
// the size of the galaxy.
func (s Size) NumSystems() int {
	mins, maxs := minSystems[s], maxSystems[s]
	return random.Intn(maxs-mins+1) + mins
}

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
		if v, ok := sv.([]byte); ok {
			// Parse the value of v
			switch string(v) {
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
