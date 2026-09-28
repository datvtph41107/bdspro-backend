package auth

import "time"

// UserSessionEntity đại diện cho phiên đăng nhập của người dùng
type UserSessionEntity struct {
	SessionID    uint64     `gorm:"primaryKey;column:session_id;autoIncrement"`
	AuthID       uint64     `gorm:"column:auth_id"`
	DeviceID     string     `gorm:"column:device_id;size:255"`
	Platform     string     `gorm:"column:platform;size:50"`
	Version      string     `gorm:"column:version;size:50"`
	OS           string     `gorm:"column:os;size:50"`
	DeviceName   string     `gorm:"column:device_name;size:255"`
	CreatedDate  *time.Time `gorm:"column:created_date;autoCreateTime"`
	FinishedDate *time.Time `gorm:"column:finished_date"`
	LastLogin    *time.Time `gorm:"column:last_login"`
	LogoutAt     *time.Time `gorm:"column:logout_at"`
	LastRequest  *time.Time `gorm:"column:last_req"`
	IPRequest    string     `gorm:"column:ip_request;size:50"`
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
