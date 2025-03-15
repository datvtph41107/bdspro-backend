package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Secret key dùng để ký token
var jwtSecret = []byte("key-generate-jwt-token-authentication")

// Claims chứa payload của token
type Claims struct {
	UserID    string `json:"sub"`
	ProfileId uint64 `json:"profile"`
	Type      string `json:"type"`
	Role      string `json:"role"`
	Session   uint   `json:"session"`
	jwt.RegisteredClaims
}

// GenerateJWT tạo JWT token với userID và thời gian hết hạn (expirationMinutes)
func GenerateJWT(userID uint, expirationMinutes int) (string, error) {
	claims := Claims{
		UserID: string(userID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT kiểm tra tính hợp lệ của token
func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Kiểm tra xem token có hợp lệ không
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ParseJWT(tokenStr string) (*Claims, error) {
	// Parse token và xác thực nó
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra tính hợp lệ của signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Nếu token hợp lệ, trả về claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func GetProfileId(c *gin.Context) uint64 {
	clms, _ := c.Get("claims")
	if claims, ok := clms.(*Claims); ok {
		return claims.ProfileId
	}

	return 0
}
