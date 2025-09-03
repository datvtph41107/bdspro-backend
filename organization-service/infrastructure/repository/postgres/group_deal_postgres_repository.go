package postgres

import (
	_dto "common/dto"
	"context"
	"time"

	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"

	"gorm.io/gorm"
)

// @bind: organization/internal/domain/repository.GroupDealRepository
type GroupDealPostgresRepository struct {
	db              *gorm.DB
	bankAccountRepo repository.BankAccountRepository
	transactionRepo iusecase.ITransaction
}

func NewGroupDealPostgresRepository(db *gorm.DB,
	bankAccountRepo repository.BankAccountRepository,
	transactionRepo iusecase.ITransaction,
) *GroupDealPostgresRepository {
	return &GroupDealPostgresRepository{
		db:              db,
		bankAccountRepo: bankAccountRepo,
		transactionRepo: transactionRepo,
	}
}

func (r *GroupDealPostgresRepository) CreateDealProduct(ctx context.Context, dealProduct []*entity.DealProduct) ([]*entity.DealProduct, error) {
	if err := GetDB(ctx, r.db).Create(dealProduct).Error; err != nil {
		return nil, err
	}
	return dealProduct, nil
}

func (r *GroupDealPostgresRepository) GetDealProducts(ctx context.Context, dealID uint64, pagable _dto.Pagable) ([]*entity.DealProduct, error) {
	var dealProducts []*entity.DealProduct
	if err := r.db.WithContext(ctx).
		Where("deal_id = ?", dealID).
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&dealProducts).Error; err != nil {
		return nil, err
	}
	return dealProducts, nil
}

func (r *GroupDealPostgresRepository) Create(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	// model := GroupDealEntityToModel(deal)
	deal.Products = make([]entity.DealProduct, len(deal.ProductIds))
	for i, productID := range deal.ProductIds {
		deal.Products[i] = entity.DealProduct{
			DealID:    deal.ID,
			ProductID: productID,
		}
	}
	if err := GetDB(ctx, r.db).Create(deal).Error; err != nil {
		return nil, err
	}
	return deal, nil
}

func (r *GroupDealPostgresRepository) Update(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	if err := r.db.WithContext(ctx).Model(&entity.Deal{}).Where("id = ?", deal.ID).Updates(deal).Error; err != nil {
		return nil, err
	}
	return deal, nil
}

func (r *GroupDealPostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *GroupDealPostgresRepository) GetByID(ctx context.Context, id uint64) (*entity.Deal, error) {
	var model entity.Deal
	if err := GetDB(ctx, r.db).
		Preload("BankAccount", "deleted_at is null").
		Where("id = ? AND is_deleted = ?", id, false).
		First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *GroupDealPostgresRepository) DetailByID(ctx context.Context, id uint64) (*entity.Deal, error) {
	var model entity.Deal
	if err := GetDB(ctx, r.db).
		Debug().
		Preload("Products").
		Preload("Members", "member_type = 30").
		Preload("Customers", "member_type = 10").
		Preload("Partners", "member_type = 20").
		Preload("BankAccount", "deleted_at is null").
		Preload("Documents", "deleted_at is null and doc_owner = 10").
		Where("id = ? AND deleted_at is null", id).
		First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *GroupDealPostgresRepository) GetByGroupID(ctx context.Context, groupID uint64, page, size int) ([]*entity.Deal, uint32, error) {
	var deals []*entity.Deal
	var total int64

	offset := (page - 1) * size

	// Count total
	if err := r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("owner_id = ? AND owner_type = ? AND is_deleted = ?", groupID, enums.OwnerTypeGroup, false).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get deals with pagination
	if err := r.db.WithContext(ctx).
		Where("owner_id = ? AND owner_type = ? AND is_deleted = ?", groupID, enums.OwnerTypeGroup, false).
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	return deals, uint32(total), nil
}

