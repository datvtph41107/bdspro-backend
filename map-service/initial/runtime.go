package initial

import (
	"fmt"
	"log"
	"map/config"
	"map/env"
)

// InitializeRuntime initializes the runtime environment
func InitializeRuntime() error {
	// Load environment variables
	env.LoadEnvironment()

	// Load configuration
	if _, err := config.LoadConfig(); err != nil {
		return fmt.Errorf("initialize map config: %w", err)
	}

	log.Println("Runtime initialized successfully")
	return nil
}
