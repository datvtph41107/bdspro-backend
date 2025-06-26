package postgres

import (
	"context"
	"log"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
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
	err := r.db.WithContext(ctx).Create(investment).Error
	if err != nil {
		return nil, err
	}
	return investment, nil
}

func (r *InvestmentPostgres) UpdateInvestment(ctx context.Context, investment *entity.Investment) (*entity.Investment, error) {
	return nil, nil
}

func (r *InvestmentPostgres) DeleteInvestment(ctx context.Context, id uint64) error {
	return nil
}

func (r *InvestmentPostgres) GetByID(ctx context.Context, id uint64) (*entity.Investment, error) {
	var investment entity.Investment
	err := r.db.WithContext(ctx).
		Model(&entity.Investment{}).
		Where("id = ? and deleted_at is null", id).
		Preload("AttachDocument", "doc_owner = ?", enums.DocOwnerInvest).
		First(&investment).Error
	if err != nil {
		return nil, err
	}
	return &investment, nil
}

func (r *InvestmentPostgres) GetInvestments(ctx context.Context, req *dto.InvestmentDTO) ([]entity.Investment, int64, error) {
	var investments []entity.Investment
	var total int64

	err := r.db.WithContext(ctx).
		Model(&entity.Investment{}).
		Where("deal_id = ? and deleted_at is null", req.DealID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Model(&entity.Investment{}).
		Where("deal_id = ? and deleted_at is null", req.DealID).
		Preload("AttachDocument", "doc_owner = ?", enums.DocOwnerInvest).
		Find(&investments).Error
	if err != nil {
		return nil, 0, err
	}

	return investments, total, nil
}

func (r *InvestmentPostgres) ChangeStatus(ctx context.Context, id uint64, status enums.ApprovedStatus) (uint64, error) {
	err := r.db.WithContext(ctx).
		Model(&entity.Investment{}).
		Where("id = ? and deleted_at is null", id).
		Update("status", status).
		Update("changed_at", time.Now()).Error
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *InvestmentPostgres) CountDashboard(ctx context.Context, req *dto.InvestmentDTO) (*dto.InfoInvestmentDTO, error) {
	var result dto.InfoInvestmentDTO

	err := r.db.WithContext(ctx).Raw(`
	SELECT
		COUNT(*) AS num_investment,
		SUM(CASE WHEN status = 10 THEN amount ELSE 0 END) AS amount_pending,
		SUM(CASE WHEN status = 20 THEN amount ELSE 0 END) AS amount_approved,
		SUM(CASE WHEN status = 30 THEN amount ELSE 0 END) AS amount_rejected,
		COUNT(CASE WHEN status = 10 THEN 1 ELSE NULL END) AS num_pending,
		COUNT(CASE WHEN status = 20 THEN 1 ELSE NULL END) AS num_approved,
		COUNT(CASE WHEN status = 30 THEN 1 ELSE NULL END) AS num_rejected,
		COUNT(DISTINCT id) AS total_investment,
		COUNT(DISTINCT created_by) AS num_member_submit,
		0 AS amount_target,
		0 AS percent_done
	FROM investments
	WHERE deal_id = ? and deleted_at is null
`, req.DealID).Scan(&result).Error

	if err != nil {
		log.Println("query investment stats failed:", err)
	}

	return &result, nil
}
