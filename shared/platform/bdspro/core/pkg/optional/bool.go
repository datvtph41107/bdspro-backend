// pkg/optional/bool.go
package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Bool là optional bool với 3 trạng thái
type Bool struct {
	state State
	value bool
}

func AbsentBool() Bool {
	return Bool{state: StateAbsent}
}

func NullBool() Bool {
	return Bool{state: StateNull}
}

func NewBool(v bool) Bool {
	return Bool{state: StatePresent, value: v}
}

func FromBoolPtr(v *bool) Bool {
	if v == nil {
		return AbsentBool()
	}
	return NewBool(*v)
}

func FromBoolWrapper(v *wrapperspb.BoolValue) Bool {
	if v == nil {
		return AbsentBool()
	}
	return NewBool(v.Value)
}

func FromNullBool(v sql.NullBool) Bool {
	if !v.Valid {
		return NullBool()
	}
	return NewBool(v.Bool)
}

func (s Bool) ApplyTo(fields map[string]interface{}, col string) {
	switch s.state {
	case StatePresent:
		fields[col] = s.value
	case StateNull:
		fields[col] = nil
	}
}

func (b Bool) State() State    { return b.state }
func (b Bool) IsAbsent() bool  { return b.state == StateAbsent }
func (b Bool) IsNull() bool    { return b.state == StateNull }
func (b Bool) IsPresent() bool { return b.state == StatePresent }
func (b Bool) Interface() interface{} {
	if !b.IsPresent() {
		return nil
	}
	return b.value
}

func (b Bool) Value() bool {
	if !b.IsPresent() {
		panic("optional.Bool: Value called on non-present value")
	}
	return b.value
}

func (b Bool) OrElse(def bool) bool {
	if b.IsPresent() {
		return b.value
	}
	return def
}

func (b Bool) Ptr() *bool {
	if !b.IsPresent() {
		return nil
	}
	return &b.value
}

func (b Bool) Map(fn func(bool) bool) Bool {
	if !b.IsPresent() {
		return b
	}
	return NewBool(fn(b.value))
}

// JSON
func (b Bool) MarshalJSON() ([]byte, error) {
	switch b.state {
	case StateAbsent, StateNull:
		return []byte("null"), nil
	case StatePresent:
		return json.Marshal(b.value)
	}
	return []byte("null"), nil
}

func (b *Bool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		b.state = StateNull
		return nil
	}

	var v bool
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	b.state = StatePresent
	b.value = v
	return nil
}

// SQL
func (b Bool) SQLValue() (driver.Value, error) {
	if !b.IsPresent() {
		return nil, nil
	}
	return b.value, nil
}

func (b *Bool) Scan(value interface{}) error {
	if value == nil {
		b.state = StateNull
		return nil
	}

	switch v := value.(type) {
	case bool:
		b.state = StatePresent
		b.value = v
	case int64:
		b.state = StatePresent
		b.value = v != 0
	default:
		return fmt.Errorf("cannot scan %T into optional.Bool", value)
	}

	return nil
}
