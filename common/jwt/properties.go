package _jwt

import (
	"time"
)

// TokenType định nghĩa kiểu token (access/refresh)
type TokenType string

const (
	AccessToken  TokenType = "ACCESS"
	RefreshToken TokenType = "REFRESH"
)

// JwtTokenProperties chứa thông tin JWT
type JwtTokenProperties struct {
	AuthID         uint64                 `json:"authId"`
	ProfileID      uint64                 `json:"profileId"`
	Role           string                 `json:"role"`
	SessionID      uint64                 `json:"sessionId"`
	MinuteExpired  int64                  `json:"minuteExpired"`
	Type           TokenType              `json:"type"`
	IssuedAt       time.Time              `json:"issuedAt"`
	Privileges     []string               `json:"privileges"`
	AdditionalInfo map[string]interface{} `json:"additionalInformation"`
}

// NewJwtTokenProperties tạo một JWT token struct mới (builder pattern)
func NewJwtTokenProperties(authId, profileId, sessionId uint64, role string, minuteExpired int64, tokenType TokenType, issuedAt time.Time, privileges []string, additionalInfo map[string]interface{}) *JwtTokenProperties {
	return &JwtTokenProperties{
		AuthID:         authId,
		ProfileID:      profileId,
		Role:           role,
		SessionID:      sessionId,
		MinuteExpired:  minuteExpired,
		Type:           tokenType,
		IssuedAt:       issuedAt,
		Privileges:     privileges,
		AdditionalInfo: additionalInfo,
	}
}
