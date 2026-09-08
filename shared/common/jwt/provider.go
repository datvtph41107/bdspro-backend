package _jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// Secret key dùng để ký token
// var jwtSecret = []byte("key-generate-jwt-token-authentication")

// Principal chứa payload của token
type Principal struct {
	AuthID         uint64           `json:"sub"`
	ProfileId      uint64           `json:"profile"`
	OriginId       uint64           `json:"origin"`
	OrganizationId *uint64          `json:"organization"`
	PlanId         *uint64          `json:"plan"`
	Type           string           `json:"type"`
	Role           string           `json:"role"`
	RoleIds        []uint64         `json:"roleIds"`
	Session        uint64           `json:"session"`
	PlanFrom       *jwt.NumericDate `json:"start"`
	jwt.RegisteredClaims
}

func GenerateToken(jwtTokenProperties JwtTokenProperties) (string, error) {
	claims := jwt.MapClaims{
		"profile":      jwtTokenProperties.ProfileID,
		"origin":       jwtTokenProperties.OriginID,
		"organization": jwtTokenProperties.OrganizationID,
		"role":         jwtTokenProperties.Role,
		"roleIds":      jwtTokenProperties.RoleIds,
		"type":         jwtTokenProperties.Type,
		"session":      jwtTokenProperties.SessionID,
		"authorities":  jwtTokenProperties.Privileges,
		"sub":          jwtTokenProperties.AuthID,
		"plan":         jwtTokenProperties.PlanID,
		"planAt":       jwtTokenProperties.PlanAt,
		"iat":          time.Now().Unix(),
	}

	// Nếu có thêm thông tin, merge vào claims
	if jwtTokenProperties.AdditionalInfo != nil {
		for key, value := range jwtTokenProperties.AdditionalInfo {
			claims[key] = value
		}
	}

	// Xác định thời gian hết hạn của token
	var expirationMinutes int64
	if jwtTokenProperties.MinuteExpired != 0 {
		expirationMinutes = jwtTokenProperties.MinuteExpired
	} else if jwtTokenProperties.Type == RefreshToken {
		expirationMinutes = viper.GetInt64("jwt.refreshExpAfterMinutes")
	} else if jwtTokenProperties.Type == TempToken {
		expirationMinutes = 10 // TEMP token chỉ có hạn 10 phút
	} else {
		expirationMinutes = viper.GetInt64("jwt.accessExpAfterMinutes")
	}

	claims["exp"] = time.Now().Add(time.Duration(expirationMinutes) * time.Minute).Unix()

	// Tạo token với thuật toán HMAC-SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(viper.GetString("jwt.key-generate")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// // GenerateJWT tạo JWT token với userID và thời gian hết hạn (expirationMinutes)
// func GenerateJWT(userID uint64, expirationMinutes int) (string, error) {
// 	claims := Principal{
// 		UserID: userID,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationMinutes) * time.Minute)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 			Issuer:    "your-app",
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(viper.GetString("jwt.key-generate"))
// }

// // ValidateJWT kiểm tra tính hợp lệ của token
// func ValidateJWT(tokenString string) (*Principal, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &Principal{}, func(token *jwt.Token) (interface{}, error) {
// 		return viper.GetString("jwt.key-generate"), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	// Kiểm tra xem token có hợp lệ không
// 	if claims, ok := token.Claims.(*Principal); ok && token.Valid {
// 		return claims, nil
// 	}

// 	return nil, errors.New("invalid token")
// }

func ParseJWT(tokenStr string) (*Principal, error) {
	// Parse token và xác thực nó
	token, err := jwt.ParseWithClaims(tokenStr, &Principal{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra tính hợp lệ của signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		key := viper.GetString("jwt.key-generate")
		return []byte(key), nil
	})

	var claims *Principal
	if token != nil && token.Claims != nil {
		claims = token.Claims.(*Principal)
	}

	// if err == nil {
	// 	return claims, nil
	// }

	// Nếu token hợp lệ, trả về claims
	// if claims, ok := token.Claims.(*Principal); ok && token.Valid {
	// 	return claims, nil
	// }

	return claims, err
}

func GetProfileId(c *gin.Context) uint64 {
	clms, _ := c.Get("claims")
	if claims, ok := clms.(*Principal); ok {
		return claims.ProfileId
	}

	return 0
}

func GetSessionId(c *gin.Context) uint64 {
	clms, _ := c.Get("claims")
	if claims, ok := clms.(*Principal); ok {
		return claims.ProfileId
	}

	return 0
}

func GetPrincipal(c *gin.Context) *Principal {
	clms, _ := c.Get("claims")
	if claims, ok := clms.(*Principal); ok {
		return claims
	}

	return nil
}

func GetPrincipalContext(c context.Context) *Principal {
	clms, ok := c.Value("claims").(*Principal)
	if !ok {
		return nil
	}

	return clms
}

func GetProperties(tokenString string) (*JwtTokenProperties, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(viper.GetString("jwt.key-generate")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Trích xuất claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}

	// Chuyển đổi dữ liệu từ claims
	authID, _ := claims["sub"].(float64)
	sessionID, _ := claims["session"].(float64)
	profileID, _ := claims["profile"].(float64)
	organizationID, _ := claims["organization"].(float64)
	planID, _ := claims["plan"].(float64)
	role, _ := claims["role"].(string)
	tokenType, _ := claims["type"].(string)
	privileges, _ := claims["authorities"].([]string)
	planFrom, _ := claims["start"].(int64)
	issuedAt, _ := claims["iat"].(float64) // Unix timestamp

	var roleIds []uint64
	if rawRoleIds, ok := claims["roleIds"].([]interface{}); ok {
		for _, item := range rawRoleIds {
			switch v := item.(type) {
			case float64:
				roleIds = append(roleIds, uint64(v))
			case int64:
				roleIds = append(roleIds, uint64(v))
			case uint64:
				roleIds = append(roleIds, v)
			}
		}
	}

	planAt := time.Unix(planFrom, 0)
	plan := uint64(planID)
	organization := uint64(organizationID)
	// Tạo đối tượng JwtTokenProperties
	return &JwtTokenProperties{
		SessionID:      uint64(sessionID),
		AuthID:         uint64(authID),
		ProfileID:      uint64(profileID),
		OrganizationID: &organization,
		PlanID:         &plan,
		Role:           role,
		RoleIds:        roleIds,
		Type:           TokenType(tokenType),
		Privileges:     privileges,
		IssuedAt:       time.Unix(int64(issuedAt), 0),
		PlanAt:         &planAt,
	}, nil
}
