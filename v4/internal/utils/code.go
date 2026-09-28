package utils

import "fmt"

func GenerateReferral(authId uint64) string {
	return "REF" + fmt.Sprintf("%06d", authId)
}
