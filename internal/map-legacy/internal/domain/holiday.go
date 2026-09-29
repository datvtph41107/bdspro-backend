package domain

import "time"

// Holiday represents a holiday for a location
type Holiday struct {
	ID            uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID    uint          `json:"location_id" gorm:"not null;index"`
	Name          string        `json:"name" gorm:"size:255;not null"`        // Tên ngày nghỉ
	Description   string        `json:"description" gorm:"size:1000"`         // Mô tả
	StartDate     time.Time     `json:"start_date" gorm:"type:date;not null"` // Ngày bắt đầu
	EndDate       time.Time     `json:"end_date" gorm:"type:date;not null"`   // Ngày kết thúc
	IsAllDay      bool          `json:"is_all_day" gorm:"default:true"`       // Cả ngày hay chỉ giờ cụ thể
	StartTime     string        `json:"start_time" gorm:"type:time"`          // Giờ bắt đầu (nếu không phải cả ngày)
	EndTime       string        `json:"end_time" gorm:"type:time"`            // Giờ kết thúc (nếu không phải cả ngày)
	IsRecurring   bool          `json:"is_recurring" gorm:"default:false"`    // Có lặp lại hàng năm không
	RecurringType RecurringType `json:"recurring_type"`                       // Loại lặp lại
	IsActive      bool          `json:"is_active" gorm:"default:true"`        // Có hoạt động không
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// RecurringType represents the type of recurring holiday
type RecurringType int

const (
	RecurringTypeNone    RecurringType = iota // Không lặp lại
	RecurringTypeYearly                       // Hàng năm
	RecurringTypeMonthly                      // Hàng tháng
	RecurringTypeWeekly                       // Hàng tuần
	RecurringTypeCustom                       // Tùy chỉnh
)

// HolidayType represents different types of holidays
type HolidayType int

const (
	HolidayTypeNational    HolidayType = iota // Ngày nghỉ quốc gia
	HolidayTypeReligious                      // Ngày nghỉ tôn giáo
	HolidayTypeCultural                       // Ngày nghỉ văn hóa
	HolidayTypeCompany                        // Ngày nghỉ công ty
	HolidayTypeMaintenance                    // Ngày bảo trì
	HolidayTypeEmergency                      // Ngày nghỉ khẩn cấp
	HolidayTypeCustom                         // Tùy chỉnh
)

// HolidayStatus represents the status of a holiday
type HolidayStatus int

const (
	HolidayStatusActive    HolidayStatus = iota // Đang hoạt động
	HolidayStatusInactive                       // Không hoạt động
	HolidayStatusExpired                        // Đã hết hạn
	HolidayStatusCancelled                      // Đã hủy
)
