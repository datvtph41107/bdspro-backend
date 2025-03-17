package models

import (
	_db "common/db"
	"time"
)

// UserSessionEntity đại diện cho phiên đăng nhập của người dùng
type UserSessionEntity struct {
	SessionID    uint64     `gorm:"primaryKey;column:session_id;autoIncrement"`
	AuthID       uint64     `gorm:"column:auth_id"`
	Platform     string     `gorm:"column:platform;size:10"`
	Version      string     `gorm:"column:version;size:20"`
	OS           string     `gorm:"column:os;size:20"`
	CreatedDate  *time.Time `gorm:"column:created_date;autoCreateTime"`
	FinishedDate *time.Time `gorm:"column:finished_date"`
	IPRequest    string     `gorm:"column:ip_request;size:20"`
	TotalRequest uint64     `gorm:"column:total_request"`
	UserAgent    string     `gorm:"column:user_agent" json:"userAgent"`
	Activate     bool       `gorm:"column:activate" json:"activate"`
	SessionKey   string     `gorm:"column:session_key" json:"sessionKey"`
	ClientID     string     `gorm:"column:client_id" json:"clientId"`
	// DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName đặt tên bảng trong DB
func (UserSessionEntity) TableName() string {
	return "user_session"
}

// GetCurrentSession lấy phiên đăng nhập gần nhất của user theo authId
func GetCurrentSession(authId uint64) (*UserSessionEntity, error) {
	var session UserSessionEntity
	err := _db.DB.Where("auth_id = ?", authId).Order("created_date desc").First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetBySessionId lấy thông tin session theo sessionId
func SessionGetBySessionId(sessionId uint64) (*UserSessionEntity, error) {
	var session UserSessionEntity
	err := _db.DB.Where("session_id = ?", sessionId).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}
