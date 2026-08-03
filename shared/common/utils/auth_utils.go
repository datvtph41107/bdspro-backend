package _utils

type AuthUtils struct {
	permissions map[string]int
}

func NewAuthUtils() *AuthUtils {
	return &AuthUtils{
		permissions: make(map[string]int),
	}
}
