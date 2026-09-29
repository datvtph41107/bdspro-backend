package zns

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"user/internal/domain/auth"
)

type authConfigRepoStub struct {
	atomicErr   error
	atomicCalls int
	values      map[string]string
}

func (s *authConfigRepoStub) GetByKey(context.Context, string) (*auth.AuthConfig, error) {
	panic("unexpected GetByKey call")
}

func (s *authConfigRepoStub) UpdateValue(context.Context, string, string) error {
	panic("unexpected UpdateValue call")
}

func (s *authConfigRepoStub) UpdateValuesAtomically(_ context.Context, values map[string]string) error {
	s.atomicCalls++
	s.values = make(map[string]string, len(values))
	for key, value := range values {
		s.values[key] = value
	}
	return s.atomicErr
}

func TestSaveTokensToDBPublishesMemoryOnlyAfterAtomicDurableWrite(t *testing.T) {
	repo := &authConfigRepoStub{}
	provider := &ZnsProvider{
		accessToken:    "old-access",
		refreshToken:   "old-refresh",
		authConfigRepo: repo,
	}

	if err := provider.saveTokensToDB(context.Background(), "new-access", "new-refresh"); err != nil {
		t.Fatalf("saveTokensToDB() error = %v", err)
	}

	want := map[string]string{
		auth.ConfigKeyZNSToken:   "new-access",
		auth.ConfigKeyZNSRefresh: "new-refresh",
	}
	if repo.atomicCalls != 1 {
		t.Fatalf("atomic calls = %d, want 1", repo.atomicCalls)
	}
	if !reflect.DeepEqual(repo.values, want) {
		t.Fatalf("atomic values = %#v, want %#v", repo.values, want)
	}
	if provider.accessToken != "new-access" || provider.refreshToken != "new-refresh" {
		t.Fatalf("memory pair = (%q, %q), want new pair", provider.accessToken, provider.refreshToken)
	}
}

func TestSaveTokensToDBKeepsPreviousMemoryWhenAtomicDurableWriteFails(t *testing.T) {
	writeErr := errors.New("database unavailable")
	repo := &authConfigRepoStub{atomicErr: writeErr}
	provider := &ZnsProvider{
		accessToken:    "old-access",
		refreshToken:   "old-refresh",
		authConfigRepo: repo,
	}

	err := provider.saveTokensToDB(context.Background(), "new-access", "new-refresh")
	if !errors.Is(err, writeErr) {
		t.Fatalf("saveTokensToDB() error = %v, want wrapped %v", err, writeErr)
	}
	if provider.accessToken != "old-access" || provider.refreshToken != "old-refresh" {
		t.Fatalf("memory pair changed on durable failure: (%q, %q)", provider.accessToken, provider.refreshToken)
	}
}
