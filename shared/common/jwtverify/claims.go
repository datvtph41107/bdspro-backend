// Package jwtverify owns process-neutral verification of the JWT contract used by QHPRO.
package jwtverify

import "github.com/golang-jwt/jwt/v5"

type TokenType string

const (
	AccessToken  TokenType = "ACCESS"
	RefreshToken TokenType = "REFRESH"
	TempToken    TokenType = "TEMP"
)

// Principal is the verified identity carried by a QHPRO JWT.
type Principal struct {
	AuthID         uint64           `json:"sub"`
	ProfileID      uint64           `json:"profile"`
	OriginID       uint64           `json:"origin"`
	OrganizationID *uint64          `json:"organization"`
	PlanID         *uint64          `json:"plan"`
	Type           string           `json:"type"`
	Role           string           `json:"role"`
	RoleIDs        []uint64         `json:"roleIds"`
	SessionID      uint64           `json:"session"`
	PlanFrom       *jwt.NumericDate `json:"start"`
	jwt.RegisteredClaims
}
