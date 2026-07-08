package domain

import (
	_models "common/models"
)

// PackageType represents the type of advertising package
type PackageType string

const (
	PackageTypeVIP       PackageType = "VIP"
	PackageTypeMetaAds   PackageType = "META_ADS"
	PackageTypeGoogleAds PackageType = "GOOGLE_ADS"
	PackageTypePR        PackageType = "PR"
)

// PackageStatus represents the status of a package
type PackageStatus string

const (
	PackageStatusActive   PackageStatus = "ACTIVE"
	PackageStatusInactive PackageStatus = "INACTIVE"
)

// Package represents an advertising package
type Package struct {
	_models.BaseEntity
	Name        string        `json:"name" gorm:"not null"`
	Description string        `json:"description"`
	Type        PackageType   `json:"type" gorm:"not null"`
	Status      PackageStatus `json:"status" gorm:"default:'ACTIVE'"`

	// Pricing
	Price    float64 `json:"price"`
	Currency string  `json:"currency" gorm:"default:'VND'"`
	Duration int     `json:"duration"` // in days

	// Features
	Features       string `json:"features"` // JSON string
	MaxImpressions *int   `json:"max_impressions"`
	MaxClicks      *int   `json:"max_clicks"`

	// D.10.2.2 specific fields
	MinBudget       *float64 `json:"min_budget"`                          // Minimum budget for Ads packages
	MaxBudget       *float64 `json:"max_budget"`                          // Maximum budget for Ads packages
	EstimatedReach  *int     `json:"estimated_reach"`                     // Estimated reach for this package
	EstimatedClicks *int     `json:"estimated_clicks"`                    // Estimated clicks for this package
	IsPopular       bool     `json:"is_popular" gorm:"default:false"`     // Popular package flag
	IsRecommended   bool     `json:"is_recommended" gorm:"default:false"` // Recommended package flag
	Icon            string   `json:"icon"`                                // Icon for UI display
	Color           string   `json:"color"`                               // Color theme for UI

	// Validation
	RequiresVerification bool `json:"requires_verification" gorm:"default:false"` // Requires verified account
	RequiresPayment      bool `json:"requires_payment" gorm:"default:true"`       // Requires payment

	// Metadata
	CreatedBy      uint64 `json:"created_by"`
	UpdatedBy      uint64 `json:"updated_by"`
	OrganizationID uint64 `json:"organization_id"`
}

func (Package) TableName() string {
	return "packages"
}