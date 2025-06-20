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

func (r *PostgreContact) QuerySearch(tx *gorm.DB, profileId uint64, ownerType enums.EOwnerType, dto dto.ContactSearchDTO) (*gorm.DB, error) {
	tx = tx.Model(&domain.ContactItem{}).
		Debug().
		Table("tb_contact").
		Joins("LEFT JOIN block bl ON (bl.profile_id = ? AND bl.blocked_id = tb_contact.profile_id)", profileId).
		Joins(`LEFT JOIN friend fr ON (
		(fr.created_by = tb_contact.profile_id and fr.receiver_id = ?) 
		OR (fr.receiver_id = tb_contact.profile_id and fr.created_by = ?)) 
		and fr.deleted_at is null`,
			profileId, profileId).
		Joins("LEFT JOIN follow fl ON fl.following_id = tb_contact.profile_id and fl.created_by = ? and fl.deleted_at is null", profileId).
		Select(`tb_contact.*, 
		fl.id IS NOT NULL AS following,
		fr.id AS friend_id,
		fr.status AS friend_status,
		fr.created_by AS friend_created_by,
		fr.receiver_id AS friend_receiver_id,
		bl.profile_id IS NOT NULL AS blocked,
		CASE 
			WHEN fr.status = 20 THEN 20
			WHEN tb_contact.profile_id is not null THEN 30
			ELSE 10
		END AS status`).
		Where("tb_contact.deleted_at IS NULL AND tb_contact.owner_id = ? and tb_contact.owner_type = ?", profileId, ownerType)

	if dto.FriendStatus != nil {
		tx = tx.Where("fr.status IN (?)", dto.FriendStatus)
	}
	// if dto.Following != nil {
	// 	tx = tx.Where("following = ?", *dto.Following)
	// }
	// if dto.Roles != nil {
	// 	tx = tx.Where("tb_contact.roles IN (?)", dto.Roles)
	// }
	// if dto.Leaded != nil {
	// 	tx = tx.Where("tb_contact.lead_id IS NOT NULL = ?", *dto.Leaded)
	// }
	// if dto.Connected != nil {
	// 	tx = tx.Where("tb_contact.connected = ?", *dto.Connected)
	// }
	if dto.AppInstalled != nil {
		if *dto.AppInstalled {
			tx = tx.Where("tb_contact.profile_id is not null")
		} else {
			tx = tx.Where("tb_contact.profile_id is null")
		}
	}

	return tx, nil
}

func (r *PostgreContact) Search(c context.Context, profileId uint64, ownerType enums.EOwnerType, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error) {
	var contacts []domain.ContactItem
	query, err := r.QuerySearch(r.DB, profileId, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}
	result := query.
		// Model(&domain.ContactEntity{}).
		Debug().
		// Where("tb_contact.deleted_at IS NULL AND tb_contact.owner_id = ? and tb_contact.owner_type = ?", profileId, ownerType).
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Scan(&contacts)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	var total int64
	query.Count(&total)
	return contacts, total, nil
}

func (r *PostgreContact) Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	result := r.DB.WithContext(c).Create(entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (r *PostgreContact) Update(c context.Context, id uint64, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	result := r.DB.WithContext(c).
		Model(&domain.ContactEntity{}).
		Where("id = ?", id).
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
	result := r.DB.Model(&domain.ContactEntity{}).Where("id = ? and deleted_at is null", id).First(&contact)
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
			DoUpdates: clause.AssignmentColumns([]string{"full_name", "profile_id"}), // thêm "updated_at" nếu cần
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

func (r *PostgreContact) UpdateNote(c context.Context, id uint64, note string) error {
	result := r.DB.Model(&domain.ContactEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("note", note)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgreContact) GetByPhone(c context.Context, profileId uint64, ownerType enums.EOwnerType, phone string) (*domain.ContactEntity, error) {
	var contact domain.ContactEntity
	result := r.DB.Model(&domain.ContactEntity{}).
		Where("phone = ? and deleted_at is null and owner_id = ? and owner_type = ?", phone, profileId, ownerType).
		First(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

func (r *PostgreContact) UpdateContactID(c context.Context, contactId uint64, leadId uint64) (*domain.ContactEntity, error) {
	result := r.DB.WithContext(c).
		Model(&domain.ContactEntity{}).
		Where("id = ? and deleted_at is null", contactId).
		Update("lead_id", leadId)
	if result.Error != nil {
		return nil, result.Error
	}
	return nil, nil
}
