package middlewares

import (
	"errors"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	JWT_SECRET = []byte("SUPER_SECRET")
	ISSUER     = "BLC TEAM"
	AUDIENCE   = "LAUNCH PAD"
)

func GenerateJWT(expiry time.Duration, additionalData map[string]interface{}) (map[string]interface{}, error) {
	expTime := time.Now().Add(expiry)

	claims := jwt.MapClaims{
		"exp": expTime.Unix(),
		"iat": time.Now().Unix(),
		"iss": ISSUER,
		"aud": AUDIENCE,
		"jti": GenerateRandomString(10),
	}

	maps.Copy(claims, additionalData)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(JWT_SECRET)
	if err != nil {
		return nil, err
	}

	claims["accessToken"] = token

	return claims, nil
}

func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
func GenerateRandomString(length int) string {
	randomUUID := uuid.New()
	return randomUUID.String()[:length]
}