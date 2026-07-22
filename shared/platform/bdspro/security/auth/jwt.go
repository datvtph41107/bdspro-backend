package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gateway/internal/principal"

	"github.com/golang-jwt/jwt/v5"
)

// ParseError phân loại lý do token không hợp lệ — dùng để log metric chính xác
type ParseError struct {
	Reason string
	Err    error
}

func (e *ParseError) Error() string { return fmt.Sprintf("auth: %s: %v", e.Reason, e.Err) }
func (e *ParseError) Unwrap() error { return e.Err }

var (
	ErrMissingToken = &ParseError{Reason: "missing_token"}
	ErrExpiredToken = &ParseError{Reason: "expired_token"}
	ErrInvalidToken = &ParseError{Reason: "invalid_token"}
)

// Parser parse JWT và trả về Principal — stateless, injectable, testable
type Parser struct {
	secretKey []byte
}

// NewParser tạo parser với secret từ config — không dùng viper global
func NewParser(secretKey string) *Parser {
	return &Parser{secretKey: []byte(secretKey)}
}

// claims là private — chỉ dùng để parse, không expose ra handler
type claims struct {
	AuthID         uint64  `json:"sub"`
	ProfileID      uint64  `json:"profile"`
	OriginID       uint64  `json:"origin"`
	OrganizationID *uint64 `json:"organization"`
	PlanID         *uint64 `json:"plan"`
	Role           string  `json:"role"`
	Type           string  `json:"type"`
	Session        uint64  `json:"session"`
	jwt.RegisteredClaims
}

// Parse nhận Authorization header value, trả Principal hoặc ParseError
// Input: "Bearer eyJhb..." hoặc "eyJhb..."
// Được thiết kế để gọi một lần duy nhất tại middleware auth
func (p *Parser) Parse(authHeader string) (*principal.Principal, error) {
	tokenStr := extractBearerToken(authHeader)
	if tokenStr == "" {
		return nil, ErrMissingToken
	}

	token, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return p.secretKey, nil
	})

	if err != nil {
		if isExpiredError(err) {
			return nil, ErrExpiredToken
		}
		return nil, &ParseError{Reason: "invalid_token", Err: err}
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return &principal.Principal{
		AuthID:         c.AuthID,
		ProfileID:      c.ProfileID,
		OriginID:       c.OriginID,
		OrganizationID: c.OrganizationID,
		PlanID:         c.PlanID,
		Role:           c.Role,
		Type:           c.Type,
		Session:        c.Session,
	}, nil
}

// ParseAllowExpired parse token kể cả khi expired — dùng cho refresh token flow
func (p *Parser) ParseAllowExpired(authHeader string) (*principal.Principal, error) {
	tokenStr := extractBearerToken(authHeader)
	if tokenStr == "" {
		return nil, ErrMissingToken
	}

	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(tokenStr, &claims{})
	if err != nil {
		return nil, &ParseError{Reason: "malformed_token", Err: err}
	}

	c, ok := token.Claims.(*claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return &principal.Principal{
		AuthID:    c.AuthID,
		ProfileID: c.ProfileID,
		Type:      c.Type,
	}, nil
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, "Bearer ", 2)
	return strings.TrimSpace(parts[len(parts)-1])
}

func isExpiredError(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired) ||
		strings.Contains(err.Error(), "token is expired") ||
		func() bool {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				return ve.Errors&jwt.ValidationErrorExpired != 0
			}
			return false
		}()
}

// TokenType định nghĩa loại token — dùng để kiểm tra policy
type TokenType string

const (
	AccessToken  TokenType = "ACCESS"
	RefreshToken TokenType = "REFRESH"
	TempToken    TokenType = "TEMP"
)

// IsTempToken kiểm tra principal có phải temp token không
func IsTempToken(p *principal.Principal) bool {
	return p != nil && p.Type == string(TempToken)
}

// IsExpired kiểm tra thời gian thủ công (cho trường hợp cần check mà không parse lại)
func IsTokenFresh(issuedAt time.Time, maxAge time.Duration) bool {
	return time.Since(issuedAt) < maxAge
}
