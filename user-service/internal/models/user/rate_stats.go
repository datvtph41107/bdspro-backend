package user_model

// RateStats represents rating statistics
type RateStats struct {
	TotalReviews int64   `json:"totalReviews"`
	AverageScore float64 `json:"averageScore"`
	Star1Count   int64   `json:"star1Count"`
	Star2Count   int64   `json:"star2Count"`
	Star3Count   int64   `json:"star3Count"`
	Star4Count   int64   `json:"star4Count"`
	Star5Count   int64   `json:"star5Count"`
}
