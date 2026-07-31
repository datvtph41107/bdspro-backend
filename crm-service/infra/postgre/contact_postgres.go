package postgre

import (
	base_enum "base/enum"
	base_util "base/utils"
	_db "common/db"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.ContactRepo
type PostgreContact struct {
	*_db.BaseRepo
}

func NewContactRepo(db *gorm.DB) *PostgreContact {
	return &PostgreContact{
		BaseRepo: _db.NewBaseRepo(db),
	}
}

func (r *PostgreContact) UpdatePinnedAtDynamic(
	ctx context.Context,
	meta enums.PinTargetMeta,
	targetIDs map[string]uint64,
	pinAt *time.Time,
) error {
	db := r.GetDB(ctx).Table(meta.TableName)
	for _, col := range meta.IDColumns {
		val, ok := targetIDs[col]
		if !ok {
			return errors.New("missing target id: " + col)
		}
		db = db.Where(fmt.Sprintf("%s = ?", col), val)
	}
	return db.Update(meta.PinColumn, pinAt).Error
}

func (r *PostgreContact) GetInterestedContactsByProduct(
	ctx context.Context,
	userId *uint64,
	req *dto.ContactProductInterestedFilter,
) ([]*dto.ContactProductInterestedDTO, int64, error) {
	var (
		items []*dto.ContactProductInterestedDTO
		total int64
	)
	db := r.GetDB(ctx).
		Table("tb_contact_product cp").
		Where("cp.product_id = ?", req.ProductID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Select(`
			cp.contact_id,
			cp.origin_profile_id,

			COALESCE(c.full_name, op.display_name) AS display_name,
			COALESCE(c.phone, op.phone) AS phone,
			COALESCE(c.avatar, op.avatar) AS avatar,

			CASE
				WHEN c.id IS NULL THEN ?
				ELSE c.status
			END AS status,

			CASE
				WHEN c.id IS NULL THEN false
				ELSE true
			END AS has_contact,

			cp.priority_pin_at IS NOT NULL AS pinned,
			cp.created_at
		`, enums.StatusContactNotSaved).
		Joins(`LEFT JOIN tb_contact c 
			ON c.id = cp.contact_id 
			AND c.owner_id = ?`, *userId).
		Joins(`LEFT JOIN origin_profile op 
			ON op.origin_id = cp.origin_profile_id`).
		Order("cp.priority_pin_at DESC NULLS LAST").
		Order("cp.last_interaction_at DESC NULLS LAST").
		Order("cp.created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Scan(&items).Error
	return items, total, err
}

func (r *PostgreContact) QuerySearch(tx *gorm.DB, profileId uint64, ownerType base_enum.EOwnerOf, dto dto.ContactSearchDTO) (*gorm.DB, error) {
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
		COALESCE(fr.status, 40) AS friend_status,
		fr.created_by AS friend_created_by,
		fr.receiver_id AS friend_receiver_id,
		bl.profile_id IS NOT NULL AS blocked,
		CASE 
			WHEN fr.status = 20 THEN 20
			WHEN tb_contact.profile_id is not null THEN 30
			ELSE 10
		END AS status`).
		Where("tb_contact.deleted_at IS NULL AND tb_contact.owner_id = ? and tb_contact.owner_of = ?", profileId, ownerType)

	if dto.Text != "" {
		tx = tx.Where("(tb_contact.full_name ILIKE ? OR tb_contact.phone ILIKE ? OR tb_contact.email ILIKE ?)",
			"%"+dto.Text+"%", "%"+dto.Text+"%", "%"+dto.Text+"%")
	}

	if len(dto.FriendStatus) > 0 {
		tx = tx.Where("COALESCE(fr.status, 40) IN (?)", dto.FriendStatus)
	}
	if dto.Following != nil {
		if *dto.Following {
			tx = tx.Where("fl.id IS NOT NULL")
		} else {
			tx = tx.Where("fl.id IS NULL")
		}
	}
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

func (r *PostgreContact) Search(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.ContactSearchDTO) ([]domain.ContactItem, int64, error) {
	var contacts []domain.ContactItem
	db := r.GetDB(c)
	profileId := _utils.GetProfileIdWithContext(c)

	baseQuery := db.Table("tb_contact tb").
		Joins("LEFT JOIN origin_profile op ON tb.origin_profile_id = op.origin_id").
		Joins("LEFT JOIN friend fr ON fr.id = tb.friend_id").
		Joins("LEFT JOIN follow fl ON fl.following_id = tb.profile_id and fl.created_by = ? and fl.deleted_at is null", profileId).
		Where("tb.deleted_at IS NULL AND tb.owner_id = ? AND tb.owner_of = ?", ownerId, ownerType)

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Select(`tb.*,
			op.origin_id as origin_id,
			op.display_name as origin_display_name,
			op.phone as origin_phone,
			op.avatar as origin_avatar,
			op.owner_id as origin_owner_id,
			op.owner_of as origin_owner_of,
			fr.id as friend_id,
			fr.status as friend_status,
			fr.created_by as friend_created_by,
			fr.receiver_id as friend_receiver_id,
			fl.id IS NOT NULL AS following
			`).
		Order("CASE WHEN tb.profile_id IS NOT NULL THEN 0 ELSE 1 END, tb.created_at DESC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Scan(&contacts).Error
	if err != nil {
		return nil, 0, err
	}

	for idx, ct := range contacts {
		status := enums.StatusContactUnknown
		if ct.FriendStatus == uint64(enums.FriendStatusAccepted) {
			status = enums.StatusContactFriend
		} else if ct.ProfileID != nil {
			status = enums.StatusContactAppInstalled
		}
		contacts[idx].Status = status
	}

	return contacts, total, nil
}

func (r *PostgreContact) SearchByOwnerOf(c context.Context, ownerType base_enum.EOwnerOf, dto dto.ContactSearchDTO) ([]domain.ContactEntity, int64, error) {
	query := r.GetDB(c).Model(&domain.ContactEntity{}).
		Where("deleted_at is null and owner_of = ?", ownerType)
	if strings.TrimSpace(dto.Text) != "" {
		q := "%" + strings.TrimSpace(dto.Text) + "%"
		query = query.Where("(full_name ILIKE ? OR phone ILIKE ? OR email ILIKE ? OR company ILIKE ?)", q, q, q, q)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var contacts []domain.ContactEntity
	err := query.Order("updated_at desc, id desc").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&contacts).Error
	return contacts, total, err
}

func (r *PostgreContact) GetListContactByProductId(
	c context.Context,
	productId uint64,
	ownerId uint64,
	ownerType base_enum.EOwnerOf,
	dto dto.ContactSearchDTO,
) ([]domain.ContactItem, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)

	db := r.GetDB(c)
	baseQuery := db.
		Table("tb_contact").
		Joins(`
			INNER JOIN tb_contact_product tcp 
				ON tcp.contact_id = tb_contact.id 
				AND tcp.product_id = ?
		`, productId).
		Where(`
			tb_contact.deleted_at IS NULL
			AND tb_contact.owner_id = ?
			AND tb_contact.owner_of = ?
		`, ownerId, ownerType)

	if dto.Text != "" {
		baseQuery = baseQuery.Where(`
			(tb_contact.full_name ILIKE ?
			OR tb_contact.phone ILIKE ?
			OR tb_contact.email ILIKE ?)
		`,
			"%"+dto.Text+"%",
			"%"+dto.Text+"%",
			"%"+dto.Text+"%",
		)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := baseQuery.
		Joins(`
			LEFT JOIN LATERAL (
				SELECT *
				FROM friend fr
				WHERE (
					(created_by = tb_contact.profile_id AND receiver_id = ?)
					OR
					(receiver_id = tb_contact.profile_id AND created_by = ?)
				)
				AND deleted_at IS NULL
				LIMIT 1
			) fr ON true
		`, profileId, profileId).
		Joins(`
			LEFT JOIN LATERAL (
				SELECT id
				FROM follow
				WHERE following_id = tb_contact.profile_id
				AND created_by = ?
				AND deleted_at IS NULL
				LIMIT 1
			) fl ON true
		`, profileId).
		Joins(`
			LEFT JOIN LATERAL (
				SELECT profile_id
				FROM block
				WHERE profile_id = ?
				AND blocked_id = tb_contact.profile_id
				LIMIT 1
			) bl ON true
		`, profileId).
		Select(`
			tb_contact.*,
			tcp.updated_at AS contact_product_updated_at,
			(fl.id IS NOT NULL) AS following,
			fr.id AS friend_id,
			COALESCE(fr.status, 40) AS friend_status,
			fr.created_by AS friend_created_by,
			fr.receiver_id AS friend_receiver_id,
			(bl.profile_id IS NOT NULL) AS blocked,
			CASE 
				WHEN fr.status = 20 THEN 20
				WHEN tb_contact.profile_id IS NOT NULL THEN 30
				ELSE 10
			END AS status
		`).
		Order(`
			tcp.updated_at DESC
		`).
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit())

	var contacts []domain.ContactItem
	if err := query.Scan(&contacts).Error; err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func (r *PostgreContact) Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	result := r.GetDB(c).Create(entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (r *PostgreContact) Update(c context.Context, id uint64, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	result := r.GetDB(c).
		Model(&domain.ContactEntity{}).
		Where("id = ?", id).
		Update("full_name", entity.FullName).
		Update("phone", entity.Phone).
		Update("avatar", entity.Avatar).
		Update("email", entity.Email).
		Update("zalo", entity.Zalo).
		Update("company", entity.Company).
		Update("address", entity.Address).
		Update("note", entity.Note).
		Update("visibility", entity.Visibility)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (r *PostgreContact) Delete(c context.Context, id uint64) error {
	result := r.GetDB(c).
		Model(&domain.ContactEntity{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgreContact) GetByID(
	c context.Context,
	id uint64,
) (*domain.ContactEntity, error) {

	var contact domain.ContactEntity
	result := r.GetDB(c).
		Model(&domain.ContactEntity{}).
		Preload("OriginProfile").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&contact)

	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

func (r *PostgreContact) ExistsContactRelation(c context.Context, ownerID uint64, contactOriginProfileID uint64) (bool, error) {
	var exists bool

	err := r.GetDB(c).
		Raw(`
            SELECT EXISTS (
                SELECT 1
                FROM tb_contact
                WHERE owner_id = ?
                  AND origin_profile_id = ?
                  AND origin_profile_id IS NOT NULL
                  AND deleted_at IS NULL
            )
        `, ownerID, contactOriginProfileID).
		Scan(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

// ExistsBy kiểm tra xem một ContactEntity có tồn tại không dựa trên điều kiện
func (r *PostgreContact) ExistByProfile(c context.Context, profileId uint64) (bool, error) {
	var count int64
	result := r.GetDB(c).Model(&domain.Profile{}).Where("profile_id = ?", profileId).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *PostgreContact) Existed(c context.Context, followId uint64) error {
	var count int64
	result := r.GetDB(c).Model(&domain.Profile{}).Where("profile_id = ?", followId).Count(&count)
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
func (r *PostgreContact) Bulk(
	ctx context.Context,
	entities []domain.ContactEntity,
) ([]domain.ContactEntity, error) {

	if len(entities) == 0 {
		return entities, nil
	}

	// INSERT ... ON CONFLICT (phone) DO UPDATE SET name = EXCLUDED.name, updated_at = NOW()
	// result := r.GetDB(ctx).Clauses(
	// 	clause.OnConflict{
	// 		Columns: []clause.Column{
	// 			{Name: "phone"},
	// 			{Name: "owner_id"},
	// 			{Name: "owner_of"},
	// 		},
	// 		DoUpdates: clause.AssignmentColumns([]string{"full_name", "profile_id"}), // thêm "updated_at" nếu cần
	// 	},
	// ).Create(&entities)
	result := r.GetDB(ctx).Create(&entities)

	if result.Error != nil {
		return nil, result.Error
	}

	return entities, nil
}

func (r *PostgreContact) GetByProfileID(c context.Context, profileId uint64, currentUserId uint64, ownerType base_enum.EOwnerOf) (*domain.ContactEntity, error) {
	var contact domain.ContactEntity
	result := r.GetDB(c).Model(&domain.ContactEntity{}).
		Where("profile_id = ? and deleted_at is null and owner_id = ? and owner_of = ?",
			profileId, currentUserId, ownerType).
		First(&contact)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contact, nil
}

// restoreAndUpdateContact khôi phục và cập nhật thông tin contact nếu cần
func (r *PostgreContact) restoreAndUpdateContact(c context.Context, contact *domain.ContactEntity, fullName, phone, avatar string) (*domain.ContactEntity, error) {
	// Khôi phục nếu đã bị xóa
	if contact.DeletedAt != nil {
		if err := r.Restore(c, contact.ID); err != nil {
			return nil, err
		}
		contact.DeletedAt = nil
	}

	// Cập nhật thông tin nếu có thay đổi
	needUpdate := false
	if contact.FullName != fullName && fullName != "" {
		contact.FullName = fullName
		needUpdate = true
	}
	if contact.Avatar != avatar && avatar != "" {
		contact.Avatar = avatar
		needUpdate = true
	}
	if contact.Phone != phone && phone != "" {
		contact.Phone = phone
		needUpdate = true
	}

	if needUpdate {
		return r.Update(c, contact.ID, contact)
	}
	return contact, nil
}

// getContactByOwnerAndIdentifier lấy contact theo ownerId, ownerOf và (phone hoặc profileId) (bao gồm cả đã bị xóa)
func (r *PostgreContact) getContactByOwnerAndIdentifier(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, phone string, profileId *uint64) (*domain.ContactEntity, error) {
	var contact domain.ContactEntity
	query := r.GetDB(c).Model(&domain.ContactEntity{}).
		Where("owner_id = ? and owner_of = ?", ownerId, ownerType)

	// Query theo phone hoặc profileId
	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	if phone != "" {
		clearPhone := base_util.ClearPhone(phone)
		conditions = append(conditions, "phone = ?")
		args = append(args, clearPhone)
	}

	if profileId != nil && *profileId > 0 {
		conditions = append(conditions, "profile_id = ?")
		args = append(args, *profileId)
	}

	if len(conditions) == 0 {
		return nil, nil
	}

	// Sử dụng OR để match phone hoặc profileId
	query = query.Where("("+strings.Join(conditions, " OR ")+")", args...)

	err := query.First(&contact).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &contact, nil
}

func (r *PostgreContact) GetOrCreateByProfileID(c context.Context, profileId uint64, ownerId uint64, ownerType base_enum.EOwnerOf, fullName string, phone string, avatar string) (*domain.ContactEntity, error) {
	// Kiểm tra contact đã tồn tại chưa theo ownerId, ownerOf và (phone hoặc profileId)
	profileIdPtr := &profileId
	contact, err := r.getContactByOwnerAndIdentifier(c, ownerId, ownerType, phone, profileIdPtr)
	if err != nil {
		return nil, err
	}
	if contact != nil && contact.ID > 0 {
		return r.restoreAndUpdateContact(c, contact, fullName, phone, avatar)
	}

	// Nếu không tồn tại, tạo mới
	newContact := &domain.ContactEntity{
		FullName:  fullName,
		Phone:     phone,
		Avatar:    avatar,
		ProfileID: profileIdPtr,
		OwnerID:   ownerId,
		OwnerOf:   ownerType,
	}

	err = r.GetDB(c).Create(newContact).Error
	if err != nil {
		return nil, err
	}

	return newContact, nil
}

func (r *PostgreContact) UpdateNote(c context.Context, id uint64, note string) error {
	result := r.GetDB(c).Model(&domain.ContactEntity{}).
		Where("id = ? and deleted_at is null", id).
		Update("note", note)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgreContact) GetByPhone(c context.Context, profileId uint64, ownerType base_enum.EOwnerOf, phone string) (*domain.ContactEntity, error) {
	clearPhone := base_util.ClearPhone(phone)
	var contact domain.ContactEntity
	result := r.GetDB(c).Model(&domain.ContactEntity{}).
		Where("phone = ? and deleted_at is null and owner_id = ? and owner_of = ?", clearPhone, profileId, ownerType).
		First(&contact)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contact, nil
}

func (r *PostgreContact) GetByPhoneIncludingDeleted(c context.Context, profileId uint64, ownerType base_enum.EOwnerOf, phone string) (*domain.ContactEntity, error) {
	clearPhone := base_util.ClearPhone(phone)
	var contact domain.ContactEntity
	result := r.GetDB(c).Model(&domain.ContactEntity{}).
		Where("phone = ? and owner_id = ? and owner_of = ?", clearPhone, profileId, ownerType).
		First(&contact)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contact, nil
}

func (r *PostgreContact) Restore(c context.Context, id uint64) error {
	result := r.GetDB(c).
		Model(&domain.ContactEntity{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": nil,
		})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgreContact) UpdateContactID(c context.Context, contactId uint64, leadId uint64) (*domain.ContactEntity, error) {
	result := r.GetDB(c).
		Model(&domain.ContactEntity{}).
		Where("id = ? and deleted_at is null", contactId).
		Update("lead_id", leadId)
	if result.Error != nil {
		return nil, result.Error
	}
	return nil, nil
}

func (r *PostgreContact) GetByIDs(c context.Context, ids []uint64) ([]domain.ContactEntity, error) {
	var contacts []domain.ContactEntity
	result := r.GetDB(c).Model(&domain.ContactEntity{}).Where("id IN (?) and deleted_at is null", ids).Find(&contacts)
	if result.Error != nil {
		return nil, result.Error
	}
	return contacts, nil
}

func (r *PostgreContact) GetByProfileIDAndOwnerIdNull(c context.Context, profileId uint64) (*domain.ContactEntity, error) {
	var contact domain.ContactEntity
	result := r.GetDB(c).Model(&domain.ContactEntity{}).
		Where("profile_id = ? and deleted_at is null and owner_id = 0", profileId).
		First(&contact)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contact, nil
}

func (r *PostgreContact) GetByOriginProfileId(c context.Context, ownerId uint64, originProfileId uint64) (*domain.ContactEntity, error) {
	var contact domain.ContactEntity
	result := r.GetDB(c).
		Preload("OriginProfile").
		Where("owner_id = ? AND origin_profile_id = ? AND deleted_at IS NULL", ownerId, originProfileId).
		First(&contact)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contact, nil
}