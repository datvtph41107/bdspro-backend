package seeders

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Seeder is one deterministic data responsibility owned by the User database.
// Scenario identities, operator bootstrap, and migration-owned reference data
// intentionally do not implement this interface.
type Seeder interface {
	Name() string
	Seed(context.Context, *gorm.DB) error
}

type Result struct {
	SeederCount int
}

// DatabaseSeeders is the Go equivalent of Laravel's DatabaseSeeder call list.
//
// The current User schema has no additional default dataset whose authority is
// not already owned by versioned migrations. Keep this list explicit rather
// than inventing development/acceptance users as database data. Add a seeder
// here only when a real User-owned dataset exists independently of migrations,
// test fixtures, and operator bootstrap.
func DatabaseSeeders() []Seeder {
	return nil
}

func RunDatabaseSeeder(ctx context.Context, db *gorm.DB) (Result, error) {
	if db == nil {
		return Result{}, errors.New("User seed database is not configured")
	}

	registered := DatabaseSeeders()
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, seeder := range registered {
			if seeder == nil {
				return errors.New("User DatabaseSeeder contains a nil seeder")
			}
			if err := seeder.Seed(ctx, tx); err != nil {
				return fmt.Errorf("seed %s: %w", seeder.Name(), err)
			}
		}
		return nil
	})
	if err != nil {
		return Result{}, err
	}

	return Result{SeederCount: len(registered)}, nil
}
