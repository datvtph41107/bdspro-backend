package _jwt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var ErrVerificationKeyRequired = errors.New("jwt verification key is required")

// TokenVerifier verifies one JWT and returns its canonical principal.
// Implementations must not depend on process-global configuration.
type TokenVerifier interface {
	Parse(token string) (*Principal, error)
}

// HMACVerifier verifies the HMAC JWTs currently issued by QHPRO.
// The verification key is explicit so edge processes do not depend on Viper/global state.
type HMACVerifier struct {
	key []byte
}

func NewHMACVerifier(secret string) (*HMACVerifier, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, ErrVerificationKeyRequired
	}
	return &HMACVerifier{key: []byte(secret)}, nil
}

func (v *HMACVerifier) Parse(tokenString string) (*Principal, error) {
	if v == nil || len(v.key) == 0 {
		return nil, ErrVerificationKeyRequired
	}
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, errors.New("jwt token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Principal{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Principal)
	if !ok || claims == nil || !token.Valid {
		return nil, errors.New("invalid jwt claims")
	}
	return claims, nil
}
