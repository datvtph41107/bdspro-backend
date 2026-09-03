package application

import (
	"context"

	"crm/internal/modules/contentcatalog/domain"
)

// Repository is the application-facing catalog persistence contract.
type Repository interface {
	ListCategories(context.Context) ([]domain.ContentCategory, error)
	GetCategory(context.Context, uint64) (*domain.ContentCategory, error)
	CreateCategory(context.Context, domain.CreateSeoCategoryInput) (*domain.ContentCategory, error)
	UpdateCategory(context.Context, uint64, domain.UpdateSeoCategoryInput) (*domain.ContentCategory, error)
	DeleteCategory(context.Context, uint64) error
	AssignPlanningNewsCategory(context.Context, uint64, uint64) (*domain.PlanningNewsItem, error)
	ListPlanningNews(context.Context, domain.PlanningNewsListFilter) (*domain.PlanningNewsListResult, error)
	GetPlanningNews(context.Context, uint64) (*domain.PlanningNewsDetailResult, error)
	ListSavedPlanningNews(context.Context, uint64, domain.SavedPlanningNewsListFilter) (*domain.PlanningNewsListResult, error)
	SavePlanningNews(context.Context, uint64, uint64) (*domain.PlanningNewsItem, error)
	UnsavePlanningNews(context.Context, uint64, uint64) error
	SearchPublishedContent(context.Context, string, []string, int) (*domain.ContentSearchResult, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) ListCategories(ctx context.Context) ([]domain.ContentCategory, error) {
	return s.repository.ListCategories(ctx)
}
func (s *Service) GetCategory(ctx context.Context, id uint64) (*domain.ContentCategory, error) {
	return s.repository.GetCategory(ctx, id)
}
func (s *Service) CreateCategory(ctx context.Context, input domain.CreateSeoCategoryInput) (*domain.ContentCategory, error) {
	return s.repository.CreateCategory(ctx, input)
}
func (s *Service) UpdateCategory(ctx context.Context, id uint64, input domain.UpdateSeoCategoryInput) (*domain.ContentCategory, error) {
	return s.repository.UpdateCategory(ctx, id, input)
}
func (s *Service) DeleteCategory(ctx context.Context, id uint64) error {
	return s.repository.DeleteCategory(ctx, id)
}
func (s *Service) AssignPlanningNewsCategory(ctx context.Context, id, categoryID uint64) (*domain.PlanningNewsItem, error) {
	return s.repository.AssignPlanningNewsCategory(ctx, id, categoryID)
}

func (s *Service) ListPlanningNews(ctx context.Context, filter domain.PlanningNewsListFilter) (*domain.PlanningNewsListResult, error) {
	return s.repository.ListPlanningNews(ctx, filter)
}
func (s *Service) GetPlanningNews(ctx context.Context, id uint64) (*domain.PlanningNewsDetailResult, error) {
	return s.repository.GetPlanningNews(ctx, id)
}
func (s *Service) ListSavedPlanningNews(ctx context.Context, userID uint64, filter domain.SavedPlanningNewsListFilter) (*domain.PlanningNewsListResult, error) {
	return s.repository.ListSavedPlanningNews(ctx, userID, filter)
}
func (s *Service) SavePlanningNews(ctx context.Context, userID, newsID uint64) (*domain.PlanningNewsItem, error) {
	return s.repository.SavePlanningNews(ctx, userID, newsID)
}
func (s *Service) UnsavePlanningNews(ctx context.Context, userID, newsID uint64) error {
	return s.repository.UnsavePlanningNews(ctx, userID, newsID)
}
func (s *Service) SearchPublishedContent(ctx context.Context, query string, kinds []string, limit int) (*domain.ContentSearchResult, error) {
	return s.repository.SearchPublishedContent(ctx, query, kinds, limit)
}
