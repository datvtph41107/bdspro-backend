package repo

import (
	_dto "common/domain/dto"
	"context"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
)

type DealRepository interface {
	Create(ctx context.Context, deal *domain.Deal) (*domain.Deal, error)
	Update(ctx context.Context, deal *domain.Deal) (*domain.Deal, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.Deal, error)
	DetailByID(ctx context.Context, id uint64) (*domain.Deal, error)
	GetIdsByGroupID(ctx context.Context, groupID uint64) ([]uint64, error)
	GetIdsByOrganizationID(ctx context.Context, bdsproID uint64) ([]uint64, error)
	GetDeals(ctx context.Context, userID uint64, searchRequest *dto.DealSearchRequest) ([]*domain.Deal, uint32, error)
	UpdateStatus(ctx context.Context, id uint64, status enums.DealStatus) error
	UpdateStatusAndCancelReason(ctx context.Context, id uint64, status enums.DealStatus, cancelReason string) error
	UpdateSetting(ctx context.Context, id uint64, setting *dto.DealSetting) error
	CreateDealProduct(ctx context.Context, dealProduct []*domain.DealProduct) ([]*domain.DealProduct, error)
	GetDealProducts(ctx context.Context, dealID uint64, pagable _dto.Pagable) ([]*domain.DealProduct, error)
	UpdateAllowManualInput(ctx context.Context, dealID uint64, allowManualInput bool) error
	GetDealProcess(ctx context.Context, dealID uint64, bdsproID uint64) ([]*dto.DealProcessItem, error)

	// Thống kê thương vụ
	GetDealOverviewByOrganization(ctx context.Context, bdsproID uint64) (*dto.DealOverview, error)
	GetDealOverviewByGroup(ctx context.Context, groupID uint64) (*dto.DealOverview, error)

	// Admin API
	GetAllDeals(ctx context.Context, searchRequest *dto.AdminDealSearchRequest) ([]*domain.Deal, int64, error)
}
