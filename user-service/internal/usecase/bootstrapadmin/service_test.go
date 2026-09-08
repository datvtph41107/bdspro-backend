package bootstrapadmin

import (
	"context"
	"errors"
	"testing"
)

type recordingStore struct {
	input Input
	hash  string
	err   error
}

func (s *recordingStore) CreateRoot(_ context.Context, input Input, hash string) (Result, error) {
	s.input, s.hash = input, hash
	return Result{ProfileID: 7, AuthID: 8, RoleID: 9, PermissionCount: 12}, s.err
}

func TestBootstrapValidatesAndHashesBeforeStore(t *testing.T) {
	store := &recordingStore{}
	service := NewService(store, func(password string) (string, error) {
		if password != "correct-horse-battery" {
			t.Fatalf("unexpected password passed to hasher: %q", password)
		}
		return "bcrypt-hash", nil
	})

	result, err := service.Bootstrap(context.Background(), Input{
		Username: " root.operator ", Password: "correct-horse-battery",
		FullName: " Root Operator ", Email: "root@example.com", Phone: " 0390000000 ",
	})
	if err != nil {
		t.Fatalf("Bootstrap() error = %v", err)
	}
	if result.ProfileID != 7 || store.hash != "bcrypt-hash" {
		t.Fatalf("unexpected result/store: %+v hash=%q", result, store.hash)
	}
	if store.input.Username != "root.operator" || store.input.FullName != "Root Operator" || store.input.Phone != "0390000000" {
		t.Fatalf("input was not normalized: %+v", store.input)
	}
}

func TestBootstrapRejectsInvalidInputBeforeHash(t *testing.T) {
	called := false
	service := NewService(&recordingStore{}, func(string) (string, error) {
		called = true
		return "", nil
	})

	_, err := service.Bootstrap(context.Background(), Input{
		Username: "root operator", Password: "short", FullName: "R", Email: "invalid",
	})
	if err == nil {
		t.Fatal("Bootstrap() expected validation error")
	}
	if called {
		t.Fatal("hasher called for invalid input")
	}
}

func TestBootstrapPreservesStoreError(t *testing.T) {
	store := &recordingStore{err: ErrAlreadyBootstrapped}
	service := NewService(store, func(string) (string, error) { return "hash", nil })
	_, err := service.Bootstrap(context.Background(), Input{
		Username: "root.operator", Password: "correct-horse-battery", FullName: "Root Operator",
	})
	if !errors.Is(err, ErrAlreadyBootstrapped) {
		t.Fatalf("Bootstrap() error = %v, want ErrAlreadyBootstrapped", err)
	}
}
