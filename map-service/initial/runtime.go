package initial

import (
	"log"
	"map/config"
	"map/env"
)

// InitializeRuntime initializes the runtime environment
func InitializeRuntime() {
	// Load environment variables
	env.LoadEnvironment()
	
	// Load configuration
	config.LoadConfig()
	
	log.Println("Runtime initialized successfully")
}
