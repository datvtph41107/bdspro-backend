package jwtverify

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var ErrVerificationKeyRequired = errors.New("jwt verification key is required")

// Verifier verifies a JWT and returns its verified principal.
type Verifier interface {
	Parse(token string) (*Principal, error)
}

// HMACVerifier verifies the HS256 JWTs issued by QHPRO.
type HMACVerifier struct{ key []byte }

func NewHMACVerifier(secret string) (*HMACVerifier, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, ErrVerificationKeyRequired
	}
	return &HMACVerifier{key: []byte(secret)}, nil
}

func (v *HMACVerifier) Parse(raw string) (*Principal, error) {
	if v == nil || len(v.key) == 0 {
		return nil, ErrVerificationKeyRequired
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("jwt token is empty")
	}
	claims := &Principal{}
	token, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return v.key, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}
	if token == nil || !token.Valid {
		return nil, errors.New("invalid jwt claims")
	}
	return claims, nil
}
