/*
Package jcode is about creating random strings as game join codes.
*/
package jcode

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base32"
	"encoding/json"
	"errors"
	"strings"
)

const (
	b32Alphabet = "0123456789ABCDEFGHJKLMNPQRSTUVWX"
	sep         = "-"
)

// CosmosEncoding uses an all upper case alphabet with digits 0-9 and characters A-X,
// omitting I, O because of their visual similarity to 1 and 0. This is the opposite of
// standard base32 encoding which omits 0, 1, and 8 because of their numeric similarity.
var CosmosEncoding = base32.NewEncoding(b32Alphabet).WithPadding(base32.NoPadding)

type JoinCode string

var (
	ErrInvalidJoinCode = errors.New("invalid join code")
	ErrScanJoinCode    = errors.New("failed to scan join code")
)

// New creates a random 16 character join code.
func New() JoinCode {
	data := make([]byte, 10)
	if _, err := rand.Read(data); err != nil {
		panic(err)
	}

	code := CosmosEncoding.EncodeToString(data)
	return JoinCode(code)
}

// =====================================================================================
// Stringer interface
// =====================================================================================

// String returns a readable representation of the join code with separators.
func (j JoinCode) String() string {
	parts := make([]string, 0, 4)
	for i := 0; i < len(j); i += 4 {
		parts = append(parts, string(j)[i:i+4])
	}
	return strings.Join(parts, sep)
}

//=====================================================================================
// JSON Marshaler interface
//=====================================================================================

// Marshal a join code as JSON
func (j JoinCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.String())
}

//=====================================================================================
// JSON Unmarshaler interface
//=====================================================================================

// Unmarshal a join code from JSON
func (j *JoinCode) UnmarshalJSON(data []byte) (err error) {
	// Unmarshal the join code as a string
	var js string
	if err = json.Unmarshal(data, &js); err != nil {
		return ErrInvalidJoinCode
	}

	// Make upper case, trim any whitespace, and remove the separator character
	js = strings.Replace(strings.ToUpper(strings.TrimSpace(js)), sep, "", -1)
	if len(js) != 16 {
		return ErrInvalidJoinCode
	}

	if _, err = CosmosEncoding.DecodeString(js); err != nil {
		return ErrInvalidJoinCode
	}

	*j = JoinCode(js)
	return nil
}

//=====================================================================================
// Valuer interface
//=====================================================================================

func (j JoinCode) Value() (driver.Value, error) {
	return string(j), nil
}

//=====================================================================================
// Scanner interface
//=====================================================================================

func (j *JoinCode) Scan(value interface{}) error {
	// If value is nil, empty string
	if value == nil {
		*j = JoinCode("")
		return nil
	}

	if bv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := bv.(string); ok {
			*j = JoinCode(v)
			return nil
		}
	}

	return ErrScanJoinCode
}

//=====================================================================================
// Nullable Type
//=====================================================================================

type NullJoinCode struct {
	JoinCode JoinCode
	Valid    bool
}

func (j *NullJoinCode) Scan(value any) (err error) {
	// If value is nil then this is a null valued join code.
	if value == nil {
		j.JoinCode = JoinCode("")
		j.Valid = false
		return nil
	}

	// Scan the internal join code if the value is not null.
	j.Valid = true
	return j.JoinCode.Scan(value)
}

func (j NullJoinCode) Value() (driver.Value, error) {
	// If not valid then this value is null
	if !j.Valid {
		return nil, nil
	}
	return j.JoinCode.Value()
}
