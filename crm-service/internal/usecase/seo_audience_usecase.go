package usecase

import (
	_errors "common/errors"
	"context"
	"errors"
	"strings"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/dto"
	"crm/internal/repo"
)

type SeoAudienceUsecase struct{ repo repo.SeoAudienceRepo }

func NewSeoAudienceUsecase(repository repo.SeoAudienceRepo) *SeoAudienceUsecase {
	return &SeoAudienceUsecase{repo: repository}
}

func (u *SeoAudienceUsecase) Create(ctx context.Context, req *dto.CreateSeoAudienceRequest) (*seo_strategy.SeoAudienceAggregate, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không được để trống")
	}
	audience, err := seo_strategy.CreateSeoAudienceSegmentDraft(req.Code, req.Name, time.Now())
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	existing, err := u.repo.GetByCode(ctx, audience.Code)
	if err == nil && existing != nil {
		return nil, _errors.ReturnError(409, "audience code đã tồn tại")
	}
	if err != nil && !errors.Is(err, repo.ErrSeoAudienceNotFound) {
		return nil, err
	}
	if err := u.repo.Create(ctx, audience); err != nil {
		return nil, err
	}
	return u.Get(ctx, audience.ID)
}

func (u *SeoAudienceUsecase) Update(ctx context.Context, req *dto.UpdateSeoAudienceRequest) (*seo_strategy.SeoAudienceAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "audience update request không hợp lệ")
	}
	audience, err := u.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, mapAudienceRepoError(err)
	}
	if audience.Status == seo_strategy.AudienceStatusArchived {
		return nil, _errors.ReturnError(409, "không thể sửa audience đã archived")
	}
	audience.Name = strings.TrimSpace(req.Name)
	audience.Description = strings.TrimSpace(req.Description)
	audience.DecisionRole = seo_strategy.AudienceDecisionRole(req.DecisionRole)
	audience.KnowledgeLevel = seo_strategy.AudienceKnowledgeLevel(req.KnowledgeLevel)
	audience.OwnerCode = strings.TrimSpace(req.OwnerCode)
	audience.ReviewDueAt = req.ReviewDueAt
	if err := seo_strategy.ValidateAudienceDefinition(audience); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	if err := u.repo.UpdateDefinition(ctx, audience, req.ExpectedVersion); err != nil {
		return nil, mapAudienceRepoError(err)
	}
	return u.Get(ctx, audience.ID)
}

func (u *SeoAudienceUsecase) Get(ctx context.Context, id uint64) (*seo_strategy.SeoAudienceAggregate, error) {
	if id == 0 {
		return nil, _errors.ReturnError(400, "audience id không hợp lệ")
	}
	audience, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapAudienceRepoError(err)
	}
	return &seo_strategy.SeoAudienceAggregate{
		Audience:  audience,
		Readiness: seo_strategy.EvaluateAudienceReadiness(audience, time.Now()),
	}, nil
}

func (u *SeoAudienceUsecase) List(ctx context.Context, req *dto.SeoAudienceListRequest) ([]seo_strategy.SeoAudienceSegment, int64, error) {
	if req == nil {
		req = &dto.SeoAudienceListRequest{}
	}
	return u.repo.List(ctx, repo.SeoAudienceFilter{
		Keyword: req.Keyword, Status: req.Status, DecisionRole: req.DecisionRole, Page: req.Page, Size: req.Size,
	})
}

func (u *SeoAudienceUsecase) Activate(ctx context.Context, req *dto.SeoAudienceCommandRequest) (*seo_strategy.SeoAudienceAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "audience command không hợp lệ")
	}
	aggregate, err := u.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if aggregate.Audience.Status == seo_strategy.AudienceStatusArchived {
		return nil, _errors.ReturnError(409, "không thể activate audience đã archived")
	}
	if !aggregate.Readiness.Ready {
		codes := make([]string, 0, len(aggregate.Readiness.Blockers))
		for _, item := range aggregate.Readiness.Blockers {
			codes = append(codes, item.Code)
		}
		return nil, _errors.ReturnError(422, "audience chưa sẵn sàng: "+strings.Join(codes, ", "))
	}
	if err := u.repo.UpdateStatus(ctx, req.ID, seo_strategy.AudienceStatusActive, req.ExpectedVersion); err != nil {
		return nil, mapAudienceRepoError(err)
	}
	return u.Get(ctx, req.ID)
}

func (u *SeoAudienceUsecase) Pause(ctx context.Context, req *dto.SeoAudienceCommandRequest) (*seo_strategy.SeoAudienceAggregate, error) {
	return u.changeStatus(ctx, req, seo_strategy.AudienceStatusPaused)
}

func (u *SeoAudienceUsecase) Archive(ctx context.Context, req *dto.SeoAudienceCommandRequest) (*seo_strategy.SeoAudienceAggregate, error) {
	return u.changeStatus(ctx, req, seo_strategy.AudienceStatusArchived)
}

func (u *SeoAudienceUsecase) changeStatus(ctx context.Context, req *dto.SeoAudienceCommandRequest, status seo_strategy.AudienceStatus) (*seo_strategy.SeoAudienceAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "audience command không hợp lệ")
	}
	if err := u.repo.UpdateStatus(ctx, req.ID, status, req.ExpectedVersion); err != nil {
		return nil, mapAudienceRepoError(err)
	}
	return u.Get(ctx, req.ID)
}

func mapAudienceRepoError(err error) error {
	switch {
	case errors.Is(err, repo.ErrSeoAudienceNotFound):
		return _errors.ReturnError(404, "seo audience không tồn tại")
	case errors.Is(err, repo.ErrSeoAudienceConflict):
		return _errors.ReturnError(409, "audience đã được cập nhật; hãy tải lại dữ liệu")
	default:
		return err
	}
}
