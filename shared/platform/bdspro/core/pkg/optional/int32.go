// pkg/optional/int32.go
package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Int32 là optional int32 với 3 trạng thái
type Int32 struct {
	state State
	value int32
}

func AbsentInt32() Int32 {
	return Int32{state: StateAbsent}
}

func NullInt32() Int32 {
	return Int32{state: StateNull}
}

func NewInt32(v int32) Int32 {
	return Int32{state: StatePresent, value: v}
}

func FromInt32Ptr(v *int32) Int32 {
	if v == nil {
		return AbsentInt32()
	}
	return NewInt32(*v)
}

func FromInt32Wrapper(v *wrapperspb.Int32Value) Int32 {
	if v == nil {
		return AbsentInt32()
	}
	return NewInt32(v.Value)
}

func FromNullInt32(v sql.NullInt32) Int32 {
	if !v.Valid {
		return NullInt32()
	}
	return NewInt32(v.Int32)
}

func (s Int32) ApplyTo(fields map[string]interface{}, col string) {
	switch s.state {
	case StatePresent:
		fields[col] = s.value
	case StateNull:
		fields[col] = nil
	}
}

func (i Int32) State() State    { return i.state }
func (i Int32) IsAbsent() bool  { return i.state == StateAbsent }
func (i Int32) IsNull() bool    { return i.state == StateNull }
func (i Int32) IsPresent() bool { return i.state == StatePresent }
func (i Int32) Interface() interface{} {
	if !i.IsPresent() {
		return nil
	}
	return i.value
}

func (i Int32) Value() int32 {
	if !i.IsPresent() {
		panic("optional.Int32: Value called on non-present value")
	}
	return i.value
}

func (i Int32) OrElse(def int32) int32 {
	if i.IsPresent() {
		return i.value
	}
	return def
}

func (i Int32) Ptr() *int32 {
	if !i.IsPresent() {
		return nil
	}
	return &i.value
}

func (i Int32) Map(fn func(int32) int32) Int32 {
	if !i.IsPresent() {
		return i
	}
	return NewInt32(fn(i.value))
}

// JSON
func (i Int32) MarshalJSON() ([]byte, error) {
	switch i.state {
	case StateAbsent, StateNull:
		return []byte("null"), nil
	case StatePresent:
		return json.Marshal(i.value)
	}
	return []byte("null"), nil
}

func (i *Int32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		i.state = StateNull
		return nil
	}

	var v int32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	i.state = StatePresent
	i.value = v
	return nil
}

// SQL
func (i Int32) SQLValue() (driver.Value, error) {
	if !i.IsPresent() {
		return nil, nil
	}
	return int64(i.value), nil
}

func (i *Int32) Scan(value interface{}) error {
	if value == nil {
		i.state = StateNull
		return nil
	}

	switch v := value.(type) {
	case int64:
		i.state = StatePresent
		i.value = int32(v)
	case int32:
		i.state = StatePresent
		i.value = v
	default:
		return fmt.Errorf("cannot scan %T into optional.Int32", value)
	}

	return nil
}
