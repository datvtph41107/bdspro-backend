package services

import (
	"context"

	crmpb "pb/types/crm"

	"userv4/modules/crm/dto"
)

// ContactService owns user-facing CRM contact flows for modules/crm.
type ContactService interface {
	// ListUser returns contacts owned by the authenticated user (ownerOf=user).
	// Path ownerId is ignored when 0; owner resolves from JWT profile.
	ListUser(ctx context.Context, ownerID uint64, req dto.ContactSearchRequest) (*crmpb.ContactListDTO, error)
	// CreateContact creates a personal contact (ownerOf=member). Returns originID, contactID, profileID.
	// linkedProfileID may be 0 when the contact has no linked profile.
	CreateContact(ctx context.Context, fullName, phone, email, avatar string, linkedProfileID uint64) (originID, contactID, profileID uint64, err error)
	// FindContactIDByPhone finds the contact id owned by the current user for the given phone.
	FindContactIDByPhone(ctx context.Context, phone string) (uint64, error)
	// GetContactProfileID returns the profileId linked to a contact.
	GetContactProfileID(ctx context.Context, contactID uint64) (uint64, error)
	// GetOriginID returns the origin_id for a profile from origin_profile.
	GetOriginID(ctx context.Context, profileID uint64) (uint64, error)
	// DeleteContact soft-deletes a contact by id.
	DeleteContact(ctx context.Context, contactID uint64) error
}
