package postgres

import (
	"time"

	"context"

	"github.com/tuandq2112/payment-service/internal/domain/entity"
	"github.com/tuandq2112/payment-service/internal/domain/repository"
	"gorm.io/gorm"
	paymentpb "pb/types/payment"
)

type PaymentMethodModel struct {
	ID             uint32 `gorm:"primaryKey"`
	OrganizationID uint32
	Name           string
	Code           string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (m *PaymentMethodModel) TableName() string {
	return "payment_methods"
}

func PaymentMethodModelToEntity(m *PaymentMethodModel) *entity.PaymentMethod {
	return &entity.PaymentMethod{
		Id:             m.ID,
		OrganizationId: m.OrganizationID,
		Name:           m.Name,
		Code:           m.Code,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func PaymentMethodEntityToModel(e *entity.PaymentMethod) *PaymentMethodModel {
	return &PaymentMethodModel{
		ID:             e.Id,
		OrganizationID: e.OrganizationId,
		Name:           e.Name,
		Code:           e.Code,
		IsActive:       e.IsActive,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

type paymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) repository.PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

func (r *paymentMethodRepository) GetPaymentMethods(ctx context.Context, organizationId uint32, page int, size int) ([]*paymentpb.PaymentMethod, int64, error) {
	var models []*PaymentMethodModel
	var total int64

	query := r.db.WithContext(ctx).Where("organization_id = ?", organizationId)
	if err := query.Limit(size).Offset(page * size).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	paymentMethods := make([]*paymentpb.PaymentMethod, len(models))
	for i, model := range models {
		paymentMethods[i] = &paymentpb.PaymentMethod{
			Id:             model.ID,
			OrganizationId: model.OrganizationID,
			Name:           model.Name,
			Code:           model.Code,
			IsActive:       model.IsActive,
		}
		paymentMethods[i].CreatedAt = model.CreatedAt.Format(time.RFC3339)
		paymentMethods[i].UpdatedAt = model.UpdatedAt.Format(time.RFC3339)
	}
	return paymentMethods, total, nil
}

func (r *paymentMethodRepository) CreatePaymentMethod(ctx context.Context, req *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	model := PaymentMethodEntityToModel(req)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return PaymentMethodModelToEntity(model), nil
}

func (r *paymentMethodRepository) UpdatePaymentMethod(ctx context.Context, req *entity.PaymentMethod) (*entity.PaymentMethod, error) {
	model := PaymentMethodEntityToModel(req)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, err
	}
	return PaymentMethodModelToEntity(model), nil
}

func (r *paymentMethodRepository) DeletePaymentMethod(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&PaymentMethodModel{}, id).Error
}

func (r *paymentMethodRepository) GetPaymentMethodById(ctx context.Context, id uint32) (*entity.PaymentMethod, error) {
	var model PaymentMethodModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return PaymentMethodModelToEntity(&model), nil
}

