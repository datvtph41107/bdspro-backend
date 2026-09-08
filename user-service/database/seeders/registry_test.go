package seeders

import (
	"errors"
	"testing"
)

func TestParseName(t *testing.T) {
	for _, want := range Names() {
		got, err := ParseName(string(want))
		if err != nil {
			t.Fatalf("ParseName(%q): %v", want, err)
		}
		if got != want {
			t.Fatalf("ParseName(%q) = %q, want %q", want, got, want)
		}
	}
	if _, err := ParseName("unknown"); err == nil {
		t.Fatal("ParseName(unknown) unexpectedly succeeded")
	}
}

func TestRequireNonProduction(t *testing.T) {
	for _, environment := range []string{"development", "dev", "local", "test", "testing", "acceptance"} {
		if err := requireNonProduction(environment); err != nil {
			t.Fatalf("requireNonProduction(%q): %v", environment, err)
		}
	}
	for _, environment := range []string{"", "production", "prod", "staging"} {
		err := requireNonProduction(environment)
		if !errors.Is(err, ErrUnsafeEnvironment) {
			t.Fatalf("requireNonProduction(%q) error = %v, want ErrUnsafeEnvironment", environment, err)
		}
	}
}

func TestAcceptanceAdminPermissionSet(t *testing.T) {
	if got, want := len(acceptanceAdminPermissionKeys), 15; got != want {
		t.Fatalf("acceptance admin permission count = %d, want %d", got, want)
	}
	seen := make(map[string]struct{}, len(acceptanceAdminPermissionKeys))
	for _, key := range acceptanceAdminPermissionKeys {
		if _, exists := seen[key]; exists {
			t.Fatalf("duplicate acceptance admin permission %q", key)
		}
		seen[key] = struct{}{}
	}
}

func TestAcceptanceIdentityPatterns(t *testing.T) {
	if !acceptanceUsernamePattern.MatchString("qhpro_acceptance_0000000000") {
		t.Fatal("canonical acceptance username rejected")
	}
	if !acceptancePhonePattern.MatchString("0390000000") {
		t.Fatal("canonical acceptance phone rejected")
	}
	if !acceptanceEmailPattern.MatchString("qhpro-0000000000@e.invalid") {
		t.Fatal("canonical acceptance email rejected")
	}
}
