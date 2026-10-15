package entity

// AuthMethod is the credential identity used by login flows.
// Full persistence shape lives in modules/auth/domain/auth.
type AuthMethod struct {
	ID       uint64
	UserID   uint64
	AuthName string
	Phone    string
	FullName string
	Provider string
	Password string
}
