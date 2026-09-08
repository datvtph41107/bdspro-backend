package seeders

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	_utils "common/utils"

	"gorm.io/gorm"
)

var (
	acceptanceUsernamePattern = regexp.MustCompile(`^qhpro_acceptance_[0-9]{10,20}$`)
	acceptancePhonePattern    = regexp.MustCompile(`^039[0-9]{7}$`)
	acceptanceEmailPattern    = regexp.MustCompile(`^qhpro-[0-9]{10,20}@e[.]invalid$`)
)

func seedAcceptanceClient(ctx context.Context, db *gorm.DB, options Options) (Result, error) {
	if err := requireNonProduction(options.Environment); err != nil {
		return Result{}, err
	}

	username := strings.TrimSpace(options.AcceptanceUsername)
	phone := strings.TrimSpace(options.AcceptancePhone)
	email := strings.TrimSpace(options.AcceptanceEmail)
	password := strings.TrimSpace(options.AcceptanceUserPassword)
	if len(password) < 12 {
		return Result{}, fmt.Errorf("QHPRO_ACCEPTANCE_USER_PASSWORD must contain at least 12 characters")
	}
	if !acceptanceUsernamePattern.MatchString(username) ||
		!acceptancePhonePattern.MatchString(phone) ||
		!acceptanceEmailPattern.MatchString(email) {
		return Result{}, fmt.Errorf("invalid acceptance identity")
	}

	passwordHash, err := _utils.HashPassword(password)
	if err != nil {
		return Result{}, fmt.Errorf("hash acceptance client password: %w", err)
	}

	result, err := seedClient(ctx, db, clientSeedSpec{
		Username:     username,
		PasswordHash: passwordHash,
		FullName:     "QHPRO Acceptance Commercial User",
		Email:        email,
		Phone:        phone,
		RoleCode:     10,
	})
	if err != nil {
		return Result{}, err
	}
	result.Name = AcceptanceClient
	return result, nil
}
