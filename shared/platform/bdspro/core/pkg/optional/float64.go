// pkg/optional/float64.go
package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Float64 là optional float64 với 3 trạng thái
type Float64 struct {
	state State
	value float64
}

func AbsentFloat64() Float64 {
	return Float64{state: StateAbsent}
}

func NullFloat64() Float64 {
	return Float64{state: StateNull}
}

func NewFloat64(v float64) Float64 {
	return Float64{state: StatePresent, value: v}
}

func FromFloat64Ptr(v *float64) Float64 {
	if v == nil {
		return AbsentFloat64()
	}
	return NewFloat64(*v)
}

func FromFloat64Wrapper(v *wrapperspb.DoubleValue) Float64 {
	if v == nil {
		return AbsentFloat64()
	}
	return NewFloat64(v.Value)
}

func FromNullFloat64(v sql.NullFloat64) Float64 {
	if !v.Valid {
		return NullFloat64()
	}
	return NewFloat64(v.Float64)
}

func (s Float64) ApplyTo(fields map[string]interface{}, col string) {
	switch s.state {
	case StatePresent:
		fields[col] = s.value
	case StateNull:
		fields[col] = nil
	}
}

func (f Float64) State() State    { return f.state }
func (f Float64) IsAbsent() bool  { return f.state == StateAbsent }
func (f Float64) IsNull() bool    { return f.state == StateNull }
func (f Float64) IsPresent() bool { return f.state == StatePresent }
func (f Float64) Interface() interface{} {
	if !f.IsPresent() {
		return nil
	}
	return f.value
}

func (f Float64) Value() float64 {
	if !f.IsPresent() {
		panic("optional.Float64: Value called on non-present value")
	}
	return f.value
}

func (f Float64) OrElse(def float64) float64 {
	if f.IsPresent() {
		return f.value
	}
	return def
}

func (f Float64) Ptr() *float64 {
	if !f.IsPresent() {
		return nil
	}
	return &f.value
}

func (f Float64) Map(fn func(float64) float64) Float64 {
	if !f.IsPresent() {
		return f
	}
	return NewFloat64(fn(f.value))
}

// JSON
func (f Float64) MarshalJSON() ([]byte, error) {
	switch f.state {
	case StateAbsent, StateNull:
		return []byte("null"), nil
	case StatePresent:
		return json.Marshal(f.value)
	}
	return []byte("null"), nil
}

func (f *Float64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		f.state = StateNull
		return nil
	}

	var v float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	f.state = StatePresent
	f.value = v
	return nil
}

// SQL
func (f Float64) SQLValue() (driver.Value, error) {
	if !f.IsPresent() {
		return nil, nil
	}
	return f.value, nil
}

func (f *Float64) Scan(value interface{}) error {
	if value == nil {
		f.state = StateNull
		return nil
	}

	switch v := value.(type) {
	case float64:
		f.state = StatePresent
		f.value = v
	case float32:
		f.state = StatePresent
		f.value = float64(v)
	default:
		return fmt.Errorf("cannot scan %T into optional.Float64", value)
	}

	return nil
}
