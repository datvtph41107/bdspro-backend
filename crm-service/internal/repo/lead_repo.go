package repo

import (
	base_enum "base/enum"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
)

type LeadRepo interface {
	Search(c context.Context, ownerId uint64, ownerType base_enum.EOwnerOf, dto dto.LeadSearchDTO) ([]domain.LeadEntity, int64, error)
	// SearchByOwnerOf lists leads whose contact has the given owner_of (shared pool; no owner_id filter).
	SearchByOwnerOf(c context.Context, ownerType base_enum.EOwnerOf, dto dto.LeadSearchDTO) ([]domain.LeadEntity, int64, error)
	AdminOpportunitySummary(c context.Context, ownerType base_enum.EOwnerOf) (*dto.AdminOpportunitiesSummaryDTO, error)
	AdminFunnel(c context.Context, ownerType base_enum.EOwnerOf, search dto.LeadSearchDTO) ([]dto.AdminFunnelBucketDTO, error)
	Create(c context.Context, entity *domain.LeadEntity) (*domain.LeadEntity, error)
	Update(c context.Context, id uint64, entity *domain.LeadEntity) (*domain.LeadEntity, error)
	Delete(c context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.LeadEntity, error)
	GetByContactID(ctx context.Context, contactID uint64) (*domain.LeadEntity, error)
	// DetailByID(ctx context.Context, contactID uint64) (*domain.LeadEntity, error)
	Existed(ctx context.Context, entity *domain.LeadEntity) (bool, error)
	ExistedWithStageID(ctx context.Context, stageID uint64) (bool, error)
	ExistedWithPipelineID(ctx context.Context, pipelineID uint64) (bool, error)
	UpdateNote(ctx context.Context, id uint64, note string) (*domain.LeadEntity, error)
	Assign(c context.Context, id uint64, dto *dto.LeadDTO) (*domain.LeadEntity, error)
	SwitchStage(c context.Context, id uint64, stageID *uint64, note string) (*domain.LeadEntity, error)
	GetOwnerWithRule(c context.Context) ([]dto.CustomerWithRuleDTO, error)
	// UpdateAutoTrigger(c context.Context, id uint64) (*domain.CustomerEntity, error)
	CareExpiredTime(c context.Context, organizationId uint64) (int64, error)
}