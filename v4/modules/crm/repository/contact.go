package repository

import (
	"context"
	"fmt"
	"time"

	base_enum "base/enum"
	_utils "common/utils"
	"mdn4/modules/crm/entity"

	"gorm.io/gorm"
)

// ContactRepository is the CRM contact store on the CRM DSN.
type ContactRepository interface {
	SearchByOwner(ctx context.Context, ownerID uint64, ownerOf base_enum.EOwnerOf, page, size uint32, text string) ([]entity.ContactListItem, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.Contact, error)
	GetFriendBetween(ctx context.Context, profileA, profileB uint64) (*entity.FriendRelation, error)
	ListProductIDs(ctx context.Context, contactID uint64) ([]uint64, error)
	Create(ctx context.Context, contact *entity.Contact) error
	SoftDelete(ctx context.Context, id uint64) error
	FindIDByPhone(ctx context.Context, ownerID uint64, ownerOf base_enum.EOwnerOf, phone string) (uint64, error)
	GetOriginIDByProfileID(ctx context.Context, profileID uint64) (uint64, error)
}

type contactPostgres struct{ db *gorm.DB }

func NewContactPostgres(db *gorm.DB) ContactRepository {
	if db == nil {
		return &contactPostgres{}
	}
	return &contactPostgres{db: db}
}

func (r *contactPostgres) SearchByOwner(
	ctx context.Context,
	ownerID uint64,
	ownerOf base_enum.EOwnerOf,
	page, size uint32,
	text string,
) ([]entity.ContactListItem, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	profileID := _utils.GetProfileIdWithContext(ctx)

	baseQuery := r.db.WithContext(ctx).Table("tb_contact tb").
		Joins("LEFT JOIN friend fr ON fr.id = tb.friend_id").
		Joins("LEFT JOIN follow fl ON fl.following_id = tb.profile_id and fl.created_by = ? and fl.deleted_at is null", profileID).
		Where("tb.deleted_at IS NULL AND tb.owner_id = ? AND tb.owner_of = ?", ownerID, ownerOf)

	if text != "" {
		like := "%" + text + "%"
		baseQuery = baseQuery.Where(
			"(tb.full_name ILIKE ? OR tb.phone ILIKE ? OR tb.email ILIKE ?)",
			like, like, like,
		)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count contacts: %w", err)
	}

	offset := 0
	if page >= 1 {
		offset = int(page) * int(size)
	}
	limit := int(size)
	if limit <= 0 {
		limit = 20
	}

	var contacts []entity.ContactListItem
	err := baseQuery.
		Select(`tb.*,
			fr.id as friend_id,
			fr.status as friend_status,
			fr.created_by as friend_created_by,
			fr.receiver_id as friend_receiver_id,
			fl.id IS NOT NULL AS following
		`).
		Order("CASE WHEN tb.profile_id IS NOT NULL THEN 0 ELSE 1 END, tb.created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(&contacts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list contacts: %w", err)
	}
	return contacts, total, nil
}

func (r *contactPostgres) GetByID(ctx context.Context, id uint64) (*entity.Contact, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var contact entity.Contact
	err := r.db.WithContext(ctx).
		Table("tb_contact").
		Joins("LEFT JOIN follow fl ON fl.following_id = tb_contact.profile_id and fl.created_by = ? and fl.deleted_at is null", profileID).
		Where("tb_contact.id = ? AND tb_contact.deleted_at IS NULL", id).
		First(&contact).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get contact: %w", err)
	}
	return &contact, nil
}

func (r *contactPostgres) GetFriendBetween(ctx context.Context, profileA, profileB uint64) (*entity.FriendRelation, error) {
	if r.db == nil || profileA == 0 || profileB == 0 {
		return nil, nil
	}
	var rel entity.FriendRelation
	err := r.db.WithContext(ctx).
		Table("friend").
		Select("id, status, created_by, receiver_id").
		Where("((created_by = ? AND receiver_id = ?) OR (created_by = ? AND receiver_id = ?)) AND deleted_at IS NULL",
			profileA, profileB, profileB, profileA).
		Limit(1).
		Scan(&rel).Error
	if err != nil {
		return nil, fmt.Errorf("get friend: %w", err)
	}
	if rel.ID == 0 {
		return nil, nil
	}
	return &rel, nil
}

func (r *contactPostgres) ListProductIDs(ctx context.Context, contactID uint64) ([]uint64, error) {
	if r.db == nil || contactID == 0 {
		return nil, nil
	}
	var productIDs []uint64
	err := r.db.WithContext(ctx).
		Table("tb_contact_product").
		Select("product_id").
		Where("contact_id = ?", contactID).
		Pluck("product_id", &productIDs).Error
	if err != nil {
		return nil, fmt.Errorf("list contact products: %w", err)
	}
	return productIDs, nil
}

// Create inserts a new contact row. If ProfileID is set and OriginProfileID is empty,
// it fills OriginProfileID from origin_profile.
func (r *contactPostgres) Create(ctx context.Context, contact *entity.Contact) error {
	if r.db == nil {
		return nil
	}
	// Backfill OriginProfileID from origin_profile when linked profile is known.
	if contact.ProfileID != nil && *contact.ProfileID > 0 && (contact.OriginProfileID == nil || *contact.OriginProfileID == 0) {
		var originID uint64
		if err := r.db.WithContext(ctx).
			Table("origin_profile").
			Select("origin_id").
			Where("owner_id = ? AND owner_of = 10 AND deleted_at IS NULL", *contact.ProfileID).
			Limit(1).
			Pluck("origin_id", &originID).Error; err == nil && originID > 0 {
			contact.OriginProfileID = &originID
		}
	}
	if err := r.db.WithContext(ctx).Table("tb_contact").Create(contact).Error; err != nil {
		return fmt.Errorf("create contact: %w", err)
	}
	return nil
}

// SoftDelete sets deleted_at on a contact row.
func (r *contactPostgres) SoftDelete(ctx context.Context, id uint64) error {
	if r.db == nil || id == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).
		Table("tb_contact").
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

// FindIDByPhone returns the contact id for the owner+phone combination.
func (r *contactPostgres) FindIDByPhone(ctx context.Context, ownerID uint64, ownerOf base_enum.EOwnerOf, phone string) (uint64, error) {
	if r.db == nil || ownerID == 0 || phone == "" {
		return 0, nil
	}
	var id uint64
	err := r.db.WithContext(ctx).
		Table("tb_contact").
		Select("id").
		Where("owner_id = ? AND owner_of = ? AND phone = ? AND deleted_at IS NULL", ownerID, ownerOf, phone).
		Limit(1).
		Pluck("id", &id).Error
	if err != nil {
		return 0, fmt.Errorf("find contact by phone: %w", err)
	}
	return id, nil
}

// GetOriginIDByProfileID returns the origin_id from origin_profile for a given profile_id.
func (r *contactPostgres) GetOriginIDByProfileID(ctx context.Context, profileID uint64) (uint64, error) {
	if r.db == nil || profileID == 0 {
		return 0, nil
	}
	var originID uint64
	err := r.db.WithContext(ctx).
		Table("origin_profile").
		Select("origin_id").
		Where("owner_id = ? AND owner_of = 10 AND deleted_at IS NULL", profileID).
		Limit(1).
		Pluck("origin_id", &originID).Error
	if err != nil {
		return 0, fmt.Errorf("get origin id: %w", err)
	}
	return originID, nil
}
