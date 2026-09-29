package _jwt

import (
	"time"
)

// TokenType định nghĩa kiểu token (access/refresh/temp)
type TokenType string

const (
	AccessToken  TokenType = "ACCESS"
	RefreshToken TokenType = "REFRESH"
	TempToken    TokenType = "TEMP" // Token tạm thời cho user chưa có fullname, hạn 10 phút
)

type UserMetadata struct {
	ProfileID      *uint64
	OrganizationID *uint64
	PlanID         *uint64
	Role           []string
}

// JwtTokenProperties chứa thông tin JWT
type JwtTokenProperties struct {
	AuthID         uint64                 `json:"authId"`
	ProfileID      uint64                 `json:"profileId"`
	OriginID       uint64                 `json:"originId"`
	OrganizationID *uint64                `json:"organizationId"`
	PlanID         *uint64                `json:"planId"`
	Role           string                 `json:"role"`
	RoleIds        []uint64               `json:"roleIds"`
	SessionID      uint64                 `json:"sessionId"`
	MinuteExpired  int64                  `json:"minuteExpired"`
	Type           TokenType              `json:"type"`
	IssuedAt       time.Time              `json:"issuedAt"`
	Privileges     []string               `json:"privileges"`
	PlanAt         *time.Time             `json:"planAt"`
	AdditionalInfo map[string]interface{} `json:"additionalInformation"`
}

// NewJwtTokenProperties tạo một JWT token struct mới (builder pattern)
func MakeJwtTokenProperties(authId, profileId uint64, planId *uint64, sessionId uint64, role string, minuteExpired int64, tokenType TokenType, issuedAt time.Time, privileges []string, additionalInfo map[string]interface{}) *JwtTokenProperties {
	return &JwtTokenProperties{
		AuthID:         authId,
		ProfileID:      profileId,
		PlanID:         planId,
		Role:           role,
		SessionID:      sessionId,
		MinuteExpired:  minuteExpired,
		Type:           tokenType,
		IssuedAt:       issuedAt,
		Privileges:     privileges,
		AdditionalInfo: additionalInfo,
	}
}
