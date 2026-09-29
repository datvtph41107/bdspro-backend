package _util

import "fmt"

func GetProductCode(productId uint64) string {
	return fmt.Sprintf("SP%06d", productId)
}
