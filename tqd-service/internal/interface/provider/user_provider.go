package provider

import (
	"context"
	"tqd/internal/dto"
)

// UserProvider defines the interface for user-related operations
type UserProvider interface {
	// GetProfileByID gets user profile by ID
	GetProfileByID(ctx context.Context, userID uint64) (*dto.UserProfileDTO, error)

	// GetProfilesByIDs gets multiple user profiles by IDs
	GetProfilesByIDs(ctx context.Context, userIDs []uint64) (map[uint64]*dto.UserProfileDTO, error)

	// GetProfileIdWithContext gets profile ID from context
	GetProfileIdWithContext(ctx context.Context) uint64

	// GetOrganizationIdFromContext gets organization ID from context
	GetOrganizationIdFromContext(ctx context.Context) uint64

	// GetUserInfo gets basic user information
	GetUserInfo(ctx context.Context, userID uint64) (*dto.UserInfoDTO, error)

	// GetUserOrganization gets user's organization information
	GetUserOrganization(ctx context.Context, userID uint64) (*dto.UserOrganizationDTO, error)
}
