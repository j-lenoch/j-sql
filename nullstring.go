package jsql

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Basically a clone of the sql.NullString, but with
// additional functionality like JSON marshalling.
type NullString struct {
	sql.NullString
}

// NewNullString returns a new NullString and ignores any errors.
func NewNullString(i interface{}) NullString {
	n, _ := TryNullString(i)
	return n
}

// NewTrimmedNullString returns a new NullString, but trims all spaces and invalidates
// the object if string is empty.
func NewTrimmedNullString(i interface{}) NullString {
	n, _ := TryNullString(i)

	if n.Valid {
		n.String = strings.TrimSpace(n.String)
		n.Valid = len(n.String) > 0
	}

	return n
}

// TryNullString tries to create a new NullString
func TryNullString(i interface{}) (NullString, error) {
	ns := NullString{}
	return ns, ns.TrySet(i)
}

// Set tries to update the objects value and ignores any errors
func (ns *NullString) Set(i interface{}) {
	ns.TrySet(i)
}

// TrySet tries to update the objects value
func (ns *NullString) TrySet(i interface{}) error {

	if i == nil {
		ns.Valid = false
		return nil
	}

	var val string
	var err error

	switch i.(type) {
	case string:
		val = i.(string)
	case []byte:
		val = string(i.([]byte))
	default:
		err = errors.New(fmt.Sprintf("given value '%s' is not en explicit string: please cast it to ensure that this behaviour is expected", i))
	}

	if err != nil {
		ns.Valid = false
		return err
	}

	ns.String = val
	ns.Valid = true
	return nil
}

// IsTrimmed returns whether the given string is trimmed. Returns true if the string is invalid
func (ns NullString) IsTrimmed() bool {
	return !ns.Valid || ns.String == strings.TrimSpace(ns.String)
}

// IsEmpty checks whether this NullString contains any data
func (ns NullString) IsEmpty() bool {
	return ns.Valid && len(strings.TrimSpace(ns.String)) > 0
}

// ToValue transform the current value into nil or string
func (ns NullString) ToValue() interface{} {
	if !ns.Valid {
		return nil
	}

	return ns.String
}

// MarshalJSON transforms the current string in either "null" or a byte representation of the string
func (ns NullString) MarshalJSON() ([]byte, error) {

	if !ns.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(ns.String)
}

// UnmarshalJSON transforms
func (ns *NullString) UnmarshalJSON(b []byte) error {

	ns.String = ""
	ns.Valid = false

	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	if len(b) >= 0 {

		// Try to extract the 'string'. If this failed we simply
		// use the base value as string.
		if err := json.Unmarshal(b, &ns.String); err != nil {
			ns.String = string(b)
		}

		ns.Valid = true
	}

	return nil
}
