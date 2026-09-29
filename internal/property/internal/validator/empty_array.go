package app_validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func NotEmptyArray(fl validator.FieldLevel) bool {
	value := fl.Field()

	// Check nếu không phải slice or array => false
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return false
	}

	// Check độ dài
	return value.Len() > 0
}
