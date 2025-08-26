package main

import (
	"testing"
)

func TestMain(m *testing.M) {
	// Setup trước khi chạy test
	setupTestEnvironment()

	// Chạy tất cả test
	m.Run()

	// Cleanup sau khi test xong
	cleanupTestEnvironment()
}

func setupTestEnvironment() {
	// Khởi tạo môi trường test
	// Có thể start các service cần thiết ở đây
}

func cleanupTestEnvironment() {
	// Dọn dẹp sau khi test
}
