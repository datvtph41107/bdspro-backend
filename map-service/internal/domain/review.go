package domain

import "time"

// Review represents a review for a location
type Review struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	LocationID uint      `json:"location_id" gorm:"not null;index"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	Rating     int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"` // 1-5 stars
	Title      string    `json:"title" gorm:"size:255"`
	Comment    string    `json:"comment" gorm:"type:text"`
	IsVerified bool      `json:"is_verified" gorm:"default:false"` // Verified purchase/visit
	IsPublic   bool      `json:"is_public" gorm:"default:true"`    // Public or private review
	IsActive   bool      `json:"is_active" gorm:"default:true"`    // Active or deleted
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ReviewImage represents images attached to a review
type ReviewImage struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ReviewID  uint      `json:"review_id" gorm:"not null;index"`
	ImageURL  string    `json:"image_url" gorm:"size:500;not null"`
	AltText   string    `json:"alt_text" gorm:"size:255"`
	Order     int       `json:"order" gorm:"default:0"` // Display order
	CreatedAt time.Time `json:"created_at"`
}

// ReviewLike represents likes/dislikes for a review
type ReviewLike struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ReviewID  uint      `json:"review_id" gorm:"not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	IsLike    bool      `json:"is_like" gorm:"not null"` // true = like, false = dislike
	CreatedAt time.Time `json:"created_at"`

	// Unique constraint on (review_id, user_id)
	_ struct{} `gorm:"uniqueIndex:idx_review_user"`
}

// ReviewReport represents reports for inappropriate reviews
type ReviewReport struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ReviewID    uint      `json:"review_id" gorm:"not null;index"`
	ReporterID  uint      `json:"reporter_id" gorm:"not null;index"`
	Reason      string    `json:"reason" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Status      int       `json:"status" gorm:"default:0"` // 0=pending, 1=approved, 2=rejected
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ReviewRating represents different rating types
type ReviewRating int

const (
	ReviewRating1Star ReviewRating = 1
	ReviewRating2Star ReviewRating = 2
	ReviewRating3Star ReviewRating = 3
	ReviewRating4Star ReviewRating = 4
	ReviewRating5Star ReviewRating = 5
)

// ReviewStatus represents the status of a review
type ReviewStatus int

const (
	ReviewStatusPending  ReviewStatus = iota // Pending approval
	ReviewStatusApproved                     // Approved and visible
	ReviewStatusRejected                     // Rejected
	ReviewStatusHidden                       // Hidden by admin
	ReviewStatusDeleted                      // Soft deleted
)

// ReviewReportStatus represents the status of a review report
type ReviewReportStatus int

const (
	ReviewReportStatusPending  ReviewReportStatus = iota // Pending review
	ReviewReportStatusApproved                           // Report approved
	ReviewReportStatusRejected                           // Report rejected
)
