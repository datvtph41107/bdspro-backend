package base_util

import "strings"

func ClearPhone(phone string) string {
	return strings.ReplaceAll(phone, " ", "")
}
