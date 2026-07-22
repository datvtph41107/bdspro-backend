package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Uint64 struct {
	state State
	value uint64
}

func AbsentUint64() Uint64      { return Uint64{state: StateAbsent} }
func NullUint64() Uint64        { return Uint64{state: StateNull} }
func NewUint64(v uint64) Uint64 { return Uint64{state: StatePresent, value: v} }

func FromUint64Ptr(v *uint64) Uint64 {
	if v == nil {
		return AbsentUint64()
	}
	return NewUint64(*v)
}

func FromNullInt64(v sql.NullInt64) Uint64 {
	if !v.Valid {
		return NullUint64()
	}
	return NewUint64(uint64(v.Int64))
}

func (u Uint64) State() State    { return u.state }
func (u Uint64) IsAbsent() bool  { return u.state == StateAbsent }
func (u Uint64) IsNull() bool    { return u.state == StateNull }
func (u Uint64) IsPresent() bool { return u.state == StatePresent }
func (u Uint64) Interface() interface{} {
	if !u.IsPresent() {
		return nil
	}
	return u.value
}

func (u Uint64) Value() uint64 {
	if !u.IsPresent() {
		panic("optional.Uint64: Value called on non-present value")
	}
	return u.value
}

func (u Uint64) OrElse(def uint64) uint64 {
	if u.IsPresent() {
		return u.value
	}
	return def
}

func (u Uint64) Ptr() *uint64 {
	if !u.IsPresent() {
		return nil
	}
	return &u.value
}

func (u Uint64) ApplyTo(fields map[string]interface{}, col string) {
	switch u.state {
	case StatePresent:
		fields[col] = u.value
	case StateNull:
		fields[col] = nil
	}
}

func (u Uint64) MarshalJSON() ([]byte, error) {
	switch u.state {
	case StateAbsent, StateNull:
		return []byte("null"), nil
	default:
		return json.Marshal(u.value)
	}
}

func (u *Uint64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		u.state = StateNull
		return nil
	}
	var v uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	u.state = StatePresent
	u.value = v
	return nil
}

func (u Uint64) SQLValue() (driver.Value, error) {
	if !u.IsPresent() {
		return nil, nil
	}
	return int64(u.value), nil
}

func (u *Uint64) Scan(value interface{}) error {
	if value == nil {
		u.state = StateNull
		return nil
	}
	switch v := value.(type) {
	case int64:
		u.state = StatePresent
		u.value = uint64(v)
	case uint64:
		u.state = StatePresent
		u.value = v
	default:
		return fmt.Errorf("cannot scan %T into optional.Uint64", value)
	}
	return nil
}
