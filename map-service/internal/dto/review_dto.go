package dto

import (
	"map/internal/domain"
	"time"
)

// ReviewWithDetails represents a review with additional details
type ReviewWithDetails struct {
	domain.Review
	User          *User            `json:"user,omitempty"`
	Location      *domain.Location `json:"location,omitempty"`
	Images        []ReviewImage    `json:"images,omitempty"`
	LikesCount    int              `json:"likes_count"`
	DislikesCount int              `json:"dislikes_count"`
	UserLiked     *bool            `json:"user_liked,omitempty"` // null if not logged in
}

// ReviewRequest represents a request to create/update review
type ReviewRequest struct {
	LocationID uint   `json:"location_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Title      string `json:"title" binding:"max=255"`
	Comment    string `json:"comment"`
	IsPublic   bool   `json:"is_public"`
}

// ReviewUpdateRequest represents a request to update review
type ReviewUpdateRequest struct {
	Rating   int    `json:"rating" binding:"min=1,max=5"`
	Title    string `json:"title" binding:"max=255"`
	Comment  string `json:"comment"`
	IsPublic bool   `json:"is_public"`
}

// ReviewFilter represents filter criteria for reviews
type ReviewFilter struct {
	LocationID *uint   `json:"location_id,omitempty"`
	UserID     *uint   `json:"user_id,omitempty"`
	Rating     *int    `json:"rating,omitempty"`
	IsVerified *bool   `json:"is_verified,omitempty"`
	IsPublic   *bool   `json:"is_public,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
	DateFrom   *string `json:"date_from,omitempty"`
	DateTo     *string `json:"date_to,omitempty"`
}

// ReviewSummary represents a summary of reviews for a location
type ReviewSummary struct {
	LocationID    uint                `json:"location_id"`
	LocationName  string              `json:"location_name"`
	TotalReviews  int                 `json:"total_reviews"`
	AverageRating float64             `json:"average_rating"`
	RatingCounts  map[int]int         `json:"rating_counts"` // {1: 5, 2: 3, 3: 10, 4: 25, 5: 57}
	VerifiedCount int                 `json:"verified_count"`
	RecentReviews []ReviewWithDetails `json:"recent_reviews"`
}

// ReviewImage represents images attached to a review
type ReviewImage struct {
	ID        uint      `json:"id"`
	ReviewID  uint      `json:"review_id"`
	ImageURL  string    `json:"image_url"`
	AltText   string    `json:"alt_text"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"created_at"`
}

// ReviewLikeRequest represents a request to like/dislike a review
type ReviewLikeRequest struct {
	ReviewID uint `json:"review_id" binding:"required"`
	IsLike   bool `json:"is_like"`
}

// ReviewReportRequest represents a request to report a review
type ReviewReportRequest struct {
	ReviewID    uint   `json:"review_id" binding:"required"`
	Reason      string `json:"reason" binding:"required,max=255"`
	Description string `json:"description"`
}

// ReviewReport represents reports for inappropriate reviews
type ReviewReport struct {
	ID          uint      `json:"id"`
	ReviewID    uint      `json:"review_id"`
	ReporterID  uint      `json:"reporter_id"`
	Reason      string    `json:"reason"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// User represents a user (simplified for review context)
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
}

// ReviewResponse represents a review response
type ReviewResponse struct {
	ID         uint      `json:"id"`
	LocationID uint      `json:"location_id"`
	UserID     uint      `json:"user_id"`
	Rating     int       `json:"rating"`
	Title      string    `json:"title"`
	Comment    string    `json:"comment"`
	IsVerified bool      `json:"is_verified"`
	IsPublic   bool      `json:"is_public"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ReviewListResponse represents a list of reviews with pagination
type ReviewListResponse struct {
	Data  []ReviewWithDetails `json:"data"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}
