package optional

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type String struct {
	state State
	value string
}

func AbsentString() String {
	return String{
		state: StateAbsent,
	}
}
func NullString() String {
	return String{
		state: StateNull,
	}
}
func NewString(v string) String {
	return String{
		state: StatePresent,
		value: v,
	}
}

// Từ con trỏ (dùng cho proto optional)
func FromStringPtr(v *string) String {
	if v == nil {
		return AbsentString()
	}
	return NewString(*v)
}

// Từ sql.NullString
func FromNullString(v sql.NullString) String {
	if !v.Valid {
		return NullString()
	}
	return NewString(v.String)
}

func (s String) State() State    { return s.state }
func (s String) IsAbsent() bool  { return s.state == StateAbsent }
func (s String) IsNull() bool    { return s.state == StateNull }
func (s String) IsPresent() bool { return s.state == StatePresent }
func (s String) Interface() interface{} {
	if !s.IsPresent() {
		return nil
	}
	return s.value
}

func (s String) Value() string {
	if !s.IsPresent() {
		panic("optional.String: Value called on non-present value")
	}
	return s.value
}

func (s String) OrElse(def string) string {
	if s.IsPresent() {
		return s.value
	}
	return def
}

func (s String) Ptr() *string {
	if !s.IsPresent() {
		return nil
	}
	return &s.value
}

// ApplyTo thêm vào map nếu không phải Absent
func (s String) ApplyTo(fields map[string]interface{}, col string) {
	switch s.state {
	case StatePresent:
		fields[col] = s.value
	case StateNull:
		fields[col] = nil
	}
}

// JSON
func (s String) MarshalJSON() ([]byte, error) {
	switch s.state {
	case StateAbsent, StateNull:
		return []byte("null"), nil
	default:
		return json.Marshal(s.value)
	}
}

func (s *String) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.state = StateNull
		return nil
	}
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s.state = StatePresent
	s.value = v
	return nil
}

// SQL
func (s String) SQLValue() (driver.Value, error) {
	if !s.IsPresent() {
		return nil, nil
	}
	return s.value, nil
}

func (s *String) Scan(value interface{}) error {
	if value == nil {
		s.state = StateNull
		return nil
	}
	switch v := value.(type) {
	case string:
		s.state = StatePresent
		s.value = v
	case []byte:
		s.state = StatePresent
		s.value = string(v)
	default:
		return fmt.Errorf("cannot scan %T into optional.String", value)
	}
	return nil
}
