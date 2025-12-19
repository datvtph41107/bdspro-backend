package utils

import "fmt"

func ConvertToBDS(number uint64) string {
	if number == 0 {
		return ""
	}
	return fmt.Sprintf("BDS%06d", number)
}

func GenerateDealCode(id uint64) string {
	return fmt.Sprintf("TV%06d", id)
}

func GenerateTransactionCode(id uint64) string {
	return fmt.Sprintf("TXT%06d", id)
}
