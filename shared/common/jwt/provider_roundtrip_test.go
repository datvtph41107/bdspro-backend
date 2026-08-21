package _jwt

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestGenerateAndParseLocalAccessToken(t *testing.T) {
	const secret = "qhpro-local-jwt-test"
	viper.Set("jwt.key-generate", secret)
	t.Cleanup(func() { viper.Set("jwt.key-generate", "") })

	token, err := GenerateToken(JwtTokenProperties{
		AuthID:        42,
		ProfileID:     42,
		OriginID:      42,
		Role:          "user",
		SessionID:     1,
		MinuteExpired: 60,
		Type:          AccessToken,
		IssuedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	principal, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("ParseJWT() error = %v", err)
	}
	if principal == nil || principal.ProfileId != 42 || principal.AuthID != 42 {
		t.Fatalf("parsed principal = %+v", principal)
	}
}

func TestHMACVerifierUsesExplicitVerificationKey(t *testing.T) {
	const secret = "qhpro-explicit-verifier-test"
	viper.Set("jwt.key-generate", secret)
	t.Cleanup(func() { viper.Set("jwt.key-generate", "") })

	token, err := GenerateToken(JwtTokenProperties{
		AuthID:        77,
		ProfileID:     88,
		OriginID:      88,
		Role:          "user",
		SessionID:     2,
		MinuteExpired: 60,
		Type:          AccessToken,
		IssuedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// Prove verification does not consult the process-global registry.
	viper.Set("jwt.key-generate", "wrong-global-key")
	verifier, err := NewHMACVerifier(secret)
	if err != nil {
		t.Fatalf("NewHMACVerifier() error = %v", err)
	}
	principal, err := verifier.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if principal.AuthID != 77 || principal.ProfileId != 88 {
		t.Fatalf("principal = %+v", principal)
	}
}

func TestNewHMACVerifierRejectsMissingKey(t *testing.T) {
	if _, err := NewHMACVerifier("  "); !errors.Is(err, ErrVerificationKeyRequired) {
		t.Fatalf("NewHMACVerifier() error = %v", err)
	}
}
