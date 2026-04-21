package optional

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Slice[T any] struct {
	state State
	value []T
}

func AbsentSlice[T any]() Slice[T]   { return Slice[T]{state: StateAbsent} }
func NullSlice[T any]() Slice[T]     { return Slice[T]{state: StateNull} }
func NewSlice[T any](v []T) Slice[T] { return Slice[T]{state: StatePresent, value: v} }

func FromSlicePtr[T any](v *[]T) Slice[T] {
	if v == nil {
		return AbsentSlice[T]()
	}
	return NewSlice(*v)
}

func (s Slice[T]) State() State    { return s.state }
func (s Slice[T]) IsAbsent() bool  { return s.state == StateAbsent }
func (s Slice[T]) IsNull() bool    { return s.state == StateNull }
func (s Slice[T]) IsPresent() bool { return s.state == StatePresent }
func (s Slice[T]) Interface() interface{} {
	if !s.IsPresent() {
		return nil
	}
	return s.value
}

func (s Slice[T]) Value() []T {
	if !s.IsPresent() {
		panic("optional.Slice: Value called on non-present value")
	}
	return s.value
}

func (s Slice[T]) OrElse(def []T) []T {
	if s.IsPresent() {
		return s.value
	}
	return def
}

func (s Slice[T]) Ptr() *[]T {
	if !s.IsPresent() {
		return nil
	}
	return &s.value
}

func (s Slice[T]) ApplyTo(fields map[string]interface{}, col string) {
	switch s.state {
	case StatePresent:
		fields[col] = s.value
	case StateNull:
		fields[col] = nil
	}
}

// JSON
func (s Slice[T]) MarshalJSON() ([]byte, error) {
	switch s.state {
	case StateAbsent, StateNull:
		return []byte("null"), nil
	default:
		return json.Marshal(s.value)
	}
}

func (s *Slice[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.state = StateNull
		return nil
	}
	var v []T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s.state = StatePresent
	s.value = v
	return nil
}

// SQL: lưu dưới dạng JSON
func (s Slice[T]) SQLValue() (driver.Value, error) {
	if !s.IsPresent() {
		return nil, nil
	}
	return json.Marshal(s.value)
}

func (s *Slice[T]) Scan(value interface{}) error {
	if value == nil {
		s.state = StateNull
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into optional.Slice", value)
	}
	var slice []T
	if err := json.Unmarshal(bytes, &slice); err != nil {
		return err
	}
	s.state = StatePresent
	s.value = slice
	return nil
}
