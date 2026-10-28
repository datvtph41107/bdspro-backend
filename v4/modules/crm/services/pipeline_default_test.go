package services

import (
	"context"
	"testing"

	_enums "common/domain/enum"
	_models "common/models"
	domain "mdn4/modules/crm/entity"
	repo "mdn4/modules/crm/repository"
)

type stubDefaultPipelineRepo struct {
	repo.PipelineRepo
	gotOrg   uint64
	pipeline *domain.PipelineEntity
}

func (s *stubDefaultPipelineRepo) GetDefault(_ context.Context, organizationId uint64) (*domain.PipelineEntity, error) {
	s.gotOrg = organizationId
	return s.pipeline, nil
}

type stubPipelineStageRepo struct {
	repo.StageRepo
	stages []domain.StageEntity
}

func (s stubPipelineStageRepo) GetByPipelineID(context.Context, uint64) ([]domain.StageEntity, error) {
	return s.stages, nil
}

func TestPipelineGetDefault_ProfileWithoutOrganization(t *testing.T) {
	pipelines := &stubDefaultPipelineRepo{pipeline: &domain.PipelineEntity{ID: 1, PipelineName: "Mặc định"}}
	stages := []domain.StageEntity{{BaseEntity: _models.BaseEntity{ID: 10}}, {BaseEntity: _models.BaseEntity{ID: 11}}}
	uc := NewPipelineUsecase(pipelines, nil, nil, nil, stubPipelineStageRepo{stages: stages}, nil)

	ctx := context.WithValue(context.Background(), _enums.ProfileIDKey, uint64(42))
	p, stage, err := uc.GetDefault(ctx)
	if err != nil {
		t.Fatalf("GetDefault: %v", err)
	}
	if pipelines.gotOrg != 0 || p.ID != 1 {
		t.Fatalf("org=%d pipeline=%+v", pipelines.gotOrg, p)
	}
	if stage == nil || stage.ID != 10 || len(p.Stages) != 2 {
		t.Fatalf("stage=%+v stages=%d", stage, len(p.Stages))
	}
}

func TestPipelineGetDefault_NoStages(t *testing.T) {
	pipelines := &stubDefaultPipelineRepo{pipeline: &domain.PipelineEntity{ID: 1}}
	uc := NewPipelineUsecase(pipelines, nil, nil, nil, stubPipelineStageRepo{}, nil)

	ctx := context.WithValue(context.Background(), _enums.ProfileIDKey, uint64(42))
	if _, stage, err := uc.GetDefault(ctx); err != nil || stage != nil {
		t.Fatalf("stage=%+v err=%v", stage, err)
	}
}

func TestPipelineGetDefault_RequiresToken(t *testing.T) {
	uc := NewPipelineUsecase(&stubDefaultPipelineRepo{}, nil, nil, nil, stubPipelineStageRepo{}, nil)
	if _, _, err := uc.GetDefault(context.Background()); err == nil {
		t.Fatal("expected unauthorized without profile")
	}
}
