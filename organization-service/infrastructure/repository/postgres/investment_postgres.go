package postgres

import (
	"context"
	"errors"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
	"strings"
	"time"

	"gorm.io/gorm"
)

// @bind: organization/internal/domain/repository.InvestmentRepository
type InvestmentPostgres struct {
	db *gorm.DB
}

func NewInvestmentPostgres(db *gorm.DB) *InvestmentPostgres {
	return &InvestmentPostgres{db: db}
}

func (r *InvestmentPostgres) CreateInvestment(ctx context.Context, investment *entity.Investment) (*entity.Investment, error) {
	if investment.Status != enums.ApprovedStatusPending &&
		investment.Status != enums.ApprovedStatusRejected &&
		investment.Status != enums.ApprovedStatusApproved {
		investment.Status = enums.ApprovedStatusPending
	}
	err := r.db.WithContext(ctx).Create(investment).Error
	if err != nil {
		return nil, err
	}
	return investment, nil
}

func (r *InvestmentPostgres) UpdateInvestment(ctx context.Context, investment *entity.Investment) (*entity.Investment, error) {
	// Get current investment to check status
	current, err := r.GetByID(ctx, investment.ID)
	if err != nil {
		return nil, err
	}

	// Cannot update if status is approved
	if current.Status == enums.ApprovedStatusApproved {
		return nil, errors.New("cannot update approved investment")
	}

	// If rejected, reset status to pending
	if current.Status == enums.ApprovedStatusRejected {
		investment.Status = enums.ApprovedStatusPending
	}

	if investment.Status.IsValid() {
		investment.Status = current.Status
	}

	updateData := map[string]interface{}{
		"member_id":            investment.MemberID,
		"amount":               investment.Amount,
		"transfer_time":        investment.TransferTime,
		"note":                 investment.Note,
		"transfer_proof_image": investment.TransferProofImage,
		"status":               investment.Status,
		"changed_at":           time.Now(),
	}

	// Only update proxy_user_id if it's provided
	if investment.ProxyUserID != nil {
		updateData["proxy_user_id"] = investment.ProxyUserID
	}

	err = r.db.WithContext(ctx).Model(&entity.Investment{}).
		Where("id = ?", investment.ID).
		Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	// Update documents if provided
	if len(investment.AttachDocument) > 0 {
		// Delete old documents
		err = r.db.WithContext(ctx).
			Where("owner_id = ? AND doc_owner = ?", investment.ID, enums.DocOwnerInvest).
			Delete(&entity.AttachDocument{}).Error
		if err != nil {
			return nil, err
		}

		// Insert new documents
		for i := range investment.AttachDocument {
			investment.AttachDocument[i].OwnerId = investment.ID
			investment.AttachDocument[i].DocOwner = enums.DocOwnerInvest
		}
		err = r.db.WithContext(ctx).Create(&investment.AttachDocument).Error
		if err != nil {
			return nil, err
		}
	}

	return investment, nil
}

func (r *InvestmentPostgres) DeleteInvestment(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Investment{}, id).Error
}

func (r *InvestmentPostgres) GetByID(ctx context.Context, id uint64) (*entity.Investment, error) {
	var investment entity.Investment
	err := r.db.WithContext(ctx).
		Model(&entity.Investment{}).
		Where("id = ? and deleted_at is null", id).
		// Preload("AttachDocument", "doc_owner = ?", enums.DocOwnerInvest).
		First(&investment).Error
	if err != nil {
		return nil, err
	}

	// Lấy các attach documents liên quan đến investment này
	var attachDocs []entity.AttachDocument
	err = r.db.WithContext(ctx).
		Model(&entity.AttachDocument{}).
		Where("owner_id = ? AND doc_owner = ?", id, enums.DocOwnerInvest).
		Find(&attachDocs).Error
	if err != nil {
		return nil, err
	}

	investment.AttachDocument = attachDocs

	return &investment, nil
}

