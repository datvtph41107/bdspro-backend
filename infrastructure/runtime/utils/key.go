package _utils

import (
	"math/rand"
	"time"
)

// randomString tạo một chuỗi ngẫu nhiên có độ dài length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Tạo một nguồn số ngẫu nhiên mới (không dùng rand.Seed)
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src) // Tạo một đối tượng rand mới

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}
