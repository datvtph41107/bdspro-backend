package validate

import "fmt"

type Violation struct {
	Field   string
	Code    Code
	Message string
}

type Violations struct {
	items []Violation
}

func (v *Violations) Add(field string, code Code, message string) {
	if v == nil {
		return
	}
	v.items = append(v.items, Violation{
		Field:   field,
		Code:    code,
		Message: message,
	})
}

func (v *Violations) Addf(field string, code Code, format string, args ...any) {
	if v == nil {
		return
	}
	v.Add(field, code, fmt.Sprintf(format, args...))
}

func (v *Violations) Merge(other *Violations) {
	if v == nil || other == nil || !other.HasErrors() {
		return
	}
	v.items = append(v.items, other.items...)
}

func (v *Violations) HasErrors() bool {
	return v != nil && len(v.items) > 0
}

func (v *Violations) Len() int {
	if v == nil {
		return 0
	}
	return len(v.items)
}

func (v *Violations) List() []Violation {
	if v == nil || len(v.items) == 0 {
		return nil
	}
	out := make([]Violation, len(v.items))
	copy(out, v.items)
	return out
}

func (v *Violations) Clear() {
	if v == nil {
		return
	}
	v.items = nil
}

func (v *Violations) Error() string {
	if !v.HasErrors() {
		return ""
	}

	out := "validation:"
	for _, item := range v.items {
		out += " " + item.Field + ": " + item.Message + ";"
	}
	return out
}

func (v *Violations) AsError() error {
	if !v.HasErrors() {
		return nil
	}
	return v
}
