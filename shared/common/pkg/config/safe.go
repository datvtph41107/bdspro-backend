// shared/pkg/config/safe.go
// Giải quyết vấn đề: muốn log config để debug, không muốn lộ secrets
package config

import (
	"reflect"
)

// SafeString: interface để struct tự biết cách hiển thị an toàn
// Bất kỳ config struct nào implement interface này sẽ được dùng khi log
type SafeStringer interface {
	SafeString() string
}

// MaskSensitive: tự động mask các field có tag `sensitive:"true"`
// Không cần developer nhớ từng field nào là secret — struct tự khai báo
//
// Dùng reflection vì chúng ta không biết trước cấu trúc của từng service config.
// Đây là trường hợp hợp lý để dùng reflection — generic behavior trên arbitrary structs.
func MaskSensitive(cfg interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	maskRecursive(reflect.ValueOf(cfg), reflect.TypeOf(cfg), result, "")
	return result
}

func maskRecursive(v reflect.Value, t reflect.Type, result map[string]interface{}, prefix string) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		mapKey := field.Tag.Get("mapstructure")
		if mapKey == "" || mapKey == ",squash" {
			// Squash: recurse vào embedded struct
			if field.Anonymous || mapKey == ",squash" {
				maskRecursive(fieldVal, field.Type, result, prefix)
			}
			continue
		}

		key := prefix + mapKey

		if field.Tag.Get("sensitive") == "true" {
			// Mask nhưng cho biết có giá trị hay không
			if fieldVal.String() != "" {
				result[key] = "***SET***"
			} else {
				result[key] = "***NOT SET***" // Cảnh báo: field bắt buộc mà không có giá trị
			}
		} else if fieldVal.Kind() == reflect.Struct {
			sub := make(map[string]interface{})
			maskRecursive(fieldVal, field.Type, sub, "")
			result[key] = sub
		} else {
			result[key] = fieldVal.Interface()
		}
	}
}
