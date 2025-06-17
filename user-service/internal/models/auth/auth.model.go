package auth_model

import (
	_db "common/db"
	"time"
)

// UserAuthEntity tương đương với Java Entity
type UserAuthEntity struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement"`
	Provider    string     `gorm:"size:10"` // enum: email, zalo, facebook, phone
	Username    string     `gorm:"size:50"`
	OAuthID     string     `gorm:"column:oauth_id"`
	Avatar      string     `gorm:"size:255"`
	IsSensitive bool       `gorm:"column:is_sensitive"`
	Name        string     `gorm:"size:255"`
	Email       string     `gorm:"size:255"`
	Phone       string     `gorm:"size:255"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	UserID      uint64     `gorm:"column:user_id"` // Không cập nhật & xóa, chỉ tạo mới
}

// TableName đặt tên bảng tương đương với Java
func (UserAuthEntity) TableName() string {
	return "user_auth"
}

// Tìm kiếm bằng số điện thoại
func AuthFindByPhone(phone string) (*UserAuthEntity, error) {
	var user UserAuthEntity
	err := _db.DB.Where("provider = ? AND oauth_id = ? AND deleted_at is null", "PHONE", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Tìm kiếm bằng provider và username
func AuthFindByOAuthId(oauthId, provider string) (*UserAuthEntity, error) {
	var user UserAuthEntity
	err := _db.DB.Where("provider = ? AND oauth_id = ? AND deleted_at is null", provider, oauthId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Lấy user theo ID
func AuthGetWithId(id uint64) (*UserAuthEntity, error) {
	var user UserAuthEntity
	err := _db.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Kiểm tra sự tồn tại của tài khoản theo OAuth ID
func AuthExistsByOAuth(phone string) (bool, error) {
	var count int64
	err := _db.DB.Model(&UserAuthEntity{}).Where("provider = ? AND oauth_id = ? AND deleted_at is null", "PHONE", phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DeleteProfile(profileId uint64) {
	_db.DB.Model(&UserAuthEntity{}).Where("user_id = ?", profileId).Update("deleted_at", time.Now())
}
