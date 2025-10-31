package mapper

import (
	_utils "common/utils"
	"hub/internal/domain"
	hubpb "pb/types/hub"
)

type FAQMapper struct{}

func NewFAQMapper() *FAQMapper {
	return &FAQMapper{}
}

// CreateRequestToEntity converts CreateFAQRequest to FAQEntity
func (m *FAQMapper) CreateRequestToEntity(req *hubpb.CreateFAQRequest) *domain.FAQEntity {
	return &domain.FAQEntity{
		Question: req.GetQuestion(),
		Answer:   req.GetAnswer(),
		GroupKey: req.GetGroupKey(),
	}
}

// UpdateRequestToEntity updates FAQEntity with UpdateFAQRequest
func (m *FAQMapper) UpdateRequestToEntity(entity *domain.FAQEntity, req *hubpb.UpdateFAQRequest) {
	if req.GetQuestion() != "" {
		entity.Question = req.GetQuestion()
	}
	if req.GetAnswer() != "" {
		entity.Answer = req.GetAnswer()
	}
	if req.GetGroupKey() != "" {
		entity.GroupKey = req.GetGroupKey()
	}
}

// EntityToDetailProto converts FAQEntity to FAQDetail
func (m *FAQMapper) EntityToDetailProto(entity *domain.FAQEntity) *hubpb.FAQDetail {
	return &hubpb.FAQDetail{
		Id:        entity.ID,
		Question:  entity.Question,
		Answer:    entity.Answer,
		GroupKey:  entity.GroupKey,
		CreatedAt: _utils.FormatTimeToString(entity.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(entity.UpdatedAt),
	}
}

// EntitiesToProto converts slice of FAQEntity to slice of FAQ
func (m *FAQMapper) EntitiesToProto(entities []*domain.FAQEntity) []*hubpb.FAQ {
	var result []*hubpb.FAQ
	for _, entity := range entities {
		result = append(result, &hubpb.FAQ{
			Id:        entity.ID,
			Question:  entity.Question,
			Answer:    entity.Answer,
			GroupKey:  entity.GroupKey,
			CreatedAt: _utils.FormatTimeToString(entity.CreatedAt),
			UpdatedAt: _utils.FormatTimeToString(entity.UpdatedAt),
		})
	}
	return result
}
