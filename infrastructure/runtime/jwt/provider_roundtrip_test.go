package _jwt

import (
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
