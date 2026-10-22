package services

import (
	"context"

	"mdn4/modules/crm/dto"
	crmpb "pb/types/crm"
)

// ContactService is the local CRM contact surface used by HTTP controllers and
// demo seeding (ListUser + create/delete helpers).
type ContactService interface {
	ListUser(ctx context.Context, ownerID uint64, req dto.ContactSearchRequest) (*crmpb.ContactListDTO, error)
	CreateContact(ctx context.Context, fullName, phone, email, avatar string, linkedProfileID uint64) (originID, contactID, profileID uint64, err error)
	FindContactIDByPhone(ctx context.Context, phone string) (uint64, error)
	GetContactProfileID(ctx context.Context, contactID uint64) (uint64, error)
	GetOriginID(ctx context.Context, profileID uint64) (uint64, error)
	DeleteContact(ctx context.Context, contactID uint64) error
}
