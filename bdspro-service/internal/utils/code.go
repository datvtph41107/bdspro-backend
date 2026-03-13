package utils

import "fmt"

func ConvertToBDS(number uint64) string {
	if number >= 1000000 {
		panic("number must be < 1,000,000")
	}

	const mod uint64 = 1000000
	const a uint64 = 7368787 // gcd(a, mod) = 1
	const b uint64 = 123457

	x := (a*number + b) % mod

	return fmt.Sprintf("BDS%06d", x)
}

func GenerateDealCode(id uint64) string {
	return fmt.Sprintf("TV%06d", id)
}

func GenerateTransactionCode(id uint64) string {
	return fmt.Sprintf("TXT%06d", id)
}
