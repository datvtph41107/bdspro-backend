package usecase

import (
	"context"

	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
)

type BusinessDomainUsecase interface {
	CreateBusinessDomain(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error)
	GetBusinessDomain(ctx context.Context, id uint32) (*entity.BusinessDomain, error)
	GetAllBusinessDomains(ctx context.Context) ([]*entity.BusinessDomain, error)
	GetActiveBusinessDomains(ctx context.Context) ([]*entity.BusinessDomain, error)
	UpdateBusinessDomain(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error)
	DeleteBusinessDomain(ctx context.Context, id uint32) error
	GetBusinessDomainsByIds(ctx context.Context, ids []uint32) ([]*entity.BusinessDomain, error)
}

type businessDomainUsecase struct {
	businessDomainRepository repository.BusinessDomainRepository
}

func NewBusinessDomainUsecase(businessDomainRepository repository.BusinessDomainRepository) BusinessDomainUsecase {
	return &businessDomainUsecase{
		businessDomainRepository: businessDomainRepository,
	}
}

func (u *businessDomainUsecase) CreateBusinessDomain(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error) {
	// Check if code already exists
	existingDomain, err := u.businessDomainRepository.FindByCode(ctx, businessDomain.Code)
	if err != nil {
		return nil, err
	}
	if existingDomain != nil {
		return nil, err // You might want to create a specific error for this
	}

	return u.businessDomainRepository.Create(ctx, businessDomain)
}

func (u *businessDomainUsecase) GetBusinessDomain(ctx context.Context, id uint32) (*entity.BusinessDomain, error) {
	return u.businessDomainRepository.FindById(ctx, id)
}

func (u *businessDomainUsecase) GetAllBusinessDomains(ctx context.Context) ([]*entity.BusinessDomain, error) {
	return u.businessDomainRepository.FindAll(ctx)
}

func (u *businessDomainUsecase) GetActiveBusinessDomains(ctx context.Context) ([]*entity.BusinessDomain, error) {
	return u.businessDomainRepository.FindActive(ctx)
}

func (u *businessDomainUsecase) UpdateBusinessDomain(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error) {
	// Check if business domain exists
	existingDomain, err := u.businessDomainRepository.FindById(ctx, businessDomain.ID)
	if err != nil {
		return nil, err
	}
	if existingDomain == nil {
		return nil, err // You might want to create a specific error for this
	}

	// Check if code already exists for another domain
	if businessDomain.Code != existingDomain.Code {
		domainWithCode, err := u.businessDomainRepository.FindByCode(ctx, businessDomain.Code)
		if err != nil {
			return nil, err
		}
		if domainWithCode != nil {
			return nil, err // You might want to create a specific error for this
		}
	}

	return u.businessDomainRepository.Update(ctx, businessDomain)
}

func (u *businessDomainUsecase) DeleteBusinessDomain(ctx context.Context, id uint32) error {
	return u.businessDomainRepository.Delete(ctx, id)
}

func (u *businessDomainUsecase) GetBusinessDomainsByIds(ctx context.Context, ids []uint32) ([]*entity.BusinessDomain, error) {
	return u.businessDomainRepository.FindByIds(ctx, ids)
}
