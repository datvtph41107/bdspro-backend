package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.CustomerRepo
type PostgreCustomer struct {
	db *gorm.DB
}

func NewPostgreCustomer(db *gorm.DB) *PostgreCustomer {
	return &PostgreCustomer{db: db}
}

func (r *PostgreCustomer) Search(c context.Context, dto dto.CustomerSearchDTO) ([]domain.CustomerEntity, int64, error) {
	return nil, 0, nil
}

func (r *PostgreCustomer) Create(c context.Context, entity *domain.CustomerEntity) (*domain.CustomerEntity, error) {
	return nil, nil
}

func (r *PostgreCustomer) Update(ctx context.Context, id uint64, customer *domain.CustomerEntity) (*domain.CustomerEntity, error) {
	return nil, nil
}

func (r *PostgreCustomer) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (r *PostgreCustomer) GetByID(ctx context.Context, id uint64) (*domain.CustomerEntity, error) {
	return nil, nil
}

func (r *PostgreCustomer) Existed(c context.Context, entity *domain.CustomerEntity) (bool, error) {
	return false, nil
}