func (r *InvestmentPostgres) GetInvestments(ctx context.Context, req *dto.InvestmentDTO) ([]entity.Investment, int64, error) {
	var investments []entity.Investment
	var total int64

	query := GetDB(ctx, r.db).Model(&entity.Investment{}).
		Where("deal_id = ? and deleted_at is null and deleted_at is null", req.DealID)

	// Filter by member ID if provided
	if req.MemberID > 0 {
		query = query.Where("member_id = ?", req.MemberID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		// Preload("AttachDocument", "doc_owner = ?", enums.DocOwnerInvest).
		Order("updated_at DESC").
		Find(&investments).
		Error
	if err != nil {
		return nil, 0, err
	}

	return investments, total, nil
}

func (r *InvestmentPostgres) ChangeStatus(ctx context.Context, id uint64, status enums.ApprovedStatus) (uint64, error) {
	err := r.db.WithContext(ctx).Model(&entity.Investment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"changed_at": time.Now(),
		}).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *InvestmentPostgres) ChangeConfirmation(ctx context.Context, id uint64, confirmed bool) (uint64, error) {
	err := r.db.WithContext(ctx).Model(&entity.Investment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"confirmed":  confirmed,
			"changed_at": time.Now(),
		}).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *InvestmentPostgres) CountDashboard(ctx context.Context, summary *dto.SummaryRequest) (*dto.SummaryResponse, error) {
	var result dto.SummaryResponse
	var selectFields []string
	var args []interface{}

	// Build SELECT fields and arguments dynamically
	if summary.AmountPending || summary.All {
		selectFields = append(selectFields, `COALESCE(SUM(CASE WHEN status = ? THEN amount ELSE 0 END), 0) as amount_pending`)
		args = append(args, enums.ApprovedStatusPending)
	}
	if summary.AmountApproved || summary.All {
		selectFields = append(selectFields, `COALESCE(SUM(CASE WHEN status = ? THEN amount ELSE 0 END), 0) as amount_approved`)
		args = append(args, enums.ApprovedStatusApproved)
	}
	if summary.AmountRejected || summary.All {
		selectFields = append(selectFields, `COALESCE(SUM(CASE WHEN status = ? THEN amount ELSE 0 END), 0) as amount_rejected`)
		args = append(args, enums.ApprovedStatusRejected)
	}
	if summary.NumPending || summary.All {
		selectFields = append(selectFields, `COUNT(CASE WHEN status = ? THEN 1 END) as num_pending`)
		args = append(args, enums.ApprovedStatusPending)
	}
	if summary.NumApproved || summary.All {
		selectFields = append(selectFields, `COUNT(CASE WHEN status = ? THEN 1 END) as num_approved`)
		args = append(args, enums.ApprovedStatusApproved)
	}
	if summary.NumRejected || summary.All {
		selectFields = append(selectFields, `COUNT(CASE WHEN status = ? THEN 1 END) as num_rejected`)
		args = append(args, enums.ApprovedStatusRejected)
	}
	if summary.NumInvestment || summary.All {
		selectFields = append(selectFields, `COUNT(*) as num_investment`)
	}
	if summary.NumMemberSubmit || summary.All {
		selectFields = append(selectFields, `COUNT(DISTINCT member_id) as num_member_submit`)
	}
	if summary.AmountCost || summary.All {
		selectFields = append(selectFields, `COALESCE(SUM(CASE WHEN status = ? THEN amount ELSE 0 END), 0) as amount_cost`)
		args = append(args, enums.ApprovedStatusApproved)
	}
	if summary.AmountCost || summary.All {
		selectFields = append(selectFields,
			`
			(SELECT sum(dc.amount) FROM deal_costs dc
			left join cost_types ct on ct.id = dc.cost_type_id 
			WHERE deal_id = ? and ct.cost_type = 10 AND dc.deleted_at IS NULL 
		) as amount_cost`)
		args = append(args, summary.DealId)
	}

	// Nếu không có field nào được chọn, tránh query rỗng
	if len(selectFields) == 0 {
		return &result, nil
	}

	// Build query
	err := r.db.WithContext(ctx).
		Debug().
		Model(&entity.Investment{}).
		Where("deal_id = ? AND deleted_at IS NULL", summary.DealId).
		Select(strings.Join(selectFields, ", "), args...).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *InvestmentPostgres) GetConfirmedInvestmentsByDealId(ctx context.Context, dealId uint64, memberId uint64) ([]entity.Investment, error) {
	var investments []entity.Investment

	err := r.db.WithContext(ctx).
		Model(&entity.Investment{}).
		// Joins("JOIN organization_members om ON om.user_id = investments.member_id").
		Where("investments.deal_id = ? AND investments.member_id = ? AND investments.status = ? AND investments.deleted_at IS NULL",
			dealId, memberId, enums.ApprovedStatusApproved).
		Order("investments.transfer_time ASC").
		Find(&investments).Error

	if err != nil {
		return nil, err
	}

	return investments, nil
}
