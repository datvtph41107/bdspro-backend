package services

import (
	"context"
	"strconv"

	_errors "common/errors"
	hubpb "pb/types/hub"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

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
