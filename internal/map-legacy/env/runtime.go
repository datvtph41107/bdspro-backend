package env

import (
	"log"
	"os"
)

// LoadEnvironment loads environment variables
func LoadEnvironment() {
	// Set default timezone
	if tz := os.Getenv("TZ"); tz == "" {
		os.Setenv("TZ", "Asia/Ho_Chi_Minh")
	}

	log.Println("Environment loaded successfully")
}
