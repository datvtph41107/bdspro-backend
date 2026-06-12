package usecase

import (
	_errors "common/errors"
	"context"
	"crm/config"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func (u *SeoDomainUsecase) GetInternalSeoRenderDataBySlug(ctx context.Context, slug string) (*dto.InternalSeoRenderDataResponse, error) {
	slug = strings.Trim(strings.TrimSpace(slug), "/")
	if slug == "" {
		return nil, _errors.ReturnError(400, "slug không được để trống")
	}

	seoDomain, err := u.seoDomainRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if seoDomain == nil {
		return nil, _errors.ReturnError(404, "seo domain không tồn tại")
	}
	if !seoDomain.Published {
		return nil, _errors.ReturnError(404, "seo domain chưa publish")
	}
	if !seoDomain.IsIndex {
		return nil, _errors.ReturnError(404, "seo domain không cho index")
	}

	return &dto.InternalSeoRenderDataResponse{
		SeoDomain:     seoDomain,
		Entity:        nil,
		InternalLinks: seoDomain.InternalLinks,
		Relatives:     seoDomain.Relatives,
	}, nil
}

func (u *SeoDomainUsecase) GetInternalSeoStaticPaths(ctx context.Context, page uint32, size uint32) ([]seo_domain.SeoDomain, int64, error) {
	return u.seoDomainRepo.GetPublishedStaticPaths(ctx, page, size)
}

func (u *SeoDomainUsecase) TriggerSeoGenerationJob(ctx context.Context) error {
	batchSize := config.AppProperties.Seo.GenerationJob.BatchSize
	if batchSize <= 0 {
		batchSize = 200
	}

	items, err := u.seoDomainRepo.GetNeedGenerateItems(ctx, batchSize)
	if err != nil {
		return err
	}

	for i := range items {
		item := &items[i]
		now := time.Now()
		hash := buildSeoStaticHash(item)
		staticPath := item.StaticHtmlPath
		if strings.TrimSpace(staticPath) == "" {
			staticPath = item.CanonicalURL
		}
		if err := u.seoDomainRepo.MarkGenerated(ctx, item.ID, staticPath, hash, now); err != nil {
			continue
		}
	}
	return nil
}

func buildSeoStaticHash(item *seo_domain.SeoDomain) string {
	if item == nil {
		return ""
	}
	sum := sha256.Sum256([]byte(fmt.Sprintf("%d|%s|%s|%s|%s|%s", item.ID, item.Slug, item.CanonicalURL, item.Title, item.Description, item.UpdatedAt.Format(time.RFC3339Nano))))
	return hex.EncodeToString(sum[:])
}