func (r *GroupDealPostgresRepository) GetByOrganizationID(ctx context.Context, organizationID uint64, page, size int) ([]*entity.Deal, uint32, error) {
	var deals []*entity.Deal
	var total int64

	offset := (page - 1) * size

	// Count total
	if err := r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("owner_id = ? AND owner_type = ? AND is_deleted = ?", organizationID, enums.OwnerTypeOrganization, false).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get deals with pagination
	if err := r.db.WithContext(ctx).
		Where("owner_id = ? AND owner_type = ? AND is_deleted = ?", organizationID, enums.OwnerTypeOrganization, false).
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	return deals, uint32(total), nil
}

func (r *GroupDealPostgresRepository) UpdateStatus(ctx context.Context, id uint64, status enums.DealStatus) error {
	return r.db.WithContext(ctx).Model(&entity.Deal{}).Where("id = ? AND is_deleted = ?", id, false).Update("status", status).Error
}

func (r *GroupDealPostgresRepository) UpdateStatusAndCancelReason(ctx context.Context, id uint64, status enums.DealStatus, cancelReason string) error {
	return r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Updates(map[string]interface{}{
			"status":        status,
			"cancel_reason": cancelReason,
		}).Error
}

func (r *GroupDealPostgresRepository) UpdateSetting(ctx context.Context, id uint64, setting *dto.DealSetting) error {
	// tx := r.db.WithContext(ctx).Begin()
	err := r.transactionRepo.WithTransaction(ctx, func(ctx context.Context) error {
		var model *entity.Deal
		err := GetDB(ctx, r.db).Model(&entity.Deal{}).
			Preload("BankAccount", "deleted_at is null").
			Where("id = ? AND is_deleted = ?", id, false).
			First(&model).Error
		if err != nil {
			return err
		}

		bankAccount, err := r.bankAccountRepo.DealUpdateAccount(ctx, model.BankAccountID, &entity.BankAccount{
			BankId:        setting.BankId,
			AccountNumber: setting.AccountNumber,
			AccountName:   setting.AccountName,
		})
		if err != nil {
			return err
		}

		err = GetDB(ctx, r.db).
			Model(&entity.Deal{}).
			Where("id = ? AND is_deleted = ?", id, false).
			Updates(map[string]interface{}{
				"from_date":                  setting.FromDate,
				"to_date":                    setting.ToDate,
				"allow_sharing":              setting.AllowSharing,
				"member_can_add_transaction": setting.MemberCanAddTransaction,
				"only_owner_get_commission":  setting.OnlyOwnerGetCommission,
				"internal_note":              setting.InternalNote,
				"owner_id":                   setting.OwnerId,
				"status":                     setting.Status,
				"allow_manual_input":         setting.AllowManualInput,
				"bank_account_id":            bankAccount.ID,
			}).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *GroupDealPostgresRepository) UpdateAllowManualInput(ctx context.Context, dealID uint64, allowManualInput bool) error {
	return r.db.WithContext(ctx).Model(&entity.Deal{}).Where("id = ?", dealID).Update("allow_manual_input", allowManualInput).Error
}

// GetDealProcess lấy thông tin process của deal
func (r *GroupDealPostgresRepository) GetDealProcess(ctx context.Context, dealID uint64, organizationID uint64) ([]*dto.DealProcessItem, error) {
	// TODO: Implement logic để lấy thông tin process từ database
	// Tạm thời trả về dữ liệu mẫu
	processItems := []*dto.DealProcessItem{
		{
			StatusName:    "Đang xử lý",
			Timestamp:     "2024-01-15T10:30:00Z",
			Code:          "PROC001",
			TransactionId: "TXN123456",
			Color:         "#FFA500",
		},
		{
			StatusName:    "Hoàn thành",
			Timestamp:     "2024-01-15T11:00:00Z",
			Code:          "COMP002",
			TransactionId: "TXN123457",
			Color:         "#00FF00",
		},
	}

	return processItems, nil
}
