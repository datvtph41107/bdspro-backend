package jwtverify

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestHMACVerifier_AcceptsOnlyHS256(t *testing.T) {
	t.Parallel()
	const secret = "test-verification-secret"
	verifier, err := NewHMACVerifier(secret)
	require.NoError(t, err)

	claims := Principal{Type: string(AccessToken), RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}}
	hs256, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	require.NoError(t, err)
	_, err = verifier.Parse(hs256)
	require.NoError(t, err)

	hs384, err := jwt.NewWithClaims(jwt.SigningMethodHS384, claims).SignedString([]byte(secret))
	require.NoError(t, err)
	_, err = verifier.Parse(hs384)
	require.Error(t, err)
}
