package services

import (
	"context"
	"strconv"

	_errors "common/errors"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/entity"
	"mdn4/modules/hub/repository"
	hubpb "pb/types/hub"
)

// FAQService owns user-facing and admin FAQ flows for modules/hub.
type FAQService interface {
	GetFAQSimple(ctx context.Context, req dto.GetFAQSimpleRequest) (*hubpb.GetFAQSimpleResponse, error)
	GetFAQSimpleDetail(ctx context.Context, id uint64) (*hubpb.FAQDetail, error)
	AdminCreate(ctx context.Context, req dto.CreateFAQRequest) (*hubpb.FAQDetail, error)
	AdminUpdate(ctx context.Context, id uint64, req dto.CreateFAQRequest) (*hubpb.FAQDetail, error)
	AdminDelete(ctx context.Context, id uint64) error
}

// @bind: mdn4/modules/hub/services.FAQService
type LocalFAQService struct {
	repo repository.FAQRepository
}

func NewLocalFAQService(repo repository.FAQRepository) *LocalFAQService {
	return &LocalFAQService{repo: repo}
}

func (s *LocalFAQService) GetFAQSimple(ctx context.Context, req dto.GetFAQSimpleRequest) (*hubpb.GetFAQSimpleResponse, error) {
	list, total, err := s.repo.ListSimple(ctx, req)
	if err != nil {
		return nil, err
	}
	data := make([]*hubpb.FAQSimple, 0, len(list))
	for _, item := range list {
		data = append(data, &hubpb.FAQSimple{
			Id:       item.ID,
			Question: item.Question,
			Answer:   item.Answer,
		})
	}
	return &hubpb.GetFAQSimpleResponse{Data: data, Total: total}, nil
}

func (s *LocalFAQService) GetFAQSimpleDetail(ctx context.Context, id uint64) (*hubpb.FAQDetail, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("faq not found")
	}
	return mapFAQDetail(row), nil
}

func (s *LocalFAQService) AdminCreate(ctx context.Context, req dto.CreateFAQRequest) (*hubpb.FAQDetail, error) {
	if req.Question == "" || req.Answer == "" {
		return nil, _errors.BadRequestException("question and answer are required")
	}
	row := &entity.FAQ{Question: req.Question, Answer: req.Answer, GroupKey: req.GroupKey}
	if err := s.repo.Create(ctx, row); err != nil {
		return nil, err
	}
	return mapFAQDetail(row), nil
}

func (s *LocalFAQService) AdminUpdate(ctx context.Context, id uint64, req dto.CreateFAQRequest) (*hubpb.FAQDetail, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, _errors.NotFoundException("faq not found")
	}
	existing.Question = req.Question
	existing.Answer = req.Answer
	existing.GroupKey = req.GroupKey
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return mapFAQDetail(existing), nil
}

func (s *LocalFAQService) AdminDelete(ctx context.Context, id uint64) error {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return _errors.NotFoundException("faq not found")
	}
	return s.repo.Delete(ctx, id)
}

func mapFAQDetail(row *entity.FAQ) *hubpb.FAQDetail {
	out := &hubpb.FAQDetail{
		Id:       row.ID,
		Question: row.Question,
		Answer:   row.Answer,
		GroupKey: row.GroupKey,
	}
	if row.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(row.CreatedAt.UnixMilli(), 10)
	}
	if row.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(row.UpdatedAt.UnixMilli(), 10)
	}
	return out
}
