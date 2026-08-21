package application

import (
	"context"

	"crm/internal/modules/contentcatalog/domain"
)

// Repository is the application-facing catalog persistence contract.
type Repository interface {
	ListDomains(context.Context) ([]domain.ContentDomain, error)
	ListCategories(context.Context, string) ([]domain.ContentCategory, error)
	GetCategory(context.Context, uint64) (*domain.ContentCategory, error)
	CreateCategory(context.Context, domain.CreateSeoCategoryInput) (*domain.ContentCategory, error)
	UpdateCategory(context.Context, uint64, domain.UpdateSeoCategoryInput) (*domain.ContentCategory, error)
	DeleteCategory(context.Context, uint64) error
	ListFeed(context.Context, string, string, string, int) (*domain.FeedResult, error)
	ListPlanningNews(context.Context, domain.PlanningNewsListFilter) (*domain.PlanningNewsListResult, error)
	GetPlanningNews(context.Context, string) (*domain.PlanningNewsDetailResult, error)
	ListSavedPlanningNews(context.Context, uint64, domain.SavedPlanningNewsListFilter) (*domain.PlanningNewsListResult, error)
	SavePlanningNews(context.Context, uint64, uint64) (*domain.PlanningNewsItem, error)
	UnsavePlanningNews(context.Context, uint64, uint64) error
	ListFollowedSeoCategories(context.Context, uint64, domain.FollowedSeoCategoryListFilter) (*domain.FollowedSeoCategoryListResult, error)
	FollowSeoCategory(context.Context, uint64, uint64) (*domain.ContentCategory, error)
	UnfollowSeoCategory(context.Context, uint64, uint64) error
	UpdatePlanningNewsVisibility(context.Context, uint64, uint32) (*domain.PlanningNewsItem, error)
	AssignPlanningNewsCategory(context.Context, uint64, uint64) (*domain.PlanningNewsItem, error)
	SetPlanningNewsTop(context.Context, uint64) (*domain.PlanningNewsItem, error)
	SearchPublishedContent(context.Context, string, []string, int) (*domain.ContentSearchResult, error)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func (s *Service) ListDomains(ctx context.Context) ([]domain.ContentDomain, error) {
	return s.repository.ListDomains(ctx)
}
func (s *Service) ListCategories(ctx context.Context, domainCode string) ([]domain.ContentCategory, error) {
	return s.repository.ListCategories(ctx, domainCode)
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
func (s *Service) ListFeed(ctx context.Context, feed, domainCode, categoryCode string, limit int) (*domain.FeedResult, error) {
	return s.repository.ListFeed(ctx, feed, domainCode, categoryCode, limit)
}
func (s *Service) ListPlanningNews(ctx context.Context, filter domain.PlanningNewsListFilter) (*domain.PlanningNewsListResult, error) {
	return s.repository.ListPlanningNews(ctx, filter)
}
func (s *Service) GetPlanningNews(ctx context.Context, slug string) (*domain.PlanningNewsDetailResult, error) {
	return s.repository.GetPlanningNews(ctx, slug)
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
func (s *Service) ListFollowedSeoCategories(ctx context.Context, userID uint64, filter domain.FollowedSeoCategoryListFilter) (*domain.FollowedSeoCategoryListResult, error) {
	return s.repository.ListFollowedSeoCategories(ctx, userID, filter)
}
func (s *Service) FollowSeoCategory(ctx context.Context, userID, categoryID uint64) (*domain.ContentCategory, error) {
	return s.repository.FollowSeoCategory(ctx, userID, categoryID)
}
func (s *Service) UnfollowSeoCategory(ctx context.Context, userID, categoryID uint64) error {
	return s.repository.UnfollowSeoCategory(ctx, userID, categoryID)
}
func (s *Service) UpdatePlanningNewsVisibility(ctx context.Context, id uint64, visibleOn uint32) (*domain.PlanningNewsItem, error) {
	return s.repository.UpdatePlanningNewsVisibility(ctx, id, visibleOn)
}
func (s *Service) AssignPlanningNewsCategory(ctx context.Context, id uint64, categoryID uint64) (*domain.PlanningNewsItem, error) {
	return s.repository.AssignPlanningNewsCategory(ctx, id, categoryID)
}
func (s *Service) SetPlanningNewsTop(ctx context.Context, id uint64) (*domain.PlanningNewsItem, error) {
	return s.repository.SetPlanningNewsTop(ctx, id)
}
func (s *Service) SearchPublishedContent(ctx context.Context, query string, kinds []string, limit int) (*domain.ContentSearchResult, error) {
	return s.repository.SearchPublishedContent(ctx, query, kinds, limit)
}
