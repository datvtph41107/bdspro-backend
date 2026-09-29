package _jwt

// TokenVerifier verifies one JWT and returns its canonical principal.
// Implementations must not depend on process-global configuration.
type TokenVerifier interface {
	Parse(token string) (*Principal, error)
}
