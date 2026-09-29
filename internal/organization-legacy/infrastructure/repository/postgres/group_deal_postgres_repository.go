package postgres

import (
	_dto "common/domain/dto"
	"context"
	"time"

	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
	iusecase "organization/internal/interface"

	"gorm.io/gorm"
)

// @bind: organization/internal/domain/repository.DealRepository
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
	if err := GetDB(ctx, r.db).Create(&dealProduct).Error; err != nil {
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
	deal.Members = nil
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
	return GetDB(ctx, r.db).Model(&entity.Deal{}).
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

func (r *GroupDealPostgresRepository) GetIdsByGroupID(ctx context.Context, groupID uint64) ([]uint64, error) {
	var deals []uint64
	// Get deal ids
	countQuery := `
		SELECT d.id
		FROM deals d 
		WHERE d.owner_id = ? 
		AND d.owner_type = ? 
		AND d.deleted_at IS NULL
	`
	if err := GetDB(ctx, r.db).Raw(countQuery, groupID, enums.OwnerOfGroup).Pluck("id", &deals).Error; err != nil {
		return nil, err
	}

	return deals, nil
}

func (r *GroupDealPostgresRepository) GetIdsByOrganizationID(ctx context.Context, organizationID uint64) ([]uint64, error) {
	var deals []uint64
	// Get deal ids
	countQuery := `
		SELECT d.id
		FROM deals d 
		WHERE d.owner_id = ? 
		AND d.owner_type = ? 
		AND d.deleted_at IS NULL
	`
	if err := GetDB(ctx, r.db).Raw(countQuery, organizationID, enums.OwnerOfOrganization).Pluck("id", &deals).Error; err != nil {
		return nil, err
	}

	return deals, nil
}

func (r *GroupDealPostgresRepository) GetDeals(ctx context.Context, userID uint64, searchRequest *dto.DealSearchRequest) ([]*entity.Deal, uint32, error) {
	var deals []*entity.Deal
	var total int64

	offset := searchRequest.GetOffset()
	limit := searchRequest.GetLimit()

	// Build base query with subquery to check user access
	baseQuery := `
		SELECT d.* 
		FROM deals d 
		WHERE 1 = 1
		AND (
			d.owner_id = ? 
			OR EXISTS (
				SELECT 1 FROM deal_members dm 
				WHERE dm.deal_id = d.id AND dm.member_id = ?
			)
		)
	`

	args := []interface{}{userID, userID}

	if searchRequest.OwnerID != nil && searchRequest.OwnerType != 0 {
		baseQuery += " AND d.owner_id = ? AND d.owner_type = ?"
		args = append(args, *searchRequest.OwnerID, searchRequest.OwnerType)
	}
	// Add organization filter if provided
	if searchRequest.OrganizationID != nil {
		baseQuery += " AND d.owner_id = ? AND d.owner_type = ?"
		args = append(args, *searchRequest.OrganizationID, enums.OwnerOfOrganization)
	}

	// Add group filter if provided
	if searchRequest.GroupID != nil {
		baseQuery += " AND d.owner_id = ? AND d.owner_type = ?"
		args = append(args, *searchRequest.GroupID, enums.OwnerOfGroup)
	}

	// Add text search if provided
	if searchRequest.Text != "" {
		baseQuery += " AND d.name ILIKE ?"
		args = append(args, "%"+searchRequest.Text+"%")
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") as count_query"
	if err := r.db.WithContext(ctx).Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get deals with pagination
	dataQuery := baseQuery + " ORDER BY d.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	if err := r.db.WithContext(ctx).Raw(dataQuery, args...).Scan(&deals).Error; err != nil {
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

		if !setting.Status.IsValid() {
			setting.Status = model.Status
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
				"status":                     setting.Status,
				"allow_manual_input":         setting.AllowManualInput,
				"bank_account_id":            bankAccount.ID,
				"name":                       setting.Name,
				"target_profit":              setting.TargetProfit,
				// "owner_id":                   setting.OwnerId,
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
			Color:         "#8E51FF",
			BgColor:       "#F5F3FF",
			BorderColor:   "#DDD6FF",
		},
		{
			StatusName:    "Hoàn thành",
			Timestamp:     "2024-01-15T11:00:00Z",
			Code:          "COMP002",
			TransactionId: "TXN123457",
			Color:         "#2B7FFF",
			BgColor:       "#EFF6FF",
			BorderColor:   "#BEDBFF",
		},
		{
			StatusName:  "Đặt cọc",
			Timestamp:   "16/09/2025",
			Code:        "GD0000120",
			Color:       "#FE9A00",
			BgColor:     "#FFFBEB",
			BorderColor: "#FEE685",
		},
		{
			StatusName:  "Ký kết",
			Timestamp:   "16/09/2025",
			Code:        "GD0000120",
			Color:       "#00BBA7",
			BgColor:     "#F0FDFA",
			BorderColor: "#96F7E4",
		},
		{
			StatusName:  "Hoàn tất",
			Timestamp:   "16/09/2025",
			Code:        "GD0000120",
			Color:       "#00A63E",
			BgColor:     "#F0FDF4",
			BorderColor: "#B9F8CF",
		},
	}

	return processItems, nil
}

// GetDealOverviewByOrganization lấy thống kê tổng quan thương vụ theo tổ chức
func (r *GroupDealPostgresRepository) GetDealOverviewByOrganization(ctx context.Context, organizationID uint64) (*dto.DealOverview, error) {
	var result dto.DealOverview

	// Lấy tổng số thương vụ, tổng vốn, đang xử lý, hoàn tất
	err := GetDB(ctx, r.db).
		Debug().
		Model(&entity.Deal{}).
		Select(`
			COUNT(*) as total_deals,
			COALESCE(SUM(target_profit), 0) as total_capital,
			SUM(CASE WHEN status = 10 THEN 1 ELSE 0 END) as total_deal_processing,
			SUM(CASE WHEN status = 20 THEN 1 ELSE 0 END) as total_deal_completed,
			SUM(CASE WHEN status = 30 THEN 1 ELSE 0 END) as total_deal_canceled
		`).
		Where("owner_id = ? AND owner_type = ? AND deleted_at IS NULL", organizationID, enums.OwnerOfOrganization).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// Lợi nhuận ước tính sẽ được lấy từ transaction service
	// Tạm thời set = 0, sẽ được cập nhật trong usecase
	// result.EstimatedProfit = 0

	return &result, nil
}

// GetDealOverviewByGroup lấy thống kê tổng quan thương vụ theo nhóm
func (r *GroupDealPostgresRepository) GetDealOverviewByGroup(ctx context.Context, groupID uint64) (*dto.DealOverview, error) {
	var result dto.DealOverview

	// Lấy tổng số thương vụ, tổng vốn, đang xử lý, hoàn tất
	err := r.db.WithContext(ctx).
		Debug().
		Model(&entity.Deal{}).
		Select(`
			COUNT(*) as total_deals,
			COALESCE(SUM(target_profit), 0) as total_capital,
			SUM(CASE WHEN status = 10 THEN 1 ELSE 0 END) as total_deal_processing,
			SUM(CASE WHEN status = 20 THEN 1 ELSE 0 END) as total_deal_completed,
			SUM(CASE WHEN status = 30 THEN 1 ELSE 0 END) as total_deal_canceled
		`).
		Where("owner_id = ? AND owner_type = ? AND deleted_at IS NULL", groupID, enums.OwnerOfGroup).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// Lợi nhuận ước tính sẽ được lấy từ transaction service
	// Tạm thời set = 0, sẽ được cập nhật trong usecase
	result.EstimatedProfit = 0

	return &result, nil
}

// GetAllDeals lấy toàn bộ danh sách deal trong hệ thống (dành cho admin)
func (r *GroupDealPostgresRepository) GetAllDeals(ctx context.Context, searchRequest *dto.AdminDealSearchRequest) ([]*entity.Deal, int64, error) {
	var deals []*entity.Deal
	var total int64

	offset := searchRequest.Page * searchRequest.Size
	limit := searchRequest.Size

	// Build base query
	query := r.db.WithContext(ctx).Model(&entity.Deal{}).
		Where("deleted_at IS NULL")

	// Apply filters
	if searchRequest.Keyword != nil && *searchRequest.Keyword != "" {
		keyword := "%" + *searchRequest.Keyword + "%"
		query = query.Where("name ILIKE ?", keyword)
	}

	if searchRequest.Status != nil {
		query = query.Where("status = ?", *searchRequest.Status)
	}

	if searchRequest.DealType != nil {
		query = query.Where("deal_type = ?", *searchRequest.DealType)
	}

	if searchRequest.OrganizationID != nil {
		query = query.Where("owner_id = ? AND owner_type = ?", *searchRequest.OrganizationID, enums.OwnerOfOrganization)
	}

	if searchRequest.GroupID != nil {
		query = query.Where("owner_id = ? AND owner_type = ?", *searchRequest.GroupID, enums.OwnerOfGroup)
	}

	if searchRequest.OwnerID != nil {
		query = query.Where("owner_id = ?", *searchRequest.OwnerID)
	}

	if searchRequest.OwnerType != nil {
		query = query.Where("owner_type = ?", *searchRequest.OwnerType)
	}

	if searchRequest.StartDate != nil {
		query = query.Where("created_at >= ?", *searchRequest.StartDate)
	}

	if searchRequest.EndDate != nil {
		query = query.Where("created_at <= ?", *searchRequest.EndDate)
	}

	if searchRequest.MinAmount != nil {
		query = query.Where("target_profit >= ?", *searchRequest.MinAmount)
	}

	if searchRequest.MaxAmount != nil {
		query = query.Where("target_profit <= ?", *searchRequest.MaxAmount)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get deals with pagination
	query = query.Offset(offset).Limit(limit).Order("created_at DESC")

	if err := query.Find(&deals).Error; err != nil {
		return nil, 0, err
	}

	return deals, total, nil
}
