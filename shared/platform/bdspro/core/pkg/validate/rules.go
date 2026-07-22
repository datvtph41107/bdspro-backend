package validate

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

func Required(v *Violations, field string, val *string) {
	if v == nil {
		return
	}
	if val == nil || strings.TrimSpace(*val) == "" {
		v.Add(field, CodeRequired, "required")
	}
}

func RequiredID(v *Violations, field string, val *uint64) {
	if v == nil {
		return
	}
	if val == nil || *val == 0 {
		v.Add(field, CodeRequired, "required")
	}
}

func RequiredInt32(v *Violations, field string, val *int32) {
	if v == nil {
		return
	}
	if val == nil {
		v.Add(field, CodeRequired, "required")
	}
}

func StringMaxLen(v *Violations, field string, val *string, max int) {
	if v == nil || val == nil {
		return
	}
	n := utf8.RuneCountInString(*val)
	if n > max {
		v.Addf(field, CodeTooLong, "max length is %d, got %d", max, n)
	}
}

func StringMinLen(v *Violations, field string, val *string, min int) {
	if v == nil || val == nil {
		return
	}
	n := utf8.RuneCountInString(*val)
	if n < min {
		v.Addf(field, CodeTooShort, "min length is %d, got %d", min, n)
	}
}

func StringBetween(v *Violations, field string, val *string, min, max int) {
	if v == nil || val == nil {
		return
	}
	n := utf8.RuneCountInString(*val)
	if n < min || n > max {
		v.Addf(field, CodeInvalidValue, "length must be between %d and %d, got %d", min, max, n)
	}
}

func StringOneOf(v *Violations, field string, val *string, allowed []string) {
	if v == nil || val == nil {
		return
	}
	for _, a := range allowed {
		if *val == a {
			return
		}
	}
	v.Addf(field, CodeInvalidValue, "must be one of [%s], got %q", strings.Join(allowed, ", "), *val)
}

func Int32OneOf(v *Violations, field string, val *int32, allowed []int32) {
	if v == nil || val == nil {
		return
	}
	for _, a := range allowed {
		if *val == a {
			return
		}
	}

	buf := make([]string, len(allowed))
	for i, a := range allowed {
		buf[i] = fmt.Sprintf("%d", a)
	}
	v.Addf(field, CodeInvalidValue, "must be one of [%s], got %d", strings.Join(buf, ", "), *val)
}

func Int32Between(v *Violations, field string, val *int32, min, max int32) {
	if v == nil || val == nil {
		return
	}
	if *val < min || *val > max {
		v.Addf(field, CodeOutOfRange, "must be between %d and %d, got %d", min, max, *val)
	}
}

func Float64NonNegative(v *Violations, field string, val *float64) {
	if v == nil || val == nil {
		return
	}
	if *val < 0 {
		v.Addf(field, CodeOutOfRange, "must be >= 0, got %v", *val)
	}
}

func Float64Between(v *Violations, field string, val *float64, min, max float64) {
	if v == nil || val == nil {
		return
	}
	if *val < min || *val > max {
		v.Addf(field, CodeOutOfRange, "must be between %v and %v, got %v", min, max, *val)
	}
}

func Uint32Positive(v *Violations, field string, val *uint32) {
	if v == nil || val == nil {
		return
	}
	if *val == 0 {
		v.Add(field, CodeRequired, "must be > 0")
	}
}

func DateStringOpt(v *Violations, field string, val *string) {
	if v == nil || val == nil || *val == "" {
		return
	}
	if _, err := time.Parse(time.RFC3339, *val); err != nil {
		v.Add(field, CodeInvalidDate, `invalid date, expected RFC3339 (e.g. "2024-01-15T00:00:00Z")`)
	}
}

func DateStringRequired(v *Violations, field string, val *string) {
	if v == nil {
		return
	}
	if val == nil || *val == "" {
		v.Add(field, CodeRequired, "required")
		return
	}
	if _, err := time.Parse(time.RFC3339, *val); err != nil {
		v.Add(field, CodeInvalidDate, `invalid date, expected RFC3339 (e.g. "2024-01-15T00:00:00Z")`)
	}
}

func DateOrder(v *Violations, startField, endField string, start, end *string) {
	if v == nil || start == nil || end == nil || *start == "" || *end == "" {
		return
	}

	s, err1 := time.Parse(time.RFC3339, *start)
	e, err2 := time.Parse(time.RFC3339, *end)
	if err1 != nil || err2 != nil {
		return
	}
	if !s.Before(e) {
		v.Addf(startField, CodeInvalidValue, "must be before %s", endField)
	}
}

func Latitude(v *Violations, field string, val *float64) {
	Float64Between(v, field, val, -90, 90)
}

func Longitude(v *Violations, field string, val *float64) {
	Float64Between(v, field, val, -180, 180)
}

func URLOpt(v *Violations, field string, val *string) {
	if v == nil || val == nil || *val == "" {
		return
	}
	if !strings.HasPrefix(*val, "http://") && !strings.HasPrefix(*val, "https://") {
		v.Add(field, CodeInvalidValue, "must start with http:// or https://")
	}
}

func MediaType(v *Violations, field string, val *string) {
	StringOneOf(v, field, val, []string{"image", "video"})
}
