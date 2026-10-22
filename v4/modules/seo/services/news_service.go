// Package services contains business logic for modules/seo.
package services

import (
	"context"
	"errors"

	_errors "common/errors"
	"mdn4/modules/seo/dto"
	"mdn4/modules/seo/entity"
	"mdn4/modules/seo/repository"

	"gorm.io/gorm"
)

// NewsService handles planning news, categories and saved-news flows.
type NewsService interface {
	ListPlanningNews(ctx context.Context, filter dto.PlanningNewsListFilter) (*dto.PlanningNewsListResult, error)
	GetPlanningNews(ctx context.Context, id uint64) (*dto.PlanningNewsDetailResult, error)
	RecordClick(ctx context.Context, req dto.RecordClickRequest) (*dto.RecordClickResult, error)
	ListCategories(ctx context.Context) ([]dto.ContentCategory, error)
	GetCategory(ctx context.Context, id uint64) (*dto.ContentCategory, error)
	AdminCreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.ContentCategory, error)
	AdminUpdateCategory(ctx context.Context, id uint64, req dto.UpdateCategoryRequest) (*dto.ContentCategory, error)
	AdminDeleteCategory(ctx context.Context, id uint64) error
	// Client saved news
	ListSavedPlanningNews(ctx context.Context, userID uint64, filter dto.SavedNewsFilter) (*dto.PlanningNewsListResult, error)
	SavePlanningNews(ctx context.Context, userID, newsID uint64) (*dto.PlanningNewsItem, error)
	UnsavePlanningNews(ctx context.Context, userID, newsID uint64) error
	// Content search
	SearchContent(ctx context.Context, req dto.ContentSearchRequest) (*dto.ContentSearchResult, error)
}

// @bind: mdn4/modules/seo/services.NewsService
type LocalNewsService struct {
	repo       repository.NewsRepository
	domainRepo repository.SeoDomainRepository
}

func NewLocalNewsService(repo repository.NewsRepository, domainRepo repository.SeoDomainRepository) *LocalNewsService {
	return &LocalNewsService{repo: repo, domainRepo: domainRepo}
}

func (s *LocalNewsService) ListPlanningNews(ctx context.Context, filter dto.PlanningNewsListFilter) (*dto.PlanningNewsListResult, error) {
	return s.repo.ListPlanningNews(ctx, filter)
}

func (s *LocalNewsService) GetPlanningNews(ctx context.Context, id uint64) (*dto.PlanningNewsDetailResult, error) {
	if id == 0 {
		return nil, _errors.NotFoundException("news not found")
	}
	result, err := s.repo.GetPlanningNews(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, _errors.NotFoundException("news not found")
	}
	if err != nil {
		return nil, err
	}
	if result == nil || result.SEODomain.ID == 0 {
		return nil, _errors.NotFoundException("news not found")
	}
	return result, nil
}

func (s *LocalNewsService) RecordClick(ctx context.Context, req dto.RecordClickRequest) (*dto.RecordClickResult, error) {
	var userID *uint64
	if req.UserID > 0 {
		id := req.UserID
		userID = &id
	}

	// Prefer explicit seo_domain.id (proto `id` or legacy `seoDomainId`).
	seoID := req.ID
	if seoID == 0 {
		seoID = req.SeoDomainID
	}
	if seoID > 0 {
		if s.domainRepo == nil {
			return nil, _errors.InternalServerException("seo domain repo is not configured")
		}
		page, err := s.domainRepo.GetByID(ctx, seoID)
		if err != nil {
			return nil, err
		}
		if page == nil || page.ID == 0 {
			return nil, _errors.NotFoundException("seo page not found")
		}
		if err := s.repo.RecordClick(ctx, page.ID, userID, req.IPAddr); err != nil {
			return nil, err
		}
		return &dto.RecordClickResult{Accepted: true, ID: page.ID}, nil
	}

	// Else require refType + refId (CRM TrackSeoClick parity).
	if req.RefID == 0 || req.RefType == 0 {
		return nil, _errors.BadRequestException("cần truyền id hoặc cặp refType+refId")
	}
	if !isValidSEORefType(req.RefType) {
		return nil, _errors.BadRequestException("refType không hợp lệ")
	}
	if s.domainRepo != nil {
		page, err := s.domainRepo.GetByRef(ctx, req.RefType, req.RefID)
		if err != nil {
			return nil, err
		}
		if page != nil && page.ID > 0 {
			if err := s.repo.RecordClick(ctx, page.ID, userID, req.IPAddr); err != nil {
				return nil, err
			}
		}
		// CRM accepts immediately even when the page is missing at flush time.
	}
	return &dto.RecordClickResult{
		Accepted: true,
		RefType:  req.RefType,
		RefID:    req.RefID,
	}, nil
}

