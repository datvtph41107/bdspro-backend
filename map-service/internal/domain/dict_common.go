package domain

import "time"

// DictCommon represents a dictionary/common data table
type DictCommon struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string    `json:"code" gorm:"size:50;not null;uniqueIndex"` // Unique code for the dictionary item
	Name        string    `json:"name" gorm:"size:255;not null"`            // Display name
	Description string    `json:"description" gorm:"size:500"`              // Description
	Category    string    `json:"category" gorm:"size:100;not null;index"`  // Category like 'work_schedule_type', 'holiday_type', etc.
	Value       string    `json:"value" gorm:"size:100"`                    // Value (optional)
	SortOrder   int       `json:"sort_order" gorm:"default:0"`              // Display order
	IsActive    bool      `json:"is_active" gorm:"default:true"`            // Active status
	ParentID    *uint     `json:"parent_id" gorm:"index"`                   // For hierarchical data
	Metadata    string    `json:"metadata" gorm:"type:json"`                // Additional metadata in JSON
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DictCommonCategory represents different categories of dictionary data
type DictCommonCategory string

const (
	DictCategoryWorkScheduleType DictCommonCategory = "work_schedule_type"
	DictCategoryHolidayType      DictCommonCategory = "holiday_type"
	DictCategoryRecurringType    DictCommonCategory = "recurring_type"
	DictCategoryReviewRating     DictCommonCategory = "review_rating"
	DictCategoryReviewStatus     DictCommonCategory = "review_status"
	DictCategoryCategoryType     DictCommonCategory = "category_type"
	DictCategoryAmenityCategory  DictCommonCategory = "amenity_category"
	DictCategoryEventType        DictCommonCategory = "event_type"
	DictCategoryEventStatus      DictCommonCategory = "event_status"
	DictCategoryNotificationType DictCommonCategory = "notification_type"
	DictCategoryLocationType     DictCommonCategory = "location_type"
	DictCategoryDayOfWeek        DictCommonCategory = "day_of_week"
	DictCategoryCurrency         DictCommonCategory = "currency"
	DictCategoryCountry          DictCommonCategory = "country"
	DictCategoryCity             DictCommonCategory = "city"
	DictCategoryLanguage         DictCommonCategory = "language"
	DictCategoryTimeZone         DictCommonCategory = "timezone"
)

// DictCommonWithChildren represents dictionary item with children (for hierarchical data)
type DictCommonWithChildren struct {
	DictCommon
	Children []DictCommonWithChildren `json:"children,omitempty"`
}

// DictCommonFilter represents filter criteria for dictionary data
type DictCommonFilter struct {
	Category *string `json:"category,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	ParentID *uint   `json:"parent_id,omitempty"`
	Search   string  `json:"search,omitempty"`
}

// DictCommonSummary represents summary of dictionary data
type DictCommonSummary struct {
	TotalItems    int             `json:"total_items"`
	ActiveItems   int             `json:"active_items"`
	Categories    map[string]int  `json:"categories"`
	TopCategories []CategoryCount `json:"top_categories"`
}

// CategoryCount represents count of items in a category
type CategoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}
