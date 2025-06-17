package auth_model

// UserIdEntity là struct cơ sở chứa authId
type UserIdEntity struct {
	AuthID uint64 `gorm:"primaryKey;column:auth_id"`
}
