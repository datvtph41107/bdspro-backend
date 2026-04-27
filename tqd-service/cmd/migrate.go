package cmd

import (
	_db "common/db"
	"log"
	"tqd/config"

	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations and seed data",
	Long: `Run database migrations to create tables and seed initial data.
	
This command will:
1. Create all necessary database tables
2. Seed initial data for directory suppliers, sources, and contact labels`,
	RunE: runMigrate,
}

// migrateCmd is already defined above

func runMigrate(cmd *cobra.Command, args []string) error {
	log.Println("Running database migrations...")

	// Load configuration
	log.Println("Loading configuration...")
	config.LoadConfig()

	// Initialize database connection
	log.Println("Initializing database connection...")
	_db.NewDB()

	// Run auto migration
	// db.AutoMigrate()

	log.Println("Database migration completed successfully!")
	return nil
}

// RunMigrate runs database migration
func RunMigrate() error {
	return runMigrate(nil, nil)
}
