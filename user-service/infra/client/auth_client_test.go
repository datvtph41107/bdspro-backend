package client

import (
	"testing"
	domainauth "user/internal/domain/auth"
)

func TestAuthMethodToProviderDoesNotExposePasswordHash(t *testing.T) {
	mapped := authMethodToProvider(&domainauth.AuthMethod{
		ID:       9,
		UserID:   27,
		Provider: domainauth.ProviderAdmin,
		AuthName: "operator",
		Password: "$2a$10$stored-hash",
	})

	if mapped == nil {
		t.Fatal("expected mapped auth method")
	}
	if mapped.Password != "" {
		t.Fatal("credential hash crossed the AuthProvider boundary")
	}
	if mapped.ID != 9 || mapped.UserID != 27 || mapped.AuthName != "operator" {
		t.Fatalf("identity fields changed during mapping: %+v", mapped)
	}
}
