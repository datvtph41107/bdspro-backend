package utils

import "fmt"

func ConvertToBDS(number uint64) string {
	if number == 0 {
		return ""
	}
	return fmt.Sprintf("BDS%06d", number)
}
