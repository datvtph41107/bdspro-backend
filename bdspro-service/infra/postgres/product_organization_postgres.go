package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductOrganizationRepo
type ProductOrganizationPostgres struct {
	db *gorm.DB
}

func NewProductOfOrgPostgres(db *gorm.DB) *ProductOrganizationPostgres {
	return &ProductOrganizationPostgres{db: db}
}

func (r *ProductOrganizationPostgres) CreateOwner(ctx context.Context, productID uint64, organizationID uint64) (*domain.ProductOrganization, error) {

	productOfOrg := &domain.ProductOrganization{
		ProductID:      productID,
		OrganizationID: organizationID,
		IsOwner:        true,
	}
	err := r.db.WithContext(ctx).Create(productOfOrg).Error
	if err != nil {
		return nil, err
	}
	return productOfOrg, nil
}

func (r *ProductOrganizationPostgres) GetByProductID(ctx context.Context, productID uint64) ([]*domain.ProductOrganization, error) {
	var productOfOrgs []*domain.ProductOrganization
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND deleted_at IS NULL", productID).
		Find(&productOfOrgs).Error
	if err != nil {
		return nil, err
	}
	return productOfOrgs, nil
}

func (r *ProductOrganizationPostgres) GetByOrganizationID(ctx context.Context, organizationID uint64) ([]*domain.ProductOrganization, error) {
	var productOfOrgs []*domain.ProductOrganization
	err := r.db.WithContext(ctx).
		Where("organization_id = ? AND deleted_at IS NULL", organizationID).
		Find(&productOfOrgs).Error
	if err != nil {
		return nil, err
	}
	return productOfOrgs, nil
}

func (r *ProductOrganizationPostgres) GetByProductAndOrganization(ctx context.Context, productID, organizationID uint64) (*domain.ProductOrganization, error) {
	var productOfOrg domain.ProductOrganization
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND organization_id = ? AND deleted_at IS NULL", productID, organizationID).
		First(&productOfOrg).Error
	if err != nil {
		return nil, err
	}
	return &productOfOrg, nil
}

func (r *ProductOrganizationPostgres) Update(ctx context.Context, productOfOrg *domain.ProductOrganization) (*domain.ProductOrganization, error) {
	err := r.db.WithContext(ctx).Save(productOfOrg).Error
	if err != nil {
		return nil, err
	}
	return productOfOrg, nil
}

func (r *ProductOrganizationPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.ProductOrganization{}, id).Error
}

func (r *ProductOrganizationPostgres) DeleteByProductID(ctx context.Context, productID uint64) error {
	return r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Delete(&domain.ProductOrganization{}).Error
}
