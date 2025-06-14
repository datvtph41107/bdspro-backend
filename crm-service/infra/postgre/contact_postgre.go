package postgre

import (
	_routes "common/routes"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @bind: crm/internal/repo.ContactRepo
type PostgreContact struct {
	DB *gorm.DB
}

func NewContactRepo(DB *gorm.DB) *PostgreContact {
	return &PostgreContact{
		DB: DB,
	}
}

func (r *PostgreContact) Search(c context.Context, profileId uint64, ownerType enums.EOwnerType, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error) {
	var contacts []domain.ContactItem
	result := r.DB.
		// Model(&domain.ContactEntity{}).
		Debug().
		Table("tb_contact").
		Joins("LEFT JOIN block bl ON (bl.profile_id = ? AND bl.blocked_id = tb_contact.profile_id)", profileId).
		Joins("LEFT JOIN friend fr ON (fr.created_by = tb_contact.profile_id OR fr.receiver_id = tb_contact.profile_id) and fr.deleted_at is null").
		Joins("LEFT JOIN follow fl ON fl.following_id = tb_contact.profile_id and fl.created_by = ? and fl.deleted_at is null", profileId).
		Select(`tb_contact.*, 
		fl.id IS NOT NULL AS following, 
		fr.id AS friend_id,
		fr.status AS friend_status,
		fr.created_by AS friend_created_by,
		fr.receiver_id AS friend_receiver_id,
		bl.profile_id IS NOT NULL AS blocked
		`).
		Where("tb_contact.deleted_at IS NULL AND tb_contact.owner_id = ? and tb_contact.owner_type = ?", profileId, ownerType).
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Scan(&contacts)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	var total int64
	r.DB.Model(&domain.ContactEntity{}).Where("deleted_at IS NULL AND owner_id = ? and owner_type = ?", profileId, ownerType).Count(&total)
	return contacts, total, nil
}

func (r *PostgreContact) Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	result := r.DB.Create(entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (r *PostgreContact) Update(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	result := r.DB.WithContext(c).
		Model(&domain.ContactEntity{}).
		Where("id = ?", entity.ID).
		Update("full_name", entity.FullName).
		Update("phone", entity.Phone).
		Update("avatar", entity.Avatar).
		Update("email", entity.Email).
		Update("zalo", entity.Zalo).
		Update("company", entity.Company).
		Update("note", entity.Note).
		Update("visibility", entity.Visibility)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (r *PostgreContact) Delete(c context.Context, id uint64) error {
	result := r.DB.WithContext(c).
		Model(&domain.ContactEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgreContact) GetByID(c context.Context, id uint64) (*domain.ContactEntity, error) {
	var contact domain.ContactEntity
	result := r.DB.Model(&domain.ContactEntity{}).Where("id = ?", id).First(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

// ExistsBy kiểm tra xem một ContactEntity có tồn tại không dựa trên điều kiện
func (r *PostgreContact) ExistByProfile(c context.Context, profileId uint64) (bool, error) {
	var count int64
	result := r.DB.Model(&domain.Profile{}).Where("profile_id = ?", profileId).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *PostgreContact) Existed(c context.Context, followId uint64) error {
	var count int64
	result := r.DB.Model(&domain.Profile{}).Where("profile_id = ?", followId).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	// ok, err := r.contactRepo.ExistByProfile(followId)
	if count == 0 {
		return &_routes.Except{
			Code:    400,
			Message: "Người dùng không tồn tại",
		}
	}

	// if err != nil {
	// 	return &_routes.Except{
	// 		Code:    500,
	// 		Message: err.Error(),
	// 	}
	// }
	return nil
}

// Bulk inserts all contacts.
// If a phone number already exists, we update the contact's name (and updated_at)
// instead of inserting a duplicate.
func (r *PostgreContact) Bulk(
	ctx context.Context,
	entities []domain.ContactEntity,
) ([]domain.ContactEntity, error) {

	if len(entities) == 0 {
		return entities, nil
	}

	// INSERT ... ON CONFLICT (phone) DO UPDATE SET name = EXCLUDED.name, updated_at = NOW()
	result := r.DB.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{
				{Name: "phone"},
				{Name: "owner_id"},
				{Name: "owner_type"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"full_name"}), // thêm "updated_at" nếu cần
		},
	).Create(&entities)

	if result.Error != nil {
		return nil, result.Error
	}

	return entities, nil
}

func (r *PostgreContact) GetByProfileID(c context.Context, profileId uint64, currentUserId uint64, ownerType enums.EOwnerType) (*domain.ContactEntity, error) {
	var contact domain.ContactEntity
	result := r.DB.Model(&domain.ContactEntity{}).
		Where("profile_id = ? and deleted_at is null and owner_id = ? and owner_type = ?",
			profileId, currentUserId, ownerType).
		First(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}
