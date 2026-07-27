package config

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// loadBackendLocalEnvFile nạp biến từ .env.backend.local (monorepo root) khi chạy gateway local.
// Không ghi đè biến đã có trong môi trường.
func loadBackendLocalEnvFile() {
	candidates := []string{
		".env.backend.local",
		filepath.Join("..", ".env.backend.local"),
		filepath.Join("..", "..", ".env.backend.local"),
	}

	for _, path := range candidates {
		if loadEnvFileIfExists(path) {
			log.Printf("Loaded local env from %s", path)
			return
		}
	}
}

func loadEnvFileIfExists(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		val = strings.Trim(val, `"'`)
		if key == "" {
			continue
		}
		if os.Getenv(key) != "" {
			continue
		}
		_ = os.Setenv(key, val)
	}
	return true
}
