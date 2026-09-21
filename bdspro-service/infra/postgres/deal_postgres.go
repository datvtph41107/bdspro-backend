package postgres

import (
	"bdspro/internal"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DealPostgresRepository struct {
	db              *gorm.DB
	bankAccountRepo repo.BankAccountRepository
	transactionRepo provider.TransactionProvider
}

func NewGroupDealPostgresRepository(db *gorm.DB,
	bankAccountRepo repo.BankAccountRepository,
	transactionRepo provider.TransactionProvider,
) repo.DealRepository {
	return &DealPostgresRepository{
		db:              db,
		bankAccountRepo: bankAccountRepo,
		transactionRepo: transactionRepo,
	}
}

func (r *DealPostgresRepository) IsAcceptedMember(
	ctx context.Context,
	userID uint64,
	dealID uint64,
) (bool, error) {
	var exists bool
	err := GetDB(ctx, r.db).
		Raw(`
			SELECT EXISTS (
				SELECT 1
				FROM deal_members
				WHERE deal_id = ?
				  AND member_id = ?
				  AND status = ?
				  AND deleted_at IS NULL
			)
		`,
			dealID,
			userID,
			domain.DealMemberStatusAccepted,
		).
		Scan(&exists).Error

	return exists, err
}

func (r *DealPostgresRepository) ExistsProductInDeal(
	ctx context.Context,
	dealID uint64,
	productID uint64,
) (bool, error) {

	var exists bool
	err := GetDB(ctx, r.db).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM deal_products
			WHERE deal_id = ?
			  AND product_id = ?
			  AND deleted_at IS NULL
		)
	`, dealID, productID).Scan(&exists).Error

	return exists, err
}

func (r *DealPostgresRepository) AddProductToDeal(
	ctx context.Context,
	dealID, productID, userID uint64,
) error {
	dp := &domain.DealProduct{
		DealID:    dealID,
		ProductID: productID,
	}

	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "deal_id"},
				{Name: "product_id"},
			},
			DoNothing: true,
		}).
		Create(dp)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return _errors.ReturnError(service.DealProductAlreadyExists)
	}

	return nil
}

func (r *DealPostgresRepository) GetDealsWithoutProduct(
	ctx context.Context,
	userID uint64,
	productID uint64,
	pagable _dto.Pagable,
) ([]*domain.Deal, int64, error) {
	var (
		items []*domain.Deal
		total int64
	)

	baseQuery := r.db.WithContext(ctx).
		Table("deals").
		Select(`
			deals.id,
			deals.name,
			deals.status,
			deals.owner_id,
			deals.deal_type
		`).
		Joins(`
			LEFT JOIN deal_members dm 
				ON dm.deal_id = deals.id
				AND dm.member_id = ?
				AND dm.status = ?
				AND dm.deleted_at IS NULL
		`, userID, domain.DealMemberStatusAccepted).
		Joins(`
			LEFT JOIN deal_products dp 
				ON dp.deal_id = deals.id
		`).
		Where("dp.id IS NULL").
		Where("deals.deleted_at IS NULL")

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := baseQuery.
		Limit(pagable.GetLimit()).
		Offset(pagable.GetOffset()).
		Order("deals.created_at DESC").
		Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *DealPostgresRepository) CreateDealProduct(ctx context.Context, dealProduct []*domain.DealProduct) ([]*domain.DealProduct, error) {
	// Sử dụng ON CONFLICT DO NOTHING để bỏ qua các bản ghi duplicate
	if err := GetDB(ctx, r.db).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "deal_id"},
				{Name: "product_id"},
			},
			DoNothing: true,
		}).
		Create(&dealProduct).Error; err != nil {
		return nil, err
	}
	return dealProduct, nil
}

func (r *DealPostgresRepository) GetDealProducts(ctx context.Context, dealID uint64, pagable _dto.Pagable) ([]*domain.DealProduct, error) {
	var dealProducts []*domain.DealProduct
	if err := r.db.WithContext(ctx).
		Where("deal_id = ?", dealID).
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&dealProducts).Error; err != nil {
		return nil, err
	}
	return dealProducts, nil
}

func (r *DealPostgresRepository) GetDealsByProductID(
	ctx context.Context,
	productID uint64,
	pagable _dto.Pagable,
) ([]*domain.Deal, int64, *time.Time, error) {

	var deals []*domain.Deal
	var total int64
	var productUpdatedAt *time.Time

	db := GetDB(ctx, r.db)

	if err := db.
		Model(&domain.DealProduct{}).
		Joins("JOIN deals ON deals.id = deal_products.deal_id").
		Where("deal_products.product_id = ? AND deals.deleted_at IS NULL", productID).
		Count(&total).Error; err != nil {
		return nil, 0, nil, err
	}

	if err := db.
		Model(&domain.Product{}).
		Select("updated_at").
		Where("id = ?", productID).
		Scan(&productUpdatedAt).Error; err != nil {
		return nil, 0, nil, err
	}

	if err := db.
		Model(&domain.Deal{}).
		Joins("JOIN deal_products dp ON dp.deal_id = deals.id").
		Where(`
			dp.product_id = ? 
			AND deals.deleted_at IS NULL
		`, productID).
		Preload("BankAccount", "deleted_at IS NULL").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Order("deals.created_at DESC").
		Find(&deals).Error; err != nil {
		return nil, 0, nil, err
	}

	return deals, total, productUpdatedAt, nil
}

func (r *DealPostgresRepository) CountDealsByProductID(ctx context.Context, productID uint64) (int64, error) {
	var count int64
	err := GetDB(ctx, r.db).
		Model(&domain.DealProduct{}).
		Joins("JOIN deals ON deals.id = deal_products.deal_id").
		Where("deal_products.product_id = ? AND deals.deleted_at IS NULL", productID).
		Count(&count).Error
	return count, err
}

func (r *DealPostgresRepository) Create(ctx context.Context, deal *domain.Deal) (*domain.Deal, error) {
	// model := GroupDealEntityToModel(deal)
	deal.Products = make([]domain.DealProduct, len(deal.ProductIds))
	for i, productID := range deal.ProductIds {
		deal.Products[i] = domain.DealProduct{
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

func (r *DealPostgresRepository) Update(ctx context.Context, deal *domain.Deal) (*domain.Deal, error) {
	if err := r.db.WithContext(ctx).Model(&domain.Deal{}).Where("id = ?", deal.ID).Updates(deal).Error; err != nil {
		return nil, err
	}
	return deal, nil
}

func (r *DealPostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&domain.Deal{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
		}).Error
}

func (r *DealPostgresRepository) GetByID(ctx context.Context, id uint64) (*domain.Deal, error) {
	var model domain.Deal
	if err := GetDB(ctx, r.db).
		Preload("BankAccount", "deleted_at is null").
		Where("id = ? AND deleted_at is null", id).
		First(&model).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *DealPostgresRepository) DetailByID(ctx context.Context, id uint64) (*domain.Deal, error) {
	var model domain.Deal
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

func (r *DealPostgresRepository) GetIdsByGroupID(ctx context.Context, groupID uint64) ([]uint64, error) {
	var deals []uint64
	// Get deal ids
	countQuery := `
		SELECT d.id
		FROM deals d 
		WHERE d.owner_id = ? 
		AND d.owner_type = ? 
		AND d.deleted_at IS NULL
	`
	if err := GetDB(ctx, r.db).Raw(countQuery, groupID, enums.EOwnerOfGroup).Pluck("id", &deals).Error; err != nil {
		return nil, err
	}

	return deals, nil
}

func (r *DealPostgresRepository) GetIdsByOrganizationID(ctx context.Context, organizationID uint64) ([]uint64, error) {
	var deals []uint64
	// Get deal ids
	countQuery := `
		SELECT d.id
		FROM deals d 
		WHERE d.owner_id = ? 
		AND d.owner_type = ? 
		AND d.deleted_at IS NULL
	`
	if err := GetDB(ctx, r.db).Raw(countQuery, organizationID, enums.EOwnerOfOrganization).Pluck("id", &deals).Error; err != nil {
		return nil, err
	}

	return deals, nil
}

func (r *DealPostgresRepository) GetDeals(ctx context.Context, userID uint64, searchRequest *dto.DealSearchRequest) ([]*domain.Deal, uint32, error) {
	var deals []*domain.Deal
	var total int64

	offset := searchRequest.GetOffset()
	limit := searchRequest.GetLimit()

	// Build base query starting from deal_members with status = 20 (Accepted), then join with deals
	baseQuery := `
		SELECT DISTINCT d.* 
		FROM deal_members dm
		INNER JOIN deals d ON dm.deal_id = d.id
		WHERE dm.member_id = ? 
		AND dm.status = ?
		AND d.deleted_at IS NULL
		AND dm.deleted_at IS NULL
	`

	args := []interface{}{userID, domain.DealMemberStatusAccepted}

	if searchRequest.OwnerID != nil && searchRequest.OwnerType != 0 {
		baseQuery += " AND d.owner_id = ? AND d.owner_type = ?"
		args = append(args, *searchRequest.OwnerID, searchRequest.OwnerType)
	}
	// Add organization filter if provided
	if searchRequest.OrganizationID != nil {
		baseQuery += " AND d.owner_id = ? AND d.owner_type = ?"
		args = append(args, *searchRequest.OrganizationID, enums.EOwnerOfOrganization)
	}

	// Add group filter if provided
	if searchRequest.GroupID != nil {
		baseQuery += " AND d.owner_id = ? AND d.owner_type = ?"
		args = append(args, *searchRequest.GroupID, enums.EOwnerOfGroup)
	}

	// Add text search if provided
	if searchRequest.Text != "" {
		baseQuery += " AND d.name ILIKE ?"
		args = append(args, "%"+searchRequest.Text+"%")
	}

	// Count total - build count query separately
	countQuery := `
		SELECT COUNT(DISTINCT d.id)
		FROM deal_members dm
		INNER JOIN deals d ON dm.deal_id = d.id
		WHERE dm.member_id = ? 
		AND dm.status = ?
		AND d.deleted_at IS NULL
		AND dm.deleted_at IS NULL
	`
	countArgs := []interface{}{userID, domain.DealMemberStatusAccepted}

	if searchRequest.OwnerID != nil && searchRequest.OwnerType != 0 {
		countQuery += " AND d.owner_id = ? AND d.owner_type = ?"
		countArgs = append(countArgs, *searchRequest.OwnerID, searchRequest.OwnerType)
	}
	if searchRequest.OrganizationID != nil {
		countQuery += " AND d.owner_id = ? AND d.owner_type = ?"
		countArgs = append(countArgs, *searchRequest.OrganizationID, enums.EOwnerOfOrganization)
	}
	if searchRequest.GroupID != nil {
		countQuery += " AND d.owner_id = ? AND d.owner_type = ?"
		countArgs = append(countArgs, *searchRequest.GroupID, enums.EOwnerOfGroup)
	}
	if searchRequest.Text != "" {
		countQuery += " AND d.name ILIKE ?"
		countArgs = append(countArgs, "%"+searchRequest.Text+"%")
	}

	if err := r.db.WithContext(ctx).Raw(countQuery, countArgs...).Scan(&total).Error; err != nil {
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

func (r *DealPostgresRepository) UpdateStatus(ctx context.Context, id uint64, status enums.DealStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Deal{}).Where("id = ? AND deleted_at is null", id).Update("status", status).Error
}

func (r *DealPostgresRepository) UpdateStatusAndCancelReason(ctx context.Context, id uint64, status enums.DealStatus, cancelReason string) error {
	return r.db.WithContext(ctx).Model(&domain.Deal{}).
		Where("id = ? AND deleted_at is null", id).
		Updates(map[string]interface{}{
			"status":        status,
			"cancel_reason": cancelReason,
		}).Error
}

func (r *DealPostgresRepository) UpdateSetting(ctx context.Context, id uint64, setting *dto.DealSetting) error {
	// tx := r.db.WithContext(ctx).Begin()
	err := r.transactionRepo.WithTransaction(ctx, func(ctx context.Context) error {
		var model *domain.Deal
		err := GetDB(ctx, r.db).Model(&domain.Deal{}).
			Preload("BankAccount", "deleted_at is null").
			Where("id = ? AND deleted_at is null", id).
			First(&model).Error
		if err != nil {
			return err
		}

		bankAccount, err := r.bankAccountRepo.DealUpdateAccount(ctx, model.BankAccountID, &domain.BankAccount{
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
			Model(&domain.Deal{}).
			Where("id = ? AND deleted_at is null", id).
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
				"share_visibility":           setting.ShareVisibility,
				"description":                setting.Description,
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

func (r *DealPostgresRepository) UpdateAllowManualInput(ctx context.Context, dealID uint64, allowManualInput bool) error {
	return r.db.WithContext(ctx).Model(&domain.Deal{}).Where("id = ?", dealID).Update("allow_manual_input", allowManualInput).Error
}

// GetDealProcess lấy thông tin process của deal
func (r *DealPostgresRepository) GetDealProcess(ctx context.Context, dealID uint64, organizationID uint64) ([]*dto.DealProcessItem, error) {
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
func (r *DealPostgresRepository) GetDealOverviewByOrganization(ctx context.Context, organizationID uint64) (*dto.DealOverview, error) {
	var result dto.DealOverview

	// Lấy tổng số thương vụ, tổng vốn, đang xử lý, hoàn tất
	err := GetDB(ctx, r.db).
		Debug().
		Model(&domain.Deal{}).
		Select(`
			COUNT(*) as total_deals,
			COALESCE(SUM(target_profit), 0) as total_capital,
			SUM(CASE WHEN status = 10 THEN 1 ELSE 0 END) as total_deal_processing,
			SUM(CASE WHEN status = 20 THEN 1 ELSE 0 END) as total_deal_completed,
			SUM(CASE WHEN status = 30 THEN 1 ELSE 0 END) as total_deal_canceled
		`).
		Where("owner_id = ? AND owner_type = ? AND deleted_at IS NULL", organizationID, enums.EOwnerOfOrganization).
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
func (r *DealPostgresRepository) GetDealOverviewByGroup(ctx context.Context, groupID uint64) (*dto.DealOverview, error) {
	var result dto.DealOverview

	// Lấy tổng số thương vụ, tổng vốn, đang xử lý, hoàn tất
	err := r.db.WithContext(ctx).
		Debug().
		Model(&domain.Deal{}).
		Select(`
			COUNT(*) as total_deals,
			COALESCE(SUM(target_profit), 0) as total_capital,
			SUM(CASE WHEN status = 10 THEN 1 ELSE 0 END) as total_deal_processing,
			SUM(CASE WHEN status = 20 THEN 1 ELSE 0 END) as total_deal_completed,
			SUM(CASE WHEN status = 30 THEN 1 ELSE 0 END) as total_deal_canceled
		`).
		Where("owner_id = ? AND owner_type = ? AND deleted_at IS NULL", groupID, enums.EOwnerOfGroup).
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
func (r *DealPostgresRepository) GetAllDeals(ctx context.Context, searchRequest *dto.AdminDealSearchRequest) ([]*domain.Deal, int64, error) {
	var deals []*domain.Deal
	var total int64

	offset := searchRequest.Page * searchRequest.Size
	limit := searchRequest.Size

	// Build base query
	query := r.db.WithContext(ctx).Model(&domain.Deal{}).
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
		query = query.Where("owner_id = ? AND owner_type = ?", *searchRequest.OrganizationID, enums.EOwnerOfOrganization)
	}

	if searchRequest.GroupID != nil {
		query = query.Where("owner_id = ? AND owner_type = ?", *searchRequest.GroupID, enums.EOwnerOfGroup)
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

func (r *DealPostgresRepository) GetByIds(ctx context.Context, ids []uint64) ([]*domain.Deal, error) {
	var deals []*domain.Deal
	if err := r.db.WithContext(ctx).Where("id IN (?) AND deleted_at IS NULL", ids).Find(&deals).Error; err != nil {
		return nil, err
	}
	return deals, nil
}
