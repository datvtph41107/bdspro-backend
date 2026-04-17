package config

import (
	"encoding/json"
	"reflect"
)

// MaskSensitive returns a JSON string with sensitive fields replaced by "***".
func MaskSensitive(cfg interface{}) string {
	m := maskRecursive(reflect.ValueOf(cfg), reflect.TypeOf(cfg))
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}

func maskRecursive(v reflect.Value, t reflect.Type) interface{} {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
		t = t.Elem()
	}
	if v.Kind() != reflect.Struct {
		return v.Interface()
	}
	result := map[string]interface{}{}
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)
		key := field.Tag.Get("mapstructure")
		if key == "" || key == ",squash" {
			sub := maskRecursive(fieldVal, field.Type)
			if subMap, ok := sub.(map[string]interface{}); ok {
				for k, v := range subMap {
					result[k] = v
				}
			}
			continue
		}
		if field.Tag.Get("sensitive") == "true" {
			if fieldVal.Kind() == reflect.String && fieldVal.String() != "" {
				result[key] = "***SET***"
			} else {
				result[key] = "***NOT SET***"
			}
		} else {
			result[key] = maskRecursive(fieldVal, field.Type)
		}
	}
	return result
}