func isValidSEORefType(v uint32) bool {
	switch entity.SEORefType(v) {
	case entity.SEORefTypeAdmUnit,
		entity.SEORefTypeParcel,
		entity.SEORefTypeRegion,
		entity.SEORefTypeProject,
		entity.SEORefTypeDocument,
		entity.SEORefTypeMap,
		entity.SEORefTypeReport,
		entity.SEORefTypeLegacy:
		return true
	default:
		return v == 0
	}
}

func (s *LocalNewsService) ListCategories(ctx context.Context) ([]dto.ContentCategory, error) {
	return s.repo.ListCategories(ctx)
}

func (s *LocalNewsService) GetCategory(ctx context.Context, id uint64) (*dto.ContentCategory, error) {
	if id == 0 {
		return nil, _errors.NotFoundException("category not found")
	}
	cat, err := s.repo.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, _errors.NotFoundException("category not found")
	}
	return cat, nil
}

func (s *LocalNewsService) AdminCreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.ContentCategory, error) {
	if req.Name == "" {
		return nil, _errors.BadRequestException("name is required")
	}
	if req.Code == "" {
		return nil, _errors.BadRequestException("code is required")
	}
	return s.repo.CreateCategory(ctx, req)
}

func (s *LocalNewsService) AdminUpdateCategory(ctx context.Context, id uint64, req dto.UpdateCategoryRequest) (*dto.ContentCategory, error) {
	if id == 0 {
		return nil, _errors.NotFoundException("category not found")
	}
	return s.repo.UpdateCategory(ctx, id, req)
}

func (s *LocalNewsService) AdminDeleteCategory(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.NotFoundException("category not found")
	}
	return s.repo.DeleteCategory(ctx, id)
}

func (s *LocalNewsService) ListSavedPlanningNews(ctx context.Context, userID uint64, filter dto.SavedNewsFilter) (*dto.PlanningNewsListResult, error) {
	if userID == 0 {
		return nil, _errors.BadRequestException("user id is required")
	}
	return s.repo.ListSavedPlanningNews(ctx, userID, filter)
}

func (s *LocalNewsService) SavePlanningNews(ctx context.Context, userID, newsID uint64) (*dto.PlanningNewsItem, error) {
	if userID == 0 {
		return nil, _errors.UnauthorizedException("authentication required")
	}
	if newsID == 0 {
		return nil, _errors.BadRequestException("news id is required")
	}
	return s.repo.SavePlanningNews(ctx, userID, newsID)
}

func (s *LocalNewsService) UnsavePlanningNews(ctx context.Context, userID, newsID uint64) error {
	if userID == 0 {
		return _errors.UnauthorizedException("authentication required")
	}
	if newsID == 0 {
		return _errors.BadRequestException("news id is required")
	}
	return s.repo.UnsavePlanningNews(ctx, userID, newsID)
}

func (s *LocalNewsService) SearchContent(ctx context.Context, req dto.ContentSearchRequest) (*dto.ContentSearchResult, error) {
	var kinds []string
	if req.Kind != "" {
		kinds = []string{req.Kind}
	}
	return s.repo.SearchContent(ctx, req.Query, kinds, req.Limit)
}
