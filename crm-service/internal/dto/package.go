package dto

import (
	"crm/internal/domain"
)

// PackageCreateDTO for creating new packages
type PackageCreateDTO struct {
	Name                 string             `json:"name" binding:"required"`
	Description          string             `json:"description"`
	Type                 domain.PackageType `json:"type" binding:"required"`
	Price                float64            `json:"price" binding:"min=0"`
	Currency             string             `json:"currency"`
	Duration             int                `json:"duration" binding:"min=1"`
	Features             string             `json:"features"`
	MaxImpressions       *int               `json:"max_impressions"`
	MaxClicks            *int               `json:"max_clicks"`
	MinBudget            *float64           `json:"min_budget"`
	MaxBudget            *float64           `json:"max_budget"`
	EstimatedReach       *int               `json:"estimated_reach"`
	EstimatedClicks      *int               `json:"estimated_clicks"`
	IsPopular            bool               `json:"is_popular"`
	IsRecommended        bool               `json:"is_recommended"`
	Icon                 string             `json:"icon"`
	Color                string             `json:"color"`
	RequiresVerification bool               `json:"requires_verification"`
	RequiresPayment      bool               `json:"requires_payment"`
}

// PackageUpdateDTO for updating packages
type PackageUpdateDTO struct {
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	Type                 domain.PackageType   `json:"type"`
	Status               domain.PackageStatus `json:"status"`
	Price                float64              `json:"price"`
	Currency             string               `json:"currency"`
	Duration             int                  `json:"duration"`
	Features             string               `json:"features"`
	MaxImpressions       *int                 `json:"max_impressions"`
	MaxClicks            *int                 `json:"max_clicks"`
	MinBudget            *float64             `json:"min_budget"`
	MaxBudget            *float64             `json:"max_budget"`
	EstimatedReach       *int                 `json:"estimated_reach"`
	EstimatedClicks      *int                 `json:"estimated_clicks"`
	IsPopular            bool                 `json:"is_popular"`
	IsRecommended        bool                 `json:"is_recommended"`
	Icon                 string               `json:"icon"`
	Color                string               `json:"color"`
	RequiresVerification bool                 `json:"requires_verification"`
	RequiresPayment      bool                 `json:"requires_payment"`
}

// PackageSearchDTO for searching packages
type PackageSearchDTO struct {
	Type                 domain.PackageType   `json:"type"`
	Status               domain.PackageStatus `json:"status"`
	MinPrice             *float64             `json:"min_price"`
	MaxPrice             *float64             `json:"max_price"`
	IsPopular            *bool                `json:"is_popular"`
	IsRecommended        *bool                `json:"is_recommended"`
	RequiresVerification *bool                `json:"requires_verification"`
	Page                 int                  `json:"page"`
	Limit                int                  `json:"limit"`
}

// PackageResponseDTO for API responses
type PackageResponseDTO struct {
	ID                   uint64               `json:"id"`
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	Type                 domain.PackageType   `json:"type"`
	Status               domain.PackageStatus `json:"status"`
	Price                float64              `json:"price"`
	Currency             string               `json:"currency"`
	Duration             int                  `json:"duration"`
	Features             string               `json:"features"`
	MaxImpressions       *int                 `json:"max_impressions"`
	MaxClicks            *int                 `json:"max_clicks"`
	MinBudget            *float64             `json:"min_budget"`
	MaxBudget            *float64             `json:"max_budget"`
	EstimatedReach       *int                 `json:"estimated_reach"`
	EstimatedClicks      *int                 `json:"estimated_clicks"`
	IsPopular            bool                 `json:"is_popular"`
	IsRecommended        bool                 `json:"is_recommended"`
	Icon                 string               `json:"icon"`
	Color                string               `json:"color"`
	RequiresVerification bool                 `json:"requires_verification"`
	RequiresPayment      bool                 `json:"requires_payment"`
	CreatedBy            uint64               `json:"created_by"`
	UpdatedBy            uint64               `json:"updated_by"`
	OrganizationID       uint64               `json:"organization_id"`
	CreatedAt            string               `json:"created_at"`
	UpdatedAt            string               `json:"updated_at"`
}

// PackageSelectionDTO for D.10.2.2 package selection
type PackageSelectionDTO struct {
	PackageID    uint64   `json:"package_id" binding:"required"`
	CustomBudget *float64 `json:"custom_budget"` // For Ads packages with custom budget
	Duration     *int     `json:"duration"`      // For VIP packages with custom duration
}

// PackageSelectionResponseDTO for D.10.2.2 response
type PackageSelectionResponseDTO struct {
	SelectedPackage PackageResponseDTO `json:"selected_package"`
	CustomBudget    *float64           `json:"custom_budget,omitempty"`
	CustomDuration  *int               `json:"custom_duration,omitempty"`
	TotalPrice      float64            `json:"total_price"`
	EstimatedReach  *int               `json:"estimated_reach,omitempty"`
	EstimatedClicks *int               `json:"estimated_clicks,omitempty"`
	NextStep        string             `json:"next_step"` // "payment" or "campaign_creation"
}

// PackageRecommendationDTO for package recommendations
type PackageRecommendationDTO struct {
	ProductID      uint64   `json:"product_id"`
	ProductType    string   `json:"product_type"`
	CurrentViews   *int     `json:"current_views"`
	TargetAudience string   `json:"target_audience"`
	Budget         *float64 `json:"budget"`
}

// PackageRecommendationResponseDTO for package recommendations response
type PackageRecommendationResponseDTO struct {
	RecommendedPackages []*PackageResponseDTO `json:"recommended_packages"`
	PopularPackages     []*PackageResponseDTO `json:"popular_packages"`
	BestValuePackages   []*PackageResponseDTO `json:"best_value_packages"`
	ComboSuggestions    []ComboPackageDTO     `json:"combo_suggestions,omitempty"`
}

// ComboPackageDTO for combo package suggestions
type ComboPackageDTO struct {
	ID          uint64               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Packages    []PackageResponseDTO `json:"packages"`
	TotalPrice  float64              `json:"total_price"`
	Discount    float64              `json:"discount"`
	Savings     float64              `json:"savings"`
}