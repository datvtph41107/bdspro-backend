package seo_strategy

import (
	"errors"
	"regexp"
	"strings"
)

var strategyCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{2,119}$`)

func NormalizeStrategyCode(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, "-", "_")
	value = strings.Join(strings.Fields(value), "_")
	return value
}

func ValidateStrategyCode(value string) error {
	if !strategyCodePattern.MatchString(value) {
		return errors.New("strategy code phải dùng lower_snake_case và dài từ 3 đến 120 ký tự")
	}
	return nil
}
