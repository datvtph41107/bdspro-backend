package auth

// IdDomain là struct cơ sở chứa authId
type IdDomain struct {
	AuthID uint64 `gorm:"primaryKey;column:auth_id"`
}
