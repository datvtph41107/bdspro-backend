package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DecodeJSON decodes the request body into dst.
func DecodeJSON(r *http.Request, dst any) error {
	if r == nil || r.Body == nil {
		return fmt.Errorf("request body is required")
	}
	return json.NewDecoder(r.Body).Decode(dst)
}

// DecodeQuery binds URL query parameters into dst using `form` struct tags
// (same tag convention previously used by Gin ShouldBindQuery).
func DecodeQuery(r *http.Request, dst any) error {
	if r == nil {
		return fmt.Errorf("request is required")
	}
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("dst must be a non-nil pointer")
	}
	return bindQuery(v.Elem(), r.URL.Query())
}

func bindQuery(v reflect.Value, query map[string][]string) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("dst must point to a struct")
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		valueField := v.Field(i)
		if !valueField.CanSet() {
			continue
		}
		queryKey := field.Tag.Get("form")
		if queryKey == "" {
			queryKey = field.Tag.Get("parser")
		}
		if queryKey == "" || queryKey == "-" {
			continue
		}
		if comma := strings.IndexByte(queryKey, ','); comma >= 0 {
			queryKey = queryKey[:comma]
		}
		rawValues, ok := query[queryKey]
		if !ok || len(rawValues) == 0 {
			continue
		}
		raw := strings.TrimSpace(rawValues[0])
		if raw == "" && valueField.Kind() != reflect.String {
			continue
		}
		if err := setQueryField(valueField, raw, rawValues); err != nil {
			return fmt.Errorf("invalid %s: %w", queryKey, err)
		}
	}
	return nil
}

func setQueryField(valueField reflect.Value, raw string, rawValues []string) error {
	if valueField.Kind() == reflect.Ptr {
		if valueField.IsNil() {
			valueField.Set(reflect.New(valueField.Type().Elem()))
		}
		return setQueryField(valueField.Elem(), raw, rawValues)
	}

	if valueField.Type() == reflect.TypeOf(time.Time{}) {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return err
		}
		valueField.Set(reflect.ValueOf(parsed))
		return nil
	}

	switch valueField.Kind() {
	case reflect.String:
		valueField.SetString(raw)
	case reflect.Bool:
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		valueField.SetBool(parsed)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(raw, 10, valueField.Type().Bits())
		if err != nil {
			return err
		}
		valueField.SetInt(parsed)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(raw, 10, valueField.Type().Bits())
		if err != nil {
			return err
		}
		valueField.SetUint(parsed)
	case reflect.Float32, reflect.Float64:
		parsed, err := strconv.ParseFloat(raw, valueField.Type().Bits())
		if err != nil {
			return err
		}
		valueField.SetFloat(parsed)
	case reflect.Slice:
		parts := rawValues
		if len(parts) == 1 && strings.Contains(parts[0], ",") {
			parts = strings.Split(parts[0], ",")
		}
		slice := reflect.MakeSlice(valueField.Type(), 0, len(parts))
		elemType := valueField.Type().Elem()
		for _, part := range parts {
			part = strings.TrimSpace(part)
			elem := reflect.New(elemType).Elem()
			if err := setQueryField(elem, part, []string{part}); err != nil {
				return err
			}
			slice = reflect.Append(slice, elem)
		}
		valueField.Set(slice)
	default:
		return fmt.Errorf("unsupported kind %s", valueField.Kind())
	}
	return nil
}
