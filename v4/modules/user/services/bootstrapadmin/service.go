// Package bootstrapadmin owns the one-time transition from an empty User IAM
// database to its first root operator. Normal administrator lifecycle remains
// behind the authenticated Admin API.
package bootstrapadmin

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

var (
	ErrAlreadyBootstrapped = errors.New("root operator already exists")
	ErrUsernameInUse       = errors.New("administrator username is already in use")
)

type Input struct {
	Username string
	Password string
	FullName string
	Email    string
	Phone    string
}

type Result struct {
	ProfileID       uint64
	AuthID          uint64
	RoleID          uint64
	PermissionCount int64
}

type Store interface {
	CreateRoot(ctx context.Context, input Input, passwordHash string) (Result, error)
}

type PasswordHasher func(password string) (string, error)

type Service struct {
	store Store
	hash  PasswordHasher
}

func NewService(store Store, hash PasswordHasher) *Service {
	return &Service{store: store, hash: hash}
}

func (s *Service) Bootstrap(ctx context.Context, input Input) (Result, error) {
	input.Username = strings.TrimSpace(input.Username)
	input.FullName = strings.TrimSpace(input.FullName)
	input.Email = strings.TrimSpace(input.Email)
	input.Phone = strings.TrimSpace(input.Phone)

	if s == nil || s.store == nil || s.hash == nil {
		return Result{}, errors.New("root operator bootstrap is not configured")
	}
	if err := validate(input); err != nil {
		return Result{}, err
	}

	passwordHash, err := s.hash(input.Password)
	if err != nil {
		return Result{}, fmt.Errorf("hash root operator password: %w", err)
	}
	return s.store.CreateRoot(ctx, input, passwordHash)
}

func validate(input Input) error {
	if len(input.Username) < 4 || len(input.Username) > 255 {
		return errors.New("username must contain between 4 and 255 characters")
	}
	if strings.ContainsAny(input.Username, " \t\r\n") {
		return errors.New("username must not contain whitespace")
	}
	if len(input.Password) < 12 {
		return errors.New("password must contain at least 12 characters")
	}
	if len(input.FullName) < 2 || len(input.FullName) > 255 {
		return errors.New("full name must contain between 2 and 255 characters")
	}
	if len(input.Email) > 50 {
		return errors.New("email must not exceed 50 characters")
	}
	if input.Email != "" {
		address, err := mail.ParseAddress(input.Email)
		if err != nil || !strings.EqualFold(address.Address, input.Email) {
			return errors.New("email is invalid")
		}
	}
	if len(input.Phone) > 20 {
		return errors.New("phone must not exceed 20 characters")
	}
	return nil
}
