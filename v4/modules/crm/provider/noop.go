package provider

import (
	"context"

	base_enum "base/enum"
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/crm/enums"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

// @bind: mdn4/modules/crm/provider.ITransaction
type NoopTransaction struct{}

func NewNoopTransaction() *NoopTransaction { return &NoopTransaction{} }

func (NoopTransaction) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// @bind: mdn4/modules/crm/provider.NotificationProvider
type NoopNotification struct{}

func NewNoopNotification() *NoopNotification { return &NoopNotification{} }

func (NoopNotification) CreateHistory(context.Context, *_dto.HistoryDTO) error { return nil }
func (NoopNotification) CreateNotification(context.Context, string, string, []string, _enum.ENotificationType, *uint64, uint64, _enum.EOwnerOf, []string) error {
	return nil
}
func (NoopNotification) RemoveNotification(context.Context, _enum.EOwnerOf, uint64, uint64, _enum.ENotificationType) error {
	return nil
}

type NoopUserClient struct{}

func NewNoopUserClient() *NoopUserClient { return &NoopUserClient{} }

func (NoopUserClient) GetProfileById(context.Context, uint64) (*domain.Profile, error) {
	return &domain.Profile{}, nil
}
func (NoopUserClient) GetProfileWithPhoneById(context.Context, uint64) (*domain.Profile, error) {
	return &domain.Profile{}, nil
}
func (NoopUserClient) GetProfileByPhones(context.Context, []string) (map[string]*domain.Profile, error) {
	return map[string]*domain.Profile{}, nil
}
func (NoopUserClient) ValidateProfileIds(context.Context, []uint64) (bool, error) { return true, nil }
func (NoopUserClient) GetProfileByIds(context.Context, []uint64) ([]*sharepb.ProfileItem, error) {
	return nil, nil
}

// @bind: mdn4/modules/crm/provider.BdsproProvider
type NoopBdspro struct{}

func NewNoopBdspro() *NoopBdspro { return &NoopBdspro{} }

func (NoopBdspro) GetDealById(context.Context, uint64) (*bdspropb.GroupDeal, error) { return nil, nil }
func (NoopBdspro) GetDealsByIds(context.Context, []uint64) ([]*bdspropb.GroupDeal, error) {
	return nil, nil
}
func (NoopBdspro) GetProductAttachmentByIds(context.Context, []uint64) (*bdspropb.ProductAttachmentResponse, error) {
	return nil, nil
}
func (NoopBdspro) GetDealContractsByIds(context.Context, []uint64) ([]*bdspropb.DealContractItem, error) {
	return nil, nil
}
func (NoopBdspro) UpdateScheduleCount(context.Context, uint64, uint32) error { return nil }
func (NoopBdspro) UpdateContactCount(context.Context, uint64, uint32) error  { return nil }
func (NoopBdspro) SearchProjects(context.Context, string, uint32, uint32) ([]*bdspropb.Project, int64, error) {
	return nil, 0, nil
}
func (NoopBdspro) SearchRegions(context.Context, string, uint32, uint32) ([]*bdspropb.Region, int64, error) {
	return nil, 0, nil
}
func (NoopBdspro) GetRegionByID(context.Context, uint64) (*bdspropb.Region, error) { return nil, nil }
func (NoopBdspro) ListAreaRegions(context.Context) ([]*bdspropb.AreaRegion, error) { return nil, nil }
func (NoopBdspro) GetAreaRegionByID(context.Context, uint64) (*bdspropb.AreaRegion, error) {
	return nil, nil
}

// @bind: mdn4/modules/crm/services.PermissionUsecase
type AllowAllPermission struct{}

func NewAllowAllPermission() *AllowAllPermission { return &AllowAllPermission{} }

func (AllowAllPermission) UserCanAccessTarget(context.Context, uint64, base_enum.EOwnerOf, uint64, int8) error {
	return nil
}
func (AllowAllPermission) UserInOwner(context.Context, uint64, base_enum.EOwnerOf) error { return nil }
func (AllowAllPermission) HasRoleWithOwner(context.Context, uint64, base_enum.EOwnerOf, enums.EAuth) error {
	return nil
}
