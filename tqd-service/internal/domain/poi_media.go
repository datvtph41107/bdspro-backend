package domain

import (
	"time"
	"tqd/internal/enums"
)

// POIMedia represents media files associated with a Point of Interest
type POIMedia struct {
	ID          uint                 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string               `json:"name" gorm:"size:255;not null"`
	POIID       uint                 `json:"poi_id" gorm:"not null;index"`
	MediaType   string               `json:"media_type" gorm:"size:50;not null"` // image, video, audio, document
	URL         string               `json:"url" gorm:"size:500;not null"`
	Thumbnail   string               `json:"thumbnail" gorm:"size:500"` // Thumbnail URL for videos/images
	Description string               `json:"description" gorm:"type:text"`
	AltText     string               `json:"alt_text" gorm:"size:255"`    // Alt text for accessibility
	FileSize    int64                `json:"file_size"`                   // File size in bytes
	Duration    int                  `json:"duration"`                    // Duration in seconds (for videos/audio)
	Width       int                  `json:"width"`                       // Image/video width
	Height      int                  `json:"height"`                      // Image/video height
	SortOrder   int                  `json:"sort_order" gorm:"default:0"` // Display order
	Status      enums.POIMediaStatus `json:"status" gorm:"default:0"`
	Type        enums.POIMediaType   `json:"type" gorm:"default:0"`
	IsActive    bool                 `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}
