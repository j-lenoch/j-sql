package jsql

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

type NullArray struct {
	Valid bool
	Array []interface{}
}

func NewNullArray(b interface{}) NullArray {

	raw := reflect.ValueOf(b)

	if raw.Kind() != reflect.Slice {

		log.Fatal("Expected a slice, got a ", raw.Kind())

		return NullArray{
			Valid: false,
			Array: nil,
		}
	}

	a := make([]interface{}, raw.Len())
	for i := 0; i < raw.Len(); i++ {
		a[i] = raw.Index(i).Interface()
	}

	return NullArray{
		Array: a,
		Valid: true,
	}
}

func (na NullArray) ToStringArray() []string {
	out := make([]string, len(na.Array))

	for i, val := range na.Array {
		out[i] = fmt.Sprint(val)
	}

	return out
}

func (na NullArray) ToInt64Array() ([]int64, error) {
	out := make([]int64, len(na.Array))

	var err error
	for i, val := range na.Array {
		if i, ok := val.(int64); ok {
			out[i] = i
			continue
		}

		str := fmt.Sprint(val)
		out[i], err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (na NullArray) ToFloat64Array() ([]float64, error) {
	out := make([]float64, len(na.Array))

	var err error
	for i, val := range na.Array {
		if f, ok := val.(float64); ok {
			out[i] = f
			continue
		}
		str := fmt.Sprint(val)
		out[i], err = strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (na NullArray) MarshalJSON() ([]byte, error) {
	if !na.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(na.Array)
}

func (na *NullArray) UnmarshalJSON(b []byte) error {
	na.Valid = false

	if bytes.Equal(b, []byte("null")) {
		return nil
	}

	if len(b) >= 0 {
		if err := json.Unmarshal(b, &na.Array); err != nil {
			return err
		}
		na.Valid = true
	}

	return nil
}

// Retrieve value as JSON data
func (na *NullArray) Scan(value interface{}) error {
	na.Valid = false

	if value == nil {
		return nil
	}

	data, ok := value.([]byte)
	if !ok {
		return errors.New("The given data is not a valid string")
	}

	if len(data) == 0 {
		return nil
	}

	return na.UnmarshalJSON(data)
}

func (na NullArray) Value() (driver.Value, error) {

	if !na.Valid {
		return nil, nil
	}

	data, err := na.MarshalJSON()
	log.Print("Get value ", na, string(data), err)
	return na.MarshalJSON()
}
