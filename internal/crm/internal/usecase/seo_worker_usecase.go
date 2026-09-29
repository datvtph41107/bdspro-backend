package usecase

import (
	"context"
	"fmt"
	"strings"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/repo"
)

type SeoWorkerUsecase struct {
	seoDomainRepo    repo.SeoDomainRepo
	seoRenderUsecase *SeoRenderUsecase
}

func NewSeoWorkerUsecase(
	seoDomainRepo repo.SeoDomainRepo,
	seoRenderUsecase *SeoRenderUsecase,
) *SeoWorkerUsecase {
	return &SeoWorkerUsecase{
		seoDomainRepo:    seoDomainRepo,
		seoRenderUsecase: seoRenderUsecase,
	}
}

func (u *SeoWorkerUsecase) GeneratePendingSeoPages(ctx context.Context, req *dto.SeoWorkerGenerateRequest) (*dto.SeoWorkerGenerateResult, error) {
	if req == nil {
		req = &dto.SeoWorkerGenerateRequest{}
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	triggerType := strings.TrimSpace(req.TriggerType)
	if triggerType == "" {
		triggerType = seo_domain.SeoGenerationTriggerScheduler
	}

	items, err := u.seoDomainRepo.GetNeedGenerateItems(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := &dto.SeoWorkerGenerateResult{
		Scanned: len(items),
		Errors:  make([]string, 0),
	}

	for i := range items {
		page := &items[i]

		if err := validatePageCanGenerate(page); err != nil {
			result.Skipped++
			result.Errors = append(result.Errors, fmt.Sprintf("seoId=%d skipped: %s", page.ID, err.Error()))
			continue
		}

		if u.seoRenderUsecase == nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("seoId=%d failed: SeoRenderUsecase is nil", page.ID))
			continue
		}

		generation, err := u.seoRenderUsecase.GenerateSeoPage(ctx, &dto.GenerateSeoPageRequest{
			ID:          page.ID,
			TriggerType: triggerType,
			Reason:      req.Reason,
		})
		if err != nil || generation == nil || generation.RenderStatus != seo_domain.SeoRenderStatusSuccess {
			result.Failed++
			message := "render failed"
			if err != nil {
				message = err.Error()
			} else if generation != nil && strings.TrimSpace(generation.ErrorMessage) != "" {
				message = generation.ErrorMessage
			}
			result.Errors = append(result.Errors, fmt.Sprintf("seoId=%d failed: %s", page.ID, message))
			continue
		}

		result.Success++
	}

	return result, nil
}
