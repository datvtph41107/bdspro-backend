// pkg/optional/uint32.go
package optional

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Uint32 là optional uint32 với 3 trạng thái
type Uint32 struct {
	state State
	value uint32
}

func AbsentUint32() Uint32 {
	return Uint32{state: StateAbsent}
}

func NullUint32() Uint32 {
	return Uint32{state: StateNull}
}

func NewUint32(v uint32) Uint32 {
	return Uint32{state: StatePresent, value: v}
}

// Từ pointer (proto optional)
func FromUint32Ptr(v *uint32) Uint32 {
	if v == nil {
		return AbsentUint32()
	}
	return NewUint32(*v)
}

// Từ wrapper (proto wrapper)
func FromUint32Wrapper(v *wrapperspb.UInt32Value) Uint32 {
	if v == nil {
		return AbsentUint32()
	}
	return NewUint32(v.Value)
}

func (s Uint32) ApplyTo(fields map[string]interface{}, col string) {
	switch s.state {
	case StatePresent:
		fields[col] = s.value
	case StateNull:
		fields[col] = nil
	}
}

// Từ sql.NullInt32
// func FromNullInt32(v sql.NullInt32) Uint32 {
// 	if !v.Valid {
// 		return NullUint32()
// 	}
// 	return NewUint32(uint32(v.Int32))
// }

func (u Uint32) State() State    { return u.state }
func (u Uint32) IsAbsent() bool  { return u.state == StateAbsent }
func (u Uint32) IsNull() bool    { return u.state == StateNull }
func (u Uint32) IsPresent() bool { return u.state == StatePresent }
func (u Uint32) Interface() interface{} {
	if !u.IsPresent() {
		return nil
	}
	return u.value
}

func (u Uint32) Value() uint32 {
	if !u.IsPresent() {
		panic("optional.Uint32: Value called on non-present value")
	}
	return u.value
}

func (u Uint32) OrElse(def uint32) uint32 {
	if u.IsPresent() {
		return u.value
	}
	return def
}

func (u Uint32) Ptr() *uint32 {
	if !u.IsPresent() {
		return nil
	}
	return &u.value
}

func (u Uint32) Map(fn func(uint32) uint32) Uint32 {
	if !u.IsPresent() {
		return u
	}
	return NewUint32(fn(u.value))
}

// JSON
func (u Uint32) MarshalJSON() ([]byte, error) {
	switch u.state {
	case StateAbsent, StateNull:
		return []byte("null"), nil
	case StatePresent:
		return json.Marshal(u.value)
	}
	return []byte("null"), nil
}

func (u *Uint32) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		u.state = StateNull
		return nil
	}

	var v uint32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	u.state = StatePresent
	u.value = v
	return nil
}

// SQL
func (u Uint32) SQLValue() (driver.Value, error) {
	if !u.IsPresent() {
		return nil, nil
	}
	return int64(u.value), nil
}

func (u *Uint32) Scan(value interface{}) error {
	if value == nil {
		u.state = StateNull
		return nil
	}

	switch v := value.(type) {
	case int64:
		u.state = StatePresent
		u.value = uint32(v)
	case uint32:
		u.state = StatePresent
		u.value = v
	default:
		return fmt.Errorf("cannot scan %T into optional.Uint32", value)
	}

	return nil
}
