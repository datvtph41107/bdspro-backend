package repository

import (
	base_enum "base/enum"
	"context"
	"mdn4/modules/crm/dto"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/crm/enums"
	"time"
)

type ContactRepo interface {
	Search(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error)
	// SearchByOwnerOf lists contacts for a shared owner_of pool (admin sales; no owner_id filter).
	SearchByOwnerOf(c context.Context, ownerType base_enum.EOwnerOf, dto dto.ContactSearchDTO) ([]domain.Contact, int64, error)
	Create(c context.Context, entity *domain.Contact) (*domain.Contact, error)
	Update(c context.Context, id uint64, entity *domain.Contact) (*domain.Contact, error)
	UpdateContactID(c context.Context, contactId uint64, leadId uint64) (*domain.Contact, error)
	Delete(c context.Context, id uint64) error
	GetByID(c context.Context, id uint64) (*domain.Contact, error)
	GetByProfileID(c context.Context, profileId uint64, currentUserId uint64, ownerType base_enum.EOwnerOf) (*domain.Contact, error)
	GetOrCreateByProfileID(c context.Context, profileId uint64, ownerId uint64, ownerType base_enum.EOwnerOf, fullName string, phone string, avatar string) (*domain.Contact, error)
	Existed(c context.Context, followId uint64) error
	Bulk(c context.Context, entities []domain.Contact) ([]domain.Contact, error)
	ExistByProfile(c context.Context, profileId uint64) (bool, error)
	UpdateNote(c context.Context, id uint64, note string) error
	GetByPhone(c context.Context, profileId uint64, ownerType base_enum.EOwnerOf, phone string) (*domain.Contact, error)
	GetByPhoneIncludingDeleted(c context.Context, profileId uint64, ownerType base_enum.EOwnerOf, phone string) (*domain.Contact, error)
	Restore(c context.Context, id uint64) error
	GetByOriginProfileId(c context.Context, ownerId uint64, originProfileId uint64) (*domain.Contact, error)
	GetListContactByProductId(c context.Context, productId uint64, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error)
	GetByIDs(c context.Context, ids []uint64) ([]domain.Contact, error)
	GetByProfileIDAndOwnerIdNull(c context.Context, profileId uint64) (*domain.Contact, error)

	ExistsContactRelation(c context.Context, ownerID, contactOriginProfileID uint64) (bool, error)

	GetInterestedContactsByProduct(ctx context.Context, userId *uint64, req *dto.ContactProductInterestedFilter) ([]*dto.ContactProductInterestedDTO, int64, error)
	UpdatePinnedAtDynamic(ctx context.Context, meta enums.PinTargetMeta, targetIDs map[string]uint64, pinAt *time.Time) error
}
