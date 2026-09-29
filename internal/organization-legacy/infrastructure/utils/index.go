package utils

import "fmt"

func GenerateDealCode(id uint64) string {
	return fmt.Sprintf("TV%06d", id)
}
